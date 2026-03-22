# Status

## Current Phase
Draft requirements

## Progress
- Feature snapshot created
- PRD created
- Technical design outline created
- Delivery task breakdown created

## Blockers
- No confirmed config schema for views yet
- No adapter-specific materialization strategy selected for v1
- No decision yet on whether frontend management ships in the first milestone

## Decisions
- Use a config-only authoring model for v1
- Support both single-model and relationship-based views
- Include transformation support as a first-class part of the definition
- Treat materialized views as read-only datasets for consumers
- Require refresh metadata so consumers can reason about freshness
- Integrate views into the existing UI through a dedicated `Views` menu

## Next Checkpoint
- Review the PRD and lock the config schema, refresh modes, and v1 adapter target before implementation starts
