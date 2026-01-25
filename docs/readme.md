
# Universal Data Viewer (UDV)

## 1. Overview

Universal Data Viewer (UDV) is a lightweight, self-hosted framework designed to provide fast, intuitive, and zero-boilerplate exploration of structured data.  

The project is inspired by the simplicity and power of Odoo’s list and tree views, combined with the flexibility of modern web stacks. UDV focuses on **data understanding first**, not CRUD workflows or dashboards.

The core idea is to allow developers to:
- Define data models and relationships once
- Automatically generate powerful list and tree views
- Perform filtering, grouping, and aggregation via UI
- Avoid writing custom admin or internal tooling repeatedly

---

## 2. Motivation & Thought Process

### Problem Being Solved

Most modern frameworks lack a **generic, relationship-aware data viewer**:

- Django Admin and similar tools require heavy customization
- React-based admin frameworks require explicit UI modeling
- BI tools are heavy and not developer-centric
- Internal dashboards are repeatedly rebuilt across projects

### Key Insight

> Developers want fast insight into their data without writing UI or SQL.

Odoo demonstrates that:
- A well-designed list view can replace many dashboards
- Grouping + filtering covers most exploratory needs
- ORM-level relationship awareness is critical

UDV extracts these ideas into a **framework-agnostic, developer-first tool**.

---

## 3. Product Vision

UDV aims to be:

- **Universal** – works with multiple databases
- **Config-driven** – no UI or model code required initially
- **Lightweight** – single Go binary + React UI
- **Extensible** – editing, workflows, and permissions can be added later

Primary focus in early versions is **read-only data exploration**.

---

## 4. Goals & Non-Goals

### Goals
- Zero-UI-code list and tree views
- Relationship-aware querying
- Visual filtering and grouping
- Fast iteration and easy deployment

### Non-Goals (Initial Phase)
- No workflow automation
- No complex RBAC
- No chart-heavy BI features
- No spreadsheet-style data editing

---

## 5. Target Users

- Backend / full-stack developers
- Startup engineering teams
- Internal tooling teams
- Data-heavy product teams

Secondary users may include product managers or ops teams in read-only mode.

---

## 6. High-Level Architecture

```

Browser
↓
React Frontend

* Model Explorer
* List View
* Tree (Grouped) View
* Filter Builder
  ↓
  Go Backend
* Config Parser
* Schema Registry
* ORM / Query Engine
* API Layer
* DB Adapters
  ↓
  Databases
* PostgreSQL (initial)
* MySQL (later)
* MongoDB (document mode)

```

---

## 7. Core Design Principles

### 7.1 Config-Driven Modeling

Instead of defining models in code, UDV uses declarative configuration files.

These definitions describe:
- Tables / collections
- Fields and types
- Relationships between models

This allows:
- Schema understanding without recompilation
- Automatic UI generation
- Backend-agnostic modeling

---

### 7.2 Relationship Awareness

Relationships are first-class citizens:
- One-to-one
- One-to-many
- Many-to-one

The system understands how to:
- Join data
- Traverse relationships
- Group by related fields

---

### 7.3 Query Abstraction

UDV avoids raw SQL in application logic.

Queries are defined as structured objects:
- Filters
- Grouping
- Aggregations
- Pagination

These are translated by database adapters into optimized queries.

---

## 8. Backend Architecture (Golang)

### 8.1 Key Components

#### Config Parser
- Reads model definitions (JSON/YAML)
- Validates relationships
- Builds internal schema graph

#### Schema Registry
- In-memory representation of all models
- Used by ORM and API layers
- Enables relationship traversal

#### ORM / Query Engine
- Model-aware query builder
- Supports filtering, grouping, aggregation
- Generates database-specific queries

#### Database Adapter Layer
- Abstracts SQL / NoSQL differences
- PostgreSQL first
- Extensible via adapter interface

#### API Layer
- REST-based JSON APIs
- Stateless
- UI and automation friendly

---

## 9. Frontend Architecture (React)

### 9.1 Design Philosophy

- Data-first UI
- Minimal configuration
- Fast rendering for large datasets
- Keyboard and power-user friendly

---

### 9.2 Core UI Concepts

#### Model Explorer
- Lists available models
- Displays metadata and relationships

#### List View
- Tabular representation of records
- Sorting, pagination, column toggling
- Inline aggregations

#### Tree / Group View
- Hierarchical grouping
- Expandable grouped rows
- Aggregated summaries per group

#### Filter Builder
- Visual filter creation
- Generates structured query objects
- Supports AND / OR logic

---

## 10. API Design (Conceptual)

The backend exposes generic APIs such as:

- List models
- Fetch model metadata
- Execute structured queries
- Fetch paginated records

The API is intentionally generic to support:
- UI
- Automation
- Future integrations

---

## 11. Deployment Model

### Backend
- Single statically compiled Go binary
- Configuration via environment variables
- Optional Docker support

### Frontend
- React app bundled and served by backend
- Or deployed independently via CDN

---

## 12. Security Model (Early Stage)

- Read-only access by default
- Token or basic authentication
- Network-level access control

Advanced permission systems are deferred.

---

## 13. Extensibility Roadmap

### Phase 1
- Read-only data exploration
- Filters, grouping, aggregation

### Phase 2
- Inline editing
- CRUD APIs
- Form auto-generation

### Phase 3
- Permissions & RBAC
- Saved views
- Export functionality

### Phase 4
- Workflow engine
- Triggers & automation
- Webhooks and integrations

---

## 14. Summary

Universal Data Viewer is a developer-first framework focused on **making data understandable with minimal effort**.

By combining:
- Config-driven modeling
- Relationship-aware querying
- Auto-generated list and tree views

UDV aims to become a foundational tool for internal data exploration across modern software systems.
```

---

