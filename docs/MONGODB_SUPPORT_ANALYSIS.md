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
   - Hard-coded to PostgreSQL

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

### 2. Data Model

| PostgreSQL | MongoDB |
|------------|---------|
| Table → Collection | Collection → Table |
| Row → Document | Document → Row |
| Column → Field | Field → Column |
| Primary Key → _id | _id → Primary Key |

### 3. CRUD Operations

| Operation | PostgreSQL | MongoDB |
|-----------|-----------|---------|
| **Create** | `INSERT INTO ... VALUES ...` | `collection.InsertOne()` or `InsertMany()` |
| **Read** | `SELECT ... FROM ... WHERE ...` | `collection.Find()` with filter |
| **Update** | `UPDATE ... SET ... WHERE ...` | `collection.UpdateOne()` or `UpdateMany()` |
| **Delete** | `DELETE FROM ... WHERE ...` | `collection.DeleteOne()` or `DeleteMany()` |

---

## Implementation Requirements

### Phase 1: Create Database Interface Abstraction (8-12 hours)

**Goal:** Abstract database operations so both PostgreSQL and MongoDB can be used.

#### 1.1 Create Database Interface (2-3 hours)

**File:** `internal/adapter/database.go` (new)

```go
package adapter

// Database represents a generic database connection
type Database interface {
    // Connection management
    Close() error
    Ping() error
    
    // Query execution
    ExecuteQuery(query interface{}) ([]map[string]interface{}, error)
    ExecuteCommand(command interface{}) (interface{}, error)
    
    // CRUD operations
    Create(collection string, data map[string]interface{}) (interface{}, error)
    Read(collection string, filter interface{}, options QueryOptions) ([]map[string]interface{}, error)
    Update(collection string, filter interface{}, data map[string]interface{}) (int64, error)
    Delete(collection string, filter interface{}) (int64, error)
}

// QueryOptions represents query options (pagination, sorting, etc.)
type QueryOptions struct {
    Limit  int
    Offset int
    Sort   []SortOption
    Fields []string
}

type SortOption struct {
    Field     string
    Direction string // "asc" or "desc"
}

// QueryBuilder represents a generic query builder
type QueryBuilder interface {
    BuildQuery(plan *planner.QueryPlan) (interface{}, error)
}
```

#### 1.2 Refactor PostgreSQL Adapter (4-6 hours)

**Files to Modify:**
- `internal/adapter/postgres/db.go` - Implement `adapter.Database` interface
- `internal/adapter/postgres/builder.go` - Implement `adapter.QueryBuilder` interface
- `internal/adapter/postgres/adapter.go` - Export adapter factory

**Changes:**
- Wrap existing PostgreSQL code in interface implementation
- Convert SQL strings to interface{} for abstraction
- Maintain backward compatibility

#### 1.3 Update API to Use Interface (2-3 hours)

**File:** `internal/api/api.go`

**Changes:**
```go
type API struct {
    registry  *schema.Registry
    validator *dsl.Validator
    planner   *planner.Planner
    builder   adapter.QueryBuilder  // Changed from *postgres.QueryBuilder
    db        adapter.Database      // Changed from *postgres.Database
}

func New(reg *schema.Registry, db adapter.Database, builder adapter.QueryBuilder) *API {
    // ...
}
```

---

### Phase 2: Implement MongoDB Adapter (24-32 hours)

#### 2.1 MongoDB Connection (4-6 hours)

**File:** `internal/adapter/mongodb/db.go` (new)

**Dependencies:**
- `go.mongodb.org/mongo-driver/mongo`
- `go.mongodb.org/mongo-driver/mongo/options`

**Implementation:**
```go
package mongodb

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
    client   *mongo.Client
    database *mongo.Database
    ctx      context.Context
}

func Connect(uri string, databaseName string) (*Database, error) {
    ctx := context.Background()
    clientOptions := options.Client().ApplyURI(uri)
    
    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        return nil, err
    }
    
    // Test connection
    err = client.Ping(ctx, nil)
    if err != nil {
        return nil, err
    }
    
    db := client.Database(databaseName)
    
    return &Database{
        client:   client,
        database: db,
        ctx:      ctx,
    }, nil
}

func (d *Database) Close() error {
    return d.client.Disconnect(d.ctx)
}

func (d *Database) Ping() error {
    return d.client.Ping(d.ctx, nil)
}
```

#### 2.2 MongoDB Query Builder (12-16 hours)

**File:** `internal/adapter/mongodb/builder.go` (new)

**Key Functions:**
- `BuildQuery()` - Convert QueryPlan to MongoDB aggregation pipeline
- `buildFindQuery()` - For simple SELECT queries
- `buildAggregationPipeline()` - For complex queries with grouping
- `buildFilter()` - Convert DSL filters to MongoDB $match
- `buildSort()` - Convert to MongoDB $sort
- `buildPagination()` - Convert to $limit and $skip

**Example Implementation:**
```go
func (qb *QueryBuilder) BuildQuery(plan *planner.QueryPlan) (interface{}, error) {
    switch plan.Operation {
    case "select":
        return qb.buildFindQuery(plan)
    case "create":
        return qb.buildInsert(plan)
    case "update":
        return qb.buildUpdate(plan)
    case "delete":
        return qb.buildDelete(plan)
    default:
        return nil, fmt.Errorf("unsupported operation: %s", plan.Operation)
    }
}

func (qb *QueryBuilder) buildFindQuery(plan *planner.QueryPlan) (*MongoQuery, error) {
    query := &MongoQuery{
        Collection: plan.RootModel.Table,
        Filter:     qb.buildFilter(plan.Filters),
        Options:    qb.buildOptions(plan),
    }
    return query, nil
}

func (qb *QueryBuilder) buildFilter(filter planner.FilterExpr) bson.M {
    // Convert DSL filters to MongoDB filter
    // Example: {field: "status", op: "=", value: "active"}
    // Becomes: {"status": "active"}
}
```

#### 2.3 MongoDB CRUD Operations (8-10 hours)

**File:** `internal/adapter/mongodb/db.go` (extend)

**Implement:**
- `Create()` - Use `collection.InsertOne()` or `InsertMany()`
- `Read()` - Use `collection.Find()` with filter
- `Update()` - Use `collection.UpdateOne()` or `UpdateMany()`
- `Delete()` - Use `collection.DeleteOne()` or `DeleteMany()`

**Example:**
```go
func (d *Database) Create(collection string, data map[string]interface{}) (interface{}, error) {
    coll := d.database.Collection(collection)
    
    // Convert map to bson.M
    doc := bson.M(data)
    
    result, err := coll.InsertOne(d.ctx, doc)
    if err != nil {
        return nil, err
    }
    
    return result.InsertedID, nil
}

func (d *Database) Read(collection string, filter interface{}, options adapter.QueryOptions) ([]map[string]interface{}, error) {
    coll := d.database.Collection(collection)
    
    // Build find options
    findOptions := options.Find()
    if options.Limit > 0 {
        findOptions.SetLimit(int64(options.Limit))
    }
    if options.Offset > 0 {
        findOptions.SetSkip(int64(options.Offset))
    }
    if len(options.Sort) > 0 {
        sortDoc := bson.D{}
        for _, s := range options.Sort {
            direction := 1
            if s.Direction == "desc" {
                direction = -1
            }
            sortDoc = append(sortDoc, bson.E{Key: s.Field, Value: direction})
        }
        findOptions.SetSort(sortDoc)
    }
    
    cursor, err := coll.Find(d.ctx, filter, findOptions)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(d.ctx)
    
    var results []map[string]interface{}
    if err = cursor.All(d.ctx, &results); err != nil {
        return nil, err
    }
    
    return results, nil
}
```

---

### Phase 3: Schema Mapping (8-12 hours)

#### 3.1 Model Configuration (4-6 hours)

**File:** `internal/schema/model.go` (modify)

**Add database type to model:**
```go
type Model struct {
    Name       string
    Table      string  // Collection name for MongoDB
    Database   string  // Database type: "postgres" or "mongodb"
    PrimaryKey string
    Fields     []Field
}
```

#### 3.2 Field Type Mapping (4-6 hours)

**File:** `internal/adapter/mongodb/types.go` (new)

**Map PostgreSQL types to MongoDB BSON types:**
- `string` → `string`
- `integer` → `int32` or `int64`
- `float` → `double`
- `boolean` → `bool`
- `date` → `date`
- `timestamp` → `timestamp`
- `json` → `object` or `array`
- `uuid` → `string` (UUID as string in MongoDB)

---

### Phase 4: Configuration & CLI (4-6 hours)

#### 4.1 Configuration File (2-3 hours)

**File:** `config.yaml` or environment variables

**Add MongoDB configuration:**
```yaml
database:
  type: "mongodb"  # or "postgres"
  mongodb:
    uri: "mongodb://localhost:27017"
    database: "agentp"
  postgres:
    dsn: "postgres://user:pass@localhost/dbname"
```

#### 4.2 CLI Updates (2-3 hours)

**File:** `cmd/server/main.go`

**Changes:**
```go
func main() {
    dbType := os.Getenv("DB_TYPE") // "postgres" or "mongodb"
    
    var db adapter.Database
    var builder adapter.QueryBuilder
    
    switch dbType {
    case "mongodb":
        uri := os.Getenv("MONGODB_URI")
        dbName := os.Getenv("MONGODB_DATABASE")
        db, err = mongodb.Connect(uri, dbName)
        builder = mongodb.NewQueryBuilder()
    case "postgres", "":
        dsn := os.Getenv("POSTGRES_DSN")
        db, err = postgres.Connect(dsn)
        builder = postgres.NewQueryBuilder()
    default:
        log.Fatal("Unsupported database type")
    }
    
    api := api.New(registry, db, builder)
    // ...
}
```

---

### Phase 5: Testing & Edge Cases (12-16 hours)

#### 5.1 Unit Tests (6-8 hours)
- Test MongoDB connection
- Test query builder
- Test CRUD operations
- Test filter conversion
- Test aggregation pipeline

#### 5.2 Integration Tests (4-6 hours)
- End-to-end query execution
- Test with real MongoDB instance
- Test error handling
- Test performance

#### 5.3 Edge Cases (2-2 hours)
- Handle MongoDB-specific features (ObjectId, arrays, nested documents)
- Handle schema-less nature (missing fields)
- Handle different data types
- Handle MongoDB limitations (no JOINs, use $lookup)

---

## Key Challenges

### 1. Query Language Translation

**Challenge:** SQL to MongoDB Query Language conversion

**Solution:**
- Use aggregation pipeline for complex queries
- Map SQL operators to MongoDB operators:
  - `=` → `$eq`
  - `!=` → `$ne`
  - `>` → `$gt`
  - `<` → `$lt`
  - `>=` → `$gte`
  - `<=` → `$lte`
  - `IN` → `$in`
  - `LIKE` → `$regex`
  - `IS NULL` → `$exists: false`

### 2. Joins (Relationships)

**Challenge:** MongoDB doesn't support SQL JOINs

**Solution:**
- Use `$lookup` aggregation stage
- For simple relationships, fetch related documents separately
- Consider denormalization for frequently accessed relationships

### 3. Aggregations

**Challenge:** Different aggregation syntax

**Solution:**
- Convert GROUP BY to `$group` stage
- Convert aggregate functions:
  - `COUNT` → `$sum: 1`
  - `SUM` → `$sum`
  - `AVG` → `$avg`
  - `MIN` → `$min`
  - `MAX` → `$max`

### 4. Transactions

**Challenge:** MongoDB transactions work differently

**Solution:**
- Use MongoDB sessions for transactions
- Support multi-document transactions (requires replica set)
- Handle transaction errors appropriately

### 5. Schema-less Nature

**Challenge:** MongoDB is schema-less, but Agent P expects schema

**Solution:**
- Use model configuration to define expected schema
- Handle missing fields gracefully
- Validate data on insert/update
- Use MongoDB schema validation (optional)

---

## Implementation Effort Summary

| Phase | Task | Hours | Priority |
|-------|------|-------|----------|
| **Phase 1** | Database Interface Abstraction | 8-12 | 🔴 HIGH |
| **Phase 2** | MongoDB Adapter Implementation | 24-32 | 🔴 HIGH |
| **Phase 3** | Schema Mapping | 8-12 | 🟡 MEDIUM |
| **Phase 4** | Configuration & CLI | 4-6 | 🟡 MEDIUM |
| **Phase 5** | Testing & Edge Cases | 12-16 | 🔴 HIGH |
| **TOTAL** | | **56-78 hours** | |

**Estimated Timeline:** 1.5-2 weeks for a single developer

---

## Dependencies

### New Go Packages Required

```go
go.mongodb.org/mongo-driver/mongo
go.mongodb.org/mongo-driver/mongo/options
go.mongodb.org/mongo-driver/bson
go.mongodb.org/mongo-driver/bson/primitive
```

**Add to `go.mod`:**
```bash
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/bson
```

---

## File Structure After Implementation

```
internal/adapter/
├── database.go           # NEW: Database interface
├── query_builder.go      # NEW: QueryBuilder interface
├── postgres/
│   ├── db.go            # MODIFY: Implement interface
│   ├── builder.go       # MODIFY: Implement interface
│   └── adapter.go       # NEW: Factory functions
└── mongodb/              # NEW: MongoDB adapter
    ├── db.go            # MongoDB connection & CRUD
    ├── builder.go       # MongoDB query builder
    ├── types.go         # Type mapping
    └── adapter.go       # Factory functions
```

---

## Configuration Example

### Environment Variables

```bash
# Database type
DB_TYPE=mongodb  # or "postgres"

# MongoDB configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=agentp

# PostgreSQL configuration (if using postgres)
POSTGRES_DSN=postgres://user:pass@localhost/dbname
```

### Config File (Alternative)

```yaml
database:
  type: mongodb
  mongodb:
    uri: mongodb://localhost:27017
    database: agentp
    options:
      max_pool_size: 100
      min_pool_size: 10
  postgres:
    dsn: postgres://user:pass@localhost/dbname
    max_connections: 100
```

---

## Testing Strategy

### 1. Unit Tests

```go
// Test MongoDB connection
func TestMongoDBConnect(t *testing.T) {
    db, err := mongodb.Connect("mongodb://localhost:27017", "test")
    assert.NoError(t, err)
    defer db.Close()
}

// Test query builder
func TestMongoDBQueryBuilder(t *testing.T) {
    builder := mongodb.NewQueryBuilder()
    // Test filter conversion
    // Test aggregation pipeline
}

// Test CRUD operations
func TestMongoDBCRUD(t *testing.T) {
    // Test Create
    // Test Read
    // Test Update
    // Test Delete
}
```

### 2. Integration Tests

- Test with real MongoDB instance
- Test query execution end-to-end
- Test error handling
- Test performance with large datasets

### 3. Compatibility Tests

- Test same queries on both PostgreSQL and MongoDB
- Verify results are equivalent (where possible)
- Test edge cases specific to each database

---

## Migration Path

### Option 1: Dual Support (Recommended)

- Support both databases simultaneously
- Choose at startup via configuration
- Same API, different backends

### Option 2: Separate Builds

- Build separate binaries for each database
- Smaller binary size
- More complex deployment

### Option 3: Plugin System (Future)

- Load database adapters as plugins
- Support multiple databases in one instance
- More complex architecture

---

## Limitations & Considerations

### MongoDB-Specific Limitations

1. **No SQL JOINs** - Must use `$lookup` or separate queries
2. **Schema-less** - Need to handle missing fields
3. **ObjectId** - MongoDB's default ID type (vs auto-increment integers)
4. **Transactions** - Require replica set (not available in standalone)
5. **Case Sensitivity** - Field names are case-sensitive
6. **No Foreign Keys** - Relationships must be managed in application

### PostgreSQL-Specific Features Not Available in MongoDB

1. **Complex JOINs** - MongoDB `$lookup` is less powerful
2. **Subqueries** - Not directly supported
3. **Window Functions** - Not available
4. **Full-text Search** - Different implementation
5. **ACID Transactions** - Different model (requires replica set)

---

## Recommended Approach

### Step 1: Create Interface (Week 1)
- Abstract database operations
- Refactor PostgreSQL to use interface
- Update API to use interface

### Step 2: Implement MongoDB Adapter (Week 2)
- MongoDB connection
- Query builder
- CRUD operations
- Basic testing

### Step 3: Polish & Test (Week 3)
- Edge cases
- Performance optimization
- Documentation
- Integration tests

---

## Success Criteria

✅ **MongoDB Support Complete When:**
1. Can connect to MongoDB database
2. Can execute SELECT queries (find operations)
3. Can execute CREATE, UPDATE, DELETE operations
4. Can handle filters, sorting, pagination
5. Can handle aggregations (GROUP BY)
6. Can handle relationships (via $lookup)
7. Tests pass for all CRUD operations
8. Documentation updated

---

## Next Steps

1. **Review this analysis** - Confirm approach and priorities
2. **Create database interface** - Start with abstraction
3. **Implement MongoDB adapter** - Build MongoDB support
4. **Test thoroughly** - Ensure compatibility
5. **Update documentation** - Document MongoDB usage

---

**Last Updated:** Based on current codebase analysis  
**Estimated Effort:** 56-78 hours (~1.5-2 weeks)  
**Priority:** High (enables multi-database support)

