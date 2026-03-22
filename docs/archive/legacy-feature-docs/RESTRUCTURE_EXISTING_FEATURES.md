# Restructure Existing Markdown Files into Feature Folders

Purpose
-------
This document describes the plan to convert all existing Markdown files in the repository into per-feature documentation folders under `docs/features/` following the `docs/features/README.md` guidelines.

High-level approach
-------------------
1. For every existing Markdown file in the repository, create a feature folder under `docs/features/` using a cleaned, hyphen-lowercase name derived from the filename or path.
2. Move the original Markdown content into `docs/features/<feature-folder>/detailed.md` to preserve the full content.
3. Create a `summary.md` in each feature folder containing a one-line summary (placeholder if the summary is not known yet).
4. Add an initial `changelog.md` entry referencing the original file and date.
5. After reviewing all folders, update `docs/features/summary/ALL_FEATURES_SUMMARY.md` with one-line entries and links.

Naming rules
------------
- Conversion: `path/to/FileName.md` -> `docs/features/<clean-name>/detailed.md`
- Clean-name: lowercase, spaces and slashes replaced with `-`, strip file extension, remove multiple hyphens.

Proposed mapping (source -> target folder)
-----------------------------------------
Each entry below maps an existing markdown file to a suggested feature folder (the folder name is the cleaned name) and the target `detailed.md` path.

| Source | Feature folder | Target file |
|---|---:|---|
| ARCHITECTURE_CORRECTION_SUMMARY.md | architecture-correction-summary | docs/features/architecture-correction-summary/detailed.md |
| DATA_MODELLING_PROCESSOR_COMPLETE.md | data-modelling-processor-complete | docs/features/data-modelling-processor-complete/detailed.md |
| DATA_MODELLING_PROCESSOR_INDEX.md | data-modelling-processor-index | docs/features/data-modelling-processor-index/detailed.md |
| DATA_MODELLING_PROCESSOR_VERIFICATION.md | data-modelling-processor-verification | docs/features/data-modelling-processor-verification/detailed.md |
| generate-models (cmd) — (binary source) | generate-models-cmd | docs/features/generate-models-cmd/detailed.md |
| IMPLEMENTATION_CHECKLIST.md | implementation-checklist | docs/features/implementation-checklist/detailed.md |
| MONGODB_IMPLEMENTATION_VERIFICATION.md | mongodb-implementation-verification | docs/features/mongodb-implementation-verification/detailed.md |
| MONGODB_MODELLING_UPDATE.md | mongodb-modelling-update | docs/features/mongodb-modelling-update/detailed.md |
| MONGODB_TEST_COVERAGE_SUMMARY.md | mongodb-test-coverage-summary | docs/features/mongodb-test-coverage-summary/detailed.md |
| MONGODB_TESTING_COMPLETE.md | mongodb-testing-complete | docs/features/mongodb-testing-complete/detailed.md |
| OBJECT_RENDERING_IMPLEMENTATION.md | object-rendering-implementation | docs/features/object-rendering-implementation/detailed.md |
| RELEASE_NOTES_v2.0.0.md | release-notes-v2-0-0 | docs/features/release-notes-v2-0-0/detailed.md |
| SERVER_SIDE_FILTERING_COMPLETE.md | server-side-filtering-complete | docs/features/server-side-filtering-complete/detailed.md |
| TYPE_CASTING_FIX.md | type-casting-fix | docs/features/type-casting-fix/detailed.md |

Documentation in `docs/` (these will each become features)
---------------------------------------------------------
| docs/query_dsl_spec.md | query-dsl-spec | docs/features/query-dsl-spec/detailed.md |
| docs/mvp_scope.md | mvp-scope | docs/features/mvp-scope/detailed.md |
| docs/QUICK_START.md | quick-start | docs/features/quick-start/detailed.md |
| docs/backend_progress.md | backend-progress | docs/features/backend-progress/detailed.md |
| docs/technical.md | technical | docs/features/technical/detailed.md |
| docs/frontend.md | frontend | docs/features/frontend-docs/detailed.md |
| docs/MONGODB_SUPPORT_ANALYSIS.md | mongodb-support-analysis | docs/features/mongodb-support-analysis/detailed.md |
| docs/README_PROJECT.md | readme-project | docs/features/readme-project/detailed.md |
| docs/DATA_MODELLING_PROCESSOR.md | data-modelling-processor | docs/features/data-modelling-processor/detailed.md |
| docs/development_playbook.md | development-playbook | docs/features/development-playbook/detailed.md |
| docs/postgres_adapter_skeleton.md | postgres-adapter-skeleton | docs/features/postgres-adapter-skeleton/detailed.md |
| docs/backend.md | backend | docs/features/backend/detailed.md |
| docs/MONGODB_SCHEMA_DISCOVERY.md | mongodb-schema-discovery | docs/features/mongodb-schema-discovery/detailed.md |
| docs/readme.md | docs-readme | docs/features/docs-readme/detailed.md |
| docs/configurations.md | configurations | docs/features/configurations/detailed.md |
| docs/frontend_progress.md | frontend-progress | docs/features/frontend-progress/detailed.md |
| docs/repo_strucutre.md | repo-structure | docs/features/repo-structure/detailed.md |
| docs/DELTA_IMPLEMENTATION_PLAN.md | delta-implementation-plan | docs/features/delta-implementation-plan/detailed.md |
| docs/DATA_MODELLING_PROCESSOR_QUICKSTART.md | data-modelling-processor-quickstart | docs/features/data-modelling-processor-quickstart/detailed.md |
| docs/GLOBAL_SEARCH_IMPLEMENTATION.md | global-search-implementation | docs/features/global-search-implementation/detailed.md |
| docs/CRUD_VIA_QUERY_ENDPOINT.md | crud-via-query-endpoint | docs/features/crud-via-query-endpoint/detailed.md |
| docs/FRONTEND_INTEGRATION.md | frontend-integration | docs/features/frontend-integration/detailed.md |
| docs/CRUD_IMPLEMENTATION_EFFORT.md | crud-implementation-effort | docs/features/crud-implementation-effort/detailed.md |
| docs/postgres_sql_generation.md | postgres-sql-generation | docs/features/postgres-sql-generation/detailed.md |
| docs/MONGODB_MODELLING.md | mongodb-modelling | docs/features/mongodb-modelling/detailed.md |
| docs/DARK_THEME_IMPLEMENTATION.md | dark-theme-implementation | docs/features/dark-theme-implementation/detailed.md |
| docs/MONGODB_IMPLEMENTATION_PLAN.md | mongodb-implementation-plan | docs/features/mongodb-implementation-plan/detailed.md |
| docs/INTEGRATION_COMPLETE.md | integration-complete | docs/features/integration-complete/detailed.md |
| docs/development_plan.md | development-plan | docs/features/development-plan/detailed.md |
| docs/query_planner.md | query-planner | docs/features/query-planner/detailed.md |
| docs/ODOO_VIEW_DELTA_ANALYSIS.md | odoo-view-delta-analysis | docs/features/odoo-view-delta-analysis/detailed.md |
| docs/end_to_end.md | end-to-end | docs/features/end-to-end/detailed.md |

Other top-level docs
--------------------
| DATA_MODELLING_PROCESSOR_COMPLETE.md | data-modelling-processor-complete | docs/features/data-modelling-processor-complete/detailed.md |
| IMPLEMENTATION_CHECKLIST.md | implementation-checklist | docs/features/implementation-checklist/detailed.md |
| MONGODB_IMPLEMENTATION_VERIFICATION.md | mongodb-implementation-verification | docs/features/mongodb-implementation-verification/detailed.md |
| OBJECT_RENDERING_IMPLEMENTATION.md | object-rendering-implementation | docs/features/object-rendering-implementation/detailed.md |

Notes and assumptions
---------------------
- This plan does not delete or move original files yet — first we will create the feature folders and copy contents into `detailed.md` so that there is no data loss.
- Where a one-line summary is not obvious, `summary.md` will be created with a `TODO: add summary` placeholder.
- Reviewers should confirm mapping names before any automated move.

Proposed execution steps
------------------------
1. Create feature folders under `docs/features/` and copy each source file into `detailed.md` (committed on the current feature branch).
2. Create `summary.md` placeholders and `changelog.md` with a note referencing original path.
3. Update `docs/features/summary/ALL_FEATURES_SUMMARY.md` with links.
4. Open a PR describing the restructure for review. Keep the original files in place during review (or add a `REPLACED_BY_FEATURE_FOLDER` header) to make the review diffs easier.
5. After review and approval, optionally remove original files or keep them as redirects.

Verification
------------
Run `git status` to ensure new files are created and review diffs before committing. Example:

```sh
git add docs/features/<folders> && git status --porcelain
```

Next actions I can take
----------------------
- Create all feature folders and copy existing Markdown into `detailed.md` plus `summary.md` placeholders.
- Create a small script to perform the conversion automatically and run it in a dry-run mode.
- Start a draft PR with the changes for review.

If you approve this plan I will create the feature folders and copy the existing Markdown files into `docs/features/<feature>/detailed.md` and create placeholder `summary.md` files. I will not delete or move the originals until you confirm.
