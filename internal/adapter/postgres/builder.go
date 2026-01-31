package postgres

// Package postgres implements PostgreSQL-specific query generation

import (
	"fmt"
	"strings"

	"udv/internal/dsl"
	"udv/internal/planner"
)

// getPostgreSQLType returns the PostgreSQL type cast string for a given FieldType
func getPostgreSQLType(fieldType planner.FieldType) string {
	switch fieldType {
	case planner.TypeUUID:
		return "uuid"
	case planner.TypeJSON:
		return "jsonb"
	case planner.TypeBinary:
		return "bytea"
	case planner.TypeTimestamp, planner.TypeDate, planner.TypeDateTime:
		return "timestamp"
	// For other types, PostgreSQL can usually infer from context
	default:
		return ""
	}
}

// needsTypeCasting returns true if the field type needs explicit type casting in SQL
func needsTypeCasting(fieldType planner.FieldType) bool {
	switch fieldType {
	case planner.TypeUUID, planner.TypeJSON, planner.TypeBinary, planner.TypeTimestamp:
		return true
	default:
		return false
	}
}

// addTypeCast adds PostgreSQL type casting to a parameterized value if needed
func addTypeCast(paramPlaceholder string, fieldType planner.FieldType) string {
	pgType := getPostgreSQLType(fieldType)
	if pgType == "" {
		return paramPlaceholder
	}
	return fmt.Sprintf("$%s::%s", strings.TrimPrefix(paramPlaceholder, "$"), pgType)
}

// QueryBuilder builds parameterized PostgreSQL queries from query plans
type QueryBuilder struct {
	params     []interface{}
	paramCount int
}

// BuildQuery converts a QueryPlan into a parameterized SQL query
func (qb *QueryBuilder) BuildQuery(plan *planner.QueryPlan) (string, []interface{}, error) {
	if plan == nil {
		return "", nil, fmt.Errorf("query plan is nil")
	}

	if plan.RootModel == nil {
		return "", nil, fmt.Errorf("root model is nil")
	}

	qb.params = []interface{}{}
	qb.paramCount = 0

	// Route to appropriate builder based on operation
	operation := plan.Operation
	if operation == "" {
		operation = "select" // Default to select
	}

	switch operation {
	case "create":
		return qb.buildInsert(plan)
	case "update":
		return qb.buildUpdate(plan)
	case "delete":
		return qb.buildDelete(plan)
	case "select", "":
		return qb.buildSelect(plan)
	default:
		return "", nil, fmt.Errorf("unsupported operation: %s", operation)
	}
}

// buildSelect builds a SELECT query (existing logic)
func (qb *QueryBuilder) buildSelect(plan *planner.QueryPlan) (string, []interface{}, error) {
	var parts []string

	// 1. SELECT clause
	selectPart := qb.buildSelectClause(plan)
	parts = append(parts, selectPart)

	// 2. FROM clause
	fromPart := qb.buildFromClause(plan)
	parts = append(parts, fromPart)

	// 3. WHERE clause (if filters exist)
	if plan.Filters != nil {
		wherePart, err := qb.buildWhereClause(plan.Filters)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, wherePart)
	}

	// 4. GROUP BY clause (if grouping exists)
	if len(plan.GroupBy) > 0 {
		groupByPart := qb.buildGroupByClause(plan)
		parts = append(parts, groupByPart)
	}

	// 5. ORDER BY clause (if sorting exists)
	if len(plan.Sort) > 0 {
		orderByPart := qb.buildOrderByClause(plan)
		parts = append(parts, orderByPart)
	}

	// 6. LIMIT/OFFSET clause
	paginationPart := qb.buildPaginationClause(plan)
	parts = append(parts, paginationPart)

	// Join all parts
	sql := strings.Join(parts, " ") + ";"

	return sql, qb.params, nil
}

// buildInsert builds an INSERT query
func (qb *QueryBuilder) buildInsert(plan *planner.QueryPlan) (string, []interface{}, error) {
	if plan.Data == nil || len(plan.Data) == 0 {
		return "", nil, fmt.Errorf("data is required for insert operation")
	}

	table := plan.RootModel.Table
	fields := []string{}
	placeholders := []string{}

	for field, value := range plan.Data {
		fields = append(fields, field)
		qb.paramCount++
		placeholders = append(placeholders, fmt.Sprintf("$%d", qb.paramCount))
		qb.params = append(qb.params, value)
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *;",
		table,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "),
	)

	return sql, qb.params, nil
}

// buildUpdate builds an UPDATE query
func (qb *QueryBuilder) buildUpdate(plan *planner.QueryPlan) (string, []interface{}, error) {
	if plan.Data == nil || len(plan.Data) == 0 {
		return "", nil, fmt.Errorf("data is required for update operation")
	}

	table := plan.RootModel.Table
	sets := []string{}

	for field, value := range plan.Data {
		qb.paramCount++
		sets = append(sets, fmt.Sprintf("%s = $%d", field, qb.paramCount))
		qb.params = append(qb.params, value)
	}

	where := ""
	if plan.ID != nil {
		// Use ID for WHERE clause
		qb.paramCount++
		where = fmt.Sprintf("WHERE %s = $%d", plan.RootModel.PrimaryKey.ColumnName, qb.paramCount)
		qb.params = append(qb.params, plan.ID)
	} else if plan.Filters != nil {
		// Use filters for WHERE clause
		wherePart, err := qb.buildWhereClause(plan.Filters)
		if err != nil {
			return "", nil, err
		}
		where = wherePart
	} else {
		return "", nil, fmt.Errorf("id or filters required for update operation")
	}

	sql := fmt.Sprintf(
		"UPDATE %s SET %s %s RETURNING *;",
		table,
		strings.Join(sets, ", "),
		where,
	)

	return sql, qb.params, nil
}

// buildDelete builds a DELETE query
func (qb *QueryBuilder) buildDelete(plan *planner.QueryPlan) (string, []interface{}, error) {
	table := plan.RootModel.Table

	where := ""
	if plan.ID != nil {
		// Use ID for WHERE clause
		qb.paramCount++
		where = fmt.Sprintf("WHERE %s = $%d", plan.RootModel.PrimaryKey.ColumnName, qb.paramCount)
		qb.params = append(qb.params, plan.ID)
	} else if plan.Filters != nil {
		// Use filters for WHERE clause
		wherePart, err := qb.buildWhereClause(plan.Filters)
		if err != nil {
			return "", nil, err
		}
		where = wherePart
	} else {
		return "", nil, fmt.Errorf("id or filters required for delete operation")
	}

	sql := fmt.Sprintf("DELETE FROM %s %s;", table, where)

	return sql, qb.params, nil
}

// buildSelectClause generates the SELECT part of the query
func (qb *QueryBuilder) buildSelectClause(plan *planner.QueryPlan) string {
	var columns []string

	// Add selected columns (if any)
	if len(plan.Select) > 0 {
		for _, expr := range plan.Select {
			colName := fmt.Sprintf("%s.%s", expr.Column.TableAlias, expr.Column.ColumnName)
			if expr.Alias != expr.Column.ColumnName {
				colName = fmt.Sprintf("%s AS %s", colName, expr.Alias)
			}
			columns = append(columns, colName)
		}
	}

	// Add group by columns if grouping
	if len(plan.GroupBy) > 0 && len(plan.Select) == 0 {
		for _, groupExpr := range plan.GroupBy {
			colName := fmt.Sprintf("%s.%s", groupExpr.Column.TableAlias, groupExpr.Column.ColumnName)
			columns = append(columns, colName)
		}
	}

	// Add aggregates
	for _, agg := range plan.Aggregates {
		aggStr := qb.buildAggregateExpression(agg)
		columns = append(columns, aggStr)
	}

	// If no columns selected, use *
	if len(columns) == 0 {
		return "SELECT *"
	}

	return "SELECT " + strings.Join(columns, ", ")
}

// buildFromClause generates the FROM part of the query
func (qb *QueryBuilder) buildFromClause(plan *planner.QueryPlan) string {
	return fmt.Sprintf("FROM %s %s", plan.RootModel.Table, plan.RootModel.Alias)
}

// buildWhereClause generates the WHERE part of the query
func (qb *QueryBuilder) buildWhereClause(filterExpr planner.FilterExpr) (string, error) {
	filterSQL, err := qb.buildFilterExpression(filterExpr)
	if err != nil {
		return "", err
	}
	return "WHERE " + filterSQL, nil
}

// buildFilterExpression recursively builds filter expressions
func (qb *QueryBuilder) buildFilterExpression(expr planner.FilterExpr) (string, error) {
	switch e := expr.(type) {
	case *planner.ComparisonFilterIR:
		return qb.buildComparisonFilter(e)

	case *planner.LogicalFilterIR:
		return qb.buildLogicalFilter(e)

	default:
		return "", fmt.Errorf("unknown filter expression type")
	}
}

// buildComparisonFilter builds a single comparison filter
func (qb *QueryBuilder) buildComparisonFilter(f *planner.ComparisonFilterIR) (string, error) {
	colName := fmt.Sprintf("%s.%s", f.Left.TableAlias, f.Left.ColumnName)

	switch f.Operator {
	case dsl.OpEqual:
		if f.Value == nil {
			return "", fmt.Errorf("value required for = operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		paramPlaceholder := fmt.Sprintf("$%d", qb.paramCount)
		if needsTypeCasting(f.Left.DataType) {
			paramPlaceholder = addTypeCast(paramPlaceholder, f.Left.DataType)
		}
		return fmt.Sprintf("%s = %s", colName, paramPlaceholder), nil

	case dsl.OpNotEqual:
		if f.Value == nil {
			return "", fmt.Errorf("value required for != operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		paramPlaceholder := fmt.Sprintf("$%d", qb.paramCount)
		if needsTypeCasting(f.Left.DataType) {
			paramPlaceholder = addTypeCast(paramPlaceholder, f.Left.DataType)
		}
		return fmt.Sprintf("%s != %s", colName, paramPlaceholder), nil

	case dsl.OpGT:
		if f.Value == nil {
			return "", fmt.Errorf("value required for > operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s > $%d", colName, qb.paramCount), nil

	case dsl.OpGTE:
		if f.Value == nil {
			return "", fmt.Errorf("value required for >= operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s >= $%d", colName, qb.paramCount), nil

	case dsl.OpLT:
		if f.Value == nil {
			return "", fmt.Errorf("value required for < operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s < $%d", colName, qb.paramCount), nil

	case dsl.OpLTE:
		if f.Value == nil {
			return "", fmt.Errorf("value required for <= operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s <= $%d", colName, qb.paramCount), nil

	case dsl.OpIn:
		if f.Value == nil {
			return "", fmt.Errorf("value required for in operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		paramPlaceholder := fmt.Sprintf("$%d", qb.paramCount)
		if needsTypeCasting(f.Left.DataType) {
			paramPlaceholder = addTypeCast(paramPlaceholder, f.Left.DataType)
		}
		return fmt.Sprintf("%s = ANY(%s)", colName, paramPlaceholder), nil

	case dsl.OpNotIn:
		if f.Value == nil {
			return "", fmt.Errorf("value required for not_in operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		paramPlaceholder := fmt.Sprintf("$%d", qb.paramCount)
		if needsTypeCasting(f.Left.DataType) {
			paramPlaceholder = addTypeCast(paramPlaceholder, f.Left.DataType)
		}
		return fmt.Sprintf("%s != ALL(%s)", colName, paramPlaceholder), nil

	case dsl.OpIsNull:
		return fmt.Sprintf("%s IS NULL", colName), nil

	case dsl.OpNotNull:
		return fmt.Sprintf("%s IS NOT NULL", colName), nil

	case dsl.OpLike:
		if f.Value == nil {
			return "", fmt.Errorf("value required for like operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s LIKE $%d", colName, qb.paramCount), nil

	case dsl.OpILike:
		if f.Value == nil {
			return "", fmt.Errorf("value required for ilike operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s ILIKE $%d", colName, qb.paramCount), nil

	case dsl.OpStartsWith:
		if f.Value == nil {
			return "", fmt.Errorf("value required for starts_with operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value.(string)+"%")
		return fmt.Sprintf("%s LIKE $%d", colName, qb.paramCount), nil

	case dsl.OpEndsWith:
		if f.Value == nil {
			return "", fmt.Errorf("value required for ends_with operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, "%"+f.Value.Value.(string))
		return fmt.Sprintf("%s LIKE $%d", colName, qb.paramCount), nil

	case dsl.OpContains:
		if f.Value == nil {
			return "", fmt.Errorf("value required for contains operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, "%"+f.Value.Value.(string)+"%")
		return fmt.Sprintf("%s LIKE $%d", colName, qb.paramCount), nil

	case dsl.OpBetween:
		if f.Value == nil {
			return "", fmt.Errorf("value required for between operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s BETWEEN $%d AND $%d", colName, qb.paramCount, qb.paramCount+1), nil

	case dsl.OpBefore:
		if f.Value == nil {
			return "", fmt.Errorf("value required for before operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s < $%d", colName, qb.paramCount), nil

	case dsl.OpAfter:
		if f.Value == nil {
			return "", fmt.Errorf("value required for after operator")
		}
		qb.paramCount++
		qb.params = append(qb.params, f.Value.Value)
		return fmt.Sprintf("%s > $%d", colName, qb.paramCount), nil

	default:
		return "", fmt.Errorf("unknown operator: %s", f.Operator)
	}
}

// buildLogicalFilter builds logical filter expressions (AND/OR/NOT)
func (qb *QueryBuilder) buildLogicalFilter(f *planner.LogicalFilterIR) (string, error) {
	if len(f.Nodes) == 0 {
		return "", fmt.Errorf("logical filter has no nodes")
	}

	var parts []string
	for _, node := range f.Nodes {
		nodeSql, err := qb.buildFilterExpression(node)
		if err != nil {
			return "", err
		}
		parts = append(parts, nodeSql)
	}

	switch f.Op {
	case "AND":
		return "(" + strings.Join(parts, " AND ") + ")", nil

	case "OR":
		return "(" + strings.Join(parts, " OR ") + ")", nil

	case "NOT":
		if len(parts) != 1 {
			return "", fmt.Errorf("NOT filter must have exactly one node")
		}
		return "NOT " + parts[0], nil

	default:
		return "", fmt.Errorf("unknown logical operator: %s", f.Op)
	}
}

// buildGroupByClause generates the GROUP BY part of the query
func (qb *QueryBuilder) buildGroupByClause(plan *planner.QueryPlan) string {
	var groupCols []string
	for _, groupExpr := range plan.GroupBy {
		colName := fmt.Sprintf("%s.%s", groupExpr.Column.TableAlias, groupExpr.Column.ColumnName)
		groupCols = append(groupCols, colName)
	}
	return "GROUP BY " + strings.Join(groupCols, ", ")
}

// buildOrderByClause generates the ORDER BY part of the query
func (qb *QueryBuilder) buildOrderByClause(plan *planner.QueryPlan) string {
	var sortCols []string
	for _, sortExpr := range plan.Sort {
		var colRef string
		if sortExpr.Column != nil {
			colRef = fmt.Sprintf("%s.%s", sortExpr.Column.TableAlias, sortExpr.Column.ColumnName)
		} else if sortExpr.Aggregate != nil {
			colRef = sortExpr.Aggregate.Alias
		}

		direction := "ASC"
		if sortExpr.Direction == "DESC" {
			direction = "DESC"
		}

		sortCols = append(sortCols, colRef+" "+direction)
	}
	return "ORDER BY " + strings.Join(sortCols, ", ")
}

// buildPaginationClause generates the LIMIT/OFFSET part of the query
func (qb *QueryBuilder) buildPaginationClause(plan *planner.QueryPlan) string {
	qb.paramCount++
	limitParam := qb.paramCount
	qb.params = append(qb.params, plan.Pagination.Limit)

	qb.paramCount++
	offsetParam := qb.paramCount
	qb.params = append(qb.params, plan.Pagination.Offset)

	return fmt.Sprintf("LIMIT $%d OFFSET $%d", limitParam, offsetParam)
}

// buildAggregateExpression builds an aggregate function expression
func (qb *QueryBuilder) buildAggregateExpression(agg planner.AggregateExpr) string {
	var aggSQL string

	switch agg.Function {
	case planner.AggCountFn:
		if agg.Column == nil {
			aggSQL = "COUNT(*)"
		} else {
			aggSQL = fmt.Sprintf("COUNT(%s.%s)", agg.Column.TableAlias, agg.Column.ColumnName)
		}

	case planner.AggSumFn:
		aggSQL = fmt.Sprintf("SUM(%s.%s)", agg.Column.TableAlias, agg.Column.ColumnName)

	case planner.AggAvgFn:
		aggSQL = fmt.Sprintf("AVG(%s.%s)", agg.Column.TableAlias, agg.Column.ColumnName)

	case planner.AggMinFn:
		aggSQL = fmt.Sprintf("MIN(%s.%s)", agg.Column.TableAlias, agg.Column.ColumnName)

	case planner.AggMaxFn:
		aggSQL = fmt.Sprintf("MAX(%s.%s)", agg.Column.TableAlias, agg.Column.ColumnName)

	default:
		aggSQL = "COUNT(*)"
	}

	return fmt.Sprintf("%s AS %s", aggSQL, agg.Alias)
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		params:     []interface{}{},
		paramCount: 0,
	}
}
