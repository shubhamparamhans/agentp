package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"udv/internal/config"
	"udv/internal/dsl"
	"udv/internal/schema"
)

func TestCreateEndpoint_Simple(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "users",
				Table:      "users",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "uuid", Nullable: false},
					{Name: "email", Type: "string", Nullable: false},
					{Name: "name", Type: "string", Nullable: true},
					{Name: "created_at", Type: "timestamp", Nullable: true},
				},
			},
		},
	}

	reg := schema.NewRegistry()
	reg.LoadFromConfig(cfg)

	a := New(reg, nil)
	mux := http.NewServeMux()
	a.RegisterRoutes(mux)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	q := dsl.Query{
		Operation: dsl.OpCreate,
		Model:     "users",
		Data: map[string]interface{}{
			"email": "test@example.com",
			"name":  "Test User",
		},
	}

	b, _ := json.Marshal(q)
	resp, err := http.Post(ts.URL+"/query", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST /query failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("unexpected status: %d body: %s", resp.StatusCode, string(body))
	}

	var out map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	// Check that SQL is generated
	sql, ok := out["sql"].(string)
	if !ok || sql == "" {
		t.Fatalf("no SQL generated in response: %v", out)
	}

	// Check that SQL contains INSERT
	if !contains(sql, "INSERT") {
		t.Fatalf("expected INSERT in SQL, got: %s", sql)
	}

	// Check that parameters are present
	params, ok := out["params"].([]interface{})
	if !ok {
		t.Fatalf("expected params in response")
	}

	if len(params) != 2 {
		t.Fatalf("expected 2 parameters, got %d", len(params))
	}

	t.Logf("SQL: %s", sql)
	t.Logf("Params: %v", params)
}

func TestCreateEndpoint_WithAllFields(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "orders",
				Table:      "orders",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "uuid", Nullable: false},
					{Name: "user_id", Type: "uuid", Nullable: true},
					{Name: "status", Type: "string", Nullable: false},
					{Name: "amount", Type: "decimal", Nullable: false},
					{Name: "metadata", Type: "json", Nullable: true},
					{Name: "created_at", Type: "timestamp", Nullable: true},
				},
			},
		},
	}

	reg := schema.NewRegistry()
	reg.LoadFromConfig(cfg)

	a := New(reg, nil)
	mux := http.NewServeMux()
	a.RegisterRoutes(mux)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	q := dsl.Query{
		Operation: dsl.OpCreate,
		Model:     "orders",
		Data: map[string]interface{}{
			"user_id":  "550e8400-e29b-41d4-a716-446655440000",
			"status":   "PENDING",
			"amount":   "150.50",
			"metadata": map[string]interface{}{"key": "value"},
		},
	}

	b, _ := json.Marshal(q)
	resp, err := http.Post(ts.URL+"/query", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST /query failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Fatalf("unexpected status: %d body: %s", resp.StatusCode, string(body))
	}

	var out map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	// Check that SQL is generated
	sql, ok := out["sql"].(string)
	if !ok || sql == "" {
		t.Fatalf("no SQL generated in response: %v", out)
	}

	// Verify all fields are in the INSERT
	requiredFields := []string{"user_id", "status", "amount", "metadata"}
	for _, field := range requiredFields {
		if !contains(sql, field) {
			t.Fatalf("expected field %s in SQL, got: %s", field, sql)
		}
	}

	t.Logf("SQL: %s", sql)
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
