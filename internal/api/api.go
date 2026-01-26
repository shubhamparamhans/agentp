package api

import (
    "encoding/json"
    "fmt"
    "net/http"

    "udv/internal/adapter/postgres"
    "udv/internal/dsl"
    "udv/internal/planner"
    "udv/internal/schema"
)

// API bundles dependencies for HTTP handlers
type API struct {
    registry  *schema.Registry
    validator *dsl.Validator
    planner   *planner.Planner
    builder   *postgres.QueryBuilder
    db        *postgres.Database
}

// New creates a new API instance with optional database connection
func New(reg *schema.Registry, db *postgres.Database) *API {
    return &API{
        registry:  reg,
        validator: dsl.NewValidator(reg),
        planner:   planner.NewPlanner(reg),
        builder:   postgres.NewQueryBuilder(),
        db:        db,
    }
}

// RegisterRoutes registers HTTP handlers onto the provided mux
func (a *API) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/models", a.handleModels)
    mux.HandleFunc("/query", a.handleQuery)
}

// handleModels returns a JSON list of models and their fields
func (a *API) handleModels(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    models := a.registry.ListModels()

    type fieldResp struct {
        Name string `json:"name"`
        Type string `json:"type"`
    }

    type modelResp struct {
        Name       string      `json:"name"`
        Table      string      `json:"table"`
        PrimaryKey string      `json:"primary_key"`
        Fields     []fieldResp `json:"fields"`
    }

    var out []modelResp
    for _, m := range models {
        fields, _ := a.registry.GetModelFields(m)
        md := a.registry.GetModel(m)
        fr := modelResp{
            Name:       m,
            Table:      "",
            PrimaryKey: "",
            Fields:     []fieldResp{},
        }
        if md != nil {
            fr.Table = md.Table
            fr.PrimaryKey = md.PrimaryKey
        }
        for _, f := range fields {
            fr.Fields = append(fr.Fields, fieldResp{Name: f.Name, Type: f.Type})
        }
        out = append(out, fr)
    }

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(out)
}

// handleQuery accepts a DSL query JSON, validates, plans, and returns SQL+params
func (a *API) handleQuery(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // Decode into a raw structure so we can handle the FilterExpr interface
    type rawQuery struct {
        Model      string          `json:"model"`
        Fields     []string        `json:"fields,omitempty"`
        Filters    json.RawMessage `json:"filters,omitempty"`
        GroupBy    []string        `json:"group_by,omitempty"`
        Aggregates []dsl.Aggregate `json:"aggregates,omitempty"`
        Sort       []dsl.Sort      `json:"sort,omitempty"`
        Pagination *dsl.Pagination `json:"pagination,omitempty"`
    }

    var rq rawQuery
    if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
        http.Error(w, fmt.Sprintf("invalid request body: %v", err), http.StatusBadRequest)
        return
    }

    q := dsl.Query{
        Model:      rq.Model,
        Fields:     rq.Fields,
        GroupBy:    rq.GroupBy,
        Aggregates: rq.Aggregates,
        Sort:       rq.Sort,
        Pagination: rq.Pagination,
    }

    // Parse filters if provided
    if len(rq.Filters) > 0 {
        // try ComparisonFilter first
        var cf dsl.ComparisonFilter
        if err := json.Unmarshal(rq.Filters, &cf); err == nil && cf.Field != "" {
            q.Filters = &cf
        } else {
            var lf dsl.LogicalFilter
            if err := json.Unmarshal(rq.Filters, &lf); err == nil {
                q.Filters = &lf
            } else {
                http.Error(w, "invalid filters format", http.StatusBadRequest)
                return
            }
        }
    }

    if err := a.validator.ValidateQuery(&q); err != nil {
        http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
        return
    }

    plan, err := a.planner.PlanQuery(&q)
    if err != nil {
        http.Error(w, fmt.Sprintf("planning error: %v", err), http.StatusInternalServerError)
        return
    }

    sql, params, err := a.builder.BuildQuery(plan)
    if err != nil {
        http.Error(w, fmt.Sprintf("sql build error: %v", err), http.StatusInternalServerError)
        return
    }

    resp := map[string]interface{}{
        "sql":    sql,
        "params": params,
    }

    // Execute query if database is available
    if a.db != nil {
        rows, err := a.db.ExecuteAndFetchRows(sql, params...)
        if err != nil {
            fmt.Printf("Warning: Failed to execute query: %v\n", err)
            // Don't fail the request - still return SQL and params
            // Frontend can see what query would have been executed
        } else {
            resp["data"] = rows
        }
    }

    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(resp)
}
