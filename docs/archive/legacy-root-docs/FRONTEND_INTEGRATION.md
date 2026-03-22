# Frontend Integration Guide

**Status**: ✅ Complete - Frontend fully integrated with backend API

**Last Updated**: January 26, 2026

---

## Overview

The Universal Data Viewer frontend has been fully integrated with the backend API. The React application now:
- Fetches available models from `GET /models`
- Executes queries via `POST /query`
- Displays backend-generated SQL and parameters
- Handles loading and error states gracefully

---

## Architecture

### Frontend Stack
- **Framework**: React 18.2 with TypeScript
- **Build Tool**: Vite 7.3.1
- **Styling**: Tailwind CSS
- **HTTP Client**: Fetch API (browser native)

### Backend Communication

```
Frontend Component
    ↓
  useEffect Hook (fetches on mount/update)
    ↓
  API Client (src/api/client.ts)
    ↓
  Backend HTTP Endpoint
    ↓
  Response (JSON)
    ↓
  State Update (React)
    ↓
  Re-render
```

---

## API Integration Points

### 1. App.tsx - Model Loading

**Purpose**: Fetch list of available models on application load

**Implementation**:
```typescript
useEffect(() => {
  const loadModels = async () => {
    try {
      const fetchedModels = await fetchModels()
      setModels(fetchedModels)
    } catch (err) {
      setError(err.message)
    }
  }
  loadModels()
}, [])
```

**Endpoint**: `GET http://localhost:8080/models`

**Response Format**:
```json
[
  {
    "name": "users",
    "table": "users",
    "primary_key": "id",
    "fields": [
      {"name": "id", "type": "integer"},
      {"name": "email", "type": "string"}
    ]
  }
]
```

### 2. ListView Component - Query Execution

**Purpose**: Execute simple SELECT queries with optional filters

**Process**:
1. User selects model and applies filters
2. Component builds DSL query using `buildDSLQuery()`
3. Calls `executeQuery()` with the DSL
4. Backend returns SQL and parameters
5. Currently displays SQL in console (mock data shown in UI)

**Filter Conversion** (UI → DSL):
- `equals` → `=`
- `contains` → `like`
- `startswith` → `starts_with`
- `endswith` → `ends_with`
- `gt` → `>`
- `lt` → `<`
- `gte` → `>=`
- `lte` → `<=`

**Example DSL Query**:
```javascript
{
  "model": "users",
  "filters": {
    "field": "email",
    "op": "contains",
    "value": "@example.com"
  },
  "pagination": {
    "limit": 100,
    "offset": 0
  }
}
```

**Generated SQL**:
```sql
SELECT * FROM users t0 
WHERE t0.email LIKE $1 
LIMIT $2 OFFSET $3;
```

### 3. GroupView Component - GROUP BY Queries

**Purpose**: Execute GROUP BY queries with aggregates

**Process**:
1. User selects grouping field
2. Component builds DSL with `group_by` parameter
3. Backend generates GROUP BY SQL with COUNT and SUM aggregates
4. Results displayed in collapsible groups

**Example DSL Query**:
```javascript
{
  "model": "orders",
  "group_by": ["status"],
  "aggregates": [
    {"fn": "count", "field": "", "alias": "count"},
    {"fn": "count", "field": "id", "alias": "total_rows"}
  ],
  "pagination": {"limit": 100, "offset": 0}
}
```

**Generated SQL**:
```sql
SELECT t0.status, COUNT(*) AS count, COUNT(t0.id) AS total_rows 
FROM orders t0 
GROUP BY t0.status 
LIMIT $1 OFFSET $2;
```

---

## API Client Module (src/api/client.ts)

### Exported Functions

#### `fetchModels(): Promise<Model[]>`
Fetches all available models from the backend.

```typescript
const models = await fetchModels()
// Returns: [{name: "users", table: "users", ...}, ...]
```

#### `executeQuery(query: unknown): Promise<QueryResponse>`
Executes a DSL query against the backend.

```typescript
const response = await executeQuery(dslQuery)
// Returns: {sql: "SELECT...", params: [...], error?: string}
```

#### `buildDSLQuery(...): any`
Constructs a DSL query object with optional filters and grouping.

**Parameters**:
- `modelName` (string) - Target model name
- `fields?` (string[]) - SELECT columns (empty = SELECT *)
- `filters?` (Filter[]) - WHERE conditions
- `groupByField?` (string) - GROUP BY field
- `limit?` (number) - LIMIT clause (default: 100)
- `offset?` (number) - OFFSET clause (default: 0)

**Example**:
```typescript
const query = buildDSLQuery(
  "orders",
  ["id", "total"],
  [{field: "status", op: "=", value: "PAID"}],
  undefined,
  50,
  0
)
```

---

## Component Integration

### App.tsx (Root Component)

**State Management**:
- `models`: Array of Model objects from backend
- `selectedModel`: Currently selected model name
- `loading`: Loading state while fetching models
- `error`: Error message if model loading fails
- `filters`: Array of user-applied filters
- `groupByField`: Field to group by (if any)

**Features**:
- ✅ Loads models from backend on mount
- ✅ Shows loading spinner while fetching
- ✅ Displays error message if backend is unavailable
- ✅ Dynamically extracts field names from backend model definitions
- ✅ Passes fields to FilterBuilder and GroupView components

### ListView Component

**State Management**:
- `data`: Query results (currently mock data)
- `loading`: Query execution state
- `error`: Query execution errors

**Behavior**:
- Fetches data when model or filters change
- Converts UI filters to DSL format
- Calls backend `/query` endpoint
- Logs generated SQL to browser console
- Falls back to mock data for UI display

### GroupView Component

**State Management**:
- `data`: Query results (currently mock data)
- `loading`: Aggregation query state
- `error`: Query execution errors
- `expandedGroups`: Which group headers are expanded

**Behavior**:
- Automatically fetches GROUP BY data when field selected
- Builds DSL with aggregates
- Logs generated SQL to console
- Falls back to mock data for UI display

---

## Environment Configuration

### Development

**Backend URL**: Defaults to `http://localhost:8080`

To change, add to window object in index.html:
```html
<script>
  window.REACT_APP_API_URL = 'http://your-backend-url'
</script>
```

### Production

Update the `API_BASE` constant in `src/api/client.ts`:
```typescript
const API_BASE = 'https://api.yourdomain.com'
```

---

## Running the Integration

### Prerequisites

1. **Backend Server** running on port 8080:
```bash
cd /Users/shubhamparamhans/Workspace/udv
go build -o server ./cmd/server
./server
```

2. **Frontend Dev Server**:
```bash
cd /Users/shubhamparamhans/Workspace/udv/frontend
nvm use 22
npm install
npm run dev
```

### Access Application

Open browser to: `http://localhost:5173`

### Verify Integration

**Check Browser Console** (F12 → Console tab):
- Should show "Generated SQL: SELECT..." messages
- No CORS errors
- Network tab shows requests to `http://localhost:8080/models` and `http://localhost:8080/query`

---

## Current Behavior vs Production-Ready

### ✅ Currently Working

- [x] Models fetched from `/models` endpoint
- [x] DSL queries built and sent to `/query` endpoint
- [x] Backend SQL generation working correctly
- [x] Loading states displayed
- [x] Error handling with fallbacks
- [x] Filter operator mapping complete
- [x] GROUP BY query building functional

### 🔄 Next Steps (For Production)

- [ ] **Execute Generated SQL**: Currently showing SQL in console; need to execute against actual database
- [ ] **Return Results to UI**: Backend needs to return query results (not just SQL)
- [ ] **Pagination**: Implement working pagination
- [ ] **Real Data Display**: Replace mock data with actual database query results
- [ ] **Query Caching**: Cache frequently executed queries
- [ ] **Error Recovery**: Better error handling for network failures
- [ ] **Performance**: Add debouncing for rapid filter changes

---

## Testing the Integration

### Test 1: Load Models

1. Open browser to `http://localhost:5173`
2. Verify models appear in left sidebar
3. Check browser console for no errors
4. Check Network tab for GET /models request

### Test 2: Execute Simple Query

1. Select "orders" model
2. Verify fields from backend appear in filter modal
3. Create a filter (e.g., status = PAID)
4. Check browser console for SQL output
5. Verify SQL matches expected format

### Test 3: Group By Query

1. Select "orders" model
2. Click "Group By" button
3. Select a field (e.g., "status")
4. Check console for GROUP BY SQL
5. Verify groups display with counts

### Test 4: Error Handling

1. Kill backend server
2. Reload frontend
3. Should show error message in header
4. UI should still be functional (using mock data fallback)

---

## Troubleshooting

### "Failed to fetch models" Error

**Causes**:
- Backend server not running
- Wrong API URL configured
- CORS issues (uncommon with fetch)

**Solution**:
```bash
# Check backend is running
curl http://localhost:8080/health

# Check API_BASE in src/api/client.ts
# Verify backend is accessible from frontend URL
```

### Models Display but Filters Fail

**Cause**: Field name type mismatch from backend

**Solution**:
- Backend returns `fields` array with `name` and `type`
- Frontend extracts just the names
- Verify `currentModelFields` is correctly populated

### SQL Logged but No Data Shown

**Expected**: Currently by design. Backend returns SQL; frontend uses mock data.

**To Implement Real Data**:
1. Backend needs to execute SQL and return results
2. Add result field to QueryResponse interface
3. Parse results in ListView/GroupView components
4. Display actual data instead of mock data

---

## File Structure

```
frontend/src/
├── api/
│   └── client.ts              # API integration (models, queries)
├── components/
│   ├── ListView/
│   │   └── ListView.tsx        # Backend query execution
│   ├── GroupView/
│   │   └── GroupView.tsx       # GROUP BY query execution
│   ├── FilterBuilder/
│   │   └── FilterBuilder.tsx   # (unchanged)
│   └── [other components]/
├── App.tsx                    # Model loading + state management
├── main.tsx                   # Entry point
└── [other files]/
```

---

## Summary

The UDV frontend is now **fully integrated** with the backend API:

✅ **Models** are loaded dynamically from `/models` endpoint  
✅ **Queries** are built as DSL and sent to `/query` endpoint  
✅ **SQL Generation** works correctly and is logged to console  
✅ **UI Components** properly pass field information from backend  
✅ **Error Handling** gracefully handles backend failures  
✅ **Type Safety** maintained with TypeScript interfaces  

**The next phase** is to have the backend execute the generated SQL and return actual results, which the frontend can then display in place of the current mock data.
