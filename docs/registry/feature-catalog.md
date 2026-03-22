# Feature Catalog

This is the main lookup table for what UDV already has, what is still being refined, and where the source docs live.

| Feature | Area | Status | Summary | Docs |
| --- | --- | --- | --- | --- |
| AI chatbot | Frontend, backend, query UX | Draft | Plain-text data questions translated into validated UDV queries with chat-oriented result rendering. | [ai-chatbot](/Users/shubhamparamhans/Workspace/udv/docs/features/ai-chatbot/index.md) |
| MongoDB support | Backend, schema, adapters | In progress | MongoDB adapter, schema discovery, modelling, and verification work. | [mongodb-support](/Users/shubhamparamhans/Workspace/udv/docs/features/mongodb-support/summary.md) |
| PostgreSQL support | Backend, schema, adapters | In progress | PostgreSQL adapter, SQL generation, data modelling processor, and type fixes. | [postgres-support](/Users/shubhamparamhans/Workspace/udv/docs/features/postgres-support/summary.md) |
| Frontend support | Frontend, UX | In progress | Frontend architecture, integration, theme, and setup consolidation. | [frontend-support](/Users/shubhamparamhans/Workspace/udv/docs/features/frontend-support/summary.md) |
| Query support | Query engine, API | Active | Query DSL, planner, and CRUD-via-query implementation references. | [query-support](/Users/shubhamparamhans/Workspace/udv/docs/features/query-support/summary.md) |
| Data modelling processor | Schema tooling | Needs doc backfill | CLI-based model generation and verification flow for database schemas. | [data-modelling-processor-index](/Users/shubhamparamhans/Workspace/udv/docs/features/data-modelling-processor-index/summary.md) |
| Object rendering | Frontend, MongoDB UX | Implemented | Expandable rendering for nested object values in the UI. | [object-rendering-implementation](/Users/shubhamparamhans/Workspace/udv/docs/features/object-rendering-implementation/summary.md) |
| Server-side filtering | Backend, frontend integration | Implemented | Filters now execute on the backend so pagination and large datasets behave correctly. | [server-side-filtering-complete](/Users/shubhamparamhans/Workspace/udv/docs/features/server-side-filtering-complete/summary.md) |
| Type casting fix | PostgreSQL, planner | Implemented | Explicit type handling for UUID and other special PostgreSQL types. | [type-casting-fix](/Users/shubhamparamhans/Workspace/udv/docs/features/type-casting-fix/summary.md) |

## Notes

- `Needs doc backfill` means the implementation exists, but the feature snapshot still needs a cleaner top-level summary.
- Legacy documentation-migration folders now live under `docs/archive/legacy-feature-docs/` and are no longer part of the active feature catalog.
