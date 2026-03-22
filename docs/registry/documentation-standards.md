# Documentation Standards

These rules keep the repository usable as both a PM dashboard and a technical manager's operating view.

## Source Of Truth Rules

- Keep exactly one active source of truth for each topic.
- Put cross-cutting product context in `docs/product/`.
- Put cross-cutting engineering and architecture context in `docs/architecture/` or `docs/engineering/`.
- Put feature-specific planning, status, and design in `docs/features/<feature-slug>/`.
- Move duplicate or historical material into `docs/archive/`.

## Feature Folder Standard

Every new feature folder should include:

- `index.md` for the feature snapshot
- `prd.md` for problem, scope, and acceptance criteria
- `design.md` for technical approach and code areas
- `tasks.md` for implementation checklist
- `status.md` for progress, blockers, decisions, and next steps
- `assets/` for preserved supporting material such as diagrams, screenshots, payloads, prompt context, imported notes, or legacy reference docs

## Minimum Snapshot Requirements

Every feature snapshot should answer:

- What problem does this feature solve?
- What is in scope and out of scope?
- What is its current status?
- Which code paths does it touch?
- What already exists today?
- What should happen next?

## Lifecycle Rules

- Draft new work in `prd.md` before implementation starts.
- Update `tasks.md` and `status.md` during implementation, not only after completion.
- Update the registry pages when a feature is added, completed, paused, or archived.
- If a feature is superseded, add a note in its folder and move older material to `docs/archive/`.

## Assets Rules

- Keep `assets/` in every feature folder, even if it starts nearly empty.
- Use `assets/` for preservation, not for the primary narrative.
- If an asset becomes important for current understanding, summarize or promote it into `index.md`, `prd.md`, `design.md`, or `status.md`.
- Add a short `assets/README.md` when the folder contains multiple supporting files so future readers know what is worth opening.

## Naming Rules

- Use lowercase kebab-case for folder names.
- Prefer feature names over branch names if the branch name is temporary or noisy.
- Reserve top-level `docs/*.md` files for long-lived entry points only.

## Migration Rule

Existing feature folders that still use `summary.md`, `detailed.md`, and `changelog.md` are allowed to remain in place. Normalize them gradually by either:

- adding the new files next to the legacy files, or
- folding the legacy content into the new file set during the next significant update.
