package postgres

import (
	"os"
	"testing"

	"udv/internal/config"
	"udv/internal/dsl"
	"udv/internal/planner"
	"udv/internal/schema"
)

// TestE2EQueryExecution tests the full pipeline: DSL -> SQL -> Execute on real DB
func TestE2EQueryExecution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Get connection string from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL environment variable not set; skipping E2E test")
	}

	// Connect to database
	db, err := Connect(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup schema registry from config
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "orders",
				Table:      "orders",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer"},
					{Name: "status", Type: "string"},
					{Name: "amount", Type: "decimal"},
					{Name: "created_at", Type: "timestamp"},
				},
			},
		},
	}

	registry := schema.NewRegistry()
	registry.LoadFromConfig(cfg)

	// Test 1: Simple SELECT query
	t.Run("Simple SELECT", func(t *testing.T) {
		q := &dsl.Query{
			Model: "orders",
			Pagination: &dsl.Pagination{
				Limit:  10,
				Offset: 0,
			},
		}

		validator := dsl.NewValidator(registry)
		if err := validator.ValidateQuery(q); err != nil {
			t.Fatalf("Validation error: %v", err)
		}

		planner := planner.NewPlanner(registry)
		plan, err := planner.PlanQuery(q)
		if err != nil {
			t.Fatalf("Planning error: %v", err)
		}

		builder := NewQueryBuilder()
		sqlIface, params, err := builder.BuildQuery(plan)
		if err != nil {
			t.Fatalf("SQL build error: %v", err)
		}

		sql, ok := sqlIface.(string)
		if !ok {
			t.Fatalf("expected sql string, got %T", sqlIface)
		}

		rows, err := db.Query(sql, params...)
		if err != nil {
			t.Fatalf("Query execution error: %v", err)
		}
		defer rows.Close()

		// Verify we got results
		count := 0
		for rows.Next() {
			count++
		}
		t.Logf("Simple SELECT returned %d rows", count)
	})

	// Test 2: SELECT with filter
	t.Run("SELECT with filter", func(t *testing.T) {
		q := &dsl.Query{
			Model: "orders",
			Filters: &dsl.ComparisonFilter{
				Field: "status",
				Op:    dsl.OpEqual,
				Value: "PAID",
			},
			Pagination: &dsl.Pagination{
				Limit:  10,
				Offset: 0,
			},
		}

		validator := dsl.NewValidator(registry)
		if err := validator.ValidateQuery(q); err != nil {
			t.Fatalf("Validation error: %v", err)
		}

		planner := planner.NewPlanner(registry)
		plan, err := planner.PlanQuery(q)
		if err != nil {
			t.Fatalf("Planning error: %v", err)
		}

		builder := NewQueryBuilder()
		sqlIface, params, err := builder.BuildQuery(plan)
		if err != nil {
			t.Fatalf("SQL build error: %v", err)
		}

		sql, ok := sqlIface.(string)
		if !ok {
			t.Fatalf("expected sql string, got %T", sqlIface)
		}
		t.Logf("Generated SQL: %s", sql)
		t.Logf("Parameters: %v", params)

		rows, err := db.Query(sql, params...)
		if err != nil {
			t.Fatalf("Query execution error: %v", err)
		}
		defer rows.Close()

		// Just verify it executes without error
		for rows.Next() {
			// consume rows
		}
		if err = rows.Err(); err != nil {
			t.Fatalf("Row iteration error: %v", err)
		}
	})

	// Test 3: SELECT with GROUP BY and aggregate
	t.Run("SELECT with GROUP BY and aggregate", func(t *testing.T) {
		q := &dsl.Query{
			Model:   "orders",
			GroupBy: []string{"status"},
			Aggregates: []dsl.Aggregate{
				{Function: dsl.AggCount, Field: "", Alias: "count"},
				{Function: dsl.AggSum, Field: "amount", Alias: "total"},
			},
			Pagination: &dsl.Pagination{
				Limit:  10,
				Offset: 0,
			},
		}

		validator := dsl.NewValidator(registry)
		if err := validator.ValidateQuery(q); err != nil {
			t.Fatalf("Validation error: %v", err)
		}

		planner := planner.NewPlanner(registry)
		plan, err := planner.PlanQuery(q)
		if err != nil {
			t.Fatalf("Planning error: %v", err)
		}

		builder := NewQueryBuilder()
		sqlIface, params, err := builder.BuildQuery(plan)
		if err != nil {
			t.Fatalf("SQL build error: %v", err)
		}

		sql, ok := sqlIface.(string)
		if !ok {
			t.Fatalf("expected sql string, got %T", sqlIface)
		}
		t.Logf("Generated SQL: %s", sql)

		rows, err := db.Query(sql, params...)
		if err != nil {
			t.Fatalf("Query execution error: %v", err)
		}
		defer rows.Close()

		count := 0
		for rows.Next() {
			count++
		}
		t.Logf("GROUP BY query returned %d rows", count)
	})

	// Test 4: SELECT with ORDER BY
	t.Run("SELECT with ORDER BY", func(t *testing.T) {
		q := &dsl.Query{
			Model: "orders",
			Sort: []dsl.Sort{
				{Field: "created_at", Direction: dsl.SortDesc},
			},
			Pagination: &dsl.Pagination{
				Limit:  5,
				Offset: 0,
			},
		}

		validator := dsl.NewValidator(registry)
		if err := validator.ValidateQuery(q); err != nil {
			t.Fatalf("Validation error: %v", err)
		}

		planner := planner.NewPlanner(registry)
		plan, err := planner.PlanQuery(q)
		if err != nil {
			t.Fatalf("Planning error: %v", err)
		}

		builder := NewQueryBuilder()
		sqlIface, params, err := builder.BuildQuery(plan)
		if err != nil {
			t.Fatalf("SQL build error: %v", err)
		}

		sql, ok := sqlIface.(string)
		if !ok {
			t.Fatalf("expected sql string, got %T", sqlIface)
		}

		rows, err := db.Query(sql, params...)
		if err != nil {
			t.Fatalf("Query execution error: %v", err)
		}
		defer rows.Close()

		for rows.Next() {
			// consume
		}
	})

	// Test 5: Complex query with filter, GROUP BY, aggregate, and ORDER BY
	t.Run("Complex query", func(t *testing.T) {
		q := &dsl.Query{
			Model:   "orders",
			GroupBy: []string{"status"},
			Filters: &dsl.ComparisonFilter{
				Field: "amount",
				Op:    dsl.OpGT,
				Value: 100,
			},
			Aggregates: []dsl.Aggregate{
				{Function: dsl.AggCount, Field: "", Alias: "order_count"},
				{Function: dsl.AggSum, Field: "amount", Alias: "total_amount"},
				{Function: dsl.AggAvg, Field: "amount", Alias: "avg_amount"},
			},
			Sort: []dsl.Sort{
				{Field: "status", Direction: dsl.SortDesc},
			},
			Pagination: &dsl.Pagination{
				Limit:  10,
				Offset: 0,
			},
		}

		validator := dsl.NewValidator(registry)
		if err := validator.ValidateQuery(q); err != nil {
			t.Fatalf("Validation error: %v", err)
		}

		planner := planner.NewPlanner(registry)
		plan, err := planner.PlanQuery(q)
		if err != nil {
			t.Fatalf("Planning error: %v", err)
		}

		builder := NewQueryBuilder()
		sqlIface, params, err := builder.BuildQuery(plan)
		if err != nil {
			t.Fatalf("SQL build error: %v", err)
		}

		sql, ok := sqlIface.(string)
		if !ok {
			t.Fatalf("expected sql string, got %T", sqlIface)
		}

		t.Logf("Complex SQL: %s", sql)
		t.Logf("Parameters: %v", params)

		rows, err := db.Query(sql, params...)
		if err != nil {
			t.Fatalf("Query execution error: %v", err)
		}
		defer rows.Close()

		rowCount := 0
		for rows.Next() {
			rowCount++
		}
		t.Logf("Complex query returned %d rows", rowCount)

		if err = rows.Err(); err != nil {
			t.Fatalf("Row iteration error: %v", err)
		}
	})
}

// TestE2EQueryWithDifferentOperators tests various filter operators against real data
func TestE2EQueryWithDifferentOperators(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL environment variable not set; skipping E2E test")
	}

	db, err := Connect(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "orders",
				Table:      "orders",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer"},
					{Name: "status", Type: "string"},
					{Name: "amount", Type: "decimal"},
					{Name: "created_at", Type: "timestamp"},
				},
			},
		},
	}

	registry := schema.NewRegistry()
	registry.LoadFromConfig(cfg)

	tests := []struct {
		name   string
		query  *dsl.Query
		expect string // expectation about result (e.g., "should succeed")
	}{
		{
			name: "Greater than operator",
			query: &dsl.Query{
				Model: "orders",
				Filters: &dsl.ComparisonFilter{
					Field: "amount",
					Op:    dsl.OpGT,
					Value: 1000,
				},
				Pagination: &dsl.Pagination{Limit: 10, Offset: 0},
			},
			expect: "should succeed",
		},
		{
			name: "Less than operator",
			query: &dsl.Query{
				Model: "orders",
				Filters: &dsl.ComparisonFilter{
					Field: "amount",
					Op:    dsl.OpLT,
					Value: 500,
				},
				Pagination: &dsl.Pagination{Limit: 10, Offset: 0},
			},
			expect: "should succeed",
		},
		{
			name: "Not equal operator",
			query: &dsl.Query{
				Model: "orders",
				Filters: &dsl.ComparisonFilter{
					Field: "status",
					Op:    dsl.OpNotEqual,
					Value: "PENDING",
				},
				Pagination: &dsl.Pagination{Limit: 10, Offset: 0},
			},
			expect: "should succeed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := dsl.NewValidator(registry)
			if err := validator.ValidateQuery(tt.query); err != nil {
				t.Fatalf("Validation error: %v", err)
			}

			planner := planner.NewPlanner(registry)
			plan, err := planner.PlanQuery(tt.query)
			if err != nil {
				t.Fatalf("Planning error: %v", err)
			}

			builder := NewQueryBuilder()
			sqlIface, params, err := builder.BuildQuery(plan)
			if err != nil {
				t.Fatalf("SQL build error: %v", err)
			}
			sql, ok := sqlIface.(string)
			if !ok {
				t.Fatalf("expected sql string, got %T", sqlIface)
			}

			rows, err := db.Query(sql, params...)
			if err != nil {
				t.Fatalf("Query execution error: %v\nSQL: %s\nParams: %v", err, sql, params)
			}
			defer rows.Close()

			// Consume rows
			for rows.Next() {
			}
		})
	}
}
