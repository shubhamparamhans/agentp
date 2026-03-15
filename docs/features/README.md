# Feature Documentation Guidelines

Goal
----
This folder contains a standard layout and guidelines for documenting each feature developed on its own branch. Each feature must have its documentation placed in a dedicated subfolder named after the feature branch. This makes reviews, changelogs and rollbacks clearer.

Folder structure
----------------
- `docs/features/<branch-name>/` — documentation for a single feature branch
  - `summary.md` — one-paragraph summary and status
  - `detailed.md` — implementation details, rationale, design decisions
  - `changelog.md` — incremental notes, commits or PR references
  - `assets/` — screenshots, diagrams, SQL dumps, or other binaries

- `docs/features/summary/` — aggregated summaries and indices
  - `ALL_FEATURES_SUMMARY.md` — master index of feature folders

How to use
----------
1. When creating a feature branch, add a folder under `docs/features/` that matches the branch name (replace `/` with `-` if needed).
2. Add at minimum a `summary.md` describing the feature goal and status.
3. Populate `detailed.md` as implementation progresses; keep `changelog.md` for stepwise notes.
4. Commit documentation to the same feature branch so docs travel with code.
5. When merging the feature, update `docs/features/summary/ALL_FEATURES_SUMMARY.md` with a link and one-line summary.

Naming conventions
------------------
- Use the branch name as the folder name (e.g., `feat/singlebinary` -> `feat-singlebinary`).
- Filenames must be lowercase and use hyphens for spaces.

Checklist for new feature docs
-----------------------------
- [ ] Create `docs/features/<branch-name>/summary.md` (required)
- [ ] Create `docs/features/<branch-name>/detailed.md` (recommended)
- [ ] Create `docs/features/<branch-name>/changelog.md` (optional)
- [ ] Add assets to `assets/` if applicable
- [ ] After merge, add one-line entry to `docs/features/summary/ALL_FEATURES_SUMMARY.md`

Example
-------
For branch `feat/singlebinary` create `docs/features/feat-singlebinary/` with the files above. Use the template in `docs/features/TEMPLATE.md` to populate the files.

Maintainers should follow this guide for all future feature branches.
