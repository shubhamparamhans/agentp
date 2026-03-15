# Feature Documentation Template

Fill out the sections below for each feature. Copy this file into `docs/features/<branch-name>/detailed.md` or `summary.md` as appropriate.

Feature Name: <short-name>
Branch: <git-branch-name>
Author: <name>
Date: <YYYY-MM-DD>
Status: (draft|in-progress|completed|merged)

Short Summary
-------------
Write a 1-2 sentence summary of the feature and its purpose.

Goals
-----
- Primary goal 1
- Primary goal 2

Scope
-----
- What is included
- What is explicitly out of scope

Design & Rationale
-------------------
Describe the chosen design, alternatives considered, and reasons for the decisions.

API / Interface Changes
-----------------------
- Endpoints added/changed
- Request/response examples

Database / Migration
--------------------
- Collections / tables added or modified
- Migration steps and rollbacks

Configuration
-------------
- New env vars or config flags

Build & Run
-----------
How to build/run locally to test the feature.

Testing & Verification
----------------------
- Unit tests added
- Integration tests
- Manual verification steps

Rollback Plan
-------------
Steps to revert the feature safely if needed.

Related PRs / Issues
--------------------
- Link to related PRs and issue trackers

Files included in this folder
-----------------------------
- `summary.md` — one-paragraph summary
- `detailed.md` — full implementation notes
- `changelog.md` — incremental notes
- `assets/` — diagrams, screenshots, fixtures
