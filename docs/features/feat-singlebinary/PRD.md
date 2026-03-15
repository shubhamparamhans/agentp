# PRD — Single Binary Distribution

Objective
---------
Deliver a single executable that runs the full UDV application stack (backend API, model generation CLI, and embedded frontend static server) to simplify deployment and developer experience.

Background
----------
Currently the project requires running multiple processes: backend server, model-generator CLI, and a separate frontend dev server. Packaging everything into one binary reduces operational complexity.

User Stories
------------
- As a developer, I can run `./udv serve` to start backend and see the frontend without starting another process.
- As a DevOps engineer, I can deploy a single binary to servers or containers.
- As a maintainer, I can run `./udv generate-models` to regenerate model JSON on demand.

Acceptance Criteria
-------------------
- One binary named `udv` can be built and run to serve API and frontend.
- `udv generate-models` produces identical output to existing CLI.
- Runtime flags and env vars documented and backwards-compatible.
- Automated tests validate serving and generation in CI.

Non-functional Requirements
---------------------------
- Binary size should be reasonable (strip symbols in release builds).
- Startup latency acceptable (< 2s for simple envs).
- Cross-platform builds supported for Linux and macOS.

Milestones
----------
1. Design & PRD (this doc) — complete
2. Refactor `cmd/generate-models` into internal package
3. Implement embedding of frontend assets
4. Add CLI router and subcommands
5. Tests, CI, and release build scripts

Risks & Mitigations
-------------------
- Risk: Binary bloat due to embedded assets — mitigate with compression and using minified assets.
- Risk: Increased complexity in one binary — keep a clear command structure and subcommands.

Metrics of Success
------------------
- Reduction in deployment steps documented
- Successful release builds for target platforms
