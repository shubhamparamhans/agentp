# Capabilities Matrix

Use this page to answer "do we already have this?" before starting a new feature.

| Capability | Current State | Primary Feature Docs | Primary Code Areas |
| --- | --- | --- | --- |
| Natural-language data querying | Planned | [ai-chatbot](/Users/shubhamparamhans/Workspace/udv/docs/features/ai-chatbot/index.md) | `frontend/src/pages/`, `frontend/src/components/`, `frontend/src/api/`, `internal/api/`, `internal/dsl/` |
| PostgreSQL querying | Available | [postgres-support](/Users/shubhamparamhans/Workspace/udv/docs/features/postgres-support/summary.md) | `internal/adapter/postgres/`, `internal/api/`, `internal/query/` |
| MongoDB querying | In progress | [mongodb-support](/Users/shubhamparamhans/Workspace/udv/docs/features/mongodb-support/summary.md) | `internal/adapter/mongodb/`, `internal/schema_processor/` |
| Query DSL | Available | [query-support](/Users/shubhamparamhans/Workspace/udv/docs/features/query-support/summary.md) | `internal/dsl/`, `internal/planner/` |
| Query planning | Available | [query-support](/Users/shubhamparamhans/Workspace/udv/docs/features/query-support/summary.md) | `internal/planner/`, `internal/ir/` |
| CRUD via query endpoint | Partial | [query-support](/Users/shubhamparamhans/Workspace/udv/docs/features/query-support/summary.md) | `internal/api/`, `internal/query/` |
| Schema discovery and model generation | Available | [data-modelling-processor-index](/Users/shubhamparamhans/Workspace/udv/docs/features/data-modelling-processor-index/summary.md) | `internal/schema_processor/`, `cmd/generate-models/`, `configs/` |
| Frontend list and group views | Available | [frontend-support](/Users/shubhamparamhans/Workspace/udv/docs/features/frontend-support/summary.md) | `frontend/src/components/`, `frontend/src/pages/` |
| Nested object rendering | Available | [object-rendering-implementation](/Users/shubhamparamhans/Workspace/udv/docs/features/object-rendering-implementation/summary.md) | `frontend/src/components/ObjectRenderer/`, `frontend/src/components/ListView/`, `frontend/src/components/DetailView/` |
| Server-side filtering | Available | [server-side-filtering-complete](/Users/shubhamparamhans/Workspace/udv/docs/features/server-side-filtering-complete/summary.md) | `internal/api/`, `internal/adapter/postgres/`, `frontend/src/components/` |
| PostgreSQL special type casting | Available | [type-casting-fix](/Users/shubhamparamhans/Workspace/udv/docs/features/type-casting-fix/summary.md) | `internal/adapter/postgres/`, `internal/planner/`, `internal/schema/` |

## Suggested Pre-Feature Check

Before creating a new feature folder:

1. Check this matrix for an existing capability.
2. Check the feature catalog for adjacent work.
3. Reuse or extend an existing feature folder if the new work is part of the same capability area.
