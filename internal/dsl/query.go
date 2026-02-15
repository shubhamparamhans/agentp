package dsl

import (
	"fmt"

	"udv/internal/schema"
)

// FilterOperator represents a comparison operator
type FilterOperator string

const (
	// Comparison operators
	OpEqual    FilterOperator = "="
	OpNotEqual FilterOperator = "!="
	OpGT       FilterOperator = ">"
	OpGTE      FilterOperator = ">="
	OpLT       FilterOperator = "<"
	OpLTE      FilterOperator = "<="
	OpIn       FilterOperator = "in"
	OpNotIn    FilterOperator = "not_in"
	OpIsNull   FilterOperator = "is_null"
	OpNotNull  FilterOperator = "not_null"

	// String operators
	OpLike       FilterOperator = "like"
	OpILike      FilterOperator = "ilike"
	OpStartsWith FilterOperator = "starts_with"
	OpEndsWith   FilterOperator = "ends_with"
	OpContains   FilterOperator = "contains"

	// Date operators
	OpBefore  FilterOperator = "before"
	OpAfter   FilterOperator = "after"
	OpBetween FilterOperator = "between"
)

// AggregateFunc represents an aggregate function
type AggregateFunc string

const (
	AggCount AggregateFunc = "count"
	AggSum   AggregateFunc = "sum"
	AggAvg   AggregateFunc = "avg"
	AggMin   AggregateFunc = "min"
	AggMax   AggregateFunc = "max"
)

// SortDirection represents sort order
type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// Operation represents the type of operation to perform
type Operation string

const (
	OpSelect Operation = "select" // Default operation
	OpCreate Operation = "create"
	OpUpdate Operation = "update"
	OpDelete Operation = "delete"
)

// Query represents a complete query specification
type Query struct {
	Operation  Operation              `json:"operation,omitempty"` // NEW: Operation type (defaults to "select")
	Model      string                 `json:"model"`
	Fields     []string               `json:"fields,omitempty"`
	Filters    FilterExpr             `json:"filters,omitempty"`
	GroupBy    []string               `json:"group_by,omitempty"`
	Aggregates []Aggregate            `json:"aggregates,omitempty"`
	Sort       []Sort                 `json:"sort,omitempty"`
	Pagination *Pagination            `json:"pagination,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"` // NEW: For create/update operations
	ID         interface{}            `json:"id,omitempty"`   // NEW: For update/delete operations
}

// FilterExpr represents a filter expression (can be AND, OR, NOT, or atomic)
type FilterExpr interface {
	isFilterExpr()
}

// LogicalFilter represents AND/OR/NOT operators
type LogicalFilter struct {
	And []*ComparisonFilter `json:"and,omitempty"`
	Or  []*ComparisonFilter `json:"or,omitempty"`
	Not *ComparisonFilter   `json:"not,omitempty"`
}

func (l *LogicalFilter) isFilterExpr() {}

// ComparisonFilter represents a single field comparison
type ComparisonFilter struct {
	Field string         `json:"field"`
	Op    FilterOperator `json:"op"`
	Value interface{}    `json:"value,omitempty"`
}

func (c *ComparisonFilter) isFilterExpr() {}

// Aggregate represents an aggregate function
type Aggregate struct {
	Function AggregateFunc `json:"fn"`
	Field    string        `json:"field,omitempty"`
	Alias    string        `json:"alias"`
}

// Sort represents a sort specification
type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction,omitempty"`
}

// Pagination represents pagination parameters
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset,omitempty"`
}

// Validator validates queries against schema
type Validator struct {
	registry *schema.Registry
}

// NewValidator creates a new query validator
func NewValidator(reg *schema.Registry) *Validator {
	return &Validator{registry: reg}
}

// ValidateQuery validates a complete query
func (v *Validator) ValidateQuery(q *Query) error {
	if q == nil {
		return fmt.Errorf("query is nil")
	}

	// Default to select operation if not specified (backward compatibility)
	if q.Operation == "" {
		q.Operation = OpSelect
	}

	// Validate model
	if q.Model == "" {
		return fmt.Errorf("model is required")
	}
	if !v.registry.ModelExists(q.Model) {
		return fmt.Errorf("model not found: %s", q.Model)
	}

	// Validate operation-specific requirements
	switch q.Operation {
	case OpCreate:
		return v.validateCreate(q)
	case OpUpdate:
		return v.validateUpdate(q)
	case OpDelete:
		return v.validateDelete(q)
	case OpSelect:
		// Continue with existing validation for select
	default:
		return fmt.Errorf("invalid operation: %s", q.Operation)
	}

	// Validate fields (for select operations)
	if err := v.validateFields(q.Model, q.Fields); err != nil {
		return err
	}

	// Validate filters
	if q.Filters != nil {
		if err := v.validateFilterExpr(q.Model, q.Filters); err != nil {
			return err
		}
	}

	// Validate group_by
	if err := v.validateGroupBy(q.Model, q.GroupBy); err != nil {
		return err
	}

	// Validate aggregates
	if err := v.validateAggregates(q.Model, q.Aggregates, len(q.GroupBy) > 0); err != nil {
		return err
	}

	// Validate sort
	if err := v.validateSort(q.Model, q.Sort); err != nil {
		return err
	}

	// Validate pagination
	if err := v.validatePagination(q.Pagination); err != nil {
		return err
	}

	return nil
}

// validateCreate validates a create operation
func (v *Validator) validateCreate(q *Query) error {
	if q.Data == nil || len(q.Data) == 0 {
		return fmt.Errorf("data is required for create operation")
	}

	model := v.registry.GetModel(q.Model)
	if model == nil {
		return fmt.Errorf("model not found: %s", q.Model)
	}

	// Validate all fields in data exist in model
	for fieldName := range q.Data {
		if !v.registry.FieldExists(q.Model, fieldName) {
			return fmt.Errorf("field not found in model %s: %s", q.Model, fieldName)
		}
	}

	// Check required fields (non-nullable fields that don't have defaults)
	// Skip the primary key field as it's typically auto-generated
	// Note: This is a basic check - full validation would require checking for default values
	fields, _ := v.registry.GetModelFields(q.Model)
	for _, field := range fields {
		// Skip primary key fields (they're typically auto-generated)
		if field.Name == model.PrimaryKey {
			continue
		}
		if !field.Nullable && q.Data[field.Name] == nil {
			// Field is required but not provided
			return fmt.Errorf("required field missing: %s", field.Name)
		}
	}

	return nil
}

// validateUpdate validates an update operation
func (v *Validator) validateUpdate(q *Query) error {
	if q.ID == nil && q.Filters == nil {
		return fmt.Errorf("id or filters required for update operation")
	}

	if q.Data == nil || len(q.Data) == 0 {
		return fmt.Errorf("data is required for update operation")
	}

	// Validate all fields being updated exist in model
	for fieldName := range q.Data {
		if !v.registry.FieldExists(q.Model, fieldName) {
			return fmt.Errorf("field not found in model %s: %s", q.Model, fieldName)
		}
	}

	return nil
}

// validateDelete validates a delete operation
func (v *Validator) validateDelete(q *Query) error {
	if q.ID == nil && q.Filters == nil {
		return fmt.Errorf("id or filters required for delete operation")
	}
	return nil
}

func (v *Validator) validateFields(modelName string, fields []string) error {
	if len(fields) == 0 {
		return nil // Empty fields is allowed
	}

	for _, field := range fields {
		if field == "" {
			return fmt.Errorf("field name cannot be empty")
		}
		if !v.registry.FieldExists(modelName, field) {
			return fmt.Errorf("field not found in model %s: %s", modelName, field)
		}
	}
	return nil
}

func (v *Validator) validateFilterExpr(modelName string, expr FilterExpr) error {
	switch e := expr.(type) {
	case *LogicalFilter:
		if e.And != nil {
			for _, f := range e.And {
				if err := v.validateComparisonFilter(modelName, f); err != nil {
					return err
				}
			}
		}
		if e.Or != nil {
			for _, f := range e.Or {
				if err := v.validateComparisonFilter(modelName, f); err != nil {
					return err
				}
			}
		}
		if e.Not != nil {
			if err := v.validateComparisonFilter(modelName, e.Not); err != nil {
				return err
			}
		}
		return nil

	case *ComparisonFilter:
		return v.validateComparisonFilter(modelName, e)

	default:
		return fmt.Errorf("invalid filter expression type")
	}
}

func (v *Validator) validateComparisonFilter(modelName string, f *ComparisonFilter) error {
	if f == nil {
		return nil
	}

	if f.Field == "" {
		return fmt.Errorf("filter field is required")
	}

	// Check field exists
	field, err := v.registry.GetField(modelName, f.Field)
	if err != nil {
		return fmt.Errorf("invalid filter field: %v", err)
	}

	// Check field is filterable
	if !field.Filterable {
		return fmt.Errorf("field is not filterable: %s", f.Field)
	}

	// Validate operator for field type
	if err := v.validateOperatorForType(f.Op, field.Type, f.Value); err != nil {
		return fmt.Errorf("invalid filter operator for field %s: %v", f.Field, err)
	}

	return nil
}

func (v *Validator) validateOperatorForType(op FilterOperator, fieldType string, value interface{}) error {
	// NULL operators don't need a value
	if op == OpIsNull || op == OpNotNull {
		return nil
	}

	// Operators that require array values
	if op == OpIn || op == OpNotIn {
		// In/NotIn should have array values
		return nil
	}

	if op == OpBetween {
		// Between requires array value
		return nil
	}

	// String operators
	if op == OpLike || op == OpILike || op == OpStartsWith || op == OpEndsWith || op == OpContains {
		if fieldType != "string" {
			return fmt.Errorf("string operator %s not valid for type %s", op, fieldType)
		}
		return nil
	}

	// Comparison operators
	comparisonOps := map[FilterOperator]bool{
		OpEqual:    true,
		OpNotEqual: true,
		OpGT:       true,
		OpGTE:      true,
		OpLT:       true,
		OpLTE:      true,
		OpBefore:   true,
		OpAfter:    true,
	}

	if comparisonOps[op] {
		return nil
	}

	return fmt.Errorf("unknown operator: %s", op)
}

func (v *Validator) validateGroupBy(modelName string, groupBy []string) error {
	if len(groupBy) == 0 {
		return nil
	}

	for _, field := range groupBy {
		if field == "" {
			return fmt.Errorf("group_by field cannot be empty")
		}

		f, err := v.registry.GetField(modelName, field)
		if err != nil {
			return fmt.Errorf("invalid group_by field: %v", err)
		}

		if !f.Groupable {
			return fmt.Errorf("field is not groupable: %s", field)
		}
	}

	return nil
}

func (v *Validator) validateAggregates(modelName string, aggs []Aggregate, hasGroupBy bool) error {
	if len(aggs) == 0 {
		return nil
	}

	validFuncs := map[AggregateFunc]bool{
		AggCount: true,
		AggSum:   true,
		AggAvg:   true,
		AggMin:   true,
		AggMax:   true,
	}

	for i, agg := range aggs {
		if agg.Alias == "" {
			return fmt.Errorf("aggregate[%d] alias is required", i)
		}

		if !validFuncs[agg.Function] {
			return fmt.Errorf("aggregate[%d] unknown function: %s", i, agg.Function)
		}

		// count can omit field
		if agg.Function == AggCount && agg.Field == "" {
			continue
		}

		if agg.Field == "" {
			return fmt.Errorf("aggregate[%d] field is required for function %s", i, agg.Function)
		}

		f, err := v.registry.GetField(modelName, agg.Field)
		if err != nil {
			return fmt.Errorf("aggregate[%d] invalid field: %v", i, err)
		}

		if !f.Aggregatable {
			return fmt.Errorf("aggregate[%d] field is not aggregatable: %s", i, agg.Field)
		}

		// Validate function for field type
		if err := v.validateAggregateForType(agg.Function, f.Type); err != nil {
			return fmt.Errorf("aggregate[%d] invalid for field %s: %v", i, agg.Field, err)
		}
	}

	return nil
}

func (v *Validator) validateAggregateForType(fn AggregateFunc, fieldType string) error {
	switch fn {
	case AggCount:
		return nil // count works on any type

	case AggSum, AggAvg:
		// Only numeric types
		if fieldType != "integer" && fieldType != "int" && fieldType != "float" && fieldType != "decimal" {
			return fmt.Errorf("function %s requires numeric field, got %s", fn, fieldType)
		}
		return nil

	case AggMin, AggMax:
		return nil // min/max work on comparable types

	default:
		return fmt.Errorf("unknown aggregate function: %s", fn)
	}
}

func (v *Validator) validateSort(modelName string, sort []Sort) error {
	if len(sort) == 0 {
		return nil
	}

	for i, s := range sort {
		if s.Field == "" {
			return fmt.Errorf("sort[%d] field is required", i)
		}

		if !v.registry.FieldExists(modelName, s.Field) {
			return fmt.Errorf("sort[%d] field not found: %s", i, s.Field)
		}

		// Validate direction
		if s.Direction != "" && s.Direction != SortAsc && s.Direction != SortDesc {
			return fmt.Errorf("sort[%d] invalid direction: %s", i, s.Direction)
		}
	}

	return nil
}

func (v *Validator) validatePagination(p *Pagination) error {
	if p == nil {
		return nil
	}

	if p.Limit <= 0 {
		return fmt.Errorf("pagination limit must be greater than 0")
	}

	if p.Offset < 0 {
		return fmt.Errorf("pagination offset must be non-negative")
	}

	return nil
}

// Helper function to create a simple comparison filter
func NewComparisonFilter(field string, op FilterOperator, value interface{}) *ComparisonFilter {
	return &ComparisonFilter{
		Field: field,
		Op:    op,
		Value: value,
	}
}
