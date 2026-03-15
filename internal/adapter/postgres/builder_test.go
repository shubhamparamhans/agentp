package postgres

import (
	"strings"
	"testing"

	"udv/internal/config"
	"udv/internal/dsl"
	"udv/internal/planner"
	"udv/internal/schema"
)

func setupTestRegistry() *schema.Registry {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "orders",
				Table:      "orders",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "user_id", Type: "integer", Nullable: false},
					{Name: "status", Type: "string", Nullable: false},
					{Name: "amount", Type: "decimal", Nullable: false},
					{Name: "created_at", Type: "timestamp", Nullable: false},
				},
			},
		},
	}

	reg := schema.NewRegistry()
	reg.LoadFromConfig(cfg)
	return reg
}

func TestBuildQuery_SimpleSelect(t *testing.T) {
	reg := setupTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model:  "orders",
		Fields: []string{"id", "status", "amount"},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	sqlIface, _, err := builder.BuildQuery(plan)

	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	sql, ok := sqlIface.(string)
	if !ok {
		t.Fatalf("expected sql string, got %T", sqlIface)
	}
	if !strings.Contains(sql, "SELECT") {
		t.Errorf("SQL missing SELECT: %s", sql)
	}
	if !strings.Contains(sql, "t0.id") {
		t.Errorf("SQL missing column: %s", sql)
	}
	if !strings.Contains(sql, "FROM orders") {
		t.Errorf("SQL missing FROM clause: %s", sql)
	}
	if !strings.Contains(sql, "LIMIT") {
		t.Errorf("SQL missing LIMIT: %s", sql)
	}
}

func TestBuildQuery_WithFilter(t *testing.T) {
	reg := setupTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model: "orders",
		Filters: &dsl.ComparisonFilter{
			Field: "status",
			Op:    dsl.OpEqual,
			Value: "PAID",
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	sqlIface, params, err := builder.BuildQuery(plan)

	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	sql, ok := sqlIface.(string)
	if !ok {
		t.Fatalf("expected sql string, got %T", sqlIface)
	}
	if !strings.Contains(sql, "WHERE") {
		t.Errorf("SQL missing WHERE: %s", sql)
	}
	if !strings.Contains(sql, "t0.status = $1") {
		t.Errorf("SQL missing filter condition: %s", sql)
	}

	if len(params) != 3 { // status value, limit, offset
		t.Errorf("Expected 3 params, got %d", len(params))
	}
	if params[0] != "PAID" {
		t.Errorf("First param should be 'PAID', got %v", params[0])
	}
}

func TestBuildQuery_WithLogicalFilter(t *testing.T) {
	reg := setupTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model: "orders",
		Filters: &dsl.LogicalFilter{
			And: []*dsl.ComparisonFilter{
				{Field: "status", Op: dsl.OpEqual, Value: "PAID"},
				{Field: "amount", Op: dsl.OpGT, Value: 1000},
			},
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	sqlIface, params, err := builder.BuildQuery(plan)

	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	sql, ok := sqlIface.(string)
	if !ok {
		t.Fatalf("expected sql string, got %T", sqlIface)
	}
	if !strings.Contains(sql, "AND") {
		t.Errorf("SQL missing AND: %s", sql)
	}

	if len(params) != 4 { // status, amount, limit, offset
		t.Errorf("Expected 4 params, got %d", len(params))
	}
}

func TestBuildQuery_WithGroupByAndAggregate(t *testing.T) {
	reg := setupTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model:   "orders",
		GroupBy: []string{"status"},
		Aggregates: []dsl.Aggregate{
			{Function: dsl.AggSum, Field: "amount", Alias: "total_amount"},
			{Function: dsl.AggCount, Field: "", Alias: "order_count"},
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	sqlIface, params, err := builder.BuildQuery(plan)

	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	sql, ok := sqlIface.(string)
	if !ok {
		t.Fatalf("expected sql string, got %T", sqlIface)
	}
	if !strings.Contains(sql, "GROUP BY") {
		t.Errorf("SQL missing GROUP BY: %s", sql)
	}
	if !strings.Contains(sql, "SUM(") {
		t.Errorf("SQL missing SUM: %s", sql)
	}
	if !strings.Contains(sql, "COUNT(*)") {
		t.Errorf("SQL missing COUNT(*): %s", sql)
	}
	if !strings.Contains(sql, "total_amount") {
		t.Errorf("SQL missing alias: %s", sql)
	}

	if len(params) != 2 { // limit, offset
		t.Errorf("Expected 2 params, got %d", len(params))
	}
}

func TestBuildQuery_WithSort(t *testing.T) {
	reg := setupTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model: "orders",
		Sort: []dsl.Sort{
			{Field: "created_at", Direction: dsl.SortDesc},
			{Field: "amount", Direction: dsl.SortAsc},
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	sqlIface, _, err := builder.BuildQuery(plan)

	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	sql, ok := sqlIface.(string)
	if !ok {
		t.Fatalf("expected sql string, got %T", sqlIface)
	}
	if !strings.Contains(sql, "ORDER BY") {
		t.Errorf("SQL missing ORDER BY: %s", sql)
	}
	if !strings.Contains(sql, "DESC") {
		t.Errorf("SQL missing DESC: %s", sql)
	}
	if !strings.Contains(sql, "ASC") {
		t.Errorf("SQL missing ASC: %s", sql)
	}
}

func TestBuildQuery_WithPagination(t *testing.T) {
	reg := setupTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model: "orders",
		Pagination: &dsl.Pagination{
			Limit:  25,
			Offset: 50,
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	sqlIface, params, err := builder.BuildQuery(plan)

	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	sql, ok := sqlIface.(string)
	if !ok {
		t.Fatalf("expected sql string, got %T", sqlIface)
	}
	if !strings.Contains(sql, "LIMIT $1 OFFSET $2") {
		t.Errorf("SQL missing correct LIMIT/OFFSET: %s", sql)
	}

	if len(params) != 2 {
		t.Errorf("Expected 2 params (limit, offset), got %d", len(params))
	}
	if params[0] != 25 {
		t.Errorf("First param should be 25, got %v", params[0])
	}
	if params[1] != 50 {
		t.Errorf("Second param should be 50, got %v", params[1])
	}
}

func TestBuildQuery_FilterOperators(t *testing.T) {
	tests := []struct {
		name     string
		op       dsl.FilterOperator
		expected string
	}{
		{"equals", dsl.OpEqual, " = $"},
		{"not_equal", dsl.OpNotEqual, " != $"},
		{"gt", dsl.OpGT, " > $"},
		{"gte", dsl.OpGTE, " >= $"},
		{"lt", dsl.OpLT, " < $"},
		{"lte", dsl.OpLTE, " <= $"},
		{"in", dsl.OpIn, " = ANY($"},
		{"not_in", dsl.OpNotIn, " != ALL($"},
		{"is_null", dsl.OpIsNull, " IS NULL"},
		{"not_null", dsl.OpNotNull, " IS NOT NULL"},
		{"like", dsl.OpLike, " LIKE $"},
		{"ilike", dsl.OpILike, " ILIKE $"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := setupTestRegistry()
			queryPlanner := planner.NewPlanner(reg)

			var value interface{} = "test"
			if tt.op == dsl.OpIsNull || tt.op == dsl.OpNotNull {
				value = nil
			}

			dslQuery := &dsl.Query{
				Model: "orders",
				Filters: &dsl.ComparisonFilter{
					Field: "status",
					Op:    tt.op,
					Value: value,
				},
			}

			plan, err := queryPlanner.PlanQuery(dslQuery)
			if err != nil {
				t.Fatalf("PlanQuery error: %v", err)
			}

			builder := NewQueryBuilder()
			sqlIface, _, err := builder.BuildQuery(plan)

			if err != nil {
				t.Errorf("BuildQuery error: %v", err)
			}

			sql, ok := sqlIface.(string)
			if !ok {
				t.Fatalf("expected sql string, got %T", sqlIface)
			}
			if !strings.Contains(sql, tt.expected) {
				t.Errorf("SQL missing expected operator %s: %s", tt.expected, sql)
			}
		})
	}
}

func TestBuildQuery_ComplexQuery(t *testing.T) {
	reg := setupTestRegistry()
	queryPlanner := planner.NewPlanner(reg)

	dslQuery := &dsl.Query{
		Model:   "orders",
		Fields:  []string{"status"},
		GroupBy: []string{"status"},
		Filters: &dsl.LogicalFilter{
			And: []*dsl.ComparisonFilter{
				{Field: "created_at", Op: dsl.OpAfter, Value: "2024-01-01"},
				{Field: "amount", Op: dsl.OpGTE, Value: 100},
			},
		},
		Aggregates: []dsl.Aggregate{
			{Function: dsl.AggSum, Field: "amount", Alias: "total_amount"},
			{Function: dsl.AggCount, Field: "", Alias: "order_count"},
		},
		Sort: []dsl.Sort{
			{Field: "total_amount", Direction: dsl.SortDesc},
		},
		Pagination: &dsl.Pagination{
			Limit:  100,
			Offset: 0,
		},
	}

	plan, err := queryPlanner.PlanQuery(dslQuery)
	if err != nil {
		t.Fatalf("PlanQuery error: %v", err)
	}

	builder := NewQueryBuilder()
	sqlIface, params, err := builder.BuildQuery(plan)

	if err != nil {
		t.Errorf("BuildQuery error: %v", err)
	}

	sql, ok := sqlIface.(string)
	if !ok {
		t.Fatalf("expected sql string, got %T", sqlIface)
	}

	// Verify all parts are present
	requiredParts := []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "LIMIT"}
	for _, part := range requiredParts {
		if !strings.Contains(sql, part) {
			t.Errorf("SQL missing %s: %s", part, sql)
		}
	}

	// Should have parameters for: created_at, amount, limit, offset
	if len(params) < 4 {
		t.Errorf("Expected at least 4 params, got %d", len(params))
	}
}

func TestBuildQuery_NilPlan(t *testing.T) {
	builder := NewQueryBuilder()
	_, _, err := builder.BuildQuery(nil)
	if err == nil {
		t.Errorf("BuildQuery(nil) should error")
	}
}

func TestBuildQuery_StringOperators(t *testing.T) {
	tests := []struct {
		name     string
		op       dsl.FilterOperator
		value    string
		contains string
	}{
		{"starts_with", dsl.OpStartsWith, "test", " LIKE $"},
		{"ends_with", dsl.OpEndsWith, "test", " LIKE $"},
		{"contains", dsl.OpContains, "test", " LIKE $"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := setupTestRegistry()
			queryPlanner := planner.NewPlanner(reg)

			dslQuery := &dsl.Query{
				Model: "orders",
				Filters: &dsl.ComparisonFilter{
					Field: "status",
					Op:    tt.op,
					Value: tt.value,
				},
			}

			plan, err := queryPlanner.PlanQuery(dslQuery)
			if err != nil {
				t.Fatalf("PlanQuery error: %v", err)
			}

			builder := NewQueryBuilder()
			sqlIface, params, err := builder.BuildQuery(plan)

			if err != nil {
				t.Errorf("BuildQuery error: %v", err)
			}

			sql, ok := sqlIface.(string)
			if !ok {
				t.Fatalf("expected sql string, got %T", sqlIface)
			}
			if !strings.Contains(sql, tt.contains) {
				t.Errorf("SQL missing expected operator: %s", sql)
			}

			// Verify params contain the wildcard patterns
			hasWildcard := false
			for _, p := range params {
				if s, ok := p.(string); ok && (strings.Contains(s, "%") || strings.Contains(s, tt.value)) {
					hasWildcard = true
					break
				}
			}
			if !hasWildcard {
				t.Errorf("Params should contain wildcard pattern for %s", tt.name)
			}
		})
	}
}

func TestNewQueryBuilder(t *testing.T) {
	builder := NewQueryBuilder()
	if builder == nil {
		t.Errorf("NewQueryBuilder() returned nil")
	}
	if len(builder.params) != 0 {
		t.Errorf("Initial params should be empty")
	}
	if builder.paramCount != 0 {
		t.Errorf("Initial paramCount should be 0")
	}
}
