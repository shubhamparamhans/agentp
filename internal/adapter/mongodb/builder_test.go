package mongodb

import (
	"testing"

	"udv/internal/config"
	"udv/internal/dsl"
	"udv/internal/planner"
	"udv/internal/schema"

	"go.mongodb.org/mongo-driver/bson"
)

func setupMongoDBTestRegistry() *schema.Registry {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "users",
				Table:      "users",
				PrimaryKey: "_id",
				Fields: []config.Field{
					{Name: "_id", Type: "uuid", Nullable: false},
					{Name: "name", Type: "string", Nullable: false},
					{Name: "email", Type: "string", Nullable: false},
					{Name: "age", Type: "integer", Nullable: true},
					{Name: "active", Type: "boolean", Nullable: true},
					{Name: "created_at", Type: "timestamp", Nullable: false},
				},
			},
			{
				Name:       "orders",
				Table:      "orders",
				PrimaryKey: "_id",
				Fields: []config.Field{
					{Name: "_id", Type: "uuid", Nullable: false},
					{Name: "user_id", Type: "uuid", Nullable: false},
					{Name: "status", Type: "string", Nullable: false},
					{Name: "amount", Type: "decimal", Nullable: false},
				},
			},
		},
	}

	reg := schema.NewRegistry()
	reg.LoadFromConfig(cfg)
	return reg
}

func TestBuildQuery_SimpleFind(t *testing.T) {
	reg := setupMongoDBTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model:  "users",
		Fields: []string{"name", "email"},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	query, _, err := builder.BuildQuery(plan)
	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		t.Errorf("Expected MongoQuery, got %T", query)
	}

	if mongoQuery.Collection != "users" {
		t.Errorf("Expected collection 'users', got '%s'", mongoQuery.Collection)
	}

	if mongoQuery.Operation != "find" {
		t.Errorf("Expected operation 'find', got '%s'", mongoQuery.Operation)
	}
}

func TestBuildQuery_WithFilter(t *testing.T) {
	reg := setupMongoDBTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model: "users",
		Filters: &dsl.ComparisonFilter{
			Field: "status",
			Op:    dsl.OpEqual,
			Value: "active",
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	query, _, err := builder.BuildQuery(plan)
	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		t.Fatalf("Expected MongoQuery, got %T", query)
	}

	// Filter should be a bson.M
	filter, ok := mongoQuery.Filter.(bson.M)
	if !ok {
		t.Errorf("Expected bson.M filter, got %T", mongoQuery.Filter)
	}

	if len(filter) == 0 {
		t.Error("Filter should not be empty")
	}
}

func TestBuildQuery_WithPagination(t *testing.T) {
	reg := setupMongoDBTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model: "users",
		Pagination: &dsl.Pagination{
			Limit:  10,
			Offset: 5,
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	query, _, err := builder.BuildQuery(plan)
	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		t.Fatalf("Expected MongoQuery, got %T", query)
	}

	if mongoQuery.Options == nil {
		t.Error("Options should not be nil for pagination")
	}
}

func TestBuildQuery_Insert(t *testing.T) {
	reg := setupMongoDBTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Operation: dsl.OpCreate,
		Model:     "users",
		Data: map[string]interface{}{
			"name":  "John",
			"email": "john@example.com",
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	query, _, err := builder.BuildQuery(plan)
	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		t.Fatalf("Expected MongoQuery, got %T", query)
	}

	if mongoQuery.Operation != "insert" {
		t.Errorf("Expected operation 'insert', got '%s'", mongoQuery.Operation)
	}

	doc, ok := mongoQuery.Document.(bson.M)
	if !ok {
		t.Errorf("Expected bson.M document, got %T", mongoQuery.Document)
	}

	if len(doc) == 0 {
		t.Error("Document should not be empty")
	}
}

func TestBuildQuery_Update(t *testing.T) {
	reg := setupMongoDBTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Operation: dsl.OpUpdate,
		Model:     "users",
		Filters: &dsl.ComparisonFilter{
			Field: "_id",
			Op:    dsl.OpEqual,
			Value: "user123",
		},
		Data: map[string]interface{}{
			"email": "newemail@example.com",
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	query, _, err := builder.BuildQuery(plan)
	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		t.Fatalf("Expected MongoQuery, got %T", query)
	}

	if mongoQuery.Operation != "update" {
		t.Errorf("Expected operation 'update', got '%s'", mongoQuery.Operation)
	}

	updateDoc, ok := mongoQuery.Update.(bson.M)
	if !ok {
		t.Errorf("Expected bson.M update, got %T", mongoQuery.Update)
	}

	if updateDoc["$set"] == nil {
		t.Error("Update should contain $set operator")
	}
}

func TestBuildQuery_Delete(t *testing.T) {
	reg := setupMongoDBTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Operation: dsl.OpDelete,
		Model:     "users",
		Filters: &dsl.ComparisonFilter{
			Field: "_id",
			Op:    dsl.OpEqual,
			Value: "user123",
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	query, _, err := builder.BuildQuery(plan)
	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		t.Fatalf("Expected MongoQuery, got %T", query)
	}

	if mongoQuery.Operation != "delete" {
		t.Errorf("Expected operation 'delete', got '%s'", mongoQuery.Operation)
	}

	if mongoQuery.Filter == nil {
		t.Error("Filter should not be nil for delete operation")
	}
}

func TestConvertOperator(t *testing.T) {
	builder := NewQueryBuilder()

	tests := []struct {
		op          string
		value       interface{}
		expectedOp  string
		expectedVal interface{}
		shouldError bool
		checkFn     func(expected, actual interface{}) bool
	}{
		{"=", "test", "$eq", "test", false, func(e, a interface{}) bool { return e == a }},
		{"!=", 42, "$ne", 42, false, func(e, a interface{}) bool { return e == a }},
		{">", 100, "$gt", 100, false, func(e, a interface{}) bool { return e == a }},
		{">=", 50, "$gte", 50, false, func(e, a interface{}) bool { return e == a }},
		{"<", 10, "$lt", 10, false, func(e, a interface{}) bool { return e == a }},
		{"<=", 5, "$lte", 5, false, func(e, a interface{}) bool { return e == a }},
		{"in", []string{"a", "b"}, "$in", []string{"a", "b"}, false, sliceEqual},
		{"not_in", []int{1, 2}, "$nin", []int{1, 2}, false, sliceEqual},
		{"like", "pattern", "$regex", "pattern", false, func(e, a interface{}) bool { return e == a }},
		{"is_null", nil, "$exists", false, false, func(e, a interface{}) bool { return e == a }},
		{"not_null", nil, "$exists", true, false, func(e, a interface{}) bool { return e == a }},
		{"unknown_op", "val", "", nil, true, nil},
	}

	for _, test := range tests {
		t.Run(test.op, func(t *testing.T) {
			mongoOp, mongoVal, err := builder.convertOperator(test.op, test.value)

			if test.shouldError && err == nil {
				t.Errorf("Expected error for operator %q", test.op)
			}

			if !test.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !test.shouldError {
				if mongoOp != test.expectedOp {
					t.Errorf("Expected operator %q, got %q", test.expectedOp, mongoOp)
				}
				if !test.checkFn(test.expectedVal, mongoVal) {
					t.Errorf("Expected value %v, got %v", test.expectedVal, mongoVal)
				}
			}
		})
	}
}

func sliceEqual(expected, actual interface{}) bool {
	switch e := expected.(type) {
	case []string:
		a, ok := actual.([]string)
		if !ok {
			return false
		}
		if len(e) != len(a) {
			return false
		}
		for i, v := range e {
			if v != a[i] {
				return false
			}
		}
		return true
	case []int:
		a, ok := actual.([]int)
		if !ok {
			return false
		}
		if len(e) != len(a) {
			return false
		}
		for i, v := range e {
			if v != a[i] {
				return false
			}
		}
		return true
	}
	return false
}

func TestBuildQuery_UnsupportedOperation(t *testing.T) {
	builder := NewQueryBuilder()

	// Create a plan with unsupported operation
	plan := &planner.QueryPlan{
		Operation: "unsupported",
		RootModel: &planner.ModelRef{
			Table: "users",
		},
	}

	_, _, err := builder.BuildQuery(plan)
	if err == nil {
		t.Error("Expected error for unsupported operation")
	}
}

func TestMongoQuery_TypeAssertion(t *testing.T) {
	query := &MongoQuery{
		Collection: "users",
		Operation:  "find",
		Filter:     bson.M{"name": "John"},
		Options:    nil,
	}

	if query.Collection != "users" {
		t.Errorf("Collection mismatch: expected 'users', got %q", query.Collection)
	}

	if query.Operation != "find" {
		t.Errorf("Operation mismatch: expected 'find', got %q", query.Operation)
	}
}
