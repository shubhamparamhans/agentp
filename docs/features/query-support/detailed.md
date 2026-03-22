# Query Support — Detailed

This feature area covers the query contract used between the frontend, API layer, planners, and database adapters.

## Focus Areas

- Query DSL shape and validation
- Query planning and intermediate representation
- CRUD behavior through query-oriented endpoints
- Adapter-facing query execution expectations

## Primary Code Areas

- `internal/dsl/query.go`
- `internal/planner/planner.go`
- `internal/ir/plan.go`
- `internal/query/builder.go`
- `internal/api/api.go`

## Reference Assets

Detailed historical notes are preserved under `assets/`.
