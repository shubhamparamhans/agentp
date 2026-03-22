# Feature Documentation Guide

`docs/features/` is where feature-level PM and engineering context meet.

Use a feature folder when you want one place to track:

- the user problem
- scope and acceptance criteria
- technical design
- implementation checklist
- delivery status
- follow-up work

## Recommended Structure

Create one folder per feature:

`docs/features/<feature-slug>/`

Each new feature should contain:

- `index.md`
- `prd.md`
- `design.md`
- `tasks.md`
- `status.md`
- `assets/`

`assets/` is the preservation area for useful supporting material that is not the primary source of truth.

Recommended contents:

- prompt outputs
- background analysis
- one-off exploration notes
- screenshots
- sample payloads
- imported legacy markdown files

Use [the template folder](/Users/shubhamparamhans/Workspace/udv/docs/features/_template/index.md) to start a new feature.

## Legacy Feature Folders

Some existing feature folders still use:

- `summary.md`
- `detailed.md`
- `changelog.md`

That format remains valid for older work. When you touch those features again, either:

- add the new file set alongside the legacy files, or
- migrate the old content into the new format.

## Folder Naming

- Use lowercase kebab-case.
- Prefer durable feature names such as `mongodb-support` over temporary branch names.
- If work is a continuation of an existing capability, update the existing feature folder instead of making a near-duplicate one.

## Definition Of Done For Docs

Before calling a feature documented, make sure the folder clearly shows:

- current status
- scope
- affected code areas
- important decisions
- next steps
- supporting context preserved in `assets/` when applicable

Also update:

- [feature catalog](/Users/shubhamparamhans/Workspace/udv/docs/registry/feature-catalog.md)
- [features by status](/Users/shubhamparamhans/Workspace/udv/docs/registry/features-by-status.md)
- [capabilities matrix](/Users/shubhamparamhans/Workspace/udv/docs/registry/capabilities-matrix.md)
