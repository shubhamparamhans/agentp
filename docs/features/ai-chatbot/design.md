# Design

## Overview
The AI chatbot should be implemented as a thin natural-language layer on top of the existing UDV architecture, not as a parallel query system. The frontend collects the user's prompt and displays the response. The backend enriches that prompt with schema metadata, asks an LLM to produce a constrained intermediate output, converts that output into the existing UDV DSL, validates it with the current validator and planner, executes it through the existing `/query` path, and returns both a narrative answer and structured result data.

This design keeps the current contract intact:
- The chatbot never generates raw SQL directly
- The validator and planner remain the safety boundary
- The schema registry stays the source of truth for models and fields
- The frontend remains metadata-driven and renders structured results like any other query workflow

## Code Areas
- `frontend/src/pages/DataViewer.tsx`
- `frontend/src/components/`
- `frontend/src/api/client.ts`
- `internal/api/api.go`
- `internal/dsl/query.go`
- `internal/planner/planner.go`
- `internal/schema/registry.go`
- `cmd/server/main.go`

## Data Flow
1. The user opens the data viewer and enters a plain-text question in a new chatbot panel.
2. The frontend sends the prompt, optional selected model context, and client-side chat metadata to a new backend endpoint such as `POST /chat/query`.
3. The backend loads relevant schema context from the registry, including allowed models, field names, field types, and relationship hints.
4. The backend calls the LLM with a constrained prompt that asks for:
   - intended model
   - fields
   - filters
   - grouping and aggregates
   - sort and pagination
   - a short user-facing explanation
5. The backend parses the model output into a typed intermediate structure and converts it into the existing `dsl.Query`.
6. The validator checks model names, fields, operators, limits, and allowed operations.
7. The planner and query builder generate the executable plan and database query.
8. The backend executes the query and formats the response with:
   - answer summary
   - generated DSL
   - SQL and params when available
   - returned rows
   - warnings for ambiguity, truncation, or partial confidence
9. The frontend renders the answer, preview of generated query intent, and the result set in a chat-aware data panel.

## API Or Interface Changes
- Add `POST /chat/query`
- Request shape:

```json
{
  "prompt": "Show the latest 20 failed orders from this week",
  "model_hint": "orders",
  "context": {
    "selected_model": "orders",
    "timezone": "Asia/Kolkata"
  }
}
```

- Response shape:

```json
{
  "answer": "Here are the latest 20 failed orders from this week.",
  "query": {
    "model": "orders",
    "fields": ["id", "status", "created_at"],
    "filters": {
      "and": [
        { "field": "status", "op": "=", "value": "failed" }
      ]
    },
    "sort": [{ "field": "created_at", "direction": "desc" }],
    "pagination": { "limit": 20, "offset": 0 }
  },
  "data": [],
  "warnings": []
}
```

- Frontend changes:
  - Add a chat input and transcript panel to `DataViewer`
  - Show loading, validation failure, and ambiguity states inline
  - Let users inspect the generated query details before trusting the answer
  - Allow one-click promotion of the generated query into the existing structured UI state

## Frontend Implementation
- Add a chat sidebar or top panel to [DataViewer.tsx](/Users/shubhamparamhans/Workspace/udv/frontend/src/pages/DataViewer.tsx) so users can ask questions in the same place they already explore data
- Create focused components such as `ChatPanel`, `ChatMessageList`, `ChatComposer`, and `ChatResultPreview`
- Extend [client.ts](/Users/shubhamparamhans/Workspace/udv/frontend/src/api/client.ts) with a `executeChatQuery()` helper and typed request/response interfaces
- Reuse the existing list or group rendering components where possible so chat results stay visually consistent with manual queries
- Keep chat state local at first: prompt text, response history, pending request state, selected result, and warnings
- Add a feature flag so the chatbot can be enabled only in selected environments

## Backend Implementation
- Extend [api.go](/Users/shubhamparamhans/Workspace/udv/internal/api/api.go) with a dedicated chat handler instead of overloading `/query`
- Introduce a chat service layer responsible for:
  - building schema context from the registry
  - calling the configured LLM provider
  - parsing and validating structured model output
  - converting it into `dsl.Query`
  - composing final answer payloads
- Keep execution read-only by explicitly rejecting non-select operations in the chat path
- Apply hard safety defaults:
  - maximum row limit
  - allowed aggregate set
  - timeout budget
  - optional allow-list of exposed models
- Register provider configuration in [main.go](/Users/shubhamparamhans/Workspace/udv/cmd/server/main.go) through environment variables so the feature can be disabled cleanly when no AI provider is configured

## Architecture Notes
- Frontend architecture fit:
  - The chatbot is an additional intent-entry surface, not a replacement for the existing metadata-driven UI
  - Structured query state should remain the canonical representation after interpretation
- Backend architecture fit:
  - The chatbot adds an interpretation layer before validation and planning
  - The validator, planner, query builder, and adapters remain unchanged as the core execution pipeline
- System architecture fit:
  - Browser -> chat endpoint -> schema-aware interpretation -> validator/planner -> database
  - This preserves the current separation between UX, query intelligence, and database execution

## Alternatives Considered
- Direct LLM-to-SQL generation
  - Rejected because it bypasses the existing query DSL, weakens validation, and makes adapter portability worse
- Frontend-only prompt interpretation
  - Rejected because schema context, provider credentials, and safety enforcement belong on the backend
- Reusing `/query` with a special `prompt` field
  - Rejected because chat interpretation has different concerns, telemetry, and failure modes from structured queries

## Rollout Notes
- Ship behind a feature flag
- Start with read-only select queries and explicit row limits
- Log prompt-to-query mismatches and manual review examples to improve prompts
- Add a fallback UI message when AI configuration is missing or disabled
