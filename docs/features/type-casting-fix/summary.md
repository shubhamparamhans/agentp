# Type Casting Fix — Summary

Branch: feat/postgres-support
Status: implemented

One-line summary
----------------
Adds explicit type handling for UUID, JSON, binary, and timestamp-like PostgreSQL fields so query building and grouping work reliably with generated models.

Code areas
----------
- `internal/adapter/postgres/`
- `internal/planner/`
- `internal/schema/`
