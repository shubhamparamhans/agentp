# CRUD via Existing Query Endpoint - Implementation Plan

## Overview

Instead of creating separate endpoints (`POST /models/{model}/records`, `PUT /models/{model}/records/{id}`, etc.), we can extend the existing `/query` endpoint to support CRUD operations. This approach is **simpler and more consistent** with the current architecture.

---

## Current State

**Existing `/query` Endpoint:**
- Accepts DSL query JSON
- Generates SELECT SQL
- Returns data
- Already has validation, planning, SQL generation

**Current DSL Structure:**
```json
{
  "model": "users",
  "fields": ["id", "name", "email"],
  "filters": {...},
  "group_by": [...],
  "sort": [...],
  "pagination": {...}
}
```

---

## Proposed Extension

### Add `operation` Field to DSL

```json
{
  "operation": "select",  // "select" | "create" | "update" | "delete"
  "model": "users",
  ...
}
```

**Default:** If `operation` is omitted, default to `"select"` (backward compatible)

---

## Implementation Details

### 1. Extended DSL Structure

#### Create Operation
```json
{
  "operation": "create",
  "model": "users",
  "data": {
    "name": "John Doe",
    "email": "john@example.com",
    "age": 30
  }
}
```

#### Update Operation
```json
{
  "operation": "update",
  "model": "users",
  "id": 123,  // or use filters to identify record
  "data": {
    "name": "John Smith",
    "age": 31
  }
}
```

#### Delete Operation
```json
{
  "operation": "delete",
  "model": "users",
  "id": 123  // or use filters
}
```

#### Select Operation (Existing)
```json
{
  "operation": "select",  // optional, default
  "model": "users",
  "fields": [...],
  "filters": {...}
}
```

---

## Backend Implementation

### Step 1: Extend DSL Structure (2-3 hours)

**File:** `internal/dsl/query.go`

```go
type Operation string

const (
    OpSelect Operation = "select"
    OpCreate Operation = "create"
    OpUpdate Operation = "update"
    OpDelete Operation = "delete"
)

type Query struct {
    Operation  Operation              `json:"operation,omitempty"`  // NEW
    Model      string                 `json:"model"`
    Fields     []string               `json:"fields,omitempty"`
    Filters    FilterExpr             `json:"filters,omitempty"`
    GroupBy    []string               `json:"group_by,omitempty"`
    Aggregates []Aggregate            `json:"aggregates,omitempty"`
    Sort       []Sort                 `json:"sort,omitempty"`
    Pagination *Pagination            `json:"pagination,omitempty"`
    
    // NEW: For create/update operations
    Data       map[string]interface{} `json:"data,omitempty"`
    ID         interface{}            `json:"id,omitempty"`  // For update/delete
}
```

**Default Behavior:**
- If `Operation` is empty, default to `OpSelect`
- Maintains backward compatibility

### Step 2: Extend Validator (4-6 hours)

**File:** `internal/dsl/query.go` (extend Validator)

```go
func (v *Validator) ValidateQuery(q *Query) error {
    // ... existing validation ...
    
    // NEW: Validate operation-specific requirements
    switch q.Operation {
    case OpCreate:
        return v.validateCreate(q)
    case OpUpdate:
        return v.validateUpdate(q)
    case OpDelete:
        return v.validateDelete(q)
    case OpSelect, "":
        // Existing validation (default)
        return nil
    default:
        return fmt.Errorf("invalid operation: %s", q.Operation)
    }
}

func (v *Validator) validateCreate(q *Query) error {
    if q.Data == nil || len(q.Data) == 0 {
        return fmt.Errorf("data is required for create operation")
    }
    
    // Validate required fields
    model := v.registry.GetModel(q.Model)
    if model == nil {
        return fmt.Errorf("model not found: %s", q.Model)
    }
    
    // Check required fields
    for fieldName, field := range model.Fields {
        if !field.Nullable && q.Data[fieldName] == nil {
            return fmt.Errorf("required field missing: %s", fieldName)
        }
    }
    
    // Validate field types
    for fieldName, value := range q.Data {
        if !v.registry.FieldExists(q.Model, fieldName) {
            return fmt.Errorf("field not found: %s", fieldName)
        }
        // Type validation...
    }
    
    return nil
}

func (v *Validator) validateUpdate(q *Query) error {
    if q.ID == nil && q.Filters == nil {
        return fmt.Errorf("id or filters required for update operation")
    }
    
    if q.Data == nil || len(q.Data) == 0 {
        return fmt.Errorf("data is required for update operation")
    }
    
    // Validate fields being updated
    for fieldName := range q.Data {
        if !v.registry.FieldExists(q.Model, fieldName) {
            return fmt.Errorf("field not found: %s", fieldName)
        }
    }
    
    return nil
}

func (v *Validator) validateDelete(q *Query) error {
    if q.ID == nil && q.Filters == nil {
        return fmt.Errorf("id or filters required for delete operation")
    }
    return nil
}
```

### Step 3: Extend SQL Builder (12-18 hours)

**File:** `internal/adapter/postgres/builder.go`

```go
// BuildQuery now handles all operations
func (qb *QueryBuilder) BuildQuery(plan *planner.QueryPlan) (string, []interface{}, error) {
    switch plan.Operation {
    case planner.OpCreate:
        return qb.buildInsert(plan)
    case planner.OpUpdate:
        return qb.buildUpdate(plan)
    case planner.OpDelete:
        return qb.buildDelete(plan)
    case planner.OpSelect, "":
        return qb.buildSelect(plan)  // Existing logic
    default:
        return "", nil, fmt.Errorf("unsupported operation: %s", plan.Operation)
    }
}

func (qb *QueryBuilder) buildInsert(plan *planner.QueryPlan) (string, []interface{}, error) {
    // INSERT INTO table (field1, field2) VALUES ($1, $2) RETURNING *
    table := plan.RootModel.Table
    
    fields := []string{}
    values := []string{}
    
    for field, value := range plan.Data {
        fields = append(fields, field)
        qb.paramCount++
        values = append(values, fmt.Sprintf("$%d", qb.paramCount))
        qb.params = append(qb.params, value)
    }
    
    sql := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s) RETURNING *",
        table,
        strings.Join(fields, ", "),
        strings.Join(values, ", "),
    )
    
    return sql, qb.params, nil
}

func (qb *QueryBuilder) buildUpdate(plan *planner.QueryPlan) (string, []interface{}, error) {
    // UPDATE table SET field1=$1, field2=$2 WHERE id=$3 RETURNING *
    table := plan.RootModel.Table
    
    sets := []string{}
    for field, value := range plan.Data {
        qb.paramCount++
        sets = append(sets, fmt.Sprintf("%s = $%d", field, qb.paramCount))
        qb.params = append(qb.params, value)
    }
    
    where := ""
    if plan.ID != nil {
        qb.paramCount++
        where = fmt.Sprintf("WHERE %s = $%d", plan.RootModel.PrimaryKey, qb.paramCount)
        qb.params = append(qb.params, plan.ID)
    } else if plan.Filters != nil {
        wherePart, err := qb.buildWhereClause(plan.Filters)
        if err != nil {
            return "", nil, err
        }
        where = wherePart
    }
    
    sql := fmt.Sprintf(
        "UPDATE %s SET %s %s RETURNING *",
        table,
        strings.Join(sets, ", "),
        where,
    )
    
    return sql, qb.params, nil
}

func (qb *QueryBuilder) buildDelete(plan *planner.QueryPlan) (string, []interface{}, error) {
    // DELETE FROM table WHERE id=$1
    table := plan.RootModel.Table
    
    where := ""
    if plan.ID != nil {
        qb.paramCount++
        where = fmt.Sprintf("WHERE %s = $%d", plan.RootModel.PrimaryKey, qb.paramCount)
        qb.params = append(qb.params, plan.ID)
    } else if plan.Filters != nil {
        wherePart, err := qb.buildWhereClause(plan.Filters)
        if err != nil {
            return "", nil, err
        }
        where = wherePart
    }
    
    sql := fmt.Sprintf("DELETE FROM %s %s", table, where)
    
    return sql, qb.params, nil
}
```

### Step 4: Extend Query Planner (4-6 hours)

**File:** `internal/planner/planner.go`

```go
type QueryPlan struct {
    Operation  Operation  // NEW
    RootModel  *ModelRef
    // ... existing fields ...
    Data       map[string]interface{}  // NEW
    ID         interface{}              // NEW
}

func (p *Planner) PlanQuery(q *dsl.Query) (*QueryPlan, error) {
    plan := &QueryPlan{
        Operation: q.Operation,
        // ... existing planning ...
    }
    
    // NEW: Handle operation-specific planning
    if q.Operation == dsl.OpCreate || q.Operation == dsl.OpUpdate {
        plan.Data = q.Data
    }
    
    if q.Operation == dsl.OpUpdate || q.Operation == dsl.OpDelete {
        plan.ID = q.ID
    }
    
    return plan, nil
}
```

### Step 5: Update API Handler (2-3 hours)

**File:** `internal/api/api.go`

```go
func (a *API) handleQuery(w http.ResponseWriter, r *http.Request) {
    // ... existing code ...
    
    // Parse operation (default to "select" if not provided)
    var operation dsl.Operation = dsl.OpSelect
    if rq.Operation != "" {
        operation = dsl.Operation(rq.Operation)
    }
    
    q := dsl.Query{
        Operation:  operation,  // NEW
        Model:      rq.Model,
        Fields:     rq.Fields,
        GroupBy:    rq.GroupBy,
        Aggregates: rq.Aggregates,
        Sort:       rq.Sort,
        Pagination: rq.Pagination,
        Data:       rq.Data,     // NEW
        ID:         rq.ID,       // NEW
    }
    
    // ... rest of existing code ...
    
    // Execute query
    if a.db != nil {
        if operation == dsl.OpDelete {
            // DELETE returns affected rows count
            result, err := a.db.Execute(sql, params...)
            resp["affected_rows"] = result
        } else {
            // CREATE/UPDATE/SELECT return data
            rows, err := a.db.ExecuteAndFetchRows(sql, params...)
            if err == nil {
                resp["data"] = rows
            }
        }
    }
}
```

---

## Frontend Implementation

### API Client Updates (2-3 hours)

**File:** `frontend/src/api/client.ts`

```typescript
export async function createRecord(
  modelName: string,
  data: Record<string, any>
): Promise<QueryResponse> {
  const query = {
    operation: 'create',
    model: modelName,
    data: data,
  }
  return executeQuery(query)
}

export async function updateRecord(
  modelName: string,
  id: string | number,
  data: Record<string, any>
): Promise<QueryResponse> {
  const query = {
    operation: 'update',
    model: modelName,
    id: id,
    data: data,
  }
  return executeQuery(query)
}

export async function deleteRecord(
  modelName: string,
  id: string | number
): Promise<QueryResponse> {
  const query = {
    operation: 'delete',
    model: modelName,
    id: id,
  }
  return executeQuery(query)
}
```

---

## Effort Comparison

### Using Existing Endpoint (Recommended)

| Task | Hours | Notes |
|------|-------|-------|
| Extend DSL structure | 2-3 | Add operation, data, id fields |
| Extend validator | 4-6 | Operation-specific validation |
| Extend SQL builder | 12-18 | INSERT, UPDATE, DELETE builders |
| Extend query planner | 4-6 | Handle operation in planning |
| Update API handler | 2-3 | Parse operation, handle responses |
| Frontend API client | 2-3 | Helper functions |
| **Total Backend** | **24-36 hours** | |
| **Total Frontend** | **2-3 hours** | |
| **GRAND TOTAL** | **26-39 hours** | |

### Creating New Endpoints (Original Plan)

| Task | Hours |
|------|-------|
| Create endpoint | 8-12 |
| Update endpoint | 8-12 |
| Delete endpoint | 4-6 |
| Validation layer | 8-12 |
| **Total Backend** | **28-42 hours** | |
| **Total Frontend** | **30-40 hours** | |
| **GRAND TOTAL** | **58-82 hours** | |

**Savings: 32-43 hours** by using existing endpoint!

---

## Advantages of Using Existing Endpoint

### ✅ Simplicity
- Single endpoint for all operations
- Consistent API structure
- Reuse existing validation, planning, SQL generation

### ✅ Consistency
- Same DSL format for all operations
- Same error handling
- Same response format

### ✅ Less Code
- No new endpoint handlers
- Reuse existing infrastructure
- Less duplication

### ✅ Easier Testing
- Test all operations through one endpoint
- Reuse existing test infrastructure

### ✅ Backward Compatible
- Default operation is "select"
- Existing queries continue to work
- No breaking changes

---

## Example Usage

### Create Record
```typescript
const response = await executeQuery({
  operation: 'create',
  model: 'users',
  data: {
    name: 'John Doe',
    email: 'john@example.com',
    age: 30
  }
})

// Response: { data: [{ id: 123, name: 'John Doe', ... }] }
```

### Update Record
```typescript
const response = await executeQuery({
  operation: 'update',
  model: 'users',
  id: 123,
  data: {
    name: 'John Smith',
    age: 31
  }
})

// Response: { data: [{ id: 123, name: 'John Smith', ... }] }
```

### Delete Record
```typescript
const response = await executeQuery({
  operation: 'delete',
  model: 'users',
  id: 123
})

// Response: { affected_rows: 1 }
```

### Update Multiple Records (using filters)
```typescript
const response = await executeQuery({
  operation: 'update',
  model: 'users',
  filters: {
    field: 'status',
    op: '=',
    value: 'inactive'
  },
  data: {
    status: 'archived'
  }
})

// Response: { affected_rows: 5 }
```

---

## Implementation Steps

### Phase 1: Backend DSL Extension (6-9 hours)
1. Add `Operation`, `Data`, `ID` to Query struct
2. Update validator to handle operations
3. Add default operation logic

### Phase 2: SQL Generation (12-18 hours)
1. Add `buildInsert()` method
2. Add `buildUpdate()` method
3. Add `buildDelete()` method
4. Update `BuildQuery()` to route to correct builder

### Phase 3: Query Planning (4-6 hours)
1. Extend QueryPlan with operation, data, id
2. Update planner to handle operations

### Phase 4: API Handler (2-3 hours)
1. Parse operation from request
2. Handle different response formats
3. Update error handling

### Phase 5: Frontend (2-3 hours)
1. Add helper functions (createRecord, updateRecord, deleteRecord)
2. Update form components to use new functions

**Total: 26-39 hours** (vs 58-82 hours for separate endpoints)

---

## Security Considerations

### Operation Validation
- Validate operation is allowed (could restrict DELETE in future)
- Check permissions per operation (future RBAC)

### Data Validation
- Validate all fields in create/update
- Check required fields
- Validate types
- Check constraints

### ID/Filters Validation
- Ensure ID exists for update/delete
- Validate filters don't match too many records (safety limit)

---

## Response Format

### Create Response
```json
{
  "sql": "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *",
  "params": ["John Doe", "john@example.com"],
  "data": [
    {
      "id": 123,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

### Update Response
```json
{
  "sql": "UPDATE users SET name=$1 WHERE id=$2 RETURNING *",
  "params": ["John Smith", 123],
  "data": [
    {
      "id": 123,
      "name": "John Smith",
      "email": "john@example.com"
    }
  ]
}
```

### Delete Response
```json
{
  "sql": "DELETE FROM users WHERE id=$1",
  "params": [123],
  "affected_rows": 1
}
```

---

## Testing Strategy

### Unit Tests
- Test DSL parsing for each operation
- Test validation for each operation
- Test SQL generation for each operation

### Integration Tests
- Test create → read → update → delete flow
- Test error cases (invalid data, missing fields)
- Test constraint violations

### E2E Tests
- Full CRUD workflow
- Bulk operations via filters
- Error handling

---

## Migration Path

### Backward Compatibility
- Existing queries without `operation` field default to `"select"`
- No breaking changes
- Gradual adoption

### Future Enhancements
- Add `operation: "bulk-update"` for bulk operations
- Add `operation: "bulk-delete"` for bulk deletes
- Add transaction support: `transaction: true`

---

## Conclusion

**Using the existing `/query` endpoint is significantly simpler:**

- **26-39 hours** vs **58-82 hours** (saves 32-43 hours)
- Single endpoint for all operations
- Consistent API design
- Reuse existing infrastructure
- Backward compatible

**Recommended Approach:** Extend the existing endpoint rather than creating new ones.

---

## Next Steps

1. **Extend DSL** - Add operation, data, id fields (2-3 hours)
2. **Extend Validator** - Add operation validation (4-6 hours)
3. **Extend SQL Builder** - Add INSERT, UPDATE, DELETE (12-18 hours)
4. **Update Planner** - Handle operations (4-6 hours)
5. **Update API** - Parse operations (2-3 hours)
6. **Frontend Helpers** - Add create/update/delete functions (2-3 hours)

**Total: 26-39 hours** for basic CRUD via existing endpoint!

