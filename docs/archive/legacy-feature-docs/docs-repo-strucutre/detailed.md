# Source: docs/repo_strucutre.md

👉 **`REPO_STRUCTURE.md`**

---

```markdown
# Universal Data Viewer (UDV)
## Repository Structure & Package Boundaries

---

## 1. Purpose of This Document

This document defines:
- The **repository layout**
- Clear **package boundaries**
- Ownership and responsibilities of each module
- Rules to prevent architectural erosion

The goal is to make UDV:
- Easy to onboard to
- Safe to refactor
- Friendly to AI-assisted coding
- Scalable for future features

---

## 2. High-Level Repository Layout

```

udv/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── schema/
│   ├── dsl/
│   ├── planner/
│   ├── ir/
│   ├── query/
│   ├── adapter/
│   │   └── postgres/
│   ├── api/
│   ├── limits/
│   └── common/
│
├── frontend/
│   ├── src/
│   ├── public/
│   └── package.json
│
├── configs/
│   └── models.json
│
├── scripts/
├── docs/
├── tests/
└── README.md

```

---

## 3. Top-Level Directories

### 3.1 `cmd/`
**Purpose:** Application entry points

Rules:
- Only wiring code
- No business logic
- No SQL
- No schema knowledge

```

cmd/server/main.go

```

Responsibilities:
- Load configs
- Initialize services
- Start HTTP server

---

### 3.2 `internal/`
**Purpose:** All core backend logic (not importable externally)

This is where **all real work happens**.

---

## 4. Backend Package Breakdown

---

### 4.1 `internal/config/`

**Responsibility**
- Load JSON/YAML configs
- Validate structure
- Fail fast on errors

**Must NOT**
- Know about databases
- Know about SQL
- Perform query logic

Example files:
```

config/
├── loader.go
├── validator.go
└── types.go

```

---

### 4.2 `internal/schema/`

**Responsibility**
- In-memory schema registry
- Relationship graph
- Field metadata access

**Must NOT**
- Load configs directly
- Execute queries
- Generate SQL

Example files:
```

schema/
├── registry.go
├── model.go
├── field.go
└── relation.go

```

---

### 4.3 `internal/dsl/`

**Responsibility**
- DSL structs
- DSL parsing
- DSL validation (syntax + schema-level)

**Must NOT**
- Resolve joins
- Generate SQL
- Access database

Example files:
```

dsl/
├── query.go
├── filter.go
├── aggregate.go
├── validate.go

```

---

### 4.4 `internal/planner/`

**Responsibility**
- Convert DSL → Query Planner IR
- Resolve relationships
- Enforce limits

**Must NOT**
- Generate SQL
- Talk to database
- Modify schema

Example files:
```

planner/
├── planner.go
├── joins.go
├── filters.go
├── groups.go

```

---

### 4.5 `internal/ir/`

**Responsibility**
- Define Query Planner IR structs
- Pure data definitions

**Must NOT**
- Contain logic
- Contain validation
- Know about DSL

Example files:
```

ir/
├── plan.go
├── select.go
├── join.go
├── filter.go

```

> ⚠️ This package should be extremely stable.

---

### 4.6 `internal/query/`

**Responsibility**
- Adapter-agnostic query building
- Translate IR → abstract query model

**Must NOT**
- Contain SQL strings
- Know database-specific syntax

Example files:
```

query/
├── builder.go
├── select.go
├── where.go

```

---

### 4.7 `internal/adapter/`

**Responsibility**
- Database-specific query generation & execution

Structure:
```

adapter/
├── adapter.go        // interface
└── postgres/
├── builder.go
├── executor.go
└── mapper.go

```

**Rules**
- One folder per database
- SQL allowed ONLY here
- Must accept IR as input
- Must return generic row format

---

### 4.8 `internal/api/`

**Responsibility**
- HTTP handlers
- Request/response mapping
- Error translation

**Must NOT**
- Contain business logic
- Generate SQL
- Understand joins

Example files:
```

api/
├── server.go
├── models.go
├── query.go
└── errors.go

```

---

### 4.9 `internal/limits/`

**Responsibility**
- Centralized enforcement of system limits

Why separate?
- Easy auditing
- One place to change safety rules

Example:
```

limits/
├── limits.go
└── enforce.go

```

---

### 4.10 `internal/common/`

**Responsibility**
- Shared utilities
- Logging
- Error types

**Must NOT**
- Contain domain logic

---

## 5. Frontend Structure

```

frontend/
├── src/
│   ├── api/
│   ├── components/
│   │   ├── ModelExplorer/
│   │   ├── ListView/
│   │   ├── GroupView/
│   │   └── FilterBuilder/
│   ├── state/
│   ├── types/
│   └── App.tsx
└── public/

```

### Frontend Rules
- No SQL assumptions
- No schema inference
- Everything driven by API responses

---

## 6. Configs & Docs

### `configs/`
- Runtime configuration
- Not versioned with secrets

### `docs/`
- All architecture documents you created
- Treated as first-class artifacts

---

## 7. Test Strategy by Folder

| Folder | Test Type |
|-----|----------|
| dsl | Validation tests |
| planner | DSL → IR snapshot tests |
| adapter/postgres | SQL golden tests |
| api | HTTP contract tests |

---

## 8. Dependency Rules (VERY IMPORTANT)

Allowed direction ONLY:

```

config → schema → dsl → planner → ir → query → adapter → api

```

❌ No backward imports  
❌ No circular dependencies  

Violations = refactor required.

---

## 9. Why This Structure Works

- Forces clean separation
- Enables parallel development
- Makes AI-generated code safer
- Prevents framework creep

---

## 10. Final Principle

> **If a package needs to know more than one concern, it is in the wrong place.**

This repo structure is a guardrail — respect it.
