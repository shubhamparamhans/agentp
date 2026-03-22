
👉 **`POSTGRES_SQL_GENERATION.md`**

This document assumes:

* Validated DSL input
* Relationships already resolved by the Query Planner
* Parameterized queries (`$1, $2, ...`)
* PostgreSQL dialect (JSONB, ILIKE, etc.)

---

# Universal Data Viewer (UDV)
## PostgreSQL SQL Generation – DSL Examples

---

## 1. Purpose of This Document

This document demonstrates **how Query DSL objects are translated into PostgreSQL SQL**.

It serves as:
- A reference for backend implementation
- A validation of DSL completeness
- A guide for adapter authors
- A test oracle for query planner correctness

---

## 2. Assumptions

- All DSL input has passed validation
- Field names are mapped to column names
- Relationships are already resolved into JOINs
- All queries are parameterized
- SQL generation is deterministic

---

## 3. Example 1 – Simple Select with Pagination

### DSL Input

```json
{
  "model": "Order",
  "fields": ["id", "status", "amount"],
  "pagination": { "limit": 50, "offset": 0 }
}
````

### Generated SQL

```sql
SELECT
  orders.id,
  orders.status,
  orders.amount
FROM orders
LIMIT $1 OFFSET $2;
```

### Parameters

```
$1 = 50
$2 = 0
```

---

## 4. Example 2 – Filtering with AND Conditions

### DSL Input

```json
{
  "model": "Order",
  "fields": ["id", "status", "amount"],
  "filters": {
    "and": [
      { "field": "status", "op": "=", "value": "PAID" },
      { "field": "amount", "op": ">", "value": 1000 }
    ]
  }
}
```

### Generated SQL

```sql
SELECT
  orders.id,
  orders.status,
  orders.amount
FROM orders
WHERE
  orders.status = $1
  AND orders.amount > $2;
```

### Parameters

```
$1 = 'PAID'
$2 = 1000
```

---

## 5. Example 3 – OR + NOT Filters

### DSL Input

```json
{
  "model": "Order",
  "filters": {
    "or": [
      { "field": "status", "op": "=", "value": "PAID" },
      {
        "not": {
          "field": "status",
          "op": "=",
          "value": "CANCELLED"
        }
      }
    ]
  }
}
```

### Generated SQL

```sql
SELECT *
FROM orders
WHERE
  (
    orders.status = $1
    OR orders.status != $2
  );
```

### Parameters

```
$1 = 'PAID'
$2 = 'CANCELLED'
```

---

## 6. Example 4 – IN Operator

### DSL Input

```json
{
  "model": "Order",
  "filters": {
    "field": "status",
    "op": "in",
    "value": ["PAID", "PENDING"]
  }
}
```

### Generated SQL

```sql
SELECT *
FROM orders
WHERE
  orders.status = ANY($1);
```

### Parameters

```
$1 = ['PAID', 'PENDING']
```

---

## 7. Example 5 – Date Range Filter

### DSL Input

```json
{
  "model": "Order",
  "filters": {
    "field": "created_at",
    "op": "between",
    "value": ["2024-01-01", "2024-01-31"]
  }
}
```

### Generated SQL

```sql
SELECT *
FROM orders
WHERE
  orders.created_at BETWEEN $1 AND $2;
```

### Parameters

```
$1 = '2024-01-01'
$2 = '2024-01-31'
```

---

## 8. Example 6 – Sorting

### DSL Input

```json
{
  "model": "Order",
  "sort": [
    { "field": "created_at", "direction": "desc" }
  ]
}
```

### Generated SQL

```sql
SELECT *
FROM orders
ORDER BY orders.created_at DESC;
```

---

## 9. Example 7 – Group By with Aggregation

### DSL Input

```json
{
  "model": "Order",
  "group_by": ["status"],
  "aggregates": [
    { "fn": "sum", "field": "amount", "alias": "total_amount" },
    { "fn": "count", "alias": "order_count" }
  ]
}
```

### Generated SQL

```sql
SELECT
  orders.status,
  SUM(orders.amount) AS total_amount,
  COUNT(*) AS order_count
FROM orders
GROUP BY orders.status;
```

---

## 10. Example 8 – Sorting on Aggregate

### DSL Input

```json
{
  "model": "Order",
  "group_by": ["status"],
  "aggregates": [
    { "fn": "sum", "field": "amount", "alias": "total_amount" }
  ],
  "sort": [
    { "field": "total_amount", "direction": "desc" }
  ]
}
```

### Generated SQL

```sql
SELECT
  orders.status,
  SUM(orders.amount) AS total_amount
FROM orders
GROUP BY orders.status
ORDER BY total_amount DESC;
```

---

## 11. Example 9 – Relationship Traversal (JOIN)

Assume:

* `Order.user_id → User.id`

### DSL Input

```json
{
  "model": "Order",
  "fields": ["id", "user.email", "amount"]
}
```

### Generated SQL

```sql
SELECT
  orders.id,
  users.email,
  orders.amount
FROM orders
LEFT JOIN users
  ON orders.user_id = users.id;
```

---

## 12. Example 10 – Filter on Related Field

### DSL Input

```json
{
  "model": "Order",
  "filters": {
    "field": "user.email",
    "op": "ilike",
    "value": "%gmail%"
  }
}
```

### Generated SQL

```sql
SELECT *
FROM orders
LEFT JOIN users
  ON orders.user_id = users.id
WHERE
  users.email ILIKE $1;
```

### Parameters

```
$1 = '%gmail%'
```

---

## 13. Example 11 – Grouping on Related Field

### DSL Input

```json
{
  "model": "Order",
  "group_by": ["user.email"],
  "aggregates": [
    { "fn": "count", "alias": "order_count" }
  ]
}
```

### Generated SQL

```sql
SELECT
  users.email,
  COUNT(*) AS order_count
FROM orders
LEFT JOIN users
  ON orders.user_id = users.id
GROUP BY users.email;
```

---

## 14. Example 12 – JSON Field Filter (Postgres Specific)

Assume `metadata` is JSONB.

### DSL Input

```json
{
  "model": "Order",
  "filters": {
    "field": "metadata.source",
    "op": "=",
    "value": "mobile"
  }
}
```

### Generated SQL

```sql
SELECT *
FROM orders
WHERE
  orders.metadata ->> 'source' = $1;
```

---

## 15. Safety Guarantees

Postgres adapter MUST ensure:

* All values are parameterized
* Identifiers are schema-resolved (never user input)
* JOIN count does not exceed limits
* OFFSET/LIMIT always present

---

## 16. Mapping Summary

| DSL Concept | SQL Construct     |
| ----------- | ----------------- |
| model       | FROM table        |
| fields      | SELECT columns    |
| filters     | WHERE             |
| group_by    | GROUP BY          |
| aggregates  | SUM / COUNT / AVG |
| sort        | ORDER BY          |
| pagination  | LIMIT / OFFSET    |
| relations   | JOIN              |

---

## 17. Summary

These examples demonstrate that:

* The DSL maps cleanly to PostgreSQL
* No SQL ambiguity exists
* Query generation is deterministic
* Adapter logic is straightforward

This validates the DSL as **production-viable**.

