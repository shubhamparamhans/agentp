# MongoDB Implementation Plan

## Overview

Implement MongoDB support for UDV in two phases:
1. **Phase A** -- Schema Discovery: infer `models.json` from document sampling
2. **Phase B** -- Full MongoDB Adapter: query builder, CRUD operations, and database interface abstraction

---

## Current Architecture

```
HTTP Request -> API Handler -> DSL Validator -> Query Planner -> QueryBuilder -> Database
```

**Current state (PostgreSQL only):**
- `internal/api/api.go` directly imports `postgres` package
- `cmd/server/main.go` creates `*postgres.Database`
- `internal/adapter/adapter.go` is an empty stub (just a package comment)
- `internal/schema_processor/processor.go` uses PostgreSQL `information_schema`
- `cmd/generate-models/main.go` only supports PostgreSQL

---

## Phase A: MongoDB Schema Discovery (First Priority)

**Goal:** Generate a valid `models.json` from MongoDB collections by sampling documents and inferring field types.

**Estimated effort:** 34-48 hours

### A1. Add MongoDB Go Driver Dependency

```bash
go get go.mongodb.org/mongo-driver/mongo
```

Adds the official MongoDB Go driver to `go.mod`.

---

### A2. Create MongoDB Sampler

**New file:** `internal/schema_processor/mongodb_sampler.go`

**Struct:** `MongoDBSampler` -- holds `*mongo.Client` and `*mongo.Database`

**Functions:**
| Function | Description |
|----------|-------------|
| `NewMongoDBSampler(client, dbName)` | Constructor |
| `GetAllCollections() ([]string, error)` | Calls `database.ListCollectionNames()` |
| `SampleDocuments(collection, sampleSize) ([]bson.M, error)` | Uses `$sample` aggregation pipeline for random sampling |

---

### A3. Create MongoDB Type Inferrer

**New file:** `internal/schema_processor/mongodb_infer.go`

Reuses existing `FieldType` constants from `processor.go` (`TypeInteger`, `TypeString`, etc.).

**Types:**

```go
type FieldStats struct {
    TypeCounts   map[FieldType]int  // How many times each type appears
    TotalCount   int                // Total documents analyzed
    NullCount    int                // How many times field is missing/null
    SampleValues []interface{}      // Sample values for debugging
}

type CollectionSchema struct {
    CollectionName string
    Fields         map[string]*FieldStats
    DocumentCount  int
}
```

**Functions:**
| Function | Description |
|----------|-------------|
| `InferSchema(documents []bson.M) *CollectionSchema` | Iterates documents, calls `analyzeDocument` |
| `analyzeDocument(doc, prefix, fields)` | Recursively walks fields, flattens nested docs with dot notation |
| `inferTypeFromValue(value) FieldType` | Maps Go/BSON types to `FieldType` |

**Type Mapping (BSON -> FieldType):**
| BSON Type | FieldType |
|-----------|-----------|
| `primitive.ObjectID` | `uuid` |
| `bool` | `boolean` |
| `int32`, `int64` | `integer` |
| `float64` | `decimal` |
| `string` | `string` (with UUID/date detection) |
| `primitive.DateTime` | `timestamp` |
| `bson.M` (nested doc) | `json` |
| `bson.A` (array) | `json` |
| `primitive.Binary` | `binary` |

---

### A4. Create MongoDB Type Resolver

**New file:** `internal/schema_processor/mongodb_resolver.go`

**Functions:**
| Function | Description |
|----------|-------------|
| `ResolveFieldType(stats) (FieldType, bool)` | Picks most common type; nullable if missing in >10% of docs |
| `GenerateModelFromSchema(collectionName, schema) Model` | Converts `CollectionSchema` to `Model` struct; `_id` is always primary key |

**Type Resolution Strategy:**
- Most common type wins (>50% threshold)
- If ambiguous, prefer more specific types (e.g., integer over string if strings are numeric)
- Nullability: field appears in <90% of sampled documents -> `nullable: true`

---

### A5. Create MongoDB Processor (Orchestrator)

**New file:** `internal/schema_processor/mongodb_processor.go`

**Struct:** `MongoDBProcessor` wrapping `MongoDBSampler`

**Functions:**
| Function | Description |
|----------|-------------|
| `NewMongoDBProcessor(uri, dbName) (*MongoDBProcessor, error)` | Connects, pings, returns processor |
| `GenerateModels(collectionNames, sampleSize) ([]Model, error)` | Per collection: sample -> infer -> resolve -> Model |
| `GenerateAndSaveModels(outputPath, collectionNames, sampleSize) error` | Generates and writes `models.json` |
| `Close() error` | Disconnects from MongoDB |

**Flow:**
```
Connect to MongoDB
    -> List collections (or use provided names)
    -> For each collection:
        -> Sample N documents ($sample aggregation)
        -> Infer schema from samples (type counts, nullability)
        -> Resolve final types (most common wins)
        -> Generate Model struct
    -> Write all models to models.json
```

---

### A6. Update CLI for MongoDB

**Modify:** `cmd/generate-models/main.go`

**New flags:**
| Flag | Default | Description |
|------|---------|-------------|
| `-type` | `postgres` | Database type: `postgres` or `mongodb` |
| `-mongodb-uri` | `MONGODB_URI` env | MongoDB connection URI |
| `-mongodb-db` | `MONGODB_DATABASE` env | MongoDB database name |
| `-collections` | (all) | Comma-separated collection names |
| `-sample-size` | `100` | Documents to sample per collection |

**Usage examples:**
```bash
# All collections with default sample size
generate-models -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "mydb"

# Specific collections, larger sample
generate-models -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "mydb" \
  -collections "users,orders,products" \
  -sample-size 500

# Using environment variables
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="mydb"
generate-models -type mongodb
```

---

### A7. Unit Tests for Schema Discovery

**New file:** `internal/schema_processor/mongodb_infer_test.go`

| Test | What it covers |
|------|---------------|
| `TestInferTypeFromValue` | All BSON types map correctly |
| `TestInferSchema` | Mixed types, missing fields, nested docs |
| `TestResolveFieldType` | Ambiguous types, nullability thresholds |
| `TestGenerateModelFromSchema` | Output matches expected `Model` format, `_id` as PK |

---

## Phase B: MongoDB Adapter & Full Support (Second Priority)

**Goal:** Abstract the database layer and implement a MongoDB adapter so the server can run queries against MongoDB.

**Estimated effort:** 56-78 hours

### B1. Define Database and QueryBuilder Interfaces

**Modify:** `internal/adapter/adapter.go`

```go
package adapter

import "udv/internal/planner"

// Database represents a generic database connection
type Database interface {
    Close() error
    Ping() error
    ExecuteQuery(query interface{}, args ...interface{}) ([]map[string]interface{}, error)
    Exec(query interface{}, args ...interface{}) (ExecResult, error)
}

// ExecResult wraps execution results
type ExecResult interface {
    RowsAffected() (int64, error)
}

// QueryBuilder converts a QueryPlan into a database-specific query
type QueryBuilder interface {
    BuildQuery(plan *planner.QueryPlan) (query interface{}, args []interface{}, err error)
}
```

---

### B2. Refactor PostgreSQL Adapter to Implement Interfaces

**Modify:** `internal/adapter/postgres/db.go` and `internal/adapter/postgres/builder.go`

| Change | Details |
|--------|---------|
| `postgres.Database` implements `adapter.Database` | `ExecuteQuery` wraps existing `ExecuteAndFetchRows` |
| `postgres.QueryBuilder` implements `adapter.QueryBuilder` | `BuildQuery` returns `(interface{}, []interface{}, error)` wrapping the SQL string |
| Zero behavior change | Existing PostgreSQL path must work identically |

---

### B3. Implement MongoDB Adapter

**New directory:** `internal/adapter/mongodb/`

#### `db.go` -- Connection & Execution
| Method | Description |
|--------|-------------|
| `Connect(uri, dbName) (*Database, error)` | Connect with ping test |
| `Close() error` | Disconnect |
| `Ping() error` | Connection health check |
| `ExecuteQuery(query, args...) ([]map[string]interface{}, error)` | Execute MongoQuery struct |
| `Exec(query, args...) (ExecResult, error)` | Execute insert/update/delete |

#### `builder.go` -- Query Builder
| Method | Description |
|--------|-------------|
| `BuildQuery(plan) (interface{}, []interface{}, error)` | Routes by operation type |
| `buildFindQuery(plan) (*MongoQuery, error)` | Filters -> `bson.M`, FindOptions |
| `buildFilter(filterExpr) (bson.M, error)` | DSL operators -> MongoDB operators |
| `buildAggregationPipeline(plan) (mongo.Pipeline, error)` | GROUP BY -> `$match/$group/$sort` |
| `buildInsert(plan) (*MongoQuery, error)` | InsertOne document |
| `buildUpdate(plan) (*MongoQuery, error)` | UpdateOne/Many with `$set` |
| `buildDelete(plan) (*MongoQuery, error)` | DeleteOne/Many |

**Operator mapping (DSL -> MongoDB):**
| DSL Operator | MongoDB Operator |
|-------------|-----------------|
| `=` | `$eq` |
| `!=` | `$ne` |
| `>` | `$gt` |
| `>=` | `$gte` |
| `<` | `$lt` |
| `<=` | `$lte` |
| `in` | `$in` |
| `not_in` | `$nin` |
| `like` / `contains` | `$regex` |
| `is_null` | `$exists: false` |
| `not_null` | `$exists: true` |

**Aggregate function mapping:**
| DSL Function | MongoDB Stage |
|-------------|--------------|
| `count` | `$sum: 1` |
| `sum` | `$sum: "$field"` |
| `avg` | `$avg: "$field"` |
| `min` | `$min: "$field"` |
| `max` | `$max: "$field"` |

#### `types.go` -- MongoQuery Struct
```go
type MongoQuery struct {
    Collection string
    Operation  string           // "find", "aggregate", "insert", "update", "delete"
    Filter     bson.M
    Pipeline   mongo.Pipeline
    Update     bson.M
    Document   bson.M
    Options    *options.FindOptions
}
```

---

### B4. Update API Layer to Use Interfaces

**Modify:** `internal/api/api.go`

| Change | Before | After |
|--------|--------|-------|
| Builder field | `builder *postgres.QueryBuilder` | `builder adapter.QueryBuilder` |
| DB field | `db *postgres.Database` | `db adapter.Database` |
| Import | `"udv/internal/adapter/postgres"` | `"udv/internal/adapter"` |
| Constructor | `New(reg, db *postgres.Database)` | `New(reg, db adapter.Database, builder adapter.QueryBuilder)` |

---

### B5. Update Server Entry Point

**Modify:** `cmd/server/main.go`

```go
dbType := os.Getenv("DB_TYPE") // "postgres" or "mongodb"

var db adapter.Database
var builder adapter.QueryBuilder

switch dbType {
case "mongodb":
    mongoURI := os.Getenv("MONGODB_URI")
    mongoDBName := os.Getenv("MONGODB_DATABASE")
    db, err = mongodb.Connect(mongoURI, mongoDBName)
    builder = mongodb.NewQueryBuilder()
case "postgres", "":
    dbURL := os.Getenv("DATABASE_URL")
    db, err = postgres.Connect(dbURL)
    builder = postgres.NewQueryBuilder()
}
```

**Environment variables:**
| Variable | Description |
|----------|-------------|
| `DB_TYPE` | `postgres` (default) or `mongodb` |
| `MONGODB_URI` | MongoDB connection string |
| `MONGODB_DATABASE` | MongoDB database name |
| `DATABASE_URL` | PostgreSQL connection string (existing) |

---

### B6. Tests for MongoDB Adapter

| File | Tests |
|------|-------|
| `internal/adapter/mongodb/builder_test.go` | Filter conversion, aggregation pipeline, CRUD query building |
| `internal/adapter/mongodb/db_test.go` | Integration tests with real MongoDB (build-tagged) |

---

## Key Design Decisions

1. **Same `models.json` format** -- Schema discovery outputs the same JSON structure for MongoDB as PostgreSQL. The config loader, schema registry, planner, and validator all work unchanged.

2. **`interface{}` for queries** -- PostgreSQL returns SQL strings; MongoDB returns pipeline/filter structs. Each adapter's `ExecuteQuery` knows how to interpret its own query type.

3. **Dot notation for nested documents** -- `address.city` becomes a flat field in `models.json`, consistent with MongoDB's native dot notation for queries.

4. **`_id` as primary key** -- MongoDB collections always use `_id` as the primary key in generated models.

5. **Statistical type inference** -- Sample 100 documents by default, infer types by frequency. Re-run discovery when schema evolves.

---

## File Changes Summary

### New Files (Phase A -- Schema Discovery)
| File | Description |
|------|-------------|
| `internal/schema_processor/mongodb_sampler.go` | Document sampling from collections |
| `internal/schema_processor/mongodb_infer.go` | Type inference from BSON values |
| `internal/schema_processor/mongodb_resolver.go` | Type resolution and model generation |
| `internal/schema_processor/mongodb_processor.go` | Orchestrator (connect, sample, infer, save) |
| `internal/schema_processor/mongodb_infer_test.go` | Unit tests for inference logic |

### New Files (Phase B -- Adapter)
| File | Description |
|------|-------------|
| `internal/adapter/mongodb/db.go` | MongoDB connection and execution |
| `internal/adapter/mongodb/builder.go` | QueryPlan -> MongoDB query conversion |
| `internal/adapter/mongodb/types.go` | MongoQuery struct and constants |
| `internal/adapter/mongodb/builder_test.go` | Query builder unit tests |

### Modified Files
| File | Phase | Change |
|------|-------|--------|
| `go.mod` | A | Add `go.mongodb.org/mongo-driver` |
| `cmd/generate-models/main.go` | A | Add `-type mongodb` flags |
| `internal/adapter/adapter.go` | B | Define Database and QueryBuilder interfaces |
| `internal/adapter/postgres/db.go` | B | Implement adapter.Database interface |
| `internal/adapter/postgres/builder.go` | B | Implement adapter.QueryBuilder interface |
| `internal/api/api.go` | B | Use interfaces instead of concrete postgres types |
| `cmd/server/main.go` | B | Add DB_TYPE switch for adapter selection |

---

## Implementation Order

```
Phase A (Schema Discovery)
  A1. go get mongo-driver          <- dependency
  A2. mongodb_sampler.go           <- sampling
  A3. mongodb_infer.go             <- type inference
  A4. mongodb_resolver.go          <- type resolution
  A5. mongodb_processor.go         <- orchestrator
  A6. Update CLI                   <- user-facing
  A7. Tests                        <- verification

Phase B (Full MongoDB Support)
  B1. adapter.go interfaces        <- abstraction layer
  B2. Refactor postgres adapter    <- implement interfaces
  B3. MongoDB adapter              <- new adapter
  B4. Refactor api.go              <- use interfaces
  B5. Update server main.go        <- adapter selection
  B6. Tests                        <- verification
```

---

**Last Updated:** February 15, 2026
**Total Estimated Effort:** 90-126 hours (~2.5-3.5 weeks)
**Priority:** High
