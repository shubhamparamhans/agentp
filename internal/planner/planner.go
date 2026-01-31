package planner

import (
	"fmt"

	"udv/internal/dsl"
	"udv/internal/schema"
)

// FieldType represents the data type of a field
type FieldType string

const (
	// String types
	TypeString    FieldType = "string"
	
	// Integer types
	TypeInteger   FieldType = "integer"
	TypeInt       FieldType = "int"
	
	// Numeric types
	TypeFloat     FieldType = "float"
	TypeDecimal   FieldType = "decimal"
	
	// Boolean type
	TypeBoolean   FieldType = "boolean"
	
	// Date/Time types
	TypeDateTime  FieldType = "datetime"
	TypeTimestamp FieldType = "timestamp"
	TypeDate      FieldType = "date"
	TypeTime      FieldType = "time"
	
	// Special types (from data modelling processor)
	TypeUUID      FieldType = "uuid"
	TypeJSON      FieldType = "json"
	TypeBinary    FieldType = "binary"
)

// ColumnRef represents a resolved column reference
type ColumnRef struct {
	TableAlias string
	ColumnName string
	DataType   FieldType
}

// SelectExpr represents a column in the SELECT clause
type SelectExpr struct {
	Column      ColumnRef
	Alias       string
	IsAggregate bool
}

// JoinType represents the type of join
type JoinType string

const (
	JoinInner JoinType = "INNER"
	JoinLeft  JoinType = "LEFT"
)

// JoinCondition represents the ON clause of a join
type JoinCondition struct {
	Left  ColumnRef
	Right ColumnRef
}

// JoinPlan represents a single join operation
type JoinPlan struct {
	Type      JoinType
	FromAlias string
	ToTable   string
	ToAlias   string
	On        JoinCondition
}

// FilterExpr is the interface for filter expressions in IR
type FilterExpr interface {
	isFilterExpr()
}

// ComparisonFilterIR represents an atomic filter in IR
type ComparisonFilterIR struct {
	Left     ColumnRef
	Operator dsl.FilterOperator
	Value    *ValueExpr
}

func (c *ComparisonFilterIR) isFilterExpr() {}

// LogicalFilterIR represents AND/OR/NOT filters in IR
type LogicalFilterIR struct {
	Op    string        // "AND", "OR", "NOT"
	Nodes []FilterExpr
}

func (l *LogicalFilterIR) isFilterExpr() {}

// ValueExpr represents a strongly typed value
type ValueExpr struct {
	Value any
	Type  FieldType
}

// GroupExpr represents a GROUP BY expression
type GroupExpr struct {
	Column ColumnRef
}

// AggregateFn represents an aggregate function
type AggregateFn string

const (
	AggCountFn AggregateFn = "COUNT"
	AggSumFn   AggregateFn = "SUM"
	AggAvgFn   AggregateFn = "AVG"
	AggMinFn   AggregateFn = "MIN"
	AggMaxFn   AggregateFn = "MAX"
)

// AggregateExpr represents an aggregate function in IR
type AggregateExpr struct {
	Function AggregateFn
	Column   *ColumnRef
	Alias    string
}

// SortTarget represents what we're sorting by
type SortTarget string

const (
	SortColumn    SortTarget = "COLUMN"
	SortAggregate SortTarget = "AGGREGATE"
)

// SortExpr represents a sort specification
type SortExpr struct {
	Target    SortTarget
	Column    *ColumnRef
	Aggregate *AggregateExpr
	Direction string // "ASC", "DESC"
}

// Pagination represents pagination parameters
type Pagination struct {
	Limit  int
	Offset int
}

// QueryPlan represents the complete query plan IR
type QueryPlan struct {
	Operation  dsl.Operation         // NEW: Operation type
	RootModel  *ModelRef
	Select     []SelectExpr
	Joins      []JoinPlan
	Filters    FilterExpr
	GroupBy    []GroupExpr
	Aggregates []AggregateExpr
	Sort       []SortExpr
	Pagination Pagination
	Data       map[string]interface{} // NEW: For create/update operations
	ID         interface{}            // NEW: For update/delete operations
}

// ModelRef represents a model in the query plan
type ModelRef struct {
	Name       string
	Table      string
	Alias      string
	PrimaryKey ColumnRef
}

// Planner converts DSL queries into execution plans
type Planner struct {
	registry *schema.Registry
}

// NewPlanner creates a new query planner
func NewPlanner(reg *schema.Registry) *Planner {
	return &Planner{registry: reg}
}

// PlanQuery converts a validated DSL query into a QueryPlan IR
func (p *Planner) PlanQuery(q *dsl.Query) (*QueryPlan, error) {
	if q == nil {
		return nil, fmt.Errorf("query is nil")
	}

	// Default to select operation if not specified
	operation := q.Operation
	if operation == "" {
		operation = dsl.OpSelect
	}

	plan := &QueryPlan{
		Operation:  operation, // NEW: Set operation
		Select:     []SelectExpr{},
		Joins:      []JoinPlan{},
		Aggregates: []AggregateExpr{},
		Sort:       []SortExpr{},
		Data:       q.Data, // NEW: Pass data for create/update
		ID:         q.ID,   // NEW: Pass id for update/delete
	}

	// 1. Create root model reference
	model := p.registry.GetModel(q.Model)
	if model == nil {
		return nil, fmt.Errorf("model not found: %s", q.Model)
	}

	rootPrimaryKey := p.schemaFieldToColumnRef(model.Name, model.PrimaryKey, "t0")
	plan.RootModel = &ModelRef{
		Name:       model.Name,
		Table:      model.Table,
		Alias:      "t0",
		PrimaryKey: rootPrimaryKey,
	}

	// For create/update/delete operations, we can skip some planning steps
	if operation == dsl.OpCreate || operation == dsl.OpUpdate || operation == dsl.OpDelete {
		// Set default pagination for mutation operations
		plan.Pagination = Pagination{Limit: 1, Offset: 0}
		return plan, nil
	}

	// 2. Process SELECT clause
	if len(q.Fields) > 0 {
		for _, field := range q.Fields {
			colRef := p.schemaFieldToColumnRef(model.Name, field, "t0")
			plan.Select = append(plan.Select, SelectExpr{
				Column:      colRef,
				Alias:       field,
				IsAggregate: false,
			})
		}
	}

	// 3. Process WHERE filters
	if q.Filters != nil {
		filterIR, err := p.convertFilterExpr(model.Name, "t0", q.Filters)
		if err != nil {
			return nil, fmt.Errorf("failed to convert filters: %w", err)
		}
		plan.Filters = filterIR
	}

	// 4. Process GROUP BY
	if len(q.GroupBy) > 0 {
		for _, field := range q.GroupBy {
			colRef := p.schemaFieldToColumnRef(model.Name, field, "t0")
			plan.GroupBy = append(plan.GroupBy, GroupExpr{Column: colRef})
		}
	}

	// 5. Process AGGREGATES
	if len(q.Aggregates) > 0 {
		for _, agg := range q.Aggregates {
			var colRef *ColumnRef
			if agg.Field != "" {
				ref := p.schemaFieldToColumnRef(model.Name, agg.Field, "t0")
				colRef = &ref
			}

			aggFn := p.dslAggToIRAgg(agg.Function)
			plan.Aggregates = append(plan.Aggregates, AggregateExpr{
				Function: aggFn,
				Column:   colRef,
				Alias:    agg.Alias,
			})
		}
	}

	// 6. Process SORT
	if len(q.Sort) > 0 {
		for _, sort := range q.Sort {
			direction := "ASC"
			if sort.Direction == dsl.SortDesc {
				direction = "DESC"
			}

			colRef := p.schemaFieldToColumnRef(model.Name, sort.Field, "t0")
			plan.Sort = append(plan.Sort, SortExpr{
				Target:    SortColumn,
				Column:    &colRef,
				Direction: direction,
			})
		}
	}

	// 7. Process PAGINATION
	if q.Pagination != nil {
		plan.Pagination = Pagination{
			Limit:  q.Pagination.Limit,
			Offset: q.Pagination.Offset,
		}
	} else {
		// Default pagination
		plan.Pagination = Pagination{
			Limit:  100,
			Offset: 0,
		}
	}

	return plan, nil
}

// convertFilterExpr recursively converts a DSL filter to IR format
func (p *Planner) convertFilterExpr(modelName, tableAlias string, expr dsl.FilterExpr) (FilterExpr, error) {
	switch e := expr.(type) {
	case *dsl.ComparisonFilter:
		return p.convertComparisonFilter(modelName, tableAlias, e)

	case *dsl.LogicalFilter:
		return p.convertLogicalFilter(modelName, tableAlias, e)

	default:
		return nil, fmt.Errorf("unknown filter expression type")
	}
}

// convertComparisonFilter converts a DSL comparison filter to IR
func (p *Planner) convertComparisonFilter(modelName, tableAlias string, f *dsl.ComparisonFilter) (*ComparisonFilterIR, error) {
	colRef := p.schemaFieldToColumnRef(modelName, f.Field, tableAlias)

	var valueExpr *ValueExpr
	if f.Op != dsl.OpIsNull && f.Op != dsl.OpNotNull {
		valueExpr = &ValueExpr{
			Value: f.Value,
			Type:  colRef.DataType,
		}
	}

	return &ComparisonFilterIR{
		Left:     colRef,
		Operator: f.Op,
		Value:    valueExpr,
	}, nil
}

// convertLogicalFilter converts a DSL logical filter to IR
func (p *Planner) convertLogicalFilter(modelName, tableAlias string, f *dsl.LogicalFilter) (*LogicalFilterIR, error) {
	logicalIR := &LogicalFilterIR{
		Nodes: []FilterExpr{},
	}

	if len(f.And) > 0 {
		logicalIR.Op = "AND"
		for _, cond := range f.And {
			irCond, err := p.convertComparisonFilter(modelName, tableAlias, cond)
			if err != nil {
				return nil, err
			}
			logicalIR.Nodes = append(logicalIR.Nodes, irCond)
		}
	} else if len(f.Or) > 0 {
		logicalIR.Op = "OR"
		for _, cond := range f.Or {
			irCond, err := p.convertComparisonFilter(modelName, tableAlias, cond)
			if err != nil {
				return nil, err
			}
			logicalIR.Nodes = append(logicalIR.Nodes, irCond)
		}
	} else if f.Not != nil {
		logicalIR.Op = "NOT"
		irCond, err := p.convertComparisonFilter(modelName, tableAlias, f.Not)
		if err != nil {
			return nil, err
		}
		logicalIR.Nodes = append(logicalIR.Nodes, irCond)
	}

	return logicalIR, nil
}

// schemaFieldToColumnRef converts a schema field to a ColumnRef
func (p *Planner) schemaFieldToColumnRef(modelName, fieldName, tableAlias string) ColumnRef {
	field, err := p.registry.GetField(modelName, fieldName)
	if err != nil {
		// This should not happen if validation was done correctly
		return ColumnRef{}
	}

	return ColumnRef{
		TableAlias: tableAlias,
		ColumnName: fieldName,
		DataType:   FieldType(field.Type),
	}
}

// dslAggToIRAgg converts DSL aggregate function to IR aggregate function
func (p *Planner) dslAggToIRAgg(fn dsl.AggregateFunc) AggregateFn {
	switch fn {
	case dsl.AggCount:
		return AggCountFn
	case dsl.AggSum:
		return AggSumFn
	case dsl.AggAvg:
		return AggAvgFn
	case dsl.AggMin:
		return AggMinFn
	case dsl.AggMax:
		return AggMaxFn
	default:
		return AggCountFn
	}
}
