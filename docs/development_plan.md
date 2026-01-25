

👉 **`DEVELOPMENT_PLAYBOOK.md`**

---

# Universal Data Viewer (UDV)

## Development Playbook (Docs → Code → Running System)

---

## 1. Mental Model: How These MD Files Should Be Used

Each document plays a **specific role at a specific phase**.

| Phase        | Primary Doc                    | Purpose                 |
| ------------ | ------------------------------ | ----------------------- |
| Architecture | Backend / Frontend Arch        | Boundaries & invariants |
| Contracts    | Config, DSL, IR                | Shared truth            |
| Execution    | SQL examples, Adapter skeleton | Implementation guide    |
| Delivery     | MVP Scope                      | Guardrails              |

> **Rule:** Never code without one “active” document open.

---

## 2. Repo Setup (Day 0 – 30 mins)

### 2.1 Initialize Repo

```bash
mkdir udv && cd udv
git init
```

### 2.2 Add Docs First

```bash
mkdir docs
# copy all MD files here
```

**Important:**
Commit docs first.

```bash
git add docs
git commit -m "docs: architecture, DSL, planner IR, MVP scope"
```

This gives Cursor/Copilot **ground truth**.

---

## 3. Backend Bootstrapping (Days 1–3)

### Active Docs

* `REPO_STRUCTURE.md`
* `BACKEND_ARCHITECTURE.md`
* `MVP_SCOPE.md`

---

### Step 3.1 – Create Backend Skeleton (No Logic Yet)

```bash
mkdir -p cmd/server internal/{config,schema,dsl,planner,ir,adapter/postgres,api}
```

Create empty files only:

```bash
touch cmd/server/main.go
```

Use Cursor prompt:

> “Create empty Go packages according to REPO_STRUCTURE.md.
> Do not add logic yet.”

**Why this works:**
Cursor respects file boundaries if they exist.

---

### Step 3.2 – Boot Server (Hello World API)

In `main.go`:

* Start HTTP server
* Add `/health` endpoint

This gives you:

* Running binary
* Deployment confidence
* Debug loop

Run:

```bash
go run ./cmd/server
```

---

### Step 3.3 – Add Config Loader (First Real Logic)

Active doc:

* `CONFIGURATION_SYSTEM.md`

Prompt Cursor:

> “Implement config loader + validator exactly as CONFIGURATION_SYSTEM.md.
> Start with model name, table, primary key validation only.”

Test by:

* Adding `configs/models.json`
* Booting server

❌ If server doesn’t fail on invalid config → stop and fix.

---

## 4. Frontend Bootstrapping (Days 2–4, Parallel)

### Active Docs

* `FRONTEND_ARCHITECTURE.md`
* `MVP_SCOPE.md`

---

### Step 4.1 – Scaffold Frontend

```bash
cd frontend
npm create vite@latest
npm install
```

Add folder structure from **Frontend Skeleton**.

Cursor prompt:

> “Create React component skeletons exactly as FRONTEND_ARCHITECTURE.md.
> No styling, no logic.”

---

### Step 4.2 – Mock Backend Contract

Before backend is ready:

Create `frontend/src/mock/models.json`

Render ModelExplorer from mock.

This allows:

* UI progress
* Faster feedback
* No backend dependency

---

## 5. Middle Layer: DSL → IR (Days 4–7)

### Active Docs

* `QUERY_DSL_SPEC.md`
* `QUERY_PLANNER_IR.md`

---

### Step 5.1 – Implement DSL Structs & Validation

In `internal/dsl/`

Cursor prompt:

> “Implement DSL structs and validation rules strictly following QUERY_DSL_SPEC.md.
> Return structured validation errors.”

Write **unit tests first**:

* Valid DSL
* Invalid field
* Invalid operator

---

### Step 5.2 – Implement Planner → IR (No SQL Yet)

Active doc:

* `QUERY_PLANNER_IR.md`

Prompt:

> “Convert validated DSL into QueryPlan IR.
> Do not generate SQL.”

Test:

* DSL → IR snapshot tests

This is where **most bugs will appear** — which is good.

---

## 6. SQL Generation Phase (Days 7–9)

### Active Docs

* `POSTGRES_SQL_GENERATION.md`
* `POSTGRES_ADAPTER_SKELETON.md`

---

### Step 6.1 – Implement Postgres Adapter Incrementally

Order:

1. SELECT + FROM
2. WHERE (single filter)
3. JOIN
4. GROUP BY
5. AGGREGATES
6. SORT
7. PAGINATION

After each step:

* Write golden SQL tests
* Compare against examples in doc

Cursor prompt style:

> “Implement WHERE clause builder for ComparisonFilter only.
> Do not add logical filters yet.”

---

## 7. First End-to-End Run (Day 10)

### Active Doc

* `END_TO_END_EXAMPLE.md`

---

### Step 7.1 – Seed Database

Create local Postgres:

```sql
INSERT INTO users ...
INSERT INTO orders ...
```

---

### Step 7.2 – Run Exact Example

Send the exact DSL from the doc.

Verify:

* SQL logs match doc
* UI shows grouped data

🎉 This is your **first real milestone**.

---

## 8. How to Use Cursor / Copilot Effectively

### Golden Rules

#### 1. One Document at a Time

Always say:

> “Follow QUERY_DSL_SPEC.md strictly.”

#### 2. Freeze Files

Tell Cursor:

> “Do not modify files outside `/internal/planner`.”

#### 3. Stop at Boundaries

Never let Cursor:

* Jump across packages
* Add “helpful” abstractions
* Change DSL shape

---

## 9. Testing Strategy While Building

| Layer   | Test Type             |
| ------- | --------------------- |
| Config  | Startup failure tests |
| DSL     | Validation tests      |
| Planner | DSL → IR snapshots    |
| Adapter | IR → SQL golden tests |
| API     | Black-box HTTP tests  |
| UI      | Contract mocks        |

---

## 10. Recommended Weekly Milestones

### Week 1

* Backend boots
* Config loads
* Frontend renders models

### Week 2

* DSL → IR works
* SQL generated correctly
* First end-to-end query

### Week 3

* Filter builder UI
* Group view UI
* Hardening & polish

---

## 11. Final Advice (Very Important)

> **Do not “improve” the design while coding the MVP.**

Every improvement idea goes into:

* `docs/POST_MVP_IDEAS.md`

Ship first. Improve later.

---

## TL;DR Execution Loop

```
Pick doc → Write minimal code → Test → Commit → Next doc
```
