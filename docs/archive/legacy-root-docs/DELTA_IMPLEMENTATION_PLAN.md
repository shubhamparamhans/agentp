# Agent P - Delta Implementation Plan

## Overview

This document provides a detailed analysis of what it would take to implement each delta (difference) between Agent P and Odoo views, based on the current architecture and codebase.

**Current State Assessment:**
- ✅ Backend: Query DSL, validation, planning, SQL generation complete
- ✅ Backend: Supports sorting, pagination, filtering, grouping, aggregates
- ✅ Backend: API endpoints `/models` and `/query` working
- ✅ Frontend: Basic ListView, GroupView, FilterBuilder, DetailView components
- ⚠️ Frontend: Uses mock data, backend integration partial
- ❌ Frontend: Missing pagination/sorting UI controls
- ❌ Backend: Query execution returns SQL but not always data

---

## Implementation Effort Summary

| Feature | Complexity | Backend Effort | Frontend Effort | Total Estimate | Priority |
|---------|-----------|---------------|-----------------|----------------|----------|
| **Frontend Pagination UI** | Low | 0 hours | 4-6 hours | 4-6 hours | 🔴 HIGH |
| **Frontend Sorting UI** | Low | 0 hours | 4-6 hours | 4-6 hours | 🔴 HIGH |
| **Column Visibility Toggle** | Medium | 0 hours | 8-12 hours | 8-12 hours | 🟡 MEDIUM |
| **Export (CSV/Excel)** | Medium | 2-4 hours | 6-8 hours | 8-12 hours | 🟡 MEDIUM |
| **Saved Views** | High | 8-12 hours | 12-16 hours | 20-28 hours | 🟢 LOW |
| **Global Search** | Medium | 4-6 hours | 8-12 hours | 12-18 hours | 🟡 MEDIUM |
| **Complex Filter Logic** | Medium | 2-4 hours | 8-12 hours | 10-16 hours | 🟡 MEDIUM |
| **Column Reordering** | Medium | 0 hours | 8-12 hours | 8-12 hours | 🟢 LOW |
| **Multi-Select Rows** | Low | 0 hours | 4-6 hours | 4-6 hours | 🟢 LOW |
| **Bulk Operations** | High | 16-24 hours | 12-16 hours | 28-40 hours | 🟢 LOW |
| **CRUD Operations** | Very High | 40-60 hours | 40-60 hours | 80-120 hours | 🟢 LOW |
| **Relationship Widgets** | High | 12-16 hours | 16-24 hours | 28-40 hours | 🟢 LOW |
| **Security/RBAC** | Very High | 40-60 hours | 20-30 hours | 60-90 hours | 🟢 LOW |
| **Advanced Views** | Very High | 60-100 hours | 80-120 hours | 140-220 hours | 🟢 LOW |

---

## Detailed Implementation Plans

### 1. Frontend Pagination UI ⚠️ **HIGH PRIORITY**

**Status:** Backend ready, frontend missing

**Current State:**
- ✅ Backend: Supports `pagination: {limit, offset}` in DSL
- ✅ Backend: Returns paginated results
- ❌ Frontend: No pagination controls

**Implementation Steps:**

#### Frontend Work (4-6 hours)

1. **Create Pagination Component** (2 hours)
   - File: `frontend/src/components/Pagination/Pagination.tsx`
   - Features:
     - Page size selector (10, 25, 50, 100)
     - Previous/Next buttons
     - Page number navigation
     - Total count display
   - Props:
     ```typescript
     interface PaginationProps {
       currentPage: number
       pageSize: number
       totalCount: number
       onPageChange: (page: number) => void
       onPageSizeChange: (size: number) => void
     }
     ```

2. **Update ListView Component** (1-2 hours)
   - Add pagination state to `ListView.tsx`
   - Calculate `offset = (currentPage - 1) * pageSize`
   - Pass pagination to `buildDSLQuery()`
   - Display Pagination component below table

3. **Update API Client** (1 hour)
   - Modify `buildDSLQuery()` to accept page/pageSize
   - Calculate offset from page number
   - Update `executeQuery()` to handle pagination in response

4. **Update AppContext** (1 hour)
   - Add pagination state to global state
   - Persist pagination across filter changes
   - Reset to page 1 on model change

**Dependencies:**
- None (backend already supports)

**Testing:**
- Test with different page sizes
- Test pagination with filters applied
- Test edge cases (last page, empty results)

**Files to Modify:**
- `frontend/src/components/ListView/ListView.tsx`
- `frontend/src/api/client.ts`
- `frontend/src/state/AppContext.tsx`
- New: `frontend/src/components/Pagination/Pagination.tsx`

---

### 2. Frontend Sorting UI ⚠️ **HIGH PRIORITY**

**Status:** Backend ready, frontend missing

**Current State:**
- ✅ Backend: Supports `sort: [{field, direction}]` in DSL
- ✅ Backend: SQL generation includes ORDER BY
- ❌ Frontend: No clickable column headers

**Implementation Steps:**

#### Frontend Work (4-6 hours)

1. **Update ListView Headers** (2-3 hours)
   - Make column headers clickable
   - Add sort indicators (↑ ↓) based on current sort
   - Handle click to toggle sort direction
   - Support multi-column sorting (optional)

2. **Update State Management** (1 hour)
   - Add `sort` state to AppContext
   - Format: `{field: string, direction: 'asc' | 'desc'}[]`
   - Persist sort across filter changes

3. **Update API Client** (1 hour)
   - Modify `buildDSLQuery()` to include sort
   - Convert frontend sort format to DSL format

4. **Visual Polish** (1 hour)
   - Add hover states to sortable headers
   - Add visual indicators for sort direction
   - Style active sort column

**Dependencies:**
- None (backend already supports)

**Testing:**
- Test single column sorting
- Test sort direction toggle
- Test sorting with filters
- Test sorting with pagination

**Files to Modify:**
- `frontend/src/components/ListView/ListView.tsx`
- `frontend/src/api/client.ts`
- `frontend/src/state/AppContext.tsx`

---

### 3. Column Visibility Toggle

**Status:** Not implemented

**Current State:**
- ✅ Backend: Supports field selection via `fields: []` in DSL
- ❌ Frontend: Shows all columns, no toggle

**Implementation Steps:**

#### Frontend Work (8-12 hours)

1. **Create Column Manager Component** (4-6 hours)
   - File: `frontend/src/components/ColumnManager/ColumnManager.tsx`
   - Features:
     - Checkbox list of all available columns
     - Show/hide toggle for each column
     - Column reordering (drag & drop optional)
     - Save column preferences (localStorage)
   - Modal or dropdown interface

2. **Update ListView** (2-3 hours)
   - Filter displayed columns based on visibility state
   - Add "Manage Columns" button in header
   - Persist column preferences per model

3. **Update State Management** (1-2 hours)
   - Add `visibleColumns: string[]` to AppContext
   - Load from localStorage on model change
   - Save to localStorage on change

4. **Update API Client** (1 hour)
   - Pass visible columns to `buildDSLQuery()` as `fields`

**Dependencies:**
- None

**Testing:**
- Test column show/hide
- Test persistence across sessions
- Test with different models
- Test with pagination/sorting

**Files to Modify:**
- `frontend/src/components/ListView/ListView.tsx`
- `frontend/src/api/client.ts`
- `frontend/src/state/AppContext.tsx`
- New: `frontend/src/components/ColumnManager/ColumnManager.tsx`

---

### 4. Export Functionality (CSV/Excel)

**Status:** Not implemented

**Current State:**
- ✅ Backend: Can execute queries and return data
- ❌ Frontend: No export capability

**Implementation Steps:**

#### Backend Work (2-4 hours)

1. **Add Export Endpoint** (2-4 hours)
   - File: `internal/api/api.go`
   - Endpoint: `GET /export?format=csv|excel&query={base64_encoded_dsl}`
   - Or: `POST /export` with query DSL in body
   - Execute query and format response
   - Return file with appropriate headers

2. **CSV Generation** (1-2 hours)
   - Use Go's `encoding/csv` package
   - Handle special characters and escaping
   - Include headers

3. **Excel Generation** (Optional, 2-4 hours)
   - Use library like `github.com/xuri/excelize/v2`
   - Create .xlsx file with formatting
   - Multiple sheets support (optional)

#### Frontend Work (6-8 hours)

1. **Add Export Button** (1 hour)
   - Add to ListView header
   - Dropdown for format selection (CSV, Excel)

2. **Export Service** (2-3 hours)
   - File: `frontend/src/services/export.ts`
   - Function to call export endpoint
   - Handle file download
   - Show loading state

3. **Client-Side Export (Alternative)** (3-4 hours)
   - If backend export not available
   - Use library like `papaparse` for CSV
   - Use `xlsx` library for Excel
   - Export current visible data

**Dependencies:**
- Backend: Standard library (CSV), optional library (Excel)
- Frontend: Optional libraries for client-side export

**Testing:**
- Test CSV export with special characters
- Test Excel export formatting
- Test export with filters applied
- Test export with large datasets

**Files to Modify:**
- `internal/api/api.go` (new endpoint)
- `frontend/src/components/ListView/ListView.tsx`
- New: `frontend/src/services/export.ts`

---

### 5. Saved Views

**Status:** Not implemented

**Current State:**
- ✅ Backend: Query DSL is serializable
- ❌ No persistence layer

**Implementation Steps:**

#### Backend Work (8-12 hours)

1. **Database Schema** (2 hours)
   - Create `saved_views` table:
     ```sql
     CREATE TABLE saved_views (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       model VARCHAR(255) NOT NULL,
       query_dsl JSONB NOT NULL,
       user_id VARCHAR(255), -- optional for future RBAC
       created_at TIMESTAMP DEFAULT NOW(),
       updated_at TIMESTAMP DEFAULT NOW()
     );
     ```

2. **Saved View Model** (1-2 hours)
   - File: `internal/models/saved_view.go`
   - Struct definition
   - CRUD operations

3. **API Endpoints** (4-6 hours)
   - `GET /saved-views` - List all saved views
   - `GET /saved-views/{id}` - Get specific view
   - `POST /saved-views` - Create new view
   - `PUT /saved-views/{id}` - Update view
   - `DELETE /saved-views/{id}` - Delete view

4. **Query Execution from Saved View** (1-2 hours)
   - `POST /saved-views/{id}/execute` - Execute saved query
   - Or modify `/query` to accept `saved_view_id`

#### Frontend Work (12-16 hours)

1. **Saved Views Component** (4-6 hours)
   - File: `frontend/src/components/SavedViews/SavedViews.tsx`
   - List of saved views
   - Create/edit/delete functionality
   - Load and apply saved view

2. **Save View Dialog** (2-3 hours)
   - Modal to name and save current query
   - Include filters, sorting, grouping, pagination

3. **Update AppContext** (2-3 hours)
   - Add saved views state
   - Load saved views on mount
   - Apply saved view functionality

4. **UI Integration** (2-4 hours)
   - Add "Save View" button
   - Add "Load View" dropdown
   - Show saved views in sidebar

**Dependencies:**
- Backend: Database connection (already exists)
- Frontend: None

**Testing:**
- Test CRUD operations
- Test applying saved view
- Test persistence across sessions
- Test with different users (if RBAC added)

**Files to Modify:**
- New: `internal/models/saved_view.go`
- `internal/api/api.go` (new endpoints)
- `frontend/src/components/SavedViews/SavedViews.tsx`
- `frontend/src/state/AppContext.tsx`

---

### 6. Global Search

**Status:** Not implemented

**Current State:**
- ✅ Backend: Supports filtering on specific fields
- ❌ No multi-field search

**Implementation Steps:**

#### Backend Work (4-6 hours)

1. **Search Endpoint** (2-3 hours)
   - File: `internal/api/api.go`
   - Endpoint: `POST /search`
   - Accept: `{model, query, fields: []}`
   - Build OR filter across multiple fields
   - Support text search operators (contains, starts_with, etc.)

2. **Search Query Builder** (2-3 hours)
   - File: `internal/search/builder.go`
   - Convert search query to DSL
   - Generate OR filter across specified fields
   - Support full-text search (PostgreSQL `tsvector` - optional)

#### Frontend Work (8-12 hours)

1. **Search Bar Component** (3-4 hours)
   - File: `frontend/src/components/SearchBar/SearchBar.tsx`
   - Global search input
   - Auto-complete suggestions (optional)
   - Search across all fields or selected fields

2. **Search Integration** (2-3 hours)
   - Add to header
   - Update query when search changes
   - Debounce search input

3. **Search Results Highlighting** (2-3 hours)
   - Highlight matching text in results
   - Show which field matched

4. **Advanced Search Options** (1-2 hours)
   - Field selection for search
   - Search operator selection

**Dependencies:**
- Backend: None (uses existing filter infrastructure)
- Frontend: Optional debounce library

**Testing:**
- Test search across multiple fields
- Test search with special characters
- Test search performance with large datasets
- Test search with filters combined

**Files to Modify:**
- `internal/api/api.go` (new endpoint)
- New: `internal/search/builder.go` (optional)
- New: `frontend/src/components/SearchBar/SearchBar.tsx`
- `frontend/src/App.tsx`

---

### 7. Complex Filter Logic (Advanced AND/OR)

**Status:** Partially implemented

**Current State:**
- ✅ Backend: Supports AND/OR/NOT in DSL
- ✅ Backend: LogicalFilter structure exists
- ⚠️ Frontend: Only supports simple AND filters
- ❌ Frontend: No visual AND/OR builder

**Implementation Steps:**

#### Backend Work (2-4 hours)

1. **Verify Complex Filter Support** (1-2 hours)
   - Test nested AND/OR/NOT filters
   - Ensure SQL generation handles complex logic
   - Add tests for edge cases

2. **Filter Validation Enhancement** (1-2 hours)
   - Validate nested filter depth
   - Prevent circular references
   - Limit filter complexity

#### Frontend Work (8-12 hours)

1. **Advanced Filter Builder** (6-8 hours)
   - File: `frontend/src/components/FilterBuilder/AdvancedFilterBuilder.tsx`
   - Visual tree builder for AND/OR logic
   - Drag-and-drop filter groups (optional)
   - Add/remove filter groups
   - Nested group support

2. **Filter State Management** (2-3 hours)
   - Update filter structure to support groups
   - Convert UI state to DSL format
   - Handle nested logical operators

3. **UI Polish** (1-2 hours)
   - Visual indicators for AND/OR groups
   - Collapsible filter groups
   - Clear visual hierarchy

**Dependencies:**
- Backend: Already supports (verification needed)
- Frontend: None

**Testing:**
- Test nested AND/OR combinations
- Test NOT operator
- Test complex filter with grouping
- Test filter validation

**Files to Modify:**
- `frontend/src/components/FilterBuilder/FilterBuilder.tsx` (enhance)
- `frontend/src/state/AppContext.tsx`
- `frontend/src/api/client.ts`

---

### 8. Column Reordering

**Status:** Not implemented

**Current State:**
- ✅ Backend: Field order in `fields: []` array
- ❌ Frontend: No drag-and-drop

**Implementation Steps:**

#### Frontend Work (8-12 hours)

1. **Drag-and-Drop Library** (2 hours)
   - Install: `react-beautiful-dnd` or `@dnd-kit/core`
   - Configure for table headers

2. **Column Reordering** (4-6 hours)
   - Make column headers draggable
   - Update column order state
   - Persist order (localStorage)
   - Visual feedback during drag

3. **Integration with Column Manager** (2-4 hours)
   - Combine with column visibility toggle
   - Unified column management UI

**Dependencies:**
- Frontend: Drag-and-drop library

**Testing:**
- Test drag-and-drop functionality
- Test persistence
- Test with different screen sizes

**Files to Modify:**
- `frontend/src/components/ListView/ListView.tsx`
- `frontend/src/components/ColumnManager/ColumnManager.tsx`

---

### 9. Multi-Select Rows

**Status:** Not implemented

**Current State:**
- ❌ No row selection mechanism

**Implementation Steps:**

#### Frontend Work (4-6 hours)

1. **Add Checkboxes** (2-3 hours)
   - Add checkbox column to table
   - Handle select all/none
   - Track selected row IDs

2. **Selection State** (1-2 hours)
   - Add `selectedRows: string[]` to state
   - Handle selection changes

3. **Selection Actions** (1-2 hours)
   - Show count of selected rows
   - Prepare for bulk operations (future)

**Dependencies:**
- None

**Testing:**
- Test single selection
- Test multi-selection
- Test select all
- Test with pagination

**Files to Modify:**
- `frontend/src/components/ListView/ListView.tsx`
- `frontend/src/state/AppContext.tsx`

---

### 10. Bulk Operations

**Status:** Not implemented (requires CRUD)

**Current State:**
- ❌ Read-only system
- ❌ No bulk operation support

**Implementation Steps:**

#### Backend Work (16-24 hours)

1. **Bulk Update Endpoint** (8-12 hours)
   - `POST /bulk-update`
   - Accept: `{model, filters, updates: {field: value}}`
   - Generate UPDATE SQL with WHERE clause
   - Execute and return affected rows count

2. **Bulk Delete Endpoint** (4-6 hours)
   - `POST /bulk-delete`
   - Accept: `{model, filters}`
   - Generate DELETE SQL
   - Execute and return deleted count

3. **Transaction Support** (4-6 hours)
   - Wrap bulk operations in transactions
   - Rollback on error
   - Return success/failure status

#### Frontend Work (12-16 hours)

1. **Bulk Actions UI** (6-8 hours)
   - Show when rows selected
   - Dropdown with actions (Update, Delete)
   - Confirmation dialogs

2. **Bulk Update Form** (4-6 hours)
   - Modal to specify field updates
   - Field selection and value input
   - Preview of affected rows

3. **Bulk Delete Confirmation** (2-3 hours)
   - Show count of rows to delete
   - Confirmation dialog
   - Success/error feedback

**Dependencies:**
- Requires CRUD operations (Phase 2)
- Backend: Transaction support

**Testing:**
- Test bulk update with filters
- Test bulk delete
- Test transaction rollback
- Test with large datasets

**Files to Modify:**
- `internal/api/api.go` (new endpoints)
- `frontend/src/components/BulkActions/BulkActions.tsx`
- `frontend/src/components/ListView/ListView.tsx`

---

### 11. CRUD Operations

**Status:** Not in MVP scope (Phase 2+)

**Current State:**
- ❌ Read-only system
- ❌ No create/update/delete

**Implementation Steps:**

#### Backend Work (40-60 hours)

1. **Create Endpoint** (8-12 hours)
   - `POST /models/{model}/records`
   - Accept record data
   - Validate against schema
   - Generate INSERT SQL
   - Return created record

2. **Update Endpoint** (8-12 hours)
   - `PUT /models/{model}/records/{id}`
   - Accept partial updates
   - Validate fields
   - Generate UPDATE SQL
   - Return updated record

3. **Delete Endpoint** (4-6 hours)
   - `DELETE /models/{model}/records/{id}`
   - Generate DELETE SQL
   - Handle cascading deletes (optional)

4. **Validation Layer** (8-12 hours)
   - Field-level validation
   - Required field checks
   - Type validation
   - Constraint validation

5. **Form Generation** (12-18 hours)
   - Generate form schema from model
   - Field type to input mapping
   - Relationship field handling

#### Frontend Work (40-60 hours)

1. **Create Form Component** (12-16 hours)
   - Auto-generate form from model
   - Field validation
   - Submit handling

2. **Edit Form Component** (12-16 hours)
   - Pre-populate with existing data
   - Handle updates
   - Validation

3. **Inline Editing** (8-12 hours)
   - Click cell to edit
   - Save on blur/enter
   - Cancel on escape

4. **Delete Confirmation** (4-6 hours)
   - Confirmation dialog
   - Handle delete action

5. **Error Handling** (4-6 hours)
   - Display validation errors
   - Handle API errors
   - User feedback

**Dependencies:**
- Backend: Validation system
- Frontend: Form library (optional)

**Testing:**
- Test create with all field types
- Test update partial records
- Test delete with relationships
- Test validation errors

**Files to Modify:**
- `internal/api/api.go` (new endpoints)
- `internal/validation/validator.go` (new)
- `frontend/src/components/CreateForm/CreateForm.tsx` (new)
- `frontend/src/components/EditForm/EditForm.tsx` (new)
- `frontend/src/components/ListView/ListView.tsx`

---

### 12. Relationship Widgets

**Status:** Not implemented

**Current State:**
- ✅ Backend: Relationship awareness
- ❌ Frontend: No relationship navigation

**Implementation Steps:**

#### Backend Work (12-16 hours)

1. **Relationship Data Endpoint** (4-6 hours)
   - `GET /models/{model}/records/{id}/relations/{relation}`
   - Return related records
   - Support pagination

2. **Relationship Query Enhancement** (4-6 hours)
   - Enhance query planner for relationship traversal
   - Support nested relationship queries
   - Optimize join queries

3. **Relationship Metadata** (4-6 hours)
   - Include relationship info in model metadata
   - Return available relationships per model

#### Frontend Work (16-24 hours)

1. **One2Many Widget** (6-8 hours)
   - Display related records in detail view
   - Inline list or table
   - Pagination for related records

2. **Many2One Widget** (4-6 hours)
   - Dropdown/autocomplete for selection
   - Display related record info
   - Navigate to related record

3. **Many2Many Widget** (4-6 hours)
   - Tag-based selection
   - Add/remove related records
   - Display selected records

4. **Relationship Navigation** (2-4 hours)
   - Click to navigate to related model
   - Breadcrumb navigation
   - Back button

**Dependencies:**
- Backend: Relationship query support
- Frontend: None

**Testing:**
- Test one2many display
- Test many2one selection
- Test many2many management
- Test navigation

**Files to Modify:**
- `internal/api/api.go` (new endpoints)
- `frontend/src/components/RelationshipWidgets/` (new)
- `frontend/src/components/DetailView/DetailView.tsx`

---

### 13. Security/RBAC

**Status:** Not in MVP scope (Phase 3+)

**Current State:**
- ❌ No authentication
- ❌ No authorization
- ❌ No user management

**Implementation Steps:**

#### Backend Work (40-60 hours)

1. **Authentication System** (12-18 hours)
   - JWT token generation/validation
   - Login endpoint
   - Token refresh
   - Password hashing

2. **User Management** (8-12 hours)
   - User model and storage
   - User CRUD operations
   - Password reset

3. **Role-Based Access Control** (12-18 hours)
   - Role model
   - Permission model
   - Role-permission mapping
   - Middleware for authorization

4. **Row-Level Security** (8-12 hours)
   - Filter queries by user permissions
   - Model-level access control
   - Field-level access control

#### Frontend Work (20-30 hours)

1. **Login/Auth UI** (6-8 hours)
   - Login form
   - Token storage
   - Auth state management

2. **Protected Routes** (4-6 hours)
   - Route guards
   - Redirect to login
   - Token refresh

3. **Permission-Based UI** (6-8 hours)
   - Hide/show features based on permissions
   - Disable actions user can't perform

4. **User Management UI** (4-8 hours)
   - User list
   - Role assignment
   - Permission management

**Dependencies:**
- Backend: JWT library, password hashing
- Frontend: Auth state management

**Testing:**
- Test authentication flow
- Test permission checks
- Test row-level security
- Test token expiration

**Files to Modify:**
- `internal/auth/` (new package)
- `internal/api/api.go` (add middleware)
- `frontend/src/components/Auth/` (new)
- `frontend/src/state/AuthContext.tsx` (new)

---

### 14. Advanced View Types (Kanban, Calendar, etc.)

**Status:** Not planned (out of scope)

**Current State:**
- ❌ Only list and tree views

**Implementation Complexity:** Very High (140-220 hours)

**Would Require:**
- Complete new view type implementations
- Backend support for view-specific queries
- Frontend components for each view type
- Configuration for view settings

**Recommendation:** Defer to Phase 4+ or consider as separate project

---

## Implementation Priority Roadmap

### Phase 1: MVP Completion (Immediate - 1-2 weeks)
1. ✅ Frontend Pagination UI (4-6 hours)
2. ✅ Frontend Sorting UI (4-6 hours)
3. ✅ Backend Query Execution Fix (if needed)

**Total: 8-12 hours**

### Phase 2: Essential Features (2-4 weeks)
1. Column Visibility Toggle (8-12 hours)
2. Export Functionality (8-12 hours)
3. Global Search (12-18 hours)
4. Complex Filter Logic (10-16 hours)

**Total: 38-58 hours**

### Phase 3: User Experience (4-6 weeks)
1. Saved Views (20-28 hours)
2. Column Reordering (8-12 hours)
3. Multi-Select Rows (4-6 hours)

**Total: 32-46 hours**

### Phase 4: CRUD & Advanced (8-12 weeks)
1. CRUD Operations (80-120 hours)
2. Relationship Widgets (28-40 hours)
3. Bulk Operations (28-40 hours)

**Total: 136-200 hours**

### Phase 5: Enterprise Features (12+ weeks)
1. Security/RBAC (60-90 hours)
2. Advanced View Types (140-220 hours)

**Total: 200-310 hours**

---

## Risk Assessment

### Low Risk (Quick Wins)
- Frontend Pagination UI
- Frontend Sorting UI
- Multi-Select Rows
- Column Visibility Toggle

### Medium Risk (Moderate Effort)
- Export Functionality
- Global Search
- Complex Filter Logic
- Saved Views

### High Risk (Significant Effort)
- CRUD Operations
- Relationship Widgets
- Security/RBAC
- Bulk Operations

---

## Dependencies & Prerequisites

### Backend Dependencies
- ✅ Query DSL system (complete)
- ✅ SQL generation (complete)
- ✅ Database connection (complete)
- ❌ Transaction support (needed for CRUD)
- ❌ Validation system (needed for CRUD)
- ❌ Authentication system (needed for RBAC)

### Frontend Dependencies
- ✅ React components (basic)
- ✅ State management (Context API)
- ✅ API client (basic)
- ❌ Drag-and-drop library (for column reordering)
- ❌ Export libraries (optional, for client-side export)
- ❌ Form library (optional, for CRUD)

---

## Testing Strategy

### Unit Tests
- Backend: Query DSL validation, SQL generation
- Frontend: Component rendering, state management

### Integration Tests
- Backend: API endpoint testing
- Frontend: API client testing

### E2E Tests
- Full user workflows
- Query execution with real data
- UI interactions

---

## Conclusion

The highest priority items (pagination and sorting UI) are **low-hanging fruit** with backend support already in place. These can be implemented quickly (8-12 hours total) to complete the MVP.

Medium-priority features (column management, export, search) require moderate effort (38-58 hours) and significantly improve usability.

High-priority features (CRUD, RBAC) are major undertakings (140-200+ hours) and should be planned as separate phases with proper architecture planning.

**Recommended Next Steps:**
1. Implement pagination and sorting UI (complete MVP)
2. Add column visibility toggle (high user value)
3. Add export functionality (frequently requested)
4. Plan Phase 2 architecture for CRUD operations

---

**Last Updated:** Based on current codebase analysis  
**Estimated Total Effort for All Deltas:** 400-600+ hours  
**Recommended MVP Completion:** 8-12 hours

