

👉 **`END_TO_END_EXAMPLE.md`**

This document is intentionally verbose and explicit so it can be fed into **AI IDEs, onboarding docs, or architecture reviews**.

---

````markdown
# Universal Data Viewer (UDV)
## End-to-End Example: Config → DSL → IR → SQL → UI

---

## 1. Scenario Overview

We want to explore **Orders** data and answer:

> “Show total order amount per user email for PAID orders, sorted by total amount descending.”

This will be expressed without writing:
- SQL
- ORM structs
- UI components

---

## 2. Database Schema (Assumed)

### orders table
| column | type |
|------|------|
| id | uuid |
| user_id | uuid |
| status | varchar |
| amount | numeric |
| created_at | timestamp |

### users table
| column | type |
|------|------|
| id | uuid |
| email | varchar |
| name | varchar |

---

## 3. Step 1 – Configuration (Model Definitions)

### `models.json`

```json
{
  "version": "1.0",
  "models": [
    {
      "name": "Order",
      "table": "orders",
      "primary_key": "id",
      "fields": [
        { "name": "id", "column": "id", "type": "uuid", "filterable": true },
        { "name": "status", "column": "status", "type": "string", "filterable": true, "groupable": true },
        { "name": "amount", "column": "amount", "type": "float", "aggregatable": true },
        { "name": "created_at", "column": "created_at", "type": "datetime" },
        { "name": "user_id", "column": "user_id", "type": "uuid" }
      ],
      "relations": [
        {
          "type": "many_to_one",
          "model": "User",
          "foreign_key": "user_id",
          "reference_key": "id"
        }
      ]
    },
    {
      "name": "User",
      "table": "users",
      "primary_key": "id",
      "fields": [
        { "name": "id", "column": "id", "type": "uuid" },
        { "name": "email", "column": "email", "type": "string", "groupable": true },
        { "name": "name", "column": "name", "type": "string" }
      ]
    }
  ]
}
````

✅ Backend now understands:

* Tables
* Fields
* Relationship: `Order.user_id → User.id`

---

## 4. Step 2 – Query DSL (Frontend → Backend)

### User Intent (UI)

* Filter: `status = PAID`
* Group by: `user.email`
* Aggregate: `SUM(amount)`
* Sort: `total_amount DESC`

### Query DSL Payload

```json
{
  "model": "Order",
  "filters": {
    "field": "status",
    "op": "=",
    "value": "PAID"
  },
  "group_by": ["user.email"],
  "aggregates": [
    {
      "fn": "sum",
      "field": "amount",
      "alias": "total_amount"
    }
  ],
  "sort": [
    {
      "field": "total_amount",
      "direction": "desc"
    }
  ],
  "pagination": {
    "limit": 50,
    "offset": 0
  }
}
```

✅ Frontend:

* Does NOT know SQL
* Does NOT know joins
* Only expresses intent

---

## 5. Step 3 – Query Planner IR (Backend Internal)

After validation and planning, the backend produces this **fully resolved IR**.

### QueryPlan (Conceptual)

```go
QueryPlan{
  RootModel: ModelRef{
    Name:  "Order",
    Table: "orders",
    Alias: "t0",
  },

  Joins: []JoinPlan{
    {
      Type:      LEFT,
      FromAlias:"t0",
      ToTable:  "users",
      ToAlias:  "t1",
      On: {
        Left:  ColumnRef{TableAlias:"t0", ColumnName:"user_id"},
        Right: ColumnRef{TableAlias:"t1", ColumnName:"id"},
      },
    },
  },

  Filters: ComparisonFilter{
    Left:     ColumnRef{TableAlias:"t0", ColumnName:"status"},
    Operator: EQ,
    Right:    ValueExpr{Value:"PAID", Type:String},
  },

  GroupBy: []GroupExpr{
    { Column: ColumnRef{TableAlias:"t1", ColumnName:"email"} },
  },

  Aggregates: []AggregateExpr{
    {
      Function: SUM,
      Column:   &ColumnRef{TableAlias:"t0", ColumnName:"amount"},
      Alias:    "total_amount",
    },
  },

  Sort: []SortExpr{
    {
      Target:    AGGREGATE,
      Aggregate: "total_amount",
      Direction: DESC,
    },
  },

  Pagination: Pagination{
    Limit:  50,
    Offset: 0,
  },
}
```

✅ At this stage:

* All joins are explicit
* All fields are resolved
* No DSL strings remain
* No SQL exists yet

---

## 6. Step 4 – PostgreSQL SQL Generation

### Generated SQL

```sql
SELECT
  t1.email,
  SUM(t0.amount) AS total_amount
FROM orders t0
LEFT JOIN users t1
  ON t0.user_id = t1.id
WHERE
  t0.status = $1
GROUP BY
  t1.email
ORDER BY
  total_amount DESC
LIMIT $2 OFFSET $3;
```

### Parameters

```
$1 = 'PAID'
$2 = 50
$3 = 0
```

✅ SQL is:

* Deterministic
* Parameterized
* Safe
* Index-friendly

---

## 7. Step 5 – Backend Response Shape

### API Response

```json
{
  "groups": [
    {
      "key": "alice@example.com",
      "aggregates": {
        "total_amount": 12500
      }
    },
    {
      "key": "bob@example.com",
      "aggregates": {
        "total_amount": 8400
      }
    }
  ],
  "meta": {
    "limit": 50,
    "offset": 0
  }
}
```

---

## 8. Step 6 – Frontend Rendering

### UI Interpretation

* View type: **Group View**
* Group key: `user.email`
* Metric: `total_amount`

### Rendered UI (Conceptual)

```
User Email               | Total Amount
-------------------------|--------------
alice@example.com        | 12,500
bob@example.com          | 8,400
```

If user expands a group (future feature):

* UI sends a **child query** scoped to that group
* Same DSL, additional filter injected

---

## 9. What This Example Proves

This walkthrough demonstrates:

✅ No ORM structs required
✅ No SQL written by users
✅ No UI code per model
✅ Full schema awareness
✅ Deterministic execution

Everything flows from **configuration + DSL**.

---

## 10. Mental Model Summary

```
Config defines reality
DSL defines intent
IR defines execution
SQL defines access
UI defines perception
```

Each layer:

* Knows only what it must
* Is replaceable
* Is testable in isolation

---

## 11. Why This Matters

This architecture enables:

* Rapid internal tooling
* AI-generated queries safely
* Future CRUD & workflows
* Multiple database backends
* Long-term maintainability

This is the foundation of a **serious developer platform**.
