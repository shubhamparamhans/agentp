# Status

## Current Phase
Draft

## Progress
- Feature folder created
- Product scope captured
- Frontend, backend, and architecture approach proposed

## Blockers
- LLM provider choice and operational constraints are not yet finalized
- The exact v1 model allow-list is still open

## Decisions
- Use the new structured feature-doc format
- Keep the chatbot read-only in v1
- Reuse the existing UDV query DSL and execution pipeline instead of generating SQL directly
- Add a dedicated chat endpoint rather than extending `/query` with natural-language parsing concerns

## Next Checkpoint
- Review the design doc and lock the API contract before implementation starts
