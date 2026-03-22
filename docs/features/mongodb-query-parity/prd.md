# PRD

## Objective
Bring MongoDB adapter capabilities closer to PostgreSQL for the core UDV query flows so that Mongo-backed models can participate in the same exploration experience with fewer backend-specific limitations.

## Users
- Developers exploring MongoDB-backed data
- Internal teams using UDV as a generic viewer across data sources

## Problem Statement
UDV's product model assumes a shared query contract across adapters, but the MongoDB adapter currently implements a narrower subset than PostgreSQL. This creates gaps in grouped exploration, aggregate views, selected-field rendering, and safe mutation targeting.

## Success Criteria
- MongoDB supports field projection for selected columns
- MongoDB supports grouped and aggregate queries through an aggregation pipeline
- MongoDB update/delete can target by record ID in a predictable way
- MongoDB tests cover the same important query shapes already validated for PostgreSQL
- The documented capabilities matrix can describe MongoDB more confidently after implementation

## Scope

### In Scope
- Read-path parity for projection, grouping, aggregation, sort, and pagination
- Safer mutation targeting for update/delete
- Builder and execution-layer support required for the above
- Tests and documentation updates

### Out Of Scope
- Full relation-join parity with relational adapters
- MongoDB-specific advanced operators not represented in the shared planner
- UI-only enhancements unrelated to adapter parity

## Risks
- The shared planner model may not map cleanly to all MongoDB query shapes
- Grouping semantics may differ from PostgreSQL in edge cases
- Mutation parity decisions may require product-level constraints, not only adapter code

## Dependencies
- Existing planner/query IR behavior
- MongoDB adapter query types and execution layer
- Capability documentation in `docs/registry/`
