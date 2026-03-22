# Design

## Overview
Materialized views are introduced as a new config-driven resource layered on top of existing models and relationship metadata. A view definition describes:

- the source dataset,
- the relationship traversal path if joins are needed,
- an ordered transform pipeline,
- the output schema,
- the materialization strategy and refresh policy.

The backend validates these definitions at startup, converts them into an internal plan, executes materialization through the relevant adapter, and exposes the resulting read-only dataset through the API. The frontend consumes view metadata in the same way it consumes other datasets, with additional freshness and refresh-status indicators.

## Architecture
The feature spans configuration, planning, execution, storage, API exposure, and frontend discovery:

- `configs/` declares materialized views as first-class config resources.
- `internal/config/` parses and validates view definitions.
- `internal/schema/` resolves source models, fields, and relationship paths for output schema generation.
- `internal/dsl/` and `internal/planner/` translate declarative transforms into an executable plan.
- `internal/query/` and `internal/adapter/` execute refreshes and serve read access to materialized rows.
- `internal/api/` exposes catalog, metadata, refresh, and row-reading endpoints.
- `frontend/src/` integrates views into the existing application shell through a dedicated `Views` menu and read-only view pages.

This keeps the authoring model config-driven while reusing the current query and UI architecture rather than creating a parallel reporting subsystem.

### Proposed Config Shape

```yaml
views:
  - id: sales_report
    name: Sales Report
    description: Stable dataset for dashboard and exports
    source:
      model: orders
      relationships:
        - customer
        - region
    transforms:
      - type: filter
        where:
          order_status:
            in: [paid, shipped]
      - type: derive
        fields:
          gross_amount: "quantity * unit_price"
      - type: group
        by: [region.name, customer.segment]
        metrics:
          total_orders: { op: count }
          total_revenue: { op: sum, field: gross_amount }
    output:
      fields:
        - source: region.name
          as: region_name
        - source: customer.segment
          as: customer_segment
        - source: total_orders
        - source: total_revenue
    materialization:
      mode: persisted
      refresh: scheduled
      schedule: "0 * * * *"
```

## Code Areas
- `configs/`
- `internal/config/`
- `internal/schema/`
- `internal/dsl/`
- `internal/planner/`
- `internal/query/`
- `internal/adapter/`
- `internal/api/`
- `frontend/src/pages/`
- `frontend/src/components/`

## Data Flow
1. Config loader reads view definitions alongside model definitions.
2. Config validator resolves referenced models, fields, and relationships.
3. Planner converts each view into an intermediate representation for projection, joins, filters, derivations, and aggregation.
4. Materialization executor builds or refreshes the persisted dataset through the active adapter.
5. Metadata store tracks generation, refresh timestamp, status, and row counts.
6. API exposes view catalog, schema, metadata, and rows.
7. Frontend renders views as read-only data sources with freshness information.
8. Users access views from a dedicated `Views` menu inside the existing UI shell.

## API Or Interface Changes
- Add config schema support for `views`.
- Add backend registry for loaded materialized-view definitions.
- Add read endpoints such as:
  - `GET /api/views`
  - `GET /api/views/:id`
  - `GET /api/views/:id/data`
  - `POST /api/views/:id/refresh`
- Add frontend navigation or dataset discovery support for views.
- Mark view data as read-only in frontend interactions.

## UI Integration
- Add a `Views` menu to the current navigation model.
- Populate that menu from the backend view catalog so all predefined config-driven views are discoverable.
- Reuse the existing page shell, table rendering, filters, and detail affordances where possible.
- Show view-level metadata such as description, last refresh, refresh status, and row count near the dataset header.
- Prevent create, edit, and delete record actions from view pages because views are read-only surfaces.

## Storage And Execution Notes
- Prefer adapter-native persisted tables or materialized objects where supported.
- Where adapter-native materialized views are not available, use a managed physical table plus swap-on-refresh behavior.
- Keep the logical view definition backend-agnostic even if execution differs by adapter.
- Store metadata separately from the materialized rows when practical so refresh state remains consistent.

## Validation Rules
- Source model must exist.
- Relationship path must resolve from the declared source.
- Output field names must be unique.
- Transform order must be valid, for example aggregate outputs cannot be referenced before they exist.
- Refresh configuration must match supported modes.

## Alternatives Considered
- Live virtual views only:
  rejected because BI-style usage needs stable performance and repeatable output.
- User-authored SQL definitions:
  rejected for v1 because it weakens portability, safety, and config validation.
- Frontend-only saved reports:
  rejected because this feature needs backend-managed consistency and reusable datasets.

## Rollout Notes
- Start with one adapter path that can reliably support persisted refreshes.
- Hide the feature behind a config flag until validation and refresh behavior are stable.
- Launch with manual and scheduled refresh before attempting incremental refresh.
- Document fallback behavior when a backend supports only managed tables rather than native materialized views.
- Roll out UI navigation only after the backend view catalog and metadata APIs are stable enough to avoid broken menu entries.
