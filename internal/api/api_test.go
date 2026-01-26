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

func setupRegistryForTest() *schema.Registry {
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
                },
            },
        },
    }

    reg := schema.NewRegistry()
    reg.LoadFromConfig(cfg)
    return reg
}

func TestModelsEndpoint(t *testing.T) {
    reg := setupRegistryForTest()
    a := New(reg, nil)  // ← Pass nil for database since we're testing without DB
    mux := http.NewServeMux()
    a.RegisterRoutes(mux)

    ts := httptest.NewServer(mux)
    defer ts.Close()

    resp, err := http.Get(ts.URL + "/models")
    if err != nil {
        t.Fatalf("GET /models failed: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("unexpected status: %d", resp.StatusCode)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    var out []map[string]interface{}
    if err := json.Unmarshal(body, &out); err != nil {
        t.Fatalf("invalid json response: %v", err)
    }

    if len(out) != 1 {
        t.Fatalf("expected 1 model, got %d", len(out))
    }
}

func TestQueryEndpoint_Simple(t *testing.T) {
    reg := setupRegistryForTest()
    a := New(reg, nil)  // ← Pass nil for database since we're testing without DB
    mux := http.NewServeMux()
    a.RegisterRoutes(mux)

    ts := httptest.NewServer(mux)
    defer ts.Close()

    q := dsl.Query{
        Model:  "orders",
        Fields: []string{"id", "status"},
        Filters: &dsl.ComparisonFilter{
            Field: "status",
            Op:    dsl.OpEqual,
            Value: "PAID",
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

    if _, ok := out["sql"]; !ok {
        t.Fatalf("response missing sql")
    }
    if params, ok := out["params"].([]interface{}); !ok || len(params) == 0 {
        t.Fatalf("response missing params or empty")
    }
}
