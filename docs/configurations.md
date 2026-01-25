
👉 **`CONFIGURATION_SYSTEM.md`**

---

# Universal Data Viewer (UDV)
## Configuration System – Technical Specification

---

## 1. Purpose of the Configuration System

The configuration system defines **how UDV understands data**.

It replaces:
- Hardcoded models
- ORM structs
- UI schemas
- Admin definitions

with **declarative, portable, machine-readable configuration files**.

This configuration is the **single source of truth** for:
- Data models
- Fields and types
- Relationships
- UI hints (minimal, optional)
- Query validation rules

---

## 2. Design Goals

### Primary Goals
- Declarative and explicit
- Database-agnostic
- Human-readable
- AI-friendly
- Backward-compatible

### Constraints
- No runtime code generation
- No database schema mutation
- Must map cleanly to existing DBs
- Must not encode business logic

---

## 3. Configuration Sources

| Config Type | Purpose |
|-----------|--------|
| Model Config | Define tables/collections and relationships |
| System Config | DB, limits, feature flags |
| UI Hints (Optional) | Presentation hints (non-binding) |

All configurations are **read-only at runtime**.

---

## 4. Configuration Format

### Supported Formats
- JSON (mandatory)
- YAML (optional, later)

### Loading Rules
- Loaded at startup
- Validated before server starts
- Failure = server does not boot

---

## 5. Model Configuration

### 5.1 Top-Level Structure

```json
{
  "version": "1.0",
  "models": [ ... ]
}
````

---

## 5.2 Model Definition

Each model represents **one database table or collection**.

```json
{
  "name": "Order",
  "table": "orders",
  "primary_key": "id",
  "fields": [ ... ],
  "relations": [ ... ],
  "options": { ... }
}
```

---

### 5.2.1 Model Attributes

| Field       | Required | Description                |
| ----------- | -------- | -------------------------- |
| name        | ✅        | Logical model name         |
| table       | ✅        | DB table / collection name |
| primary_key | ✅        | Primary key field          |
| fields      | ✅        | List of fields             |
| relations   | ❌        | Relationship definitions   |
| options     | ❌        | Model-level behavior flags |

---

## 6. Field Configuration

### 6.1 Field Definition

```json
{
  "name": "created_at",
  "column": "created_at",
  "type": "datetime",
  "nullable": false,
  "indexed": true,
  "filterable": true,
  "groupable": true,
  "aggregatable": false
}
```

---

### 6.2 Field Attributes

| Attribute    | Purpose               |
| ------------ | --------------------- |
| name         | Logical field name    |
| column       | Actual DB column      |
| type         | Logical data type     |
| nullable     | Validation hint       |
| indexed      | Optimization hint     |
| filterable   | Allowed in filters    |
| groupable    | Allowed in GROUP BY   |
| aggregatable | Allowed in aggregates |

---

### 6.3 Supported Field Types (Initial)

| Type     | Notes             |
| -------- | ----------------- |
| string   | Text / varchar    |
| int      | Integer           |
| float    | Decimal / numeric |
| boolean  | Boolean           |
| datetime | Timestamp         |
| date     | Date              |
| uuid     | UUID              |
| json     | JSON / JSONB      |

> Field types influence **UI controls**, **query validation**, and **SQL generation**.

---

## 7. Relationship Configuration

### 7.1 Relationship Definition

```json
{
  "type": "many_to_one",
  "model": "User",
  "foreign_key": "user_id",
  "reference_key": "id"
}
```

---

### 7.2 Supported Relationship Types

| Type         | Description    |
| ------------ | -------------- |
| one_to_one   | 1 ↔ 1          |
| one_to_many  | 1 → N          |
| many_to_one  | N → 1          |
| many_to_many | N ↔ N (future) |

---

### 7.3 Relationship Rules

* Relationships must reference valid models
* Foreign keys must exist as fields
* Cycles are allowed but must be explicit
* Join direction is always deterministic

---

## 8. Model Options

Optional flags that influence behavior.

```json
"options": {
  "read_only": true,
  "default_sort": {
    "field": "created_at",
    "direction": "desc"
  },
  "page_size": 50
}
```

---

## 9. UI Hint Configuration (Optional)

UI hints are **non-binding** and can be ignored by clients.

```json
{
  "ui": {
    "label": "Created At",
    "hidden": false,
    "width": 160
  }
}
```

### Rules

* UI hints must never affect query semantics
* Backend treats UI hints as opaque metadata

---

## 10. System Configuration

### 10.1 Database Config

```json
{
  "database": {
    "type": "postgres",
    "host": "localhost",
    "port": 5432,
    "name": "app_db",
    "user": "udv",
    "ssl": false
  }
}
```

---

### 10.2 Query Limits

```json
{
  "limits": {
    "max_limit": 500,
    "default_limit": 50,
    "max_joins": 5,
    "max_group_levels": 3
  }
}
```

Purpose:

* Prevent runaway queries
* Protect production databases

---

## 11. Validation Rules

### Config Validation

* Duplicate model names → error
* Missing primary key → error
* Invalid field types → error
* Broken relationships → error

### Runtime Enforcement

* Non-filterable fields rejected
* Non-groupable fields rejected
* Non-aggregatable fields rejected

---

## 12. Config to Runtime Mapping

```
Config Files
   ↓
Config Loader
   ↓
Schema Registry
   ↓
Query Validator
   ↓
Query Planner
```

No other component may interpret configs independently.

---

## 13. Versioning Strategy

* Config version declared at root
* Backward-compatible additions only
* Breaking changes require major version bump

---

## 14. Explicit Non-Goals

Configuration system must NOT:

* Encode business rules
* Store permissions
* Define workflows
* Define UI layouts

---

## 15. Summary

The configuration system is the **foundation of UDV**.

If the config is:

* Correct
* Explicit
* Stable

Then:

* Backend queries remain safe
* Frontend UI remains generic
* The system remains extensible
