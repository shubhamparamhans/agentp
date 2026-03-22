# PRD — Restructure existing Markdown docs into feature folders

Objective
---------
Provide a reliable, reviewable process to consolidate scattered documentation into the `docs/features/` layout so that each feature's documentation is grouped and discoverable.

Acceptance Criteria
-------------------
- All existing `.md` files have corresponding copies under `docs/features/` after running the migration script (dry-run validated first).
- No original file is deleted during the process.
- `ALL_FEATURES_SUMMARY.md` is updated with links for migrated features.

Milestones
----------
1. Documentation and review of approach (this PRD).
2. Dry-run migration and review of mapping.
3. Execution of migration and PR for review.

Risks
-----
- Collisions in cleaned names — mitigation: require manual review for duplicates.
