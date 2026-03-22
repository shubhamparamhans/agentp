# Global Search Implementation Guide

## Overview

Global search allows users to search across multiple fields simultaneously using a single search query. This document details what's needed to implement this feature.

---

## Current State Analysis

### ✅ What We Already Have

1. **Backend OR Filter Support**
   - `LogicalFilter` with `Or` field exists in `internal/dsl/query.go`
   - SQL generation supports OR filters in `internal/adapter/postgres/builder.go`
   - Query validation supports OR filters

2. **String Search Operators**
   - `contains` - Substring match (LIKE '%value%')
   - `starts_with` - Prefix match (LIKE 'value%')
   - `ends_with` - Suffix match (LIKE '%value')
   - `ilike` - Case-insensitive LIKE
   - `like` - Case-sensitive LIKE

3. **Query Infrastructure**
   - DSL query system
   - Query validation
   - SQL generation
   - API endpoints

### ❌ What's Missing

1. **Backend**
   - No dedicated search endpoint
   - No search query builder utility
   - No field type detection for search (which fields are searchable)

2. **Frontend**
   - No search bar component
   - No search state management
   - No search result highlighting
   - No debouncing mechanism

---

## Implementation Requirements

### Backend Requirements (4-6 hours)

#### 1. Search Query Builder Utility

**File:** `internal/search/builder.go` (new)

**Purpose:** Convert a search query string into a DSL query with OR filters across multiple fields.

**Function Signature:**
```go
func BuildSearchQuery(
    modelName string,
    searchTerm string,
    searchFields []string,
    operator dsl.FilterOperator, // e.g., contains, starts_with, ilike
) (*dsl.Query, error)
```

**Implementation Logic:**
```go
// Pseudo-code
1. Validate model exists
2. Validate all searchFields exist in model
3. Filter searchFields to only string/text fields (optional: make configurable)
4. Create OR filter with multiple ComparisonFilters:
   - For each searchField:
     - Create ComparisonFilter{field: searchField, op: operator, value: searchTerm}
5. Combine into LogicalFilter with Or array
6. Return Query with filters set
```

**Example Output:**
```json
{
  "model": "users",
  "filters": {
    "or": [
      {"field": "name", "op": "contains", "value": "john"},
      {"field": "email", "op": "contains", "value": "john"},
      {"field": "description", "op": "contains", "value": "john"}
    ]
  }
}
```

**Considerations:**
- Should only search string/text fields (not numbers, dates, etc.)
- Should handle empty search term gracefully
- Should limit number of fields searched (performance)
- Should support case-insensitive search by default (use `ilike`)

#### 2. Search Endpoint (Optional - Can Use Existing /query)

**Option A: New Dedicated Endpoint**

**File:** `internal/api/api.go`

**Endpoint:** `POST /search`

**Request Body:**
```json
{
  "model": "users",
  "query": "john",
  "fields": ["name", "email", "description"], // optional, defaults to all string fields
  "operator": "contains" // optional, defaults to "contains"
}
```

**Response:** Same as `/query` endpoint

**Option B: Use Existing /query Endpoint (Recommended)**

No new endpoint needed! Frontend can build the OR filter query directly using existing `/query` endpoint.

**Recommendation:** Use Option B (existing endpoint) to reduce backend complexity.

#### 3. Field Type Detection for Search

**File:** `internal/schema/registry.go` (enhancement)

**Purpose:** Identify which fields are searchable (string/text types).

**Function:**
```go
func (r *Registry) GetSearchableFields(modelName string) ([]string, error)
```

**Logic:**
- Return all fields with type: `string`, `text`, `varchar`, `char`
- Exclude: `integer`, `float`, `boolean`, `date`, `timestamp`, etc.

---

### Frontend Requirements (8-12 hours)

#### 1. Search Bar Component

**File:** `frontend/src/components/SearchBar/SearchBar.tsx` (new)

**Features:**
- Text input for search query
- Clear button (X icon)
- Loading indicator (when searching)
- Placeholder text: "Search across all fields..."
- Debounced input (wait 300-500ms after user stops typing)

**Props:**
```typescript
interface SearchBarProps {
  value: string
  onChange: (value: string) => void
  onClear: () => void
  placeholder?: string
  disabled?: boolean
  loading?: boolean
}
```

**Styling:**
- Match dark theme (gray-800 background, cyan accents)
- Icon on left (🔍)
- Clear button on right (when value exists)
- Smooth transitions

#### 2. Search State Management

**File:** `frontend/src/App.tsx` (modify)

**Add State:**
```typescript
const [searchQuery, setSearchQuery] = useState<string>('')
const [searchFields, setSearchFields] = useState<string[]>([]) // optional: specific fields
```

**Integration:**
- Add SearchBar to header (next to Filters/Group By buttons)
- Pass searchQuery to ListView
- Reset search when model changes

#### 3. Search Query Building

**File:** `frontend/src/api/client.ts` (modify)

**New Function:**
```typescript
export function buildSearchQuery(
  modelName: string,
  searchTerm: string,
  searchFields: string[],
  operator: string = 'contains'
): any {
  if (!searchTerm || searchTerm.trim() === '') {
    return null // No search filter
  }

  // Build OR filter across all search fields
  const orFilters = searchFields.map(field => ({
    field,
    op: operator,
    value: searchTerm.trim()
  }))

  return {
    or: orFilters
  }
}
```

**Update `buildDSLQuery`:**
```typescript
export function buildDSLQuery(
  modelName: string,
  fields?: string[],
  filters?: Filter[],
  groupByField?: string,
  limit: number = 100,
  offset: number = 0,
  sort?: Sort[],
  searchFilter?: any // NEW: OR filter from search
): any {
  // ... existing code ...
  
  // Combine search filter with existing filters
  if (searchFilter && filters && filters.length > 0) {
    // AND the search filter with existing filters
    query.filters = {
      and: [
        searchFilter,
        ...(filters.length === 1 
          ? [filters[0]]
          : [{ and: filters }]
        )
      ]
    }
  } else if (searchFilter) {
    query.filters = searchFilter
  }
  // ... rest of existing code ...
}
```

#### 4. ListView Integration

**File:** `frontend/src/components/ListView/ListView.tsx` (modify)

**Add Props:**
```typescript
interface ListViewProps {
  // ... existing props ...
  searchQuery?: string
  searchFields?: string[]
}
```

**Update Query Building:**
```typescript
// In useEffect
const searchFilter = searchQuery 
  ? buildSearchQuery(modelName, searchQuery, searchFields || [])
  : null

const query = buildDSLQuery(
  modelName,
  modelFields.length > 0 ? modelFields : undefined,
  dslFilters.length > 0 ? dslFilters : undefined,
  undefined,
  pageSize,
  offset,
  sortArray,
  searchFilter // NEW
)
```

**Add to useEffect dependencies:**
```typescript
}, [modelName, filters, modelFields, currentPage, pageSize, sort, searchQuery, searchFields])
```

#### 5. Search Result Highlighting (Optional Enhancement)

**File:** `frontend/src/components/ListView/ListView.tsx` (enhancement)

**Function:**
```typescript
function highlightMatch(text: string, searchTerm: string): React.ReactNode {
  if (!searchTerm) return text
  
  const regex = new RegExp(`(${escapeRegex(searchTerm)})`, 'gi')
  const parts = String(text).split(regex)
  
  return parts.map((part, i) => 
    regex.test(part) ? (
      <mark key={i} className="bg-yellow-600 text-yellow-100">
        {part}
      </mark>
    ) : (
      <span key={i}>{part}</span>
    )
  )
}
```

**Usage in Table Cells:**
```typescript
<td className="px-6 py-4 text-sm text-gray-200">
  {searchQuery 
    ? highlightMatch(String(row[column] ?? ''), searchQuery)
    : String(row[column] ?? '')
  }
</td>
```

#### 6. Debouncing Implementation

**Option A: Custom Hook**

**File:** `frontend/src/hooks/useDebounce.ts` (new)

```typescript
import { useState, useEffect } from 'react'

export function useDebounce<T>(value: T, delay: number = 500): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value)

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)

    return () => {
      clearTimeout(handler)
    }
  }, [value, delay])

  return debouncedValue
}
```

**Usage in App.tsx:**
```typescript
const [searchInput, setSearchInput] = useState('')
const debouncedSearchQuery = useDebounce(searchInput, 500)

// Use debouncedSearchQuery in queries, not searchInput
```

**Option B: Library**

Install: `npm install use-debounce`

```typescript
import { useDebounce } from 'use-debounce'

const [searchInput, setSearchInput] = useState('')
const [debouncedSearchQuery] = useDebounce(searchInput, 500)
```

---

## Implementation Steps

### Phase 1: Backend (2-3 hours)

1. **Create Search Builder Utility** (1-2 hours)
   - Create `internal/search/builder.go`
   - Implement `BuildSearchQuery` function
   - Add tests for search query building
   - Handle edge cases (empty search, no fields, invalid fields)

2. **Enhance Schema Registry** (30 minutes)
   - Add `GetSearchableFields` method
   - Filter fields by type (string/text only)
   - Add tests

3. **Testing** (30 minutes)
   - Test search query generation
   - Test with existing filter system
   - Test SQL generation for OR filters

### Phase 2: Frontend Core (4-6 hours)

1. **Create Search Bar Component** (2-3 hours)
   - Create component file
   - Add styling (dark theme)
   - Add clear button
   - Add loading state

2. **Add Debouncing** (1 hour)
   - Create `useDebounce` hook or install library
   - Test debouncing behavior

3. **Update API Client** (1 hour)
   - Add `buildSearchQuery` function
   - Update `buildDSLQuery` to accept search filter
   - Test query building

4. **Integrate with ListView** (1-2 hours)
   - Add search props to ListView
   - Update query building in useEffect
   - Test search functionality

### Phase 3: Frontend Polish (2-3 hours)

1. **Add to App Header** (30 minutes)
   - Add SearchBar to header
   - Position next to Filters/Group By
   - Handle state management

2. **Search Result Highlighting** (1-2 hours)
   - Implement highlight function
   - Update table cell rendering
   - Style highlights

3. **Advanced Options** (1 hour)
   - Field selection dropdown (optional)
   - Search operator selection (optional)
   - Search scope indicator

---

## Technical Details

### Search Query Structure

**Input:**
```json
{
  "model": "users",
  "query": "john",
  "fields": ["name", "email"]
}
```

**Generated DSL:**
```json
{
  "model": "users",
  "filters": {
    "or": [
      {"field": "name", "op": "contains", "value": "john"},
      {"field": "email", "op": "contains", "value": "john"}
    ]
  }
}
```

**Generated SQL:**
```sql
SELECT * FROM users t0 
WHERE (t0.name LIKE '%john%' OR t0.email LIKE '%john%')
LIMIT $1 OFFSET $2
```

### Combining Search with Filters

**Search + Filters:**
```json
{
  "model": "users",
  "filters": {
    "and": [
      {
        "or": [
          {"field": "name", "op": "contains", "value": "john"},
          {"field": "email", "op": "contains", "value": "john"}
        ]
      },
      {"field": "status", "op": "=", "value": "active"}
    ]
  }
}
```

**SQL:**
```sql
SELECT * FROM users t0 
WHERE (
  (t0.name LIKE '%john%' OR t0.email LIKE '%john%')
  AND t0.status = $1
)
LIMIT $2 OFFSET $3
```

### Field Type Detection

**Searchable Field Types:**
- `string`
- `text`
- `varchar` (if supported)
- `char` (if supported)

**Non-Searchable Field Types:**
- `integer`, `int`
- `float`, `decimal`
- `boolean`
- `date`, `datetime`, `timestamp`
- `uuid`
- `json`

---

## Performance Considerations

### 1. Limit Search Fields

**Problem:** Searching across 50+ fields can be slow.

**Solution:**
- Default to searching only visible/selected fields
- Limit to top 10-20 most relevant fields
- Allow user to select specific fields

### 2. Index Usage

**Problem:** LIKE queries with leading wildcards (`%value%`) can't use indexes.

**Solution:**
- Use `starts_with` operator when possible (can use indexes)
- Consider PostgreSQL full-text search (`tsvector`) for advanced use cases
- Add database indexes on commonly searched fields

### 3. Debouncing

**Problem:** Searching on every keystroke creates too many requests.

**Solution:**
- Debounce search input (300-500ms delay)
- Cancel pending requests when new search starts
- Show loading state during search

### 4. Result Limiting

**Problem:** Large result sets slow down rendering.

**Solution:**
- Use pagination (already implemented)
- Limit initial results
- Show "X results found" message

---

## Testing Checklist

### Backend Tests

- [ ] Search query builder creates correct OR filter
- [ ] Search works with empty search term (returns all)
- [ ] Search validates field existence
- [ ] Search only includes string fields
- [ ] Search combines correctly with existing filters
- [ ] SQL generation for OR filters works correctly
- [ ] Search with special characters (SQL injection safe)

### Frontend Tests

- [ ] Search bar renders correctly
- [ ] Search input updates state
- [ ] Debouncing works (delays API calls)
- [ ] Clear button clears search
- [ ] Search query builds correctly
- [ ] Search combines with filters
- [ ] Search resets on model change
- [ ] Search highlighting works
- [ ] Search works with pagination
- [ ] Search works with sorting

### Integration Tests

- [ ] End-to-end search flow
- [ ] Search + Filter combination
- [ ] Search + Sort combination
- [ ] Search + Pagination combination
- [ ] Search performance with large datasets

---

## Alternative Approaches

### Option 1: Client-Side Search (Not Recommended)

**Pros:**
- No backend changes needed
- Instant results

**Cons:**
- Only searches loaded data
- Doesn't work with pagination
- Performance issues with large datasets

### Option 2: Full-Text Search (Advanced)

**PostgreSQL Full-Text Search:**
- Use `tsvector` and `tsquery`
- Better performance for text search
- Ranking and relevance scoring
- More complex implementation

**When to Use:**
- Large text fields (descriptions, content)
- Need relevance ranking
- Advanced search features

### Option 3: Elasticsearch Integration (Future)

**Pros:**
- Best search performance
- Advanced features (fuzzy search, autocomplete)
- Relevance scoring

**Cons:**
- Additional infrastructure
- Data synchronization complexity
- Overkill for MVP

---

## Estimated Effort

| Task | Backend | Frontend | Total |
|------|---------|----------|-------|
| **Core Implementation** | 2-3 hours | 4-6 hours | 6-9 hours |
| **Polish & Enhancements** | 1-2 hours | 2-3 hours | 3-5 hours |
| **Testing** | 1 hour | 1 hour | 2 hours |
| **Total** | **4-6 hours** | **7-10 hours** | **11-16 hours** |

---

## Dependencies

### Backend
- ✅ No new dependencies (uses existing DSL system)
- ✅ No database changes required

### Frontend
- Optional: `use-debounce` library (or custom hook)
- No other dependencies

---

## Files to Create/Modify

### New Files
- `internal/search/builder.go` - Search query builder
- `frontend/src/components/SearchBar/SearchBar.tsx` - Search bar component
- `frontend/src/hooks/useDebounce.ts` - Debounce hook (if custom)

### Modified Files
- `internal/schema/registry.go` - Add `GetSearchableFields` method
- `internal/api/api.go` - Optional: Add `/search` endpoint (or use existing `/query`)
- `frontend/src/api/client.ts` - Add search query building
- `frontend/src/components/ListView/ListView.tsx` - Integrate search
- `frontend/src/App.tsx` - Add search bar and state

---

## Conclusion

Global search can be implemented using **existing backend infrastructure** (OR filters). The main work is:

1. **Backend:** Utility to build OR filter queries (2-3 hours)
2. **Frontend:** Search bar component and integration (4-6 hours)
3. **Polish:** Highlighting, debouncing, advanced options (2-3 hours)

**Total: 8-12 hours** for a complete implementation.

The implementation leverages existing OR filter support, so no major backend changes are needed. The feature integrates seamlessly with existing filters, sorting, and pagination.

