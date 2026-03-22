
👉 **`MVP_SCOPE.md`**

This document is intentionally opinionated and defensive.

---

```markdown
# Universal Data Viewer (UDV)
## MVP Scope & Explicit Exclusions

---

## 1. Purpose of the MVP

The MVP exists to validate **one core hypothesis**:

> A config-driven, schema-aware data viewer can replace custom admin dashboards for most internal data exploration needs.

The MVP must:
- Be usable end-to-end
- Be safe to deploy on production databases (read-only)
- Demonstrate clear value with minimal surface area

---

## 2. MVP Definition (One Sentence)

> A self-hosted web app that lets developers explore relational data using list and group views with filters and aggregations, without writing SQL or UI code.

---

## 3. In-Scope Features (MVP)

### 3.1 Backend – MUST HAVE

#### Configuration
- JSON-based model configuration
- Field definitions with types
- One-to-many / many-to-one relationships
- Config validation at startup

#### Query Engine
- Query DSL parsing & validation
- Query Planner IR generation
- PostgreSQL adapter
- Parameterized SQL generation

#### Query Features
- Field selection
- Filtering (AND / OR / NOT)
- Grouping (1–2 levels max)
- Aggregations (count, sum, avg)
- Sorting
- Pagination

#### API
- Fetch model metadata
- Execute query DSL
- Return grouped and flat results

---

### 3.2 Frontend – MUST HAVE

#### Core UI
- Model selector (sidebar)
- List view (table)
- Group view (tree-like grouping)
- Filter builder UI
- Aggregate selector
- Sort controls

#### UX Constraints
- Read-only UI
- Keyboard-friendly navigation
- Pagination controls
- Clear error messages

---

### 3.3 Deployment – MUST HAVE

- Single Go binary
- Embedded React build OR static hosting
- Environment-based configuration
- Docker optional (not required)

---

## 4. Explicit Exclusions (MVP WILL NOT INCLUDE)

### 4.1 Data Mutation

❌ Create / Update / Delete  
❌ Inline editing  
❌ Bulk updates  
❌ Forms or validation logic  

**Reason:** Mutation multiplies complexity, permissions, and risk.

---

### 4.2 Permissions & Auth

❌ Role-based access control  
❌ Field-level permissions  
❌ Row-level security  
❌ Multi-tenant isolation  

**Reason:** MVP is developer/internal-first.

---

### 4.3 Advanced Query Features

❌ Subqueries  
❌ HAVING clauses  
❌ Window functions  
❌ User-defined functions  
❌ Arbitrary expressions  

**Reason:** 80% use cases covered without these.

---

### 4.4 UI Complexity

❌ Charts / graphs  
❌ Dashboards  
❌ Custom layouts  
❌ Drag-and-drop builders  

**Reason:** Tables + grouping already deliver insights.

---

### 4.5 Database Support

❌ MySQL  
❌ MongoDB  
❌ Multi-database joins  

**Reason:** PostgreSQL first to validate core engine.

---

### 4.6 Performance Optimizations

❌ Materialized views  
❌ Query caching  
❌ Async execution  
❌ Background jobs  

**Reason:** Premature optimization.

---

### 4.7 Persistence Features

❌ Saved views  
❌ Query history  
❌ Favorites  
❌ Export (CSV / Excel)  

**Reason:** UX polish after validation.

---

## 5. Hard Limits (Enforced in MVP)

| Limit | Value |
|----|----|
| Max rows per page | 100 |
| Default rows per page | 50 |
| Max joins | 3 |
| Max group depth | 2 |
| Max filters | 10 |
| Query timeout | 5 seconds |

---

## 6. Acceptance Criteria (MVP is DONE when…)

### Backend
- Given valid config → backend boots
- Given valid DSL → deterministic SQL generated
- Given invalid DSL → descriptive error returned
- No SQL injection possible

### Frontend
- A new model appears without UI code
- Filters and groups work across relations
- Grouped results show correct aggregates
- UI does not crash on errors

### System
- Can point to an existing Postgres DB
- Can be deployed in < 10 minutes
- Zero DB schema changes required

---

## 7. Non-Goals (Explicitly Deferred)

These are **not missing features**, they are *conscious deferrals*:

- Editing
- Permissions
- Charts
- BI-level analytics
- Multi-DB support
- AI-generated queries

---

## 8. MVP Validation Signals

The MVP is successful if:

- Engineers use it instead of writing SQL
- New models appear instantly via config
- Internal teams ask for edit support (good sign)
- No custom dashboards are requested

---

## 9. Post-MVP Expansion Path (Locked)

MVP success unlocks:
1. Inline editing (CRUD)
2. Permissions & RBAC
3. Saved views
4. Export
5. Workflow automation

---

## 10. Final Principle

> **If a feature does not directly improve data exploration, it is out of scope for MVP.**

This document is a guardrail. Do not negotiate with it lightly.
