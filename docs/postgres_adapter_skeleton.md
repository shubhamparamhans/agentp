👉 **`POSTGRES_ADAPTER_SKELETON.md`**

---

## 1. Adapter Package Structure

```text
internal/adapter/
├── adapter.go                # DB adapter interface
└── postgres/
    ├── adapter.go            # PostgresAdapter struct
    ├── builder.go            # IR → SQL builder
    ├── executor.go           # Query execution
    ├── mapper.go             # Row mapping
    └── types.go              # Helper types
```

---

## 2. Adapter Interface (Database-Agnostic)

### `internal/adapter/adapter.go`

```go
package adapter

import "udv/internal/ir"

type Row map[string]any

type Adapter interface {
	BuildQuery(plan ir.QueryPlan) (string, []any, error)
	Execute(query string, args []any) ([]Row, error)
	Close() error
}
```

> Rule: **All adapters implement this interface**

---

## 3. Postgres Adapter Entry Point

### `internal/adapter/postgres/adapter.go`

```go
package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"

	"udv/internal/adapter"
)

type PostgresAdapter struct {
	db *sql.DB
}

func New(connString string) (*PostgresAdapter, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return &PostgresAdapter{db: db}, nil
}

func (p *PostgresAdapter) Close() error {
	return p.db.Close()
}
```

---

## 4. SQL Builder (IR → SQL)

### `internal/adapter/postgres/builder.go`

```go
package postgres

import (
	"fmt"
	"strings"

	"udv/internal/ir"
)

type sqlBuilder struct {
	args []any
}

func (b *sqlBuilder) nextArg(value any) string {
	b.args = append(b.args, value)
	return fmt.Sprintf("$%d", len(b.args))
}

func (p *PostgresAdapter) BuildQuery(plan ir.QueryPlan) (string, []any, error) {
	builder := &sqlBuilder{}

	var sb strings.Builder

	// SELECT
	sb.WriteString("SELECT ")
	sb.WriteString(buildSelectClause(plan))
	sb.WriteString("\n")

	// FROM
	sb.WriteString("FROM ")
	sb.WriteString(buildFromClause(plan))
	sb.WriteString("\n")

	// JOINS
	for _, join := range plan.Joins {
		sb.WriteString(buildJoinClause(join))
		sb.WriteString("\n")
	}

	// WHERE
	if plan.Filters != nil {
		whereSQL := buildWhereClause(plan.Filters, builder)
		if whereSQL != "" {
			sb.WriteString("WHERE ")
			sb.WriteString(whereSQL)
			sb.WriteString("\n")
		}
	}

	// GROUP BY
	if len(plan.GroupBy) > 0 {
		sb.WriteString("GROUP BY ")
		sb.WriteString(buildGroupByClause(plan))
		sb.WriteString("\n")
	}

	// ORDER BY
	if len(plan.Sort) > 0 {
		sb.WriteString("ORDER BY ")
		sb.WriteString(buildOrderByClause(plan))
		sb.WriteString("\n")
	}

	// PAGINATION
	sb.WriteString("LIMIT ")
	sb.WriteString(builder.nextArg(plan.Pagination.Limit))
	sb.WriteString(" OFFSET ")
	sb.WriteString(builder.nextArg(plan.Pagination.Offset))

	return sb.String(), builder.args, nil
}
```

---

## 5. Clause Builders

### `internal/adapter/postgres/types.go`

```go
package postgres

import "udv/internal/ir"

func buildSelectClause(plan ir.QueryPlan) string {
	parts := []string{}

	for _, sel := range plan.Select {
		col := sel.Column.TableAlias + "." + sel.Column.ColumnName
		if sel.Alias != "" {
			col += " AS " + sel.Alias
		}
		parts = append(parts, col)
	}

	for _, agg := range plan.Aggregates {
		if agg.Column == nil {
			parts = append(parts, "COUNT(*) AS "+agg.Alias)
		} else {
			col := agg.Column.TableAlias + "." + agg.Column.ColumnName
			parts = append(parts,
				string(agg.Function)+"("+col+") AS "+agg.Alias,
			)
		}
	}

	return join(parts)
}

func buildFromClause(plan ir.QueryPlan) string {
	return plan.RootModel.Table + " " + plan.RootModel.Alias
}

func buildJoinClause(join ir.JoinPlan) string {
	return "LEFT JOIN " + join.ToTable + " " + join.ToAlias +
		" ON " +
		join.On.Left.TableAlias + "." + join.On.Left.ColumnName +
		" = " +
		join.On.Right.TableAlias + "." + join.On.Right.ColumnName
}

func join(parts []string) string {
	return strings.Join(parts, ", ")
}
```

---

## 6. WHERE Clause Builder (Recursive)

### `internal/adapter/postgres/where.go`

```go
package postgres

import (
	"strings"

	"udv/internal/ir"
)

func buildWhereClause(expr ir.FilterExpr, b *sqlBuilder) string {
	switch f := expr.(type) {

	case ir.LogicalFilter:
		parts := []string{}
		for _, node := range f.Nodes {
			parts = append(parts, buildWhereClause(node, b))
		}
		return "(" + strings.Join(parts, " "+string(f.Op)+" ") + ")"

	case ir.ComparisonFilter:
		left := f.Left.TableAlias + "." + f.Left.ColumnName
		arg := b.nextArg(f.Right.Value)
		return left + " " + string(f.Operator) + " " + arg

	default:
		return ""
	}
}
```

---

## 7. GROUP BY & ORDER BY

### `internal/adapter/postgres/group_sort.go`

```go
package postgres

import "udv/internal/ir"

func buildGroupByClause(plan ir.QueryPlan) string {
	parts := []string{}
	for _, g := range plan.GroupBy {
		parts = append(parts,
			g.Column.TableAlias+"."+g.Column.ColumnName,
		)
	}
	return join(parts)
}

func buildOrderByClause(plan ir.QueryPlan) string {
	parts := []string{}
	for _, s := range plan.Sort {
		if s.Target == ir.SortTargetAggregate {
			parts = append(parts, s.AggregateAlias+" "+string(s.Direction))
		} else {
			col := s.Column.TableAlias + "." + s.Column.ColumnName
			parts = append(parts, col+" "+string(s.Direction))
		}
	}
	return join(parts)
}
```

---

## 8. Query Execution

### `internal/adapter/postgres/executor.go`

```go
package postgres

import (
	"context"
	"time"

	"udv/internal/adapter"
)

func (p *PostgresAdapter) Execute(query string, args []any) ([]adapter.Row, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return mapRows(rows)
}
```

---

## 9. Row Mapping

### `internal/adapter/postgres/mapper.go`

```go
package postgres

import (
	"database/sql"

	"udv/internal/adapter"
)

func mapRows(rows *sql.Rows) ([]adapter.Row, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []adapter.Row{}

	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))

		for i := range values {
			ptrs[i] = &values[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		row := adapter.Row{}
		for i, col := range cols {
			row[col] = values[i]
		}

		results = append(results, row)
	}

	return results, nil
}
```

---

## 10. Why This Skeleton Is Correct

✅ SQL lives **only** in adapter
✅ Planner IR is untouched
✅ Parameterization enforced
✅ Easy to test (`IR → SQL` golden tests)
✅ Easy to extend (MySQL adapter mirrors structure)

---

## 11. What to Implement Next (Recommended Order)

1. **Golden tests**: DSL → IR → SQL snapshots
2. Filter operators (`IN`, `BETWEEN`, `ILIKE`)
3. NULL handling
4. JOIN deduplication
5. Error wrapping & tracing

---
