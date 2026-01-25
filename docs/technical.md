
# Universal Data Viewer (UDV)
## Technical Architecture & System Design Document

---

## 1. Purpose of This Document

This document describes the **technical architecture, system components, data flow, and design decisions** behind Universal Data Viewer (UDV).

It is intended to:
- Provide full technical context to engineers and AI IDEs
- Act as a single source of truth for implementation
- Serve as a foundation for PRDs, tickets, and code generation

---

## 2. System Overview

UDV is a **config-driven, read-heavy data exploration framework** consisting of:

- A **Golang backend** responsible for schema understanding, query generation, and database access
- A **React frontend** responsible for rendering list/tree views and building query definitions
- A **model configuration layer** that defines data structure and relationships
- A **database adapter layer** that supports multiple databases

The system is designed to be:
- Stateless
- Self-hosted
- Extensible via configuration and adapters

---

## 3. High-Level Architecture

```

Client (Browser)
|
| HTTP / JSON
v
React Frontend
|
| REST API
v
Go Backend
├── Config Loader
├── Schema Registry
├── Query Planner
├── ORM / SQL Generator
├── API Layer
└── DB Adapter Interface
|
v
Database
(Postgres | MySQL | MongoDB)

````

---

## 4. Backend Architecture (Golang)

### 4.1 Backend Responsibilities

The backend is responsible for:

- Loading and validating model configurations
- Understanding schema and relationships
- Translating UI queries into executable database queries
- Executing queries efficiently
- Returning structured results for UI consumption

---

## 4.2 Core Backend Modules

### 4.2.1 Config Loader

**Purpose**
- Parse model definitions from JSON/YAML
- Validate field types and relationships
- Fail fast on invalid schemas

**Responsibilities**
- Load configuration at startup
- Watch for config reloads (optional)
- Produce normalized schema definitions

**Example Model Definition**
```json
{
  "name": "Order",
  "table": "orders",
  "primary_key": "id",
  "fields": [
    { "name": "id", "type": "uuid" },
    { "name": "status", "type": "string", "index": true },
    { "name": "amount", "type": "float" },
    { "name": "created_at", "type": "datetime" }
  ],
  "relations": [
    {
      "type": "many_to_one",
      "model": "User",
      "foreign_key": "user_id",
      "reference_key": "id"
    }
  ]
}
````

---

### 4.2.2 Schema Registry

**Purpose**

* Maintain an in-memory representation of all models and relationships

**Key Characteristics**

* Read-only after initialization
* Relationship graph based
* Accessible across all backend layers

**Internal Representation**

```go
type ModelSchema struct {
    Name        string
    Table       string
    PrimaryKey  string
    Fields      map[string]Field
    Relations   []Relation
}
```

---

### 4.2.3 Query Specification (Core Abstraction)

All queries are represented as structured objects.

```json
{
  "model": "Order",
  "fields": ["id", "status", "amount"],
  "filters": [
    { "field": "status", "op": "=", "value": "PAID" },
    { "field": "amount", "op": ">", "value": 1000 }
  ],
  "group_by": ["status"],
  "aggregates": [
    { "fn": "sum", "field": "amount", "alias": "total_amount" }
  ],
  "sort": [
    { "field": "created_at", "direction": "desc" }
  ],
  "limit": 50,
  "offset": 0
}
```

This structure is shared between frontend and backend.

---

### 4.2.4 Query Planner

**Purpose**

* Convert Query Specification into an executable plan

**Responsibilities**

* Validate fields against schema
* Resolve relationships and joins
* Determine grouping and aggregation strategy
* Optimize query structure (where possible)

**Output**

* QueryPlan (intermediate representation)

---

### 4.2.5 ORM / SQL Generator

**Purpose**

* Translate QueryPlan into database-specific queries

**Design Goals**

* No raw SQL in business logic
* Database-agnostic core
* Adapter-specific query generation

```go
type QueryBuilder interface {
    BuildSelect(plan QueryPlan) (Query, error)
}
```

---

### 4.2.6 Database Adapter Layer

**Purpose**

* Abstract database differences

**Adapter Interface**

```go
type DBAdapter interface {
    BuildQuery(plan QueryPlan) (string, []any, error)
    Execute(query string, args []any) ([]map[string]any, error)
}
```

**Adapters**

* PostgreSQL (initial)
* MySQL (future)
* MongoDB (document mode)

---

### 4.2.7 API Layer

**Design**

* REST-based
* Stateless
* JSON only

**Representative Endpoints**

| Method | Endpoint         | Description       |
| ------ | ---------------- | ----------------- |
| GET    | /models          | List all models   |
| GET    | /models/{name}   | Model metadata    |
| POST   | /query           | Execute query     |
| GET    | /records/{model} | Simple list fetch |

---

## 5. Frontend Architecture (React)

### 5.1 Frontend Responsibilities

* Fetch model metadata
* Render list and tree views
* Build structured queries
* Manage UI state

---

### 5.2 Frontend Core Modules

#### 5.2.1 Model Explorer

* Lists available models
* Displays relationships and fields

#### 5.2.2 List View Component

* Virtualized table rendering
* Sorting and pagination
* Column visibility control
* Inline aggregate display

#### 5.2.3 Tree / Group View

* Hierarchical grouping
* Expand / collapse nodes
* Aggregated values per group

#### 5.2.4 Filter Builder

* Visual rule builder
* AND / OR support
* Type-aware inputs

---

### 5.3 Frontend State Model

```ts
interface QueryState {
  model: string;
  fields: string[];
  filters: FilterRule[];
  groupBy?: string[];
  aggregates?: Aggregate[];
  sort?: SortRule[];
  pagination: {
    limit: number;
    offset: number;
  };
}
```

---

## 6. Configuration System

### 6.1 Configuration Sources

* Model definitions (JSON/YAML)
* Database connection config
* Feature flags (optional)

### 6.2 Design Goals

* No recompilation required
* Declarative
* Human-readable
* Schema-validated

---

## 7. Non-Functional Requirements

### Performance

* Handle large tables via pagination
* Support indexed filtering
* Avoid N+1 queries

### Scalability

* Horizontal backend scaling
* Stateless API layer

### Security

* Read-only mode by default
* Auth middleware support
* No direct DB exposure

### Observability

* Structured logging
* Query execution timing
* Error tracing

---

## 8. Deployment Model

### Backend

* Single statically compiled Go binary
* Environment-based configuration
* Optional Docker container

### Frontend

* Bundled and served by backend
* Or standalone deployment

---

## 9. Future Extensions (Non-Binding)

* Inline editing and CRUD
* Permission and RBAC system
* Saved views
* Export (CSV / JSON)
* Workflow engine
* Webhooks and triggers

---

## 10. Summary

UDV is architected as a **schema-aware, query-driven data exploration platform**.

By separating:

* Model definition
* Query specification
* Query execution
* UI rendering

The system remains flexible, extensible, and easy to reason about for both humans and AI-assisted tooling.

```

