# Feature: AI Chatbot

## Status
Draft

## Owner
Engineering

## Problem
Users can already explore data with model selection, filters, grouping, and search, but they must understand the UI and query concepts first. We want a plain-text chatbot that lets a user ask for data in natural language and receive the requested records, summaries, or counts without manually constructing the query.

## Scope
- Add a chat-based entry point in the frontend for natural-language data questions
- Translate prompts into validated UDV query DSL on the backend
- Reuse the existing schema registry, validator, planner, and `/query` execution path
- Return both a human-readable answer and structured result data the UI can render
- Add observability, safety limits, and fallback behavior for unsupported prompts

## Out Of Scope
- Autonomous write actions such as create, update, or delete through chat
- Multi-step agent workflows outside the UDV query domain
- Training or fine-tuning custom models
- Long-lived conversational memory across browser sessions

## Code Areas
- `frontend/src/pages/`
- `frontend/src/components/`
- `frontend/src/api/client.ts`
- `internal/api/`
- `internal/dsl/`
- `internal/planner/`
- `internal/schema/`
- `cmd/server/main.go`

## Related Docs
- [PRD](./prd.md)
- [Design](./design.md)
- [Tasks](./tasks.md)
- [Status](./status.md)
- [Frontend Architecture](/Users/shubhamparamhans/Workspace/udv/docs/architecture/frontend.md)
- [Backend Architecture](/Users/shubhamparamhans/Workspace/udv/docs/architecture/backend.md)
- [System Overview](/Users/shubhamparamhans/Workspace/udv/docs/architecture/system-overview.md)

## Current Capability
UDV already exposes model metadata through `/models`, validates structured queries with the query DSL, plans them against the schema registry, and executes them through the existing `/query` API. The frontend already has a metadata-driven data viewer, which makes this feature a natural extension rather than a new data plane.

## Next Steps
- Define the chat request and response contract
- Implement backend prompt-to-DSL translation with strict validation
- Add a chat panel in the data viewer and render answer plus result table together
- Gate rollout behind a feature flag until accuracy and safety are acceptable
