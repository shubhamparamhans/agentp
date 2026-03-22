

👉 **`CURSOR_PROMPT_PACK.md`**

---

# Universal Data Viewer (UDV)

## Cursor Prompt Pack (Phase-wise)

> **Core Rule for All Prompts**
>
> Always include:
>
> * The **active document**
> * The **scope boundary**
> * The **explicit stop condition**

---

## PHASE 0 — Repo & Docs Grounding

### 🎯 Goal

Make Cursor understand the project **before any code is written**.

### Prompt

```
You are working on a project called Universal Data Viewer (UDV).

Read and internalize all documents under /docs.
These documents are the source of truth.

Rules:
- Do NOT invent new architecture
- Do NOT rename concepts
- Do NOT add features not present in docs
- Prefer clarity over cleverness

Acknowledge once you understand the system.
Do not write any code yet.
```

✅ **Expected Output**
Cursor summarizes architecture in its own words.

---

## PHASE 1 — Repo Skeleton & Boundaries

### 🎯 Goal

Create file/folder structure **without logic**

### Active Doc

`REPO_STRUCTURE.md`

### Prompt

```
Using REPO_STRUCTURE.md as the only reference:

Create the folder and file skeleton for the backend and frontend.
Rules:
- Create empty files only
- No functions
- No imports
- No logic
- Respect package boundaries strictly
- Do not create extra folders

Stop after creating the structure.
```

✅ **Checkpoint**
Tree matches the doc exactly.

---

## PHASE 2 — Backend Bootstrapping

### 🎯 Goal

Get a running Go server (no business logic)

### Active Docs

* `BACKEND_ARCHITECTURE.md`
* `MVP_SCOPE.md`

### Prompt

```
Implement a minimal Go HTTP server in cmd/server/main.go.

Rules:
- Only a /health endpoint
- No config loading
- No database access
- No business logic
- Use net/http only

This is temporary bootstrap code.
Stop once server starts successfully.
```

✅ **Checkpoint**
`curl /health` returns OK.

---

## PHASE 3 — Configuration System

### 🎯 Goal

Load and validate models config

### Active Doc

`CONFIGURATION_SYSTEM.md`

### Prompt

```
Implement config loading and validation in /internal/config.

Scope:
- Parse JSON model config
- Validate: model name, table, primary_key
- Fail fast on error

Rules:
- Do not access database
- Do not resolve relationships yet
- Do not implement UI hints
- No extra features

Add unit tests for invalid configs.
Stop when server fails on invalid config.
```

✅ **Checkpoint**
Bad config = server won’t boot.

---

## PHASE 4 — Schema Registry

### 🎯 Goal

In-memory schema representation

### Active Docs

* `CONFIGURATION_SYSTEM.md`
* `BACKEND_ARCHITECTURE.md`

### Prompt

```
Implement schema registry in /internal/schema.

Scope:
- In-memory representation only
- Model, Field, Relation structs
- Read-only after initialization

Rules:
- No database access
- No query logic
- No DSL parsing
- No validation duplication

Expose lookup helpers only.
Stop after registry is usable by other packages.
```

---

## PHASE 5 — Query DSL (Parsing + Validation)

### 🎯 Goal

Accept and validate query intent

### Active Doc

`QUERY_DSL_SPEC.md`

### Prompt

```
Implement Query DSL structs and validation in /internal/dsl.

Scope:
- DSL struct definitions
- Validation of fields, operators, types
- Logical filters (AND / OR / NOT)

Rules:
- No SQL generation
- No join resolution
- No database access
- Validation errors must be explicit

Add unit tests for:
- valid DSL
- invalid operator
- invalid field
Stop after validation passes tests.
```

---

## PHASE 6 — Query Planner → IR

### 🎯 Goal

Convert DSL → executable plan

### Active Doc

`QUERY_PLANNER_IR.md`

### Prompt

```
Implement Query Planner in /internal/planner.

Scope:
- Convert validated DSL into QueryPlan IR
- Resolve relationships into explicit joins
- Assign deterministic table aliases

Rules:
- Do NOT generate SQL
- Do NOT mutate DSL
- Do NOT access database
- IR must be immutable

Add snapshot tests: DSL → IR.
Stop after IR matches the spec exactly.
```

---

## PHASE 7 — PostgreSQL SQL Generation

### 🎯 Goal

IR → SQL

### Active Docs

* `POSTGRES_ADAPTER_SKELETON.md`
* `POSTGRES_SQL_GENERATION.md`

### Prompt

```
Implement PostgreSQL adapter using the provided skeleton.

Scope (incremental):
1. SELECT + FROM
2. WHERE (single comparison)
3. JOIN
4. GROUP BY
5. AGGREGATES
6. ORDER BY
7. PAGINATION

Rules:
- SQL allowed ONLY in adapter/postgres
- Parameterized queries only
- Follow SQL examples exactly

Add golden tests comparing generated SQL with docs.
Stop after each step.
```

---

## PHASE 8 — API Wiring

### 🎯 Goal

Expose query execution over HTTP

### Active Docs

* `BACKEND_ARCHITECTURE.md`
* `MVP_SCOPE.md`

### Prompt

```
Implement API endpoints in /internal/api.

Scope:
- GET /models
- POST /query

Rules:
- No business logic in handlers
- Call DSL → Planner → Adapter only
- Return structured errors
- No auth for MVP

Stop after API can run one query end-to-end.
```

---

## PHASE 9 — Frontend Bootstrapping

### 🎯 Goal

Render models & data

### Active Docs

* `FRONTEND_ARCHITECTURE.md`
* Frontend skeleton

### Prompt

```
Implement frontend wiring using existing component skeletons.

Scope:
- Fetch models
- Select active model
- Run a hardcoded query
- Render ListView

Rules:
- No styling
- No advanced UX
- No state persistence

Stop after data renders on screen.
```

---

## PHASE 10 — End-to-End Validation

### 🎯 Goal

Prove system works

### Active Doc

`END_TO_END_EXAMPLE.md`

### Prompt

```
Run the exact end-to-end example from END_TO_END_EXAMPLE.md.

Rules:
- Do not modify the example
- SQL must match
- UI output must match conceptually

If something fails, identify which layer broke.
Do not fix more than one layer at a time.
```

---

## 🔒 Golden Cursor Rules (Pin These)

1. **One phase at a time**
2. **One doc at a time**
3. **No cross-package logic**
4. **No feature creep**
5. **If Cursor suggests “improvements”, reject unless documented**

---

## 🧠 Pro Tip: Cursor Memory Hack

At the top of each new coding session, paste:

```
Before writing any code:
- Identify the active architecture document
- State which layer you are working on
- State which layers you must NOT touch
