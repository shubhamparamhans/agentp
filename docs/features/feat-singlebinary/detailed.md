# Feature: feat-singlebinary

Feature Name: feat-singlebinary
Branch: feat/singlebinary
Author: shubhamparamhans
Date: 2026-03-15
Status: in-progress

Short Summary
-------------
Create a single, self-contained binary that runs the backend API, performs model-json generation (current `cmd/generate-models`), and serves the frontend static files — removing the need to run multiple processes in development or production.

Goals
-----
- Bundle backend server, model generator, and frontend static assets into one runnable binary.
- Provide command-line flags to run sub-commands: `serve`, `generate-models`, `migrate`, `version`.
- Support environment-driven configuration so the same binary can run in dev, staging, and production.

Scope
-----
Included:
- Main HTTP API and static file serving for frontend.
- On-demand or build-time model JSON generation integrated as a subcommand.
- Configuration via environment variables and config file.

Out of scope:
- Rewriting the frontend to be embedded differently — initial approach will embed built static files.
- Large refactors of business logic.

Design & Rationale
-------------------
Approach:
- Use Go's `embed` package to include compiled frontend assets when building release binaries.
- Convert `cmd/generate-models` logic into an internal package callable from the main binary and expose it via a subcommand.
- Use a small command router (e.g., `cobra` or `urfave/cli`) to provide subcommands.

Alternatives considered:
- Container-only distribution (keeps separate services) — rejected because single-binary simplifies edge deployments.

API / Interface Changes
-----------------------
- No external API changes required for the HTTP endpoints; optional admin endpoints added for triggering model regeneration.
- CLI: `udv serve --port=8080`, `udv generate-models --out=./configs/models.json`, `udv migrate up`.

Database / Migration
--------------------
- No new schema changes expected. Migration commands (if used) will still call the existing migration code.

Configuration
-------------
- Environment variables to control mode: `UDV_ENV=(development|staging|production)`
- `UDV_PORT`, `UDV_DB_URL`, `UDV_GENERATE_ON_START=true|false`

Build & Run
-----------
Build (dev):
```
go build -o udv ./cmd/server
```
Build (release, embedding frontend):
```
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o udv ./cmd/server
```
Run:
```
./udv serve --port=8080
./udv generate-models --out=./configs/models.json
```

Testing & Verification
----------------------
- Unit tests for any refactored packages.
- Integration test: run `udv serve` and confirm frontend is served and API endpoints respond.
- Manual: verify `./udv generate-models` produces the same JSON as current CLI.

Rollback Plan
-------------
- Revert the merge and restore the separate `cmd/` binaries if regression occurs.
- Provide a migration compatibility layer for config flags for a few releases.

Related PRs / Issues
--------------------
- TBD: link PR and issue numbers here.

Files included in this folder
-----------------------------
- `summary.md` — one-paragraph summary
- `detailed.md` — this file
- `changelog.md` — incremental notes
- `PRD.md` — product requirements document
- `assets/` — screenshots or examples
