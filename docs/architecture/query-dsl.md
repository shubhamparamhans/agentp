

👉 **`QUERY_DSL_SPEC.md`**

---
# Universal Data Viewer (UDV)
## Formal Query DSL Specification

---

## 1. Purpose of the Query DSL

The Query DSL defines **how query intent is expressed** in UDV.

It acts as:
- The contract between frontend and backend
- A safe abstraction over SQL / NoSQL queries
- A validation boundary
- A deterministic query representation

The DSL is **declarative**, **structured**, and **side-effect free**.

---

## 2. Core Design Principles

- No raw SQL or database syntax
- Fully serializable (JSON)
- Deterministic interpretation
- Schema-aware validation
- Explicit over implicit

---

## 3. High-Level Query Object

Every query is represented as a single JSON object.

```json
{
  "model": "Order",
  "fields": [],
  "filters": {},
  "group_by": [],
  "aggregates": [],
  "sort": [],
  "pagination": {}
}
````

All keys are optional except `model`.

---

## 4. Model Selection

### 4.1 `model`

```json
"model": "Order"
```

* Must match a configured model name
* Represents the root table / collection
* All fields are resolved relative to this model

---

## 5. Field Selection

### 5.1 `fields`

```json
"fields": ["id", "status", "amount"]
```

Rules:

* Empty or omitted → backend decides default fields
* Fields must exist in schema
* Relationship traversal allowed (see §11)

---

## 6. Filtering

### 6.1 Filter Structure

Filters are represented as a **logical expression tree**.

```json
"filters": {
  "and": [
    { "field": "status", "op": "=", "value": "PAID" },
    { "field": "amount", "op": ">", "value": 1000 }
  ]
}
```

---

### 6.2 Logical Operators

| Operator | Description               |
| -------- | ------------------------- |
| and      | All conditions must match |
| or       | Any condition may match   |
| not      | Negation                  |

Example:

```json
{
  "not": {
    "field": "status",
    "op": "=",
    "value": "CANCELLED"
  }
}
```

---

### 6.3 Atomic Filter Condition

```json
{
  "field": "amount",
  "op": ">=",
  "value": 500
}
```

---

### 6.4 Supported Filter Operators

#### Generic Operators

| Operator | Meaning           |
| -------- | ----------------- |
| =        | Equals            |
| !=       | Not equals        |
| >        | Greater than      |
| >=       | Greater or equal  |
| <        | Less than         |
| <=       | Less or equal     |
| in       | Value in list     |
| not_in   | Value not in list |
| is_null  | Is NULL           |
| not_null | Is NOT NULL       |

---

#### String Operators

| Operator    | Meaning               |
| ----------- | --------------------- |
| like        | SQL LIKE              |
| ilike       | Case-insensitive LIKE |
| starts_with | Prefix match          |
| ends_with   | Suffix match          |
| contains    | Substring match       |

---

#### Date / Time Operators

| Operator | Meaning         |
| -------- | --------------- |
| before   | <               |
| after    | >               |
| between  | Inclusive range |

```json
{
  "field": "created_at",
  "op": "between",
  "value": ["2024-01-01", "2024-01-31"]
}
```

---

### 6.5 Validation Rules (Filters)

* Field must be `filterable`
* Operator must be valid for field type
* `in` and `between` require array values
* NULL checks must not include `value`

---

## 7. Grouping

### 7.1 `group_by`

```json
"group_by": ["status"]
```

Rules:

* Fields must be `groupable`
* Order matters (hierarchical grouping)
* Max depth enforced by system limits

---

### 7.2 Nested Grouping Semantics

```json
"group_by": ["status", "created_at"]
```

Represents:

```
Status
 └── Created Date
```

---

## 8. Aggregations

### 8.1 Aggregate Structure

```json
{
  "fn": "sum",
  "field": "amount",
  "alias": "total_amount"
}
```

---

### 8.2 Supported Aggregate Functions

| Function | Field Types    |
| -------- | -------------- |
| count    | any            |
| sum      | int, float     |
| avg      | int, float     |
| min      | any comparable |
| max      | any comparable |

---

### 8.3 Aggregate Rules

* Field must be `aggregatable`
* `count` may omit field (`count(*)`)
* Aggregates require `group_by` unless global

---

## 9. Sorting

### 9.1 Sort Structure

```json
"sort": [
  { "field": "created_at", "direction": "desc" }
]
```

Rules:

* Sorting on aggregated fields allowed
* Sorting on non-selected fields allowed
* Direction defaults to `asc`

---

## 10. Pagination

### 10.1 Pagination Object

```json
"pagination": {
  "limit": 50,
  "offset": 0
}
```

Rules:

* Limit is mandatory (default applied if missing)
* Max limit enforced by backend
* Offset must be ≥ 0

---

## 11. Relationship Traversal

### 11.1 Field Path Syntax

```json
"user.email"
```

Meaning:

* Join `User` via configured relationship
* Select `email` field

---

### 11.2 Relationship Resolution Rules

* Only configured relationships allowed
* Join direction inferred from schema
* Max join depth enforced
* Cyclic traversal allowed but limited

---

## 12. Query Result Shape

### 12.1 Flat Query Result

```json
{
  "rows": [
    { "id": "1", "status": "PAID", "amount": 1200 }
  ],
  "meta": {
    "total": 100,
    "limit": 50,
    "offset": 0
  }
}
```

---

### 12.2 Grouped Query Result

```json
{
  "groups": [
    {
      "key": "PAID",
      "aggregates": { "total_amount": 5000 },
      "rows": [ ... ]
    }
  ]
}
```

---

## 13. Error Model

### 13.1 Validation Error

```json
{
  "error": "INVALID_FILTER",
  "message": "Field 'amount' is not filterable"
}
```

---

### 13.2 Execution Error

```json
{
  "error": "QUERY_EXECUTION_FAILED",
  "message": "Database timeout"
}
```

---

## 14. Safety & Limits

The backend MUST enforce:

* Max filters
* Max joins
* Max group depth
* Max limit
* Query timeout

Queries violating limits must fail fast.

---

## 15. Explicit Non-Goals

The DSL does NOT support:

* Arbitrary expressions
* Subqueries
* User-defined functions
* Free-form SQL

---

## 16. Summary

This Query DSL is the **core contract of UDV**.

If followed strictly:

* Frontend remains generic
* Backend remains safe
* Databases remain protected
* AI tooling can reason deterministically

All future features (CRUD, workflows, exports) build on top of this DSL.

```

