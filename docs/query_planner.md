

👉 **`QUERY_PLANNER_IR.md`**

# Universal Data Viewer (UDV)
## Query Planner – Internal Representation (IR) Specification

---

## 1. Purpose of the Query Planner IR

The Query Planner Internal Representation (IR) is a **normalized, executable query model** derived from the validated Query DSL.

It exists to:
- Decouple DSL semantics from SQL generation
- Centralize join resolution
- Enforce safety limits
- Enable deterministic query compilation
- Support multiple database adapters

The IR is **database-agnostic** and **immutable once built**.

---

## 2. Position in the Execution Pipeline

```

Query DSL (JSON)
↓
DSL Validation
↓
Query Planner
↓
Query Planner IR   ← THIS DOCUMENT
↓
SQL Generator (Postgres / MySQL / Mongo)
↓
Database Execution

````

---

## 3. Core Design Principles

- Fully resolved (no symbolic references)
- Explicit joins
- No user input strings
- Stable ordering
- Adapter-friendly

---

## 4. High-Level IR Structure

```go
type QueryPlan struct {
    RootModel    ModelRef
    Select       []SelectExpr
    Joins        []JoinPlan
    Filters      FilterExpr
    GroupBy      []GroupExpr
    Aggregates   []AggregateExpr
    Sort         []SortExpr
    Pagination   Pagination
    Limits       QueryLimits
}
````

Each component is described below.

---

## 5. Model Reference

### 5.1 ModelRef

Represents the root model of the query.

```go
type ModelRef struct {
    Name      string
    Table     string
    Alias     string
    PrimaryKey ColumnRef
}
```

Rules:

* Alias is planner-generated (e.g. `t0`)
* Alias uniqueness guaranteed
* Root model always has alias `t0`

---

## 6. Column & Field Resolution

### 6.1 ColumnRef

All fields are resolved into column references.

```go
type ColumnRef struct {
    TableAlias string
    ColumnName string
    DataType   FieldType
}
```

No string field paths remain at IR level.

---

## 7. SELECT Clause Representation

### 7.1 SelectExpr

```go
type SelectExpr struct {
    Column     ColumnRef
    Alias      string
    IsAggregate bool
}
```

Rules:

* All selected fields must be explicitly listed
* Aggregates are excluded here and handled separately
* Alias defaults to column name if omitted

---

## 8. JOIN Resolution

### 8.1 JoinPlan

Each join is explicitly defined.

```go
type JoinPlan struct {
    Type         JoinType        // LEFT, INNER
    FromAlias   string
    ToTable     string
    ToAlias     string
    On          JoinCondition
}
```

---

### 8.2 JoinCondition

```go
type JoinCondition struct {
    Left  ColumnRef
    Right ColumnRef
}
```

---

### 8.3 Join Rules

* Joins are derived strictly from schema relationships
* Planner generates deterministic aliases (t1, t2, ...)
* Max join depth enforced here
* Duplicate joins are deduplicated

---

## 9. Filter Expression Tree

### 9.1 FilterExpr (Recursive)

```go
type FilterExpr interface{}
```

---

### 9.2 Logical Filters

```go
type LogicalFilter struct {
    Op    LogicalOp   // AND, OR, NOT
    Nodes []FilterExpr
}
```

---

### 9.3 Comparison Filter

```go
type ComparisonFilter struct {
    Left     ColumnRef
    Operator FilterOp
    Right    ValueExpr
}
```

---

### 9.4 ValueExpr

```go
type ValueExpr struct {
    Value any
    Type  FieldType
}
```

Rules:

* All values are strongly typed
* No raw user input at SQL generation stage
* Arrays allowed only for IN/BETWEEN

---

## 10. Group By Representation

### 10.1 GroupExpr

```go
type GroupExpr struct {
    Column ColumnRef
}
```

Rules:

* Order matters
* Group depth validated earlier
* Fields must be groupable

---

## 11. Aggregate Representation

### 11.1 AggregateExpr

```go
type AggregateExpr struct {
    Function AggregateFn   // SUM, COUNT, AVG, MIN, MAX
    Column   *ColumnRef    // nil for COUNT(*)
    Alias    string
}
```

Rules:

* Aggregates are separate from SelectExpr
* Alias required for sorting and result mapping

---

## 12. Sorting Representation

### 12.1 SortExpr

```go
type SortExpr struct {
    Target    SortTarget   // COLUMN or AGGREGATE
    Column    *ColumnRef
    Aggregate *AggregateExpr
    Direction SortDirection // ASC, DESC
}
```

Rules:

* Sort target must exist
* Aggregates referenced by alias
* Default direction = ASC

---

## 13. Pagination

### 13.1 Pagination

```go
type Pagination struct {
    Limit  int
    Offset int
}
```

Rules:

* Limit always present
* Offset defaults to 0
* Max limits enforced earlier

---

## 14. Query Limits

### 14.1 QueryLimits

```go
type QueryLimits struct {
    MaxJoins      int
    MaxGroupBy    int
    MaxAggregates int
}
```

Limits are attached for observability and safety enforcement.

---

## 15. IR Construction Algorithm (High-Level)

1. Resolve root model → assign alias `t0`
2. Resolve all field paths → ColumnRefs
3. Build JOIN graph → deduplicate joins
4. Normalize filters → logical tree
5. Attach GROUP BY expressions
6. Attach aggregates
7. Validate sorting targets
8. Attach pagination
9. Freeze IR (immutable)

---

## 16. Guarantees Provided by IR

Once constructed, the IR guarantees:

* No invalid fields
* No invalid joins
* No unsafe operators
* No missing pagination
* No schema ambiguity

SQL generators **must not perform validation**, only translation.

---

## 17. Adapter Responsibilities

Adapters must:

* Respect join order
* Respect aliases
* Preserve filter tree semantics
* Preserve grouping hierarchy
* Parameterize all values

Adapters must NOT:

* Modify IR
* Add joins
* Infer relationships

---

## 18. Testing Strategy

IR-level tests should include:

* DSL → IR snapshot tests
* Join resolution tests
* Filter tree correctness tests
* Aggregate + group correctness
* Sorting on aggregates

This enables database-independent correctness testing.

---

## 19. Why This IR Matters

The IR:

* Makes the system debuggable
* Enables future optimizations
* Allows non-SQL backends
* Keeps frontend and backend loosely coupled

This is the **most important abstraction** in UDV.

---

## 20. Summary

The Query Planner IR is a **fully-resolved, immutable execution plan**.

It is:

* Safe
* Deterministic
* Database-agnostic
* Adapter-friendly

All correctness in UDV depends on this layer.


