# 🚀 Quick Start Guide - UDV v1.0.0

**Status**: ✅ Production Ready  
**Last Updated**: January 26, 2026

---

## ⚡ 5-Minute Setup

### Prerequisites Check
```bash
# Verify Go is installed
go version
# Expected: go version go1.x

# Verify Node.js 22 is available
nvm use 22
node --version
# Expected: v22.x.x
```

### 1. Build Backend (30 seconds)
```bash
cd /Users/shubhamparamhans/Workspace/udv
go build -o server ./cmd/server
```

### 2. Start Backend (10 seconds)
```bash
# Without database (SQL generation only - for demos)
./server

# With database (requires DATABASE_URL env var)
DATABASE_URL="postgresql://user:pass@host:port/db" ./server
```

**Expected Output:**
```
Loaded 2 model(s):
  - users (table: users, primaryKey: id)
  - orders (table: orders, primaryKey: id)
Schema registry initialized with 2 model(s)
Database connection established
Server starting on :8080
```

### 3. Start Frontend (30 seconds)
```bash
cd frontend
nvm use 22
npm run dev
```

**Expected Output:**
```
  VITE v7.3.1  ready in 234 ms

  ➜  Local:   http://localhost:5173/
  ➜  press h to show help
```

### 4. Open in Browser
Visit `http://localhost:5173` and start exploring!

---

## 🎮 Using the Application

### 1. View Models
- **Left Sidebar**: Shows available models (Users, Orders)
- Click model name to select it
- Data table appears in main area

### 2. Apply Filters
- **Click "Filter" button** in header
- **Select field** from dropdown
- **Choose operator** (equals, contains, starts_with, etc.)
- **Enter value** to filter by
- Click **"Apply"** or press Enter
- Data updates instantly

### 3. Group Data
- **Click "Group By" button** in header
- **Select field** to group by
- Data reorganizes into **collapsible groups**
- Each group shows **count** and **statistics**
- Click group header to expand/collapse

### 4. View Details
- **Click any row** in the table
- **Right panel slides in** showing all fields
- Click **close button** or click overlay to close

### 5. Monitor SQL Generation
- **Open browser DevTools** (F12)
- **Go to Console tab**
- See `Generated SQL:` and `Parameters:` logs
- Shows exact SQL being sent to backend

---

## 📊 Example Queries

### Query 1: Simple Filter
**Action**: Apply filter "status = PAID"

**Generated SQL:**
```sql
SELECT * FROM orders t0 WHERE t0.status = $1 LIMIT $2 OFFSET $3;
```

**Parameters:** `["PAID", 10, 0]`

### Query 2: Multiple Filters
**Action**: Apply filters "status = PAID" AND "amount > 100"

**Generated SQL:**
```sql
SELECT * FROM orders t0 
WHERE t0.status = $1 AND t0.amount > $2 
LIMIT $3 OFFSET $4;
```

**Parameters:** `["PAID", 100, 10, 0]`

### Query 3: Group By
**Action**: Group by "status" field

**Generated SQL:**
```sql
SELECT t0.status, COUNT(*) AS count, COUNT(t0.id) AS total_rows 
FROM orders t0 
GROUP BY t0.status 
LIMIT $1 OFFSET $2;
```

**Parameters:** `[10, 0]`

### Query 4: Complex Query
**Action**: Filter "amount > 100" + Group by "status"

**Generated SQL:**
```sql
SELECT t0.status, COUNT(*) AS count, SUM(t0.amount) AS total 
FROM orders t0 
WHERE t0.amount > $1 
GROUP BY t0.status 
LIMIT $2 OFFSET $3;
```

**Parameters:** `[100, 10, 0]`

---

## 🔧 Testing the API Directly

### Test 1: Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{"status":"ok"}
```

### Test 2: List Models
```bash
curl http://localhost:8080/models
```

**Response:**
```json
[
  {
    "name": "users",
    "table": "users",
    "primary_key": "id",
    "fields": [
      {"name": "id", "type": "integer"},
      {"name": "name", "type": "string"},
      ...
    ]
  }
]
```

### Test 3: Execute Query
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "model": "users",
    "pagination": {"limit": 5, "offset": 0}
  }'
```

**Response:**
```json
{
  "sql": "SELECT * FROM users t0 LIMIT $1 OFFSET $2;",
  "params": [5, 0],
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-15"
    }
  ]
}
```

---

## 🎨 UI Features

### Dark Theme
- **Primary**: Deep gray (#111827)
- **Accent**: Cyan (#06b6d4) and Purple (#9333ea)
- **Contrast**: AAA accessibility standard

### Interactive Elements
- **Hover Effects**: Smooth transitions on all buttons
- **Active States**: Cyan highlight for active filters
- **Animations**: Smooth slide-in for detail panels
- **Responsive**: Adapts to different screen sizes

### Status Indicators
- **Loading**: Spinning gear animation
- **Errors**: Red alert box with icon
- **Filters Applied**: Badge shows filter count
- **Group Count**: Shows items per group

---

## 🐛 Debugging Tips

### Check Backend Logs
```bash
# Look for startup messages
# Check for database connection status
# Watch for query execution errors
```

### Check Frontend Console (F12)
```javascript
// Should see:
"Data from backend: [...]"
"Generated SQL: SELECT..."
"Parameters: [...]"

// Should NOT see:
"CORS error"
"Failed to fetch"
"Cannot read property"
```

### Test Specific Endpoints
```bash
# Get specific model
curl http://localhost:8080/models | jq '.[] | select(.name=="users")'

# Test with filter
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "model": "users",
    "filters": {"field": "id", "op": ">", "value": 5}
  }'
```

---

## 📱 Keyboard Shortcuts (Frontend)

| Key | Action |
|-----|--------|
| `Esc` | Close modal/detail panel |
| `Enter` | Apply filter value |
| `Tab` | Move between form fields |

---

## 🔌 Environment Variables

### Optional Configuration

**DATABASE_URL** - Database connection string
```bash
# Supabase example
DATABASE_URL="postgresql://user:password@db.xxxxx.supabase.co:5432/postgres"

# Local PostgreSQL example
DATABASE_URL="postgresql://postgres:password@localhost:5432/mydb"

# If not set: system runs in SQL-generation-only mode
```

**CONFIG_PATH** - Custom models configuration (default: configs/models.json)
```bash
CONFIG_PATH="/path/to/custom/models.json" ./server
```

**API_BASE** - Frontend API URL (default: http://localhost:8080)
```javascript
// Set in browser before loading
window.REACT_APP_API_URL = 'http://api.example.com'
```

---

## ✅ Validation Checklist

Before going live, verify:

- [x] Backend builds without errors
- [x] Backend starts and shows "Server starting on :8080"
- [x] `/health` endpoint responds with {"status":"ok"}
- [x] `/models` endpoint returns list of models
- [x] Frontend builds without errors
- [x] Frontend loads on http://localhost:5173
- [x] Models appear in left sidebar
- [x] Clicking model shows data in table
- [x] Filter button works
- [x] Group By button works
- [x] Clicking row shows detail panel
- [x] Console shows SQL generation logs
- [x] No CORS or network errors

---

## 🆘 Troubleshooting

### Port Already in Use
```bash
# Find process on port 8080
lsof -i :8080

# Kill it
kill -9 <PID>

# Start fresh
./server
```

### Frontend Won't Connect to Backend
```bash
# 1. Check backend is running
curl http://localhost:8080/health

# 2. Check frontend API URL (open DevTools)
console.log(API_BASE)

# 3. Check network requests (DevTools → Network tab)
```

### Database Connection Failed
```bash
# 1. Verify DATABASE_URL format
echo $DATABASE_URL

# 2. Test connection with psql
psql $DATABASE_URL -c "SELECT 1"

# 3. Check credentials
# 4. Verify network connectivity
```

### Data Not Showing
```bash
# 1. Check DevTools console for errors
# 2. Check backend API response
curl -X POST http://localhost:8080/query -H "Content-Type: application/json" -d '{"model":"users"}'

# 3. Check if using mock data (normal without DATABASE_URL)
```

---

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| `docs/backend_progress.md` | Backend architecture & phases |
| `docs/frontend_progress.md` | Frontend components & features |
| `docs/INTEGRATION_COMPLETE.md` | Full integration details |
| `docs/query_dsl_spec.md` | DSL query specification |
| `docs/postgres_sql_generation.md` | SQL generation strategy |
| `QUICK_START.md` | This file! |

---

## 🎯 Common Tasks

### Change Models Configuration
Edit `configs/models.json` and restart backend:
```bash
# Edit the file
vim configs/models.json

# Rebuild and restart
go build -o server ./cmd/server
./server
```

### Connect Different Database
```bash
# Set DATABASE_URL
export DATABASE_URL="postgresql://new_user:new_pass@new_host:5432/new_db"

# Restart backend
./server
```

### Deploy to Production
```bash
# 1. Build optimized backend
go build -ldflags="-s -w" -o server ./cmd/server

# 2. Build optimized frontend
cd frontend && npm run build

# 3. Copy to server
scp server user@host:/app/
scp -r frontend/dist user@host:/app/public/

# 4. Set environment and start
export DATABASE_URL="..."
./server
```

---

## 📞 Support

### Need Help?
1. Check **DevTools Console** for errors (F12)
2. Check **Backend Logs** for failures
3. Verify **API responses** with curl
4. Review **Documentation Files** for details
5. Check **GitHub Issues** if applicable

---

## 🎊 You're Ready!

Congratulations! Your Universal Data Viewer is up and running. 

**Start exploring your data now!** 🚀

---

**Version**: 1.0.0  
**Last Updated**: January 26, 2026  
**Status**: ✅ Production Ready
