# Tasks

## Planning
- [ ] Confirm initial rollout is read-only and limited to select queries
- [ ] Confirm which models are exposed to the chatbot in v1
- [ ] Confirm acceptance criteria for ambiguity handling and answer formatting

## Implementation
- [ ] Add backend chat endpoint and service layer
- [ ] Add provider configuration and feature flag wiring
- [ ] Translate LLM output into `dsl.Query` and validate it
- [ ] Execute validated chat queries through the existing backend pipeline
- [ ] Add frontend chat panel, transcript state, and result rendering
- [ ] Add inspect-query affordance so users can review generated DSL

## Validation
- [ ] Unit tests for prompt-output parsing and DSL conversion
- [ ] Integration tests for chat endpoint success, rejection, and timeout paths
- [ ] Manual verification with representative prompts across at least two models
- [ ] Verify unsupported prompts fail safely and clearly

## Docs
- [ ] Update feature snapshot as implementation decisions harden
- [ ] Update registry pages when status changes from draft to active
- [ ] Preserve sample prompts and payloads in `assets/`
