## 1. Backend Purpose

The backend is the **core intelligence layer** of UDV.

Its responsibilities are:
- Understanding database schemas and relationships
- Translating high-level query intent into executable database queries
- Abstracting database-specific behavior
- Providing a stable, generic API for frontend and automation use cases

The backend is intentionally **read-heavy and stateless** in early versions.

---

## 2. Architectural Principles

### 2.1 Design Principles

- Schema-aware, not table-dumb
- Config-driven, not code-driven
- Database-agnostic core
- No UI or presentation logic
- Deterministic query generation

### 2.2 Constraints

- Must compile into a single Go binary
- Must not require database schema changes
- Must work with existing production databases
- Must avoid ORM magic that hides query behavior

---

## 3. High-Level Backend Architecture


API Layer
|
Query Validation & Planning
|
Schema Registry
|
Query Builder / ORM
|
Database Adapter
|
Database

---

## 4. Core Backend Components

### 4.1 Config Loader

#### Purpose
Load and validate model definitions that describe database structure and relationships.

#### Responsibilities
- Parse JSON/YAML configuration files
- Validate:
  - Field names and types
  - Primary keys
  - Relationship integrity
- Fail-fast on invalid configs

#### Requirements
- Must support hot reload (optional, non-blocking)
- Must normalize configs into internal representations
- Must be deterministic

---

### 4.2 Schema Registry

#### Purpose
Serve as the **single source of truth** for schema metadata during runtime.

#### Responsibilities
- Store model definitions in memory
- Maintain relationship graph
- Provide lookup utilities for:
  - Fields
  - Joins
  - Relationship traversal

#### Requirements
- Read-only after initialization
- Thread-safe access
- Zero database calls

---

### 4.3 Query Specification (Input Contract)

All backend queries must be expressed as **structured query specifications**, never raw SQL.

#### Characteristics
- Declarative
- Serializable (JSON)
- Frontend-agnostic
- Database-agnostic

This spec is the **core contract** between frontend and backend.

---

### 4.4 Query Validator

#### Purpose
Validate incoming query specifications before execution.

#### Responsibilities
- Verify model existence
- Verify field existence and types
- Validate operators against field types
- Enforce safe query limits

#### Requirements
- No database access
- Produce explicit validation errors
- Prevent unsafe queries

---

### 4.5 Query Planner

#### Purpose
Convert validated query specifications into an executable **Query Plan**.

#### Responsibilities
- Resolve joins via relationships
- Determine grouping hierarchy
- Decide aggregate execution order
- Prepare adapter-friendly query plan

#### Requirements
- Deterministic output
- No SQL generation
- No database access

---

### 4.6 ORM / Query Builder

#### Purpose
Translate query plans into database-specific query representations.

#### Responsibilities
- Generate SELECT statements
- Generate JOINs from relationships
- Generate WHERE, GROUP BY, ORDER BY
- Handle pagination safely

#### Requirements
- No string concatenation in business logic
- Parameterized queries only
- Explicit SQL generation per adapter

---

### 4.7 Database Adapter Layer

#### Purpose
Isolate database-specific behavior.

#### Responsibilities
- Translate generic query plans into native queries
- Execute queries
- Map results into generic row format

#### Adapter Interface (Conceptual)
```go
type DBAdapter interface {
  BuildQuery(plan QueryPlan) (string, []any, error)
  Execute(query string, args []any) ([]Row, error)
}
````

#### Initial Support

* PostgreSQL

#### Future Support

* MySQL
* MongoDB (document mode)

---

## 5. API Layer

### 5.1 Design Goals

* Stateless
* JSON-only
* Generic endpoints
* No domain-specific logic

### 5.2 Responsibilities

* Accept query specifications
* Invoke validation, planning, execution
* Return structured responses
* Handle errors consistently

---

## 6. Non-Functional Requirements

### Performance

* Must support pagination on all queries
* Must avoid N+1 queries
* Must support indexed filtering

### Scalability

* Horizontal scaling via stateless APIs
* Database connection pooling

### Security

* Read-only by default
* Authentication middleware ready
* No raw SQL exposure

### Observability

* Structured logs
* Query execution timing
* Error categorization

---

## 7. Explicit Non-Responsibilities

The backend must NOT:

* Render UI
* Store UI state
* Contain business workflows
* Implicitly mutate data (initial phase)

---

## 8. Summary

The backend of UDV is a **schema-aware query engine**, not a traditional CRUD service.

Its success depends on:

* Clear schema understanding
* Deterministic query planning
* Strong separation of concerns