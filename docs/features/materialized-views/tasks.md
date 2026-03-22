# Tasks

## Planning
- [ ] Confirm config schema for `views`
- [ ] Confirm supported transform set for v1
- [ ] Confirm refresh modes and operational expectations
- [ ] Confirm adapter support matrix and first target backend
- [ ] Confirm API surface and frontend discovery scope
- [ ] Confirm `Views` menu placement in the existing UI

## Implementation
- [ ] Add `views` definitions to runtime config loading
- [ ] Add config validation for source models, fields, relationships, and transform ordering
- [ ] Add planner support for view-specific projection, joins, derive, filter, group, and aggregate stages
- [ ] Add materialization executor and refresh orchestration
- [ ] Add adapter contracts for persisted view storage and atomic replacement
- [ ] Add metadata tracking for freshness, generation, status, and errors
- [ ] Add API endpoints for catalog, metadata, data read, and refresh trigger
- [ ] Add frontend support to list and inspect materialized views as read-only datasets
- [ ] Add a `Views` menu that lists all predefined configured views

## Validation
- [ ] Unit tests for config parsing and validation failures
- [ ] Unit tests for transform planning and output schema generation
- [ ] Integration tests for single-model view materialization
- [ ] Integration tests for relationship-based view materialization
- [ ] Integration tests for refresh success and failed-refresh fallback
- [ ] Manual verification using a representative reporting dataset

## Docs
- [x] Create feature snapshot and PRD
- [x] Create design, task list, and status docs
- [ ] Update registry pages if this feature moves from draft to active planning
- [ ] Add sample config docs once schema is finalized
