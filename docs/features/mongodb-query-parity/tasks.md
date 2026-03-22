# Tasks

## Planning
- [ ] Confirm parity target against PostgreSQL behavior
- [ ] Decide which grouped/aggregate semantics are guaranteed for MongoDB
- [ ] Decide whether update/delete should default to single-record or multi-record operations

## Implementation
- [ ] Add projection support for selected fields in MongoDB `find`
- [ ] Add aggregation pipeline generation for grouped queries
- [ ] Add aggregate function support in MongoDB builder
- [ ] Add `plan.ID` handling for update operations
- [ ] Add `plan.ID` handling for delete operations
- [ ] Ensure execution layer cleanly supports generated aggregate queries

## Validation
- [ ] Add MongoDB tests for projection
- [ ] Add MongoDB tests for group-by and aggregates
- [ ] Add MongoDB tests for ID-based update/delete behavior
- [ ] Compare parity-critical query flows with PostgreSQL expectations

## Docs
- [ ] Update `docs/registry/feature-catalog.md`
- [ ] Update `docs/registry/capabilities-matrix.md`
- [ ] Update `docs/registry/features-by-status.md`
- [ ] Preserve parity analysis notes in `assets/`
