

# 📕 FRONTEND_ARCHITECTURE.md
# Universal Data Viewer (UDV)
## Frontend Architecture & Technical Requirements

---

## 1. Frontend Purpose

The frontend is a **generic data exploration interface**.

It does not know business logic.
It does not know SQL.
It only understands:
- Models
- Fields
- Relationships
- Query intent

---

## 2. Design Philosophy

### Core Principles
- Zero hardcoded domain logic
- Metadata-driven UI
- Power-user friendly
- Fast feedback loops

### UI Bias
- Tables over charts
- Grouping over dashboards
- Exploration over presentation

---

## 3. High-Level Frontend Architecture


API Client
|
State Management
|
Query Builder
|
View Renderer
|
User Interface


---

## 4. Frontend Responsibilities

- Fetch model metadata
- Build query specifications
- Render list and tree views
- Manage local UI state
- Never infer data semantics

---

## 5. Core Frontend Modules

### 5.1 API Client

#### Purpose
Communicate with backend APIs.

#### Responsibilities
- Fetch models
- Execute queries
- Handle pagination
- Normalize responses

#### Requirements
- Typed API contracts
- Centralized error handling

---

### 5.2 Model Explorer

#### Purpose
Expose available models and relationships.

#### Responsibilities
- Display model list
- Display field metadata
- Display relationships

#### Requirements
- No assumptions about field meaning
- Render purely from metadata

---

### 5.3 Query Builder (UI Layer)

#### Purpose
Allow users to express query intent visually.

#### Responsibilities
- Filter builder
- Group-by selector
- Aggregate selector
- Sorting configuration

#### Requirements
- Type-aware controls
- AND / OR logic support
- Serialize into query specification

---

### 5.4 List View

#### Purpose
Render flat result sets.

#### Responsibilities
- Column selection
- Sorting
- Pagination
- Inline aggregates

#### Requirements
- Virtualized rendering
- Keyboard navigation
- Large dataset support

---

### 5.5 Tree / Group View

#### Purpose
Render hierarchical grouped data.

#### Responsibilities
- Expand / collapse groups
- Display aggregated summaries
- Nested group rendering

#### Requirements
- Lazy loading support
- Clear hierarchy visualization

---

## 6. Frontend State Model

### Core State

interface ViewState {
  model: string;
  query: QuerySpec;
  viewType: "list" | "group";
}


### Characteristics

* Serializable
* Shareable (future saved views)
* Deterministic

---

## 7. Error Handling

### Principles

* Backend errors are displayed verbatim
* Validation errors shown inline
* Partial failures do not crash UI

---

## 8. Non-Functional Requirements

### Performance

* Render 10k+ rows smoothly
* Avoid unnecessary re-renders
* Client-side caching of metadata

### Accessibility

* Keyboard-first navigation
* Screen-reader friendly tables

### Extensibility

* Component isolation
* Feature-flag driven UI extensions

---

## 9. Explicit Non-Responsibilities

The frontend must NOT:

* Generate SQL
* Infer joins automatically
* Store data permanently
* Enforce permissions

---

## 10. Summary

The frontend is a **query composition and visualization layer**, not a dashboard builder.

Its primary role is to:

* Help users express intent
* Visualize results efficiently
* Stay fully driven by backend metadata
