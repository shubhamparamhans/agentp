# feat-restructure-md-files — Summary

Branch: feat/restructure-md-files
Author: shubhamparamhans
Date: 2026-03-15
Status: draft

One-line summary
----------------
Restructure all existing project Markdown documentation into per-feature folders under `docs/features/` following project guidelines, without deleting originals until review.

Goals
-----
- Create a safe, auditable migration of existing `.md` docs into `docs/features/<feature>/detailed.md`.
- Add `summary.md` placeholders and `changelog.md` entries for each migrated document.
- Provide a dry-run script and clear review/rollback steps.

Acceptance criteria
-------------------
- Every existing `.md` file is copied to a corresponding `docs/features/<feature>/detailed.md`.
- `summary.md` exists for each feature folder (TODO placeholder allowed).
- No original files are removed during the initial migration.
