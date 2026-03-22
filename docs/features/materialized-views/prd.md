# PRD

## Objective
Introduce config-driven materialized views that let teams define stable reporting datasets over a base table or relationship graph, including transformations and aggregations, without writing custom code per report.

## Users
- Product and operations teams who need consistent report datasets.
- Backend developers and solution engineers configuring data products.
- Frontend consumers that need predictable, report-friendly schemas.
- External BI or analytics tools consuming stable API-backed data.

## Problem Statement
Raw entity models are useful for CRUD and exploration, but reporting workflows usually need denormalized, transformed, and consistent datasets. Without a materialized-view capability, every reporting use case must recreate joins, derived fields, filters, and aggregation logic repeatedly. That causes drift, duplicated logic, inconsistent metrics, and slower query performance for dashboards that expect stable output.

## Success Criteria
- A user can define a reporting view entirely through config, with no code changes required for the happy path.
- A view can source data from one model or a declared relationship path across multiple models.
- The system supports core transformations: select, alias, type cast, derived expression, filter, group-by, aggregate, sort, and default limit behavior.
- Invalid view configs fail validation before runtime traffic is served.
- Materialized output can be refreshed manually and by schedule in v1, with clear refresh status surfaced by the API.
- Report consumers receive a stable schema and predictable field names across refreshes unless config changes.
- Query latency for repeated dashboard reads improves relative to equivalent live relationship traversal.

## Scope

### In Scope
- New config schema for materialized view definitions.
- View definitions that reference existing models and relationships.
- Transform pipeline with deterministic ordering.
- Materialization lifecycle: build, refresh, replace, and read.
- Metadata fields such as last refresh time, source definition version, row count, and refresh status.
- API support to list views, inspect view definitions, read rows, and trigger refresh.
- Frontend support to discover views and display them similarly to read-only datasets.
- A dedicated `Views` menu in the existing UI that lists all predefined views from config.
- Adapter-aware materialization strategy for supported backends.
- Tests covering config validation, planning, execution, and refresh behavior.

### Out Of Scope
- Arbitrary user-defined SQL or JavaScript/Python transforms.
- Visual drag-and-drop report design.
- Write-back or editable materialized views.
- Cost-based optimizer work beyond basic execution planning.
- Fully incremental refresh, CDC-based refresh, or lineage visualization in v1.

## Functional Requirements

### View Definition
- Each view must have a unique id, display name, description, and owner.
- Each view must declare a source model or source model plus relationship traversal path.
- Each view must declare output fields explicitly or use a controlled include pattern.
- Each view may declare transformations in ordered stages.
- Each view must declare a materialization mode and refresh policy.

### Transformation Support
- Projection: choose source fields to expose.
- Alias: rename fields for consumer-friendly output.
- Derivation: create computed fields from supported expressions.
- Cast: normalize output types for report stability.
- Filter: apply static predicates before materialization.
- Join traversal: pull fields through configured relationships.
- Aggregate: support count, sum, avg, min, max, and grouped outputs.
- Sort: define default ordering for materialized rows.

### Materialization Behavior
- Views must materialize into a durable backend representation or managed cache abstraction.
- Refresh operations must be atomic from a consumer perspective.
- A failed refresh must not corrupt the last successful materialized dataset.
- The system must surface refresh status and error details.

### Consumption
- Consumers must be able to query a materialized view as a read-only dataset.
- Consumers must receive metadata describing freshness and schema.
- Frontend list and detail surfaces must treat views as read-only resources.
- The current UI must expose a dedicated `Views` menu that lists all configured predefined views.
- Selecting a view from the `Views` menu must open the view inside the existing UI shell, not a separate reporting application.

### UI Requirements
- The feature must integrate into the current frontend information architecture.
- The navigation must include a `Views` menu or equivalent first-class entry point.
- The `Views` menu must render all predefined config-driven views available to the current environment.
- Each view page must show:
  - view name and description,
  - refresh status and freshness metadata,
  - read-only tabular output using existing list/detail patterns where practical.
- Views must be visually consistent with the current UI rather than introducing a separate report-builder experience.

## Non-Functional Requirements
- Startup validation should reject invalid view definitions with actionable errors.
- Refresh behavior should be observable through logs and metrics.
- View reads should prioritize predictable latency over fully live source freshness.
- Config changes should be versioned and traceable to a refresh generation.
- UI integration should reuse existing frontend patterns so the feature feels native to the current product.

## User Stories
- As an operator, I want to define a daily sales view with joins and derived metrics in config so dashboards read a stable dataset.
- As a developer, I want config validation to catch invalid relationship paths before deployment.
- As a frontend user, I want materialized views to appear as read-only datasets with clear freshness information.
- As an analyst, I want consistent column names and types so external reporting tools do not break between refreshes.

## Acceptance Criteria
- Given a valid config for a single-model view, the system creates and serves the materialized dataset.
- Given a valid config for a relationship-based view, the system resolves joins and produces the expected flattened output.
- Given a config with ordered transforms, the output reflects those transforms deterministically.
- Given an invalid source field, relationship path, or transform, startup fails with a precise validation error.
- Given a scheduled refresh policy, the system records refresh attempts and exposes the latest refresh outcome.
- Given a failed refresh, consumers can still read the last successful materialized snapshot.

## Risks
- Different adapters may support materialization differently, making a single abstraction too leaky.
- Transform-expression design can become too limited or too complex if not constrained carefully.
- Relationship-driven views may create expensive plans if aggregation and filtering are not applied early.
- Refresh orchestration adds operational complexity and can create stale-data confusion if metadata is weak.

## Dependencies
- Existing model and relationship configuration in `configs/`.
- Config parsing and validation in `internal/config/`.
- Query planning and execution in `internal/dsl/`, `internal/planner/`, and `internal/query/`.
- Backend-specific persistence support in `internal/adapter/`.
- API exposure in `internal/api/`.
- Frontend dataset rendering in `frontend/src/`.
