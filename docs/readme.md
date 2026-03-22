# UDV Documentation

This `docs/` directory is organized to answer two questions quickly:

1. What does the product already support?
2. Where is the technical source of truth for a feature or subsystem?

## Start Here

- [Product vision](/Users/shubhamparamhans/Workspace/udv/docs/product/vision.md)
- [Roadmap](/Users/shubhamparamhans/Workspace/udv/docs/product/roadmap.md)
- [MVP scope](/Users/shubhamparamhans/Workspace/udv/docs/product/mvp-scope.md)
- [Open questions](/Users/shubhamparamhans/Workspace/udv/docs/product/open-questions.md)
- [Feature catalog](/Users/shubhamparamhans/Workspace/udv/docs/registry/feature-catalog.md)
- [Features by status](/Users/shubhamparamhans/Workspace/udv/docs/registry/features-by-status.md)
- [Capabilities matrix](/Users/shubhamparamhans/Workspace/udv/docs/registry/capabilities-matrix.md)
- [Documentation standards](/Users/shubhamparamhans/Workspace/udv/docs/registry/documentation-standards.md)

## Active Areas

### Product
- `docs/product/` keeps PM-facing context:
  - vision
  - scope
  - roadmap
  - open questions

### Architecture
- `docs/architecture/` keeps cross-cutting technical reference:
  - system overview
  - backend
  - frontend
  - configuration
  - query design
  - repo structure

### Engineering
- `docs/engineering/` keeps developer workflow docs:
  - quick start
  - development playbook

### Features
- `docs/features/` is the source of truth for feature-level planning and implementation tracking.
- New feature work should start from [the feature guide](/Users/shubhamparamhans/Workspace/udv/docs/features/README.md) and [the template](/Users/shubhamparamhans/Workspace/udv/docs/features/_template/index.md).
- Each feature folder should also keep an `assets/` folder for preserved prompt context, analysis notes, and supporting artifacts.

### Archive
- `docs/archive/` contains legacy notes, migration-era files, old summaries, and one-off implementation writeups that should remain discoverable but are no longer the primary source of truth.

## Working Rule

Use one active source of truth per topic:

- Cross-cutting project guidance goes in `product/`, `architecture/`, or `engineering/`.
- Feature-specific planning and delivery docs go in one feature folder under `docs/features/`.
- Old or duplicate material gets moved to `docs/archive/` instead of staying mixed with active docs.
