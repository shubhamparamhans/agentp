# Feature: MongoDB Query Parity

## Status
Draft

## Owner
Backend

## Problem
MongoDB support currently covers basic filtering, sorting, pagination, and CRUD, but it does not yet match PostgreSQL for core UDV query capabilities such as projection, grouping, aggregation, and safer record targeting.

## Scope
- Add MongoDB query projection support from planned selected fields
- Add aggregation pipeline generation for grouped and aggregated queries
- Add ID-based update and delete targeting parity with PostgreSQL
- Align MongoDB query behavior with the core planner/query model where practical
- Expand tests to cover parity-critical query flows

## Out Of Scope
- Frontend redesign
- Authentication or RBAC
- Saved views or exports
- Cross-database joins
- Advanced Mongo-only features that do not map to the shared planner model

## Code Areas
- `internal/adapter/mongodb/`
- `internal/planner/`
- `internal/dsl/`
- `internal/api/`

## Related Docs
- [PRD](./prd.md)
- [Design](./design.md)
- [Tasks](./tasks.md)
- [Status](./status.md)

## Current Capability
MongoDB currently supports basic `find`, CRUD operations, filter translation, sort, and pagination, but lacks builder support for projection and aggregation-pipeline based grouped queries.

## Next Steps
- Confirm the exact parity target versus PostgreSQL
- Implement projection support
- Implement aggregation pipeline generation
- Add parity-focused test coverage
