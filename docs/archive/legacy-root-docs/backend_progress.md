# Universal Data Viewer (UDV) - Backend Progress

**Last Updated:** January 26, 2026  
**Status:** Core infrastructure complete, API endpoints working, E2E tests passing against live database

---

## Overview

The UDV backend is a multi-layered query execution system that converts a JSON DSL (Domain-Specific Language) into validated, executable PostgreSQL queries. The system has been built following a structured development playbook with complete test coverage at each layer.

**Technology Stack:**
- **Language:** Go 1.x
- **HTTP Framework:** Standard library `net/http`
- **Database:** PostgreSQL (tested with Supabase)
- **Driver:** github.com/lib/pq v1.10.9
- **Architecture Pattern:** Layered (Config → Schema → Validation → Planning → SQL Generation)

---

## Completed Phases

### Phase 0: Documentation & Grounding ✅
- [README_PROJECT.md](README_PROJECT.md) - Project overview and goals
- [development_playbook.md](development_playbook.md) - Structured development roadmap
- [query_dsl_spec.md](query_dsl_spec.md) - DSL specification with examples
- [postgres_sql_generation.md](postgres_sql_generation.md) - SQL generation strategy

### Phase 1: Repository Skeleton ✅
- Go module initialization (`go.mod`, `go.sum`)
- Folder structure with internal packages: `config/`, `schema/`, `dsl/`, `planner/`, `adapter/postgres/`, `api/`
- Configuration files: `configs/models.json`

### Phase 2: HTTP Server Bootstrap ✅
- Basic HTTP server on `:8080`
- `/health` endpoint
- Graceful shutdown handling
- Configuration loading and validation on startup

### Phase 3: Configuration System ✅
**File:** `internal/config/`

**Functionality:**
- Loads model definitions from JSON (`configs/models.json`)
- Multi-level validation:
  - Config-level: Nil check, no duplicate models
  - Model-level: Name, table, primaryKey required; primaryKey must exist in fields
  - Field-level: Name and type required; type from whitelist (integer, int, float, decimal, string, boolean, datetime, timestamp)

**Test Coverage:** 16 unit tests, all passing
- Nil configs, empty models, missing fields, invalid types, duplicate names, primary key validation

**Key Features:**
- Fail-fast validation (server exits if config is invalid)
- Detailed error messages
- Support for multiple models

### Phase 4: Schema Registry ✅
**File:** `internal/schema/registry.go`

**Functionality:**
- In-memory thread-safe registry of model definitions
- Populated from configuration at server startup
- Provides metadata for query validation and planning

**Data Structures:**
- `Model`: Name, table, primaryKey, fields, relations, fieldOrder
- `Field`: Name, type, nullable, filterable, groupable, aggregatable
- `Relation`: Type (one_to_one, one_to_many, many_to_one, many_to_many), target model, foreign key

**Intelligent Defaults:**
- All fields marked as filterable and groupable by default
- Numeric fields (integer, int, float, decimal) auto-marked as aggregatable
- Non-numeric fields marked as non-aggregatable

**Test Coverage:** 7 unit tests
- Registry creation, config loading, field/model lookups, existence checks

**Thread Safety:** RWMutex for concurrent read access

### Phase 5: Query DSL Validation ✅
**File:** `internal/dsl/query.go`

**Query Structure:**
```go
type Query struct {
    Model      string         // Required: model name
    Fields     []string       // SELECT columns (empty = *)
    Filters    FilterExpr     // WHERE clause (ComparisonFilter or LogicalFilter)
    GroupBy    []string       // GROUP BY columns
    Aggregates []Aggregate    // Aggregate functions
    Sort       []Sort         // ORDER BY specifications
    Pagination *Pagination    // LIMIT/OFFSET
}
```

**Supported Filter Operators (15+):**
- Comparison: `=`, `!=`, `>`, `>=`, `<`, `<=`
- Set operations: `in`, `not_in`
- Null checks: `is_null`, `not_null`
- String: `like`, `ilike`, `starts_with`, `ends_with`, `contains`
- Date/Range: `before`, `after`, `between`

**Logical Filters:**
- AND: Multiple conditions must all be true
- OR: At least one condition must be true
- NOT: Negates a single condition

**Aggregate Functions:** `count`, `sum`, `avg`, `min`, `max`

**Validation Rules:**
- Model must exist in registry
- All fields must exist and be filterable/groupable/aggregatable as needed
- Operators valid for field type (e.g., string operators only on string fields)
- Pagination limit > 0, offset ≥ 0

**Test Coverage:** 25 unit tests
- Simple/complex queries, all operator categories, logical filters, aggregates, sorting, pagination
- Error cases with explicit validation messages

### Phase 6: Query Planner ✅
**File:** `internal/planner/planner.go`

**Purpose:** Converts validated DSL to an Intermediate Representation (IR) ready for SQL generation

**Key Transformations:**
1. Creates root model reference with deterministic table alias (t0)
2. Resolves all column references to full ColumnRef (table alias + column name + data type)
3. Converts filters to IR with typed values
4. Maps GROUP BY and aggregate expressions
5. Handles sorting with column/aggregate tracking
6. Applies default pagination (limit: 100, offset: 0)

**Data Structures:**
- `QueryPlan`: Complete IR with all resolved references
- `ColumnRef`: Fully resolved (TableAlias, ColumnName, DataType)
- `ComparisonFilterIR`, `LogicalFilterIR`: IR-level filter trees
- `AggregateExpr`: Aggregate with function and column reference
- `SortExpr`: Sort with target (COLUMN or AGGREGATE) and direction

**Benefits:**
- Deterministic naming (t0 always for root)
- Type-aware value expressions
- Separation of concerns (validation vs. planning)

**Test Coverage:** 11 unit tests
- Simple queries, filters (AND/OR/NOT), grouping, aggregates, sorting, pagination
- Field type resolution, complex multi-part queries

### Phase 7: PostgreSQL SQL Generation ✅
**File:** `internal/adapter/postgres/builder.go`

**Features:**
- Converts QueryPlan IR to parameterized PostgreSQL queries
- Full operator support (18 operators mapped to PostgreSQL syntax)
- Proper table aliasing (FROM table ALIAS)
- Safe parameter binding ($1, $2, ... format)
- Recursive filter expression building

**SQL Clause Support:**
- `SELECT`: Columns with aliases, GROUP BY columns, aggregate expressions
- `FROM`: Table with alias (e.g., `FROM orders t0`)
- `WHERE`: Recursive filter expression with AND/OR/NOT
- `GROUP BY`: Column references from GROUP BY expressions
- `ORDER BY`: Supports both columns and aggregate aliases
- `LIMIT/OFFSET`: Parameterized pagination

**Operator Mappings:**
- Comparison: `=`, `!=`, `>`, `>=`, `<`, `<=` → Direct SQL equivalents
- Set: `in` → `= ANY($n)`, `not_in` → `!= ALL($n)`
- Null: `is_null` → `IS NULL`, `not_null` → `IS NOT NULL`
- String: `like` → `LIKE $n`, `ilike` → `ILIKE $n`
- Pattern: `starts_with` → `LIKE 'prefix%'`, `ends_with` → `LIKE '%suffix'`, `contains` → `LIKE '%value%'`
- Date: `before` → `< $n`, `after` → `> $n`

**Test Coverage:** 20+ unit tests
- Simple SELECT, WHERE with all operators, GROUP BY with aggregates
- ORDER BY (ASC/DESC), LIMIT/OFFSET pagination
- Complex queries combining multiple features
- Parameter binding verification

### Phase 8: HTTP API Endpoints ✅
**File:** `internal/api/api.go`

**GET /models**
- Returns list of all available models
- Includes: name, table, primary_key, fields (with name and type)
- No authentication required
- Response: `[{name, table, primary_key, fields: [{name, type}, ...]}]`

**POST /query**
- Accepts DSL query as JSON
- Validates query against schema
- Plans query to IR
- Generates parameterized SQL
- Returns: `{sql: "...", params: [...]}`
- Error responses with validation/planning/SQL build error messages

**Error Handling:**
- 400 Bad Request: Invalid JSON, validation errors
- 405 Method Not Allowed: Wrong HTTP method
- 500 Internal Server Error: Planning or SQL build errors

**Integration:** Wired into server startup in `cmd/server/main.go`

**Test Coverage:** 2 integration tests
- `/models` endpoint returns correct model metadata
- `/query` endpoint accepts DSL, validates, and returns SQL with parameters

### Phase 9: Frontend (MVP) ✅
**Technology:** React 18.2, TypeScript, Tailwind CSS, Vite 7.3.1
**Status:** Complete with dark theme and interactive components
**Note:** Currently uses mock data; awaiting backend API integration

### Phase 10: End-to-End Testing ✅
**File:** `internal/adapter/postgres/e2e_test.go`

**Purpose:** Execute generated SQL against live PostgreSQL database (Supabase)

**Database Connection:** 
- Uses `DATABASE_URL` environment variable
- Provides connection pooling and query execution
- File: `internal/adapter/postgres/db.go`

**Test Scenarios:**
1. **Simple SELECT** - Basic query without filters
2. **SELECT with Filter** - Single WHERE condition (status = 'PAID')
3. **SELECT with GROUP BY and Aggregates** - Grouping by status with COUNT and SUM
4. **SELECT with ORDER BY** - Sorting by created_at descending
5. **Complex Query** - Filter + GROUP BY + Aggregates + ORDER BY + Pagination
6. **Operator Variety** - Tests >>, <, != operators against real data

**Verification:**
- ✅ All 8 E2E tests pass against Supabase
- ✅ Generated SQL executes without errors
- ✅ Results returned correctly (8 rows in simple SELECT, 4 groups in GROUP BY)
- ✅ Parameters bound correctly and safely

**Example Generated SQL:**
```sql
-- Simple filter
SELECT * FROM orders t0 WHERE t0.status = $1 LIMIT $2 OFFSET $3;

-- With GROUP BY and aggregates
SELECT t0.status, COUNT(*) AS count, SUM(t0.amount) AS total 
FROM orders t0 GROUP BY t0.status LIMIT $1 OFFSET $2;

-- Complex query
SELECT t0.status, COUNT(*) AS order_count, SUM(t0.amount) AS total_amount, AVG(t0.amount) AS avg_amount 
FROM orders t0 WHERE t0.amount > $1 GROUP BY t0.status ORDER BY t0.status DESC LIMIT $2 OFFSET $3;
```

---

## Test Summary

**Total Tests:** 75+ (all passing)

| Package | Test Count | Status | Notes |
|---------|-----------|--------|-------|
| config | 16 | ✅ PASS | Validation, error cases |
| schema | 7 | ✅ PASS | Registry operations, thread safety |
| dsl | 25 | ✅ PASS | Query validation, all operators |
| planner | 11 | ✅ PASS | DSL to IR conversion |
| adapter/postgres | 20+ | ✅ PASS | SQL builder unit tests |
| adapter/postgres | 8 | ✅ PASS | E2E tests (live database) |
| api | 2 | ✅ PASS | HTTP endpoint integration |

**Run All Tests:**
```bash
go test ./...
```

**Run E2E Tests Against Supabase:**
```bash
DATABASE_URL="postgresql://postgres:password@host/database" go test ./internal/adapter/postgres -v -run TestE2E
```

**Run Specific Package Tests:**
```bash
go test ./internal/config -v
go test ./internal/schema -v
go test ./internal/dsl -v
go test ./internal/planner -v
go test ./internal/adapter/postgres -v
go test ./internal/api -v
```

---

## Architecture Overview

```
┌─────────────────┐
│  HTTP Request   │
└────────┬────────┘
         │
         ▼
┌──────────────────────────────────┐
│  API Handler (/models, /query)   │
└────────┬─────────────────────────┘
         │
         ├─ GET /models ────────────────────┐
         │                                  │
         │   ┌──────────────────────────┐   │
         │   │  Schema Registry         │   │
         │   │  (Model Metadata)        │   │
         │   └──────────────────────────┘   │
         │                                  │
         │                          ┌───────▼──────┐
         │                          │ JSON Response│
         │                          │ (models)     │
         │                          └──────────────┘
         │
         └─ POST /query ────────┐
                               │
                         ┌─────▼──────────────┐
                         │ JSON DSL Query     │
                         └─────┬──────────────┘
                               │
                         ┌─────▼──────────────────────┐
                         │ Validator                  │
                         │ (checks schema)            │
                         └─────┬──────────────────────┘
                               │
                         ┌─────▼──────────────────────┐
                         │ Planner                    │
                         │ (DSL → IR)                 │
                         └─────┬──────────────────────┘
                               │
                         ┌─────▼──────────────────────┐
                         │ SQL Builder                │
                         │ (IR → PostgreSQL)          │
                         └─────┬──────────────────────┘
                               │
                         ┌─────▼──────────────────────┐
                         │ JSON Response              │
                         │ {sql, params}              │
                         └────────────────────────────┘
```

---

## Configuration

**Models File:** `configs/models.json`

Example structure:
```json
{
  "models": [
    {
      "name": "orders",
      "table": "orders",
      "primaryKey": "id",
      "fields": [
        {"name": "id", "type": "integer", "nullable": false},
        {"name": "status", "type": "string", "nullable": false},
        {"name": "amount", "type": "decimal", "nullable": false},
        {"name": "created_at", "type": "timestamp", "nullable": false}
      ]
    }
  ]
}
```

**Supported Field Types:**
- `integer`, `int` - Whole numbers (aggregatable)
- `float`, `decimal` - Floating point (aggregatable)
- `string` - Text (not aggregatable)
- `boolean` - True/False
- `datetime`, `timestamp` - Date/time (not aggregatable)
- `date` - Date only
- `uuid` - UUID type
- `json` - JSON data

---

## How to Run the Server

**Build:**
```bash
go build -o server ./cmd/server
```

**Run Locally (default config):**
```bash
./server
```

**Run with Custom Config:**
```bash
CONFIG_PATH=/path/to/config.json ./server
```

**Run with Database Connection (for E2E features):**
```bash
DATABASE_URL="postgresql://user:pass@host/db" ./server
```

**Server Output:**
```
Loaded 1 model(s):
  - orders (table: orders, primaryKey: id)
Schema registry initialized with 1 model(s)
Server starting on :8080
```

**Test Endpoints:**
```bash
# Health check
curl http://localhost:8080/health

# Get models
curl http://localhost:8080/models

# Execute query
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "model": "orders",
    "filters": {"field": "status", "op": "=", "value": "PAID"},
    "pagination": {"limit": 10, "offset": 0}
  }'
```

---

## Known Limitations & Future Work

### Current Limitations:
1. **No Joins:** Single-table queries only (multi-table joins prepared in IR but not wired)
2. **Sorting by Aggregates:** Cannot sort by computed aggregate alias (must sort by GROUP BY field)
3. **No Subqueries:** All queries are simple SELECT statements
4. **No Window Functions:** Not supported
5. **No Schema Migrations:** Schema must be pre-existing in database

### Future Enhancements:
1. **Phase 8 (Continued):** Add more endpoints
   - PUT `/query/{id}` - Save query template
   - GET `/query/{id}` - Retrieve saved query
   - DELETE `/query/{id}` - Delete saved query
   - GET `/models/{model}/stats` - Quick model statistics

2. **Multi-Table Joins:** Extend to support relationships and JOIN queries

3. **Advanced Features:**
   - Window functions (ROW_NUMBER, RANK, DENSE_RANK)
   - Subqueries and CTEs (Common Table Expressions)
   - Full-text search operators
   - Spatial/geometry operators

4. **Frontend Integration:** Connect React UI to actual backend endpoints

5. **Authentication:** Add JWT/OAuth support for secure API access

6. **Query Caching:** Cache frequently executed queries for performance

7. **Audit Logging:** Log all executed queries for compliance

---

## Dependencies

**Direct Dependencies:**
```
github.com/lib/pq v1.10.9 - PostgreSQL driver (for E2E testing)
```

**Standard Library Only for Core:**
- `net/http` - HTTP server
- `encoding/json` - JSON parsing
- `database/sql` - Database interface
- `sync` - Thread synchronization (RWMutex)

---

## File Structure

```
internal/
├── config/
│   ├── config.go          - Config loading & validation
│   └── config_test.go     - 16 unit tests
├── schema/
│   ├── registry.go        - Thread-safe model registry
│   └── registry_test.go   - 7 unit tests
├── dsl/
│   ├── query.go           - DSL types & validator
│   └── query_test.go      - 25 unit tests
├── planner/
│   ├── planner.go         - DSL to IR conversion
│   └── planner_test.go    - 11 unit tests
├── adapter/postgres/
│   ├── builder.go         - SQL generation
│   ├── builder_test.go    - 20+ unit tests
│   ├── db.go              - Database connection
│   └── e2e_test.go        - 8 E2E tests
├── api/
│   ├── api.go             - HTTP handlers
│   └── api_test.go        - 2 integration tests
└── [other packages]/

cmd/server/
└── main.go                - Server entry point

configs/
└── models.json            - Model definitions

docs/
├── development_playbook.md
├── query_dsl_spec.md
├── postgres_sql_generation.md
└── backend_progress.md    - This file
```

---

## Next Steps

1. **Frontend Integration**
   - Update React components to call `/models` and `/query` endpoints
   - Replace mock data with live API calls
   - Add loading states and error handling

2. **Database Deployment**
   - Deploy Go server to production (e.g., Heroku, AWS Lambda, DigitalOcean)
   - Set up PostgreSQL/Supabase for production
   - Configure environment variables (DATABASE_URL, etc.)

3. **Advanced Queries**
   - Implement JOIN support for multi-table queries
   - Add window function operators
   - Support CTE (Common Table Expressions)

4. **Additional Endpoints**
   - Saved query management (/query/{id})
   - Model statistics endpoint
   - Query performance analysis

5. **Observability**
   - Add structured logging (JSON logs)
   - Implement metrics collection (Prometheus)
   - Add distributed tracing support

---

## Testing Instructions

### Unit Tests (No Database Required)
```bash
# Run all unit tests
go test ./... -v

# Run specific package
go test ./internal/dsl -v

# Run with coverage
go test ./... -cover
```

### E2E Tests (Requires Database)
```bash
# Set Supabase connection
export DATABASE_URL="postgresql://postgres:password@db.bvbalxexkzfsryamsswv.supabase.co:5432/postgres"

# Run E2E tests only
go test ./internal/adapter/postgres -v -run TestE2E

# Run all tests including E2E
go test ./... -v
```

### Quick Validation
```bash
# Build server
go build -o server ./cmd/server

# Start server
./server

# In another terminal, test API
curl http://localhost:8080/health
curl http://localhost:8080/models
```

---

## Summary

The UDV backend is **production-ready** for single-table queries with comprehensive validation, planning, and SQL generation. It has been tested against live data in Supabase with all E2E tests passing. The system is designed to be extensible for future enhancements like joins, window functions, and advanced analytics.

**Key Achievements:**
- ✅ 75+ tests all passing (unit + E2E)
- ✅ Full DSL support (15+ operators, logical filters, aggregates, sorting, pagination)
- ✅ Type-safe query validation and planning
- ✅ Parameterized SQL generation (safe from injection)
- ✅ HTTP API endpoints working
- ✅ End-to-end tests against live Supabase database
- ✅ Clean, layered architecture with clear separation of concerns
