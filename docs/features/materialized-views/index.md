# Feature: Materialized Views

## Status
Draft

## Owner
Platform / Backend

## Problem
Users need a stable, report-friendly dataset that can be defined without writing custom code for every dashboard or recurring analysis workflow. Today the platform exposes raw models and relationships, but it does not provide a config-driven way to define a reusable, transformed, materialized view over one table or a relationship graph.

## Scope
- Define materialized views entirely through configuration.
- Support a single base model or a relationship-driven source graph.
- Support transformations such as field selection, renaming, casting, derived fields, filtering, grouping, and aggregations.
- Materialize view output into a consistent schema suitable for reporting and BI-like consumption.
- Support refresh strategies such as manual, on-read, and scheduled refresh.
- Expose view metadata and view data through existing API and frontend surfaces.
- Surface predefined views inside the existing UI through a dedicated `Views` menu.
- Validate view configs at startup and fail fast on invalid definitions.

## Out Of Scope
- Building a full Power BI-style visualization builder.
- User-authored SQL or arbitrary scripting in v1.
- Cross-database materialization in a single view.
- Row-level security policy design beyond existing platform controls.
- Incremental refresh optimization in the first release.

## Code Areas
- `configs/`
- `internal/config/`
- `internal/schema/`
- `internal/dsl/`
- `internal/planner/`
- `internal/query/`
- `internal/adapter/`
- `internal/api/`
- `frontend/src/`

## Related Docs
- [PRD](./prd.md)
- [Design](./design.md)
- [Tasks](./tasks.md)
- [Status](./status.md)

## Current Capability
The platform already supports config-driven models, schema discovery, relationship-aware querying, and frontend list/detail flows. It does not yet have a first-class concept for a reusable, transformed, persisted reporting dataset that is declared once in config and consumed as a stable view. It also does not yet expose a dedicated UI entry point for report-style predefined views.

## Next Steps
- Confirm the config schema and lifecycle for view refresh.
- Decide whether v1 ships as backend-only or includes frontend discovery and management surfaces.
- Align materialization storage strategy with adapter capabilities.
- Break implementation into parser, planner, executor, refresh, and API milestones.
