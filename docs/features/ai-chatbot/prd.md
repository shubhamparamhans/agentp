# PRD

## Objective
Enable users to ask questions about configured data models in plain text and get reliable answers backed by the existing UDV query engine, without requiring them to manually build filters or understand the underlying schema.

## Users
- Data analysts exploring unfamiliar models
- Internal operators who need quick answers from connected datasets
- Engineers validating model coverage and query behavior through a simpler interface

## Problem Statement
The current product is powerful for structured exploration, but it still assumes that users can choose models, fields, filters, grouping, and search patterns themselves. That creates friction for first-time users and slows down common questions like "show the latest failed orders" or "how many active customers signed up this month?".

## Success Criteria
- A user can submit a plain-text question from the frontend and receive a result without using the manual filter builder
- The backend converts supported prompts into valid UDV DSL with schema-aware validation before execution
- Unsupported or ambiguous prompts fail safely with a clear explanation instead of running a risky query
- Chat results render in under 5 seconds for normal requests excluding external model latency outliers
- Prompt, generated DSL, execution outcome, and errors are logged for debugging and iteration

## Scope

### In Scope
- Read-only natural-language querying
- Schema-aware prompt construction using model and field metadata
- New backend endpoint for chat query interpretation and execution
- Frontend chat experience embedded in the existing data viewer surface
- Structured response payload containing answer text, generated DSL, SQL summary when available, and returned rows
- Guardrails for row limits, model allow-listing, and unsupported operation handling

### Out Of Scope
- Insert, update, or delete operations from chat
- Voice input, file attachments, or image understanding
- Cross-session chat history persistence
- User-specific permission models beyond existing backend access controls

## Risks
- Model hallucinations may reference non-existent fields or unsupported aggregations
- Broad prompts may produce expensive or low-signal queries without good default limits
- Users may trust generated answers without reviewing source rows or query rationale
- Provider-specific AI dependencies can increase latency, cost, and operational complexity

## Dependencies
- Existing schema registry and query DSL remain the execution contract
- A configurable LLM provider integration for prompt interpretation
- Frontend components for chat input, streaming or pending states, and result rendering
- Basic observability for prompt, validation, execution, and error tracing
