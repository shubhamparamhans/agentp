````markdown
# MongoDB Support Implementation Analysis

## Overview

This document analyzes what's required to add MongoDB support to Agent P, which currently only supports PostgreSQL.

---

## Current Architecture

### Database Adapter Structure

```
internal/adapter/
├── adapter.go          # (if exists - interface definition)
├── postgres/
│   ├── db.go          # Database connection and execution
│   ├── builder.go     # SQL query builder
│   └── builder_test.go
```

### Current Dependencies

1. **PostgreSQL Adapter** (`internal/adapter/postgres/`)
	- Uses `database/sql` with `github.com/lib/pq` driver
	- Generates SQL queries
	- Executes SQL and returns results as `[]map[string]interface{}`

2. **Query Builder** (`internal/adapter/postgres/builder.go`)
	- Converts QueryPlan IR to SQL
	- Handles SELECT, INSERT, UPDATE, DELETE
	- PostgreSQL-specific syntax (type casting, etc.)

3. **API Integration** (`internal/api/api.go`)
	- Directly uses `postgres.Database` and `postgres.QueryBuilder`

---

## Key Differences: PostgreSQL vs MongoDB

### 1. Query Language

| Aspect | PostgreSQL | MongoDB |
|--------|------------|---------|
| **Query Language** | SQL | MongoDB Query Language (MQL) |
| **Structure** | Tables, Rows, Columns | Collections, Documents, Fields |
| **Schema** | Fixed schema (with flexibility) | Schema-less (flexible) |
| **Joins** | SQL JOINs | $lookup aggregation |
| **Aggregations** | GROUP BY, aggregate functions | Aggregation pipeline |
| **Filtering** | WHERE clause | $match stage |
| **Sorting** | ORDER BY | $sort stage |
| **Pagination** | LIMIT/OFFSET | limit()/skip() |

---

... (full analysis content copied from `docs/MONGODB_SUPPORT_ANALYSIS.md`)

````
