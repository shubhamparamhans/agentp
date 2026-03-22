# Design

## Overview
The MongoDB adapter should support the same core UDV exploration contract as PostgreSQL wherever the shared planner model already expresses that capability. The implementation should favor shared semantics over Mongo-specific shortcuts so the rest of the stack can stay adapter-agnostic.

## Code Areas
- `internal/adapter/mongodb/builder.go`
- `internal/adapter/mongodb/db.go`
- `internal/adapter/mongodb/types.go`
- `internal/adapter/mongodb/builder_test.go`
- `internal/planner/`

## Data Flow
- DSL query is validated and planned into a shared `QueryPlan`
- MongoDB builder decides whether the plan should become:
  - a `find` query with filter, projection, sort, and pagination, or
  - an `aggregate` query with a pipeline for group and aggregate flows
- MongoDB DB execution layer runs the query and returns generic row maps

## API Or Interface Changes
- No public API contract change is required if parity is implemented behind the shared planner model
- Internal `MongoQuery` usage may need clearer support for projection and aggregate pipelines
- Update/delete builder behavior should accept `plan.ID` consistently, similar to PostgreSQL

## Alternatives Considered
- Keep MongoDB intentionally limited to document browsing only
  - Rejected because UDV aims for a generic exploration surface across supported backends
- Build Mongo-only semantics that bypass the planner model
  - Rejected because it increases adapter divergence and frontend inconsistency

## Rollout Notes
- Start with read-path parity before broadening mutation semantics
- Keep tests focused on planner-driven parity scenarios
- Update the capabilities matrix after implementation to reflect the new MongoDB support level
