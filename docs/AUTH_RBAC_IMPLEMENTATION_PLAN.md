# Authentication & RBAC Implementation Plan

## Overview

This document provides a detailed analysis and implementation plan for adding SQLite-based authentication and Role-Based Access Control (RBAC) to Agent P (Universal Data Viewer). The implementation will support config-driven initialization of admin users and role definitions, enabling fine-grained access control over models and operations.

---

## Current Architecture

### Authentication State

**Frontend:**
- ✅ Basic login UI component (`frontend/src/components/Login/Login.tsx`)
- ✅ Auth context with localStorage-based state (`frontend/src/contexts/AuthContext.tsx`)
- ✅ Protected route wrapper (`frontend/src/components/ProtectedRoute/ProtectedRoute.tsx`)
- ❌ No backend authentication endpoints
- ❌ No real credential validation
- ❌ No session management

**Backend:**
- ❌ No authentication middleware
- ❌ No user management
- ❌ No authorization checks
- ❌ All endpoints are publicly accessible

### Current API Structure

```
internal/api/
├── api.go          # Main API handlers
│   ├── /models     # GET - List all models
│   └── /query      # POST - Execute queries
└── server.go       # Server setup
```

### Current Configuration

- JSON-based model configuration (`configs/models.json`)
- No user/role configuration
- No permission definitions

---

## Key Requirements

### 1. Authentication
- SQLite database for user storage
- Password hashing (bcrypt)
- JWT-based session tokens
- Config file for initial admin users
- Console command to initialize admin users

### 2. Authorization (RBAC)
- Role-based access control
- Model-level permissions (read, write, view)
- Operation-level permissions (select, create, update, delete)
- Config file for default roles
- User-role assignment

### 3. Permission Model

**Access Types:**
- **read**: Can execute SELECT queries
- **write**: Can execute CREATE, UPDATE, DELETE operations
- **view**: Can access model metadata (fields, structure)

**Scope:**
- Model-level: Permissions apply to specific models
- Global: Permissions apply to all models

---

## Implementation Requirements

### Phase 1: SQLite Auth Database Setup (8-12 hours)

**Goal:** Create SQLite database schema and adapter for authentication and authorization.

#### 1.1 Database Schema Design (2-3 hours)

**File:** `internal/auth/schema.sql` (new)

```sql
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN DEFAULT 1,
    is_admin BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- User-Role assignments
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Permissions table
CREATE TABLE IF NOT EXISTS permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    model_name TEXT,  -- NULL means global permission
    access_type TEXT NOT NULL,  -- 'read', 'write', 'view'
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

-- Sessions table (optional, for token blacklisting)
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token_hash TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_permissions_role_model ON permissions(role_id, model_name);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);
```

#### 1.2 SQLite Adapter (4-6 hours)

**File:** `internal/auth/sqlite.go` (new)

**Dependencies:**
- `github.com/mattn/go-sqlite3` (SQLite driver)
- `golang.org/x/crypto/bcrypt` (password hashing)

**Implementation:**
```go
package auth

import (
    "database/sql"
    "time"
    _ "github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/bcrypt"
)

type SQLiteAuth struct {
    db *sql.DB
}

func NewSQLiteAuth(dbPath string) (*SQLiteAuth, error) {
    db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
    if err != nil {
        return nil, err
    }
    
    // Initialize schema
    if err := initSchema(db); err != nil {
        return nil, err
    }
    
    return &SQLiteAuth{db: db}, nil
}

func (a *SQLiteAuth) Close() error {
    return a.db.Close()
}

// User management
func (a *SQLiteAuth) CreateUser(username, email, password string, isAdmin bool) (*User, error)
func (a *SQLiteAuth) GetUserByUsername(username string) (*User, error)
func (a *SQLiteAuth) ValidatePassword(username, password string) (*User, error)
func (a *SQLiteAuth) UpdateUserPassword(userID int, newPassword string) error

// Role management
func (a *SQLiteAuth) CreateRole(name, description string) (*Role, error)
func (a *SQLiteAuth) GetRoleByName(name string) (*Role, error)
func (a *SQLiteAuth) AssignRoleToUser(userID, roleID int) error
func (a *SQLiteAuth) GetUserRoles(userID int) ([]*Role, error)

// Permission management
func (a *SQLiteAuth) GrantPermission(roleID int, modelName string, accessType string) error
func (a *SQLiteAuth) CheckPermission(userID int, modelName, accessType string) (bool, error)
func (a *SQLiteAuth) GetRolePermissions(roleID int) ([]*Permission, error)
```

#### 1.3 Data Models (2-3 hours)

**File:** `internal/auth/models.go` (new)

```go
package auth

import "time"

type User struct {
    ID           int
    Username     string
    Email        string
    PasswordHash string
    IsActive     bool
    IsAdmin      bool
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type Role struct {
    ID          int
    Name        string
    Description string
    CreatedAt   time.Time
}

type Permission struct {
    ID         int
    RoleID     int
    ModelName  *string  // nil = global permission
    AccessType string   // "read", "write", "view"
}

type Session struct {
    ID        int
    UserID    int
    TokenHash string
    ExpiresAt time.Time
    CreatedAt time.Time
}
```

---

### Phase 2: Configuration System (6-8 hours)

**Goal:** Support config file-based initialization of users and roles.

#### 2.1 Auth Configuration Schema (2-3 hours)

**File:** `configs/auth.json` (new)

```json
{
  "database_path": "./data/auth.db",
  "initial_admin_users": [
    {
      "username": "admin",
      "email": "admin@example.com",
      "password": "changeme123",
      "is_admin": true
    }
  ],
  "default_roles": [
    {
      "name": "admin",
      "description": "Full access to all models",
      "permissions": [
        {
          "model_name": null,
          "access_type": "read"
        },
        {
          "model_name": null,
          "access_type": "write"
        },
        {
          "model_name": null,
          "access_type": "view"
        }
      ]
    },
    {
      "name": "viewer",
      "description": "Read-only access to all models",
      "permissions": [
        {
          "model_name": null,
          "access_type": "read"
        },
        {
          "model_name": null,
          "access_type": "view"
        }
      ]
    },
    {
      "name": "editor",
      "description": "Read and write access to all models",
      "permissions": [
        {
          "model_name": null,
          "access_type": "read"
        },
        {
          "model_name": null,
          "access_type": "write"
        },
        {
          "model_name": null,
          "access_type": "view"
        }
      ]
    },
    {
      "name": "orders_manager",
      "description": "Full access to orders model only",
      "permissions": [
        {
          "model_name": "orders",
          "access_type": "read"
        },
        {
          "model_name": "orders",
          "access_type": "write"
        },
        {
          "model_name": "orders",
          "access_type": "view"
        }
      ]
    }
  ]
}
```

#### 2.2 Config Loader (2-3 hours)

**File:** `internal/auth/config.go` (new)

```go
package auth

type AuthConfig struct {
    DatabasePath      string           `json:"database_path"`
    InitialAdminUsers []AdminUser      `json:"initial_admin_users"`
    DefaultRoles      []RoleDefinition `json:"default_roles"`
}

type AdminUser struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    IsAdmin  bool   `json:"is_admin"`
}

type RoleDefinition struct {
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Permissions []PermissionDef `json:"permissions"`
}

type PermissionDef struct {
    ModelName  *string `json:"model_name"`  // null = global
    AccessType string  `json:"access_type"`  // "read", "write", "view"
}

func LoadAuthConfig(filePath string) (*AuthConfig, error)
func (c *AuthConfig) InitializeDatabase(auth *SQLiteAuth) error
```

#### 2.3 CLI Initialization Command (2-2 hours)

**File:** `cmd/init-auth/main.go` (new)

```go
package main

import (
    "flag"
    "fmt"
    "os"
    "udv/internal/auth"
)

func main() {
    configPath := flag.String("config", "configs/auth.json", "Path to auth config file")
    dbPath := flag.String("db", "./data/auth.db", "Path to SQLite database")
    flag.Parse()
    
    // Load config
    cfg, err := auth.LoadAuthConfig(*configPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
        os.Exit(1)
    }
    
    // Initialize database
    authDB, err := auth.NewSQLiteAuth(*dbPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
        os.Exit(1)
    }
    defer authDB.Close()
    
    // Initialize from config
    if err := cfg.InitializeDatabase(authDB); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Authentication database initialized successfully!")
    fmt.Printf("Created %d admin user(s)\n", len(cfg.InitialAdminUsers))
    fmt.Printf("Created %d role(s)\n", len(cfg.DefaultRoles))
}
```

**Usage:**
```bash
go run cmd/init-auth/main.go -config configs/auth.json -db ./data/auth.db
```

---

### Phase 3: JWT Authentication (8-10 hours)

**Goal:** Implement JWT-based session management.

#### 3.1 JWT Token Management (4-5 hours)

**File:** `internal/auth/jwt.go` (new)

**Dependencies:**
- `github.com/golang-jwt/jwt/v5`

**Implementation:**
```go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secretKey     []byte
    tokenDuration time.Duration
}

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    IsAdmin  bool   `json:"is_admin"`
    Roles    []string `json:"roles"`
    jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager
func (m *JWTManager) GenerateToken(user *User, roles []*Role) (string, error)
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error)
```

#### 3.2 Authentication Middleware (4-5 hours)

**File:** `internal/auth/middleware.go` (new)

```go
package auth

import (
    "net/http"
    "strings"
)

func AuthMiddleware(authDB *SQLiteAuth, jwtMgr *JWTManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing authorization header", http.StatusUnauthorized)
                return
            }
            
            // Parse Bearer token
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "invalid authorization header", http.StatusUnauthorized)
                return
            }
            
            token := parts[1]
            claims, err := jwtMgr.ValidateToken(token)
            if err != nil {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }
            
            // Attach user info to request context
            ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
            ctx = context.WithValue(ctx, "username", claims.Username)
            ctx = context.WithValue(ctx, "is_admin", claims.IsAdmin)
            ctx = context.WithValue(ctx, "roles", claims.Roles)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

---

### Phase 4: Authorization Middleware (10-12 hours)

**Goal:** Implement RBAC-based authorization checks.

#### 4.1 Permission Checker (4-6 hours)

**File:** `internal/auth/authorizer.go` (new)

```go
package auth

import (
    "context"
    "net/http"
)

type Authorizer struct {
    authDB *SQLiteAuth
}

func NewAuthorizer(authDB *SQLiteAuth) *Authorizer {
    return &Authorizer{authDB: authDB}
}

// CheckPermission verifies if user has required permission
func (a *Authorizer) CheckPermission(userID int, modelName, accessType string) (bool, error) {
    // Get user roles
    roles, err := a.authDB.GetUserRoles(userID)
    if err != nil {
        return false, err
    }
    
    // Check if user is admin (admins have all permissions)
    user, err := a.authDB.GetUserByID(userID)
    if err != nil {
        return false, err
    }
    if user.IsAdmin {
        return true, nil
    }
    
    // Check role permissions
    for _, role := range roles {
        // Check global permissions first
        hasGlobal, err := a.authDB.HasPermission(role.ID, nil, accessType)
        if err != nil {
            continue
        }
        if hasGlobal {
            return true, nil
        }
        
        // Check model-specific permissions
        hasModel, err := a.authDB.HasPermission(role.ID, &modelName, accessType)
        if err != nil {
            continue
        }
        if hasModel {
            return true, nil
        }
    }
    
    return false, nil
}

// RequirePermission middleware
func (a *Authorizer) RequirePermission(modelName, accessType string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID, ok := r.Context().Value("user_id").(int)
            if !ok {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            
            allowed, err := a.CheckPermission(userID, modelName, accessType)
            if err != nil {
                http.Error(w, "authorization error", http.StatusInternalServerError)
                return
            }
            
            if !allowed {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

#### 4.2 Operation-to-Permission Mapping (2-3 hours)

**File:** `internal/auth/permissions.go` (new)

```go
package auth

// Map DSL operations to access types
func OperationToAccessType(operation string) string {
    switch operation {
    case "select":
        return "read"
    case "create", "update", "delete":
        return "write"
    default:
        return "read"  // Default to read for safety
    }
}

// CheckModelAccess checks if user can access model metadata
func CheckModelAccess(accessType string) bool {
    return accessType == "view" || accessType == "read" || accessType == "write"
}
```

#### 4.3 API Integration (4-3 hours)

**File:** `internal/api/api.go` (modify)

**Changes:**
- Add authorization checks to `/models` endpoint (requires "view" permission)
- Add authorization checks to `/query` endpoint (requires "read" or "write" based on operation)
- Extract model name from query and check permissions

```go
// Modified handleModels
func (a *API) handleModels(w http.ResponseWriter, r *http.Request) {
    // Get user from context
    userID := r.Context().Value("user_id").(int)
    
    // Check view permission for each model
    models := a.registry.ListModels()
    var allowedModels []modelResp
    
    for _, m := range models {
        // Check if user has view permission
        allowed, _ := a.authorizer.CheckPermission(userID, m, "view")
        if allowed {
            // Add model to response
        }
    }
    
    // Return filtered models
}

// Modified handleQuery
func (a *API) handleQuery(w http.ResponseWriter, r *http.Request) {
    // Get user from context
    userID := r.Context().Value("user_id").(int)
    
    // Parse query
    // ...
    
    // Determine required permission
    accessType := auth.OperationToAccessType(operation)
    
    // Check permission
    allowed, _ := a.authorizer.CheckPermission(userID, q.Model, accessType)
    if !allowed {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }
    
    // Process query
    // ...
}
```

---

### Phase 5: Authentication Endpoints (6-8 hours)

**Goal:** Add login, logout, and user management endpoints.

#### 5.1 Auth API Handlers (4-5 hours)

**File:** `internal/api/auth.go` (new)

```go
package api

import (
    "encoding/json"
    "net/http"
    "time"
    "udv/internal/auth"
)

type AuthAPI struct {
    authDB *auth.SQLiteAuth
    jwtMgr *auth.JWTManager
}

func NewAuthAPI(authDB *auth.SQLiteAuth, jwtMgr *auth.JWTManager) *AuthAPI {
    return &AuthAPI{
        authDB: authDB,
        jwtMgr: jwtMgr,
    }
}

// POST /auth/login
func (a *AuthAPI) handleLogin(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    
    // Validate credentials
    user, err := a.authDB.ValidatePassword(req.Username, req.Password)
    if err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }
    
    // Get user roles
    roles, _ := a.authDB.GetUserRoles(user.ID)
    
    // Generate token
    token, err := a.jwtMgr.GenerateToken(user, roles)
    if err != nil {
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }
    
    // Return token and user info
    json.NewEncoder(w).Encode(map[string]interface{}{
        "token": token,
        "user": map[string]interface{}{
            "id":       user.ID,
            "username": user.Username,
            "email":    user.Email,
            "is_admin": user.IsAdmin,
        },
    })
}

// POST /auth/logout
func (a *AuthAPI) handleLogout(w http.ResponseWriter, r *http.Request) {
    // Optional: Blacklist token
    // For now, just return success
    w.WriteHeader(http.StatusOK)
}

// GET /auth/me
func (a *AuthAPI) handleMe(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("user_id").(int)
    user, _ := a.authDB.GetUserByID(userID)
    roles, _ := a.authDB.GetUserRoles(userID)
    
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user": user,
        "roles": roles,
    })
}
```

#### 5.2 Route Registration (2-3 hours)

**File:** `internal/api/api.go` (modify)

```go
func (a *API) RegisterRoutes(mux *http.ServeMux) {
    // Public routes
    if a.authAPI != nil {
        mux.HandleFunc("/auth/login", a.authAPI.handleLogin)
        mux.HandleFunc("/auth/logout", a.authAPI.handleLogout)
    }
    
    // Protected routes
    var protectedMux http.Handler = mux
    if a.authMiddleware != nil {
        protectedMux = a.authMiddleware(protectedMux)
    }
    
    if a.authAPI != nil {
        mux.HandleFunc("/auth/me", a.authAPI.handleMe)
    }
    
    mux.HandleFunc("/models", a.handleModels)
    mux.HandleFunc("/query", a.handleQuery)
}
```

---

### Phase 6: Frontend Integration (8-10 hours)

**Goal:** Update frontend to use real authentication.

#### 6.1 API Client Updates (3-4 hours)

**File:** `frontend/src/api/client.ts` (modify)

```typescript
// Add token management
let authToken: string | null = localStorage.getItem('auth_token')

export function setAuthToken(token: string) {
  authToken = token
  localStorage.setItem('auth_token', token)
}

export function clearAuthToken() {
  authToken = null
  localStorage.removeItem('auth_token')
}

// Update fetch calls to include Authorization header
function fetchWithAuth(url: string, options: RequestInit = {}) {
  const headers = new Headers(options.headers)
  if (authToken) {
    headers.set('Authorization', `Bearer ${authToken}`)
  }
  return fetch(url, { ...options, headers })
}

// Add login function
export async function login(username: string, password: string) {
  const response = await fetchWithAuth('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
  if (!response.ok) throw new Error('Login failed')
  const data = await response.json()
  setAuthToken(data.token)
  return data
}
```

#### 6.2 Auth Context Updates (2-3 hours)

**File:** `frontend/src/contexts/AuthContext.tsx` (modify)

```typescript
import { login as apiLogin, logout as apiLogout } from '../api/client'

const login = async (username: string, password: string): Promise<boolean> => {
  try {
    const data = await apiLogin(username, password)
    setIsAuthenticated(true)
    setUser(data.user.username)
    return true
  } catch (error) {
    return false
  }
}

const logout = () => {
  apiLogout()
  setIsAuthenticated(false)
  setUser(null)
}
```

#### 6.3 Protected Route Updates (1-2 hours)

**File:** `frontend/src/components/ProtectedRoute/ProtectedRoute.tsx` (modify)

- Add token refresh logic
- Handle 401 responses
- Redirect to login on auth failure

#### 6.4 Error Handling (2-1 hours)

- Handle 401/403 responses globally
- Show permission denied messages
- Redirect to login on token expiry

---

### Phase 7: Testing & Edge Cases (12-16 hours)

#### 7.1 Unit Tests (6-8 hours)
- Test SQLite adapter functions
- Test JWT token generation/validation
- Test permission checking logic
- Test config loading and initialization

#### 7.2 Integration Tests (4-6 hours)
- Test login flow end-to-end
- Test authorization on API endpoints
- Test role assignment and permissions
- Test admin user initialization

#### 7.3 Edge Cases (2-2 hours)
- Handle expired tokens
- Handle invalid tokens
- Handle missing permissions
- Handle concurrent user creation
- Handle role deletion with assigned users

---

## Key Challenges

### 1. Permission Resolution

**Challenge:** Determining if a user has permission when both global and model-specific permissions exist.

**Solution:**
- Check global permissions first (more permissive)
- Fall back to model-specific permissions
- Admin users bypass all checks

### 2. Operation-to-Permission Mapping

**Challenge:** Mapping DSL operations (select, create, update, delete) to access types (read, write, view).

**Solution:**
- `select` → `read`
- `create`, `update`, `delete` → `write`
- Model metadata access → `view`

### 3. Token Management

**Challenge:** Handling token expiration and refresh.

**Solution:**
- Use JWT with expiration claims
- Optional: Implement refresh tokens
- Store token in localStorage (frontend)
- Validate token on each request

### 4. Config-Driven Initialization

**Challenge:** Ensuring idempotent initialization (can run multiple times safely).

**Solution:**
- Check if users/roles exist before creating
- Use UNIQUE constraints in database
- Handle conflicts gracefully

### 5. Backward Compatibility

**Challenge:** Making authentication optional during migration.

**Solution:**
- Make auth middleware optional (check if authDB is nil)
- Allow running without authentication for development
- Environment variable to enable/disable auth

---

## Implementation Effort Summary

| Phase | Task | Hours | Priority |
|-------|------|-------|----------|
| **Phase 1** | SQLite Auth Database Setup | 8-12 | 🔴 HIGH |
| **Phase 2** | Configuration System | 6-8 | 🔴 HIGH |
| **Phase 3** | JWT Authentication | 8-10 | 🔴 HIGH |
| **Phase 4** | Authorization Middleware | 10-12 | 🔴 HIGH |
| **Phase 5** | Authentication Endpoints | 6-8 | 🟡 MEDIUM |
| **Phase 6** | Frontend Integration | 8-10 | 🟡 MEDIUM |
| **Phase 7** | Testing & Edge Cases | 12-16 | 🔴 HIGH |
| **TOTAL** | | **58-76 hours** | |

**Estimated Timeline:** 1.5-2 weeks for a single developer

---

## Dependencies

### New Go Packages Required

```go
github.com/mattn/go-sqlite3          // SQLite driver
golang.org/x/crypto/bcrypt          // Password hashing
github.com/golang-jwt/jwt/v5        // JWT tokens
```

**Add to `go.mod`:**
```bash
go get github.com/mattn/go-sqlite3
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v5
```

---

## File Structure After Implementation

```
internal/
├── auth/                          # NEW: Authentication module
│   ├── sqlite.go                  # SQLite adapter
│   ├── models.go                  # Data models
│   ├── jwt.go                     # JWT management
│   ├── middleware.go              # Auth middleware
│   ├── authorizer.go              # Authorization logic
│   ├── permissions.go             # Permission utilities
│   ├── config.go                  # Config loading
│   └── schema.sql                 # Database schema
├── api/
│   ├── api.go                     # MODIFY: Add auth checks
│   └── auth.go                    # NEW: Auth endpoints
└── ...

cmd/
├── server/
│   └── main.go                    # MODIFY: Initialize auth
└── init-auth/                     # NEW: Auth initialization CLI
    └── main.go

configs/
├── models.json                    # Existing
└── auth.json                      # NEW: Auth configuration
```

---

## Configuration Example

### Auth Configuration File

**File:** `configs/auth.json`

```json
{
  "database_path": "./data/auth.db",
  "initial_admin_users": [
    {
      "username": "admin",
      "email": "admin@example.com",
      "password": "changeme123",
      "is_admin": true
    }
  ],
  "default_roles": [
    {
      "name": "admin",
      "description": "Full access to all models",
      "permissions": [
        { "model_name": null, "access_type": "read" },
        { "model_name": null, "access_type": "write" },
        { "model_name": null, "access_type": "view" }
      ]
    },
    {
      "name": "viewer",
      "description": "Read-only access",
      "permissions": [
        { "model_name": null, "access_type": "read" },
        { "model_name": null, "access_type": "view" }
      ]
    }
  ]
}
```

### Environment Variables

```bash
# Enable authentication (optional, defaults to false)
ENABLE_AUTH=true

# JWT secret key (required if auth enabled)
JWT_SECRET=your-secret-key-here

# JWT token duration (optional, defaults to 24h)
JWT_TOKEN_DURATION=24h

# Auth database path (optional, defaults to ./data/auth.db)
AUTH_DB_PATH=./data/auth.db

# Auth config path (optional, defaults to configs/auth.json)
AUTH_CONFIG_PATH=configs/auth.json
```

---

## Testing Strategy

### 1. Unit Tests

```go
// Test SQLite adapter
func TestSQLiteAuth_CreateUser(t *testing.T)
func TestSQLiteAuth_ValidatePassword(t *testing.T)
func TestSQLiteAuth_CheckPermission(t *testing.T)

// Test JWT manager
func TestJWTManager_GenerateToken(t *testing.T)
func TestJWTManager_ValidateToken(t *testing.T)

// Test permission checking
func TestAuthorizer_CheckPermission(t *testing.T)
```

### 2. Integration Tests

- Test login flow with real database
- Test authorization on protected endpoints
- Test role assignment and permission inheritance
- Test config initialization

### 3. Manual Testing Checklist

- [ ] Initialize auth database from config
- [ ] Login with admin user
- [ ] Access models endpoint (should work)
- [ ] Execute query (should work)
- [ ] Create viewer user and assign role
- [ ] Login with viewer user
- [ ] Try to execute write operation (should fail)
- [ ] Try to execute read operation (should succeed)
- [ ] Test token expiration
- [ ] Test invalid token handling

---

## Migration Path

### Option 1: Gradual Rollout (Recommended)

1. **Phase 1:** Implement auth system but keep it optional
2. **Phase 2:** Enable auth in development environment
3. **Phase 3:** Initialize admin users via config
4. **Phase 4:** Enable auth in production
5. **Phase 5:** Create roles and assign users

### Option 2: Big Bang

- Implement all phases at once
- Deploy with authentication enabled from start
- Initialize all users and roles via config

---

## Security Considerations

### 1. Password Security

- Use bcrypt with appropriate cost factor (10-12)
- Never store plaintext passwords
- Enforce password complexity (optional, future enhancement)

### 2. Token Security

- Use strong JWT secret key
- Set appropriate token expiration
- Consider implementing refresh tokens
- Optional: Token blacklisting for logout

### 3. SQL Injection Prevention

- Use parameterized queries
- Validate all inputs
- Use SQLite driver's built-in protections

### 4. Authorization

- Always check permissions server-side
- Never trust client-side permission checks
- Default to deny if permission unclear

### 5. Admin Access

- Limit admin user creation
- Require strong passwords for admin accounts
- Consider 2FA for admin accounts (future enhancement)

---

## Limitations & Considerations

### Current Limitations

1. **No Password Reset:** Users cannot reset passwords (future enhancement)
2. **No User Management UI:** Users must be managed via config/CLI (future enhancement)
3. **No Role Management UI:** Roles must be managed via config (future enhancement)
4. **No Audit Logging:** No logging of who accessed what (future enhancement)
5. **No Session Management:** No way to revoke active sessions (future enhancement)

### Future Enhancements

1. **Password Reset Flow:** Email-based password reset
2. **User Management API:** CRUD operations for users
3. **Role Management API:** CRUD operations for roles
4. **Audit Logging:** Track all access and changes
5. **Session Management:** View and revoke active sessions
6. **2FA Support:** Two-factor authentication
7. **LDAP Integration:** Support LDAP authentication
8. **OAuth Integration:** Support OAuth providers

---

## Recommended Approach

### Step 1: Database & Config (Week 1, Days 1-2)
- Create SQLite schema
- Implement SQLite adapter
- Create config system
- Build initialization CLI

### Step 2: Authentication (Week 1, Days 3-5)
- Implement JWT management
- Create auth middleware
- Add login/logout endpoints
- Test authentication flow

### Step 3: Authorization (Week 2, Days 1-3)
- Implement permission checking
- Add authorization middleware
- Integrate with API endpoints
- Test authorization flow

### Step 4: Frontend & Polish (Week 2, Days 4-5)
- Update frontend to use real auth
- Handle errors and edge cases
- Write tests
- Documentation

---

## Success Criteria

✅ **Authentication & RBAC Complete When:**

1. ✅ Can initialize admin users from config file
2. ✅ Can login and receive JWT token
3. ✅ Protected endpoints require valid token
4. ✅ Can assign roles to users
5. ✅ Can define permissions per role
6. ✅ Model access is controlled by permissions
7. ✅ Operation access is controlled by permissions
8. ✅ Admin users have full access
9. ✅ Frontend integrates with auth system
10. ✅ Tests pass for all auth/authz flows

---

## Next Steps

1. **Review this plan** - Confirm approach and priorities
2. **Create SQLite schema** - Start with database structure
3. **Implement auth adapter** - Build SQLite authentication layer
4. **Add config system** - Support config-driven initialization
5. **Implement JWT auth** - Add token-based authentication
6. **Add authorization** - Implement RBAC checks
7. **Update API** - Integrate auth/authz into endpoints
8. **Update frontend** - Connect to real auth system
9. **Test thoroughly** - Ensure security and functionality
10. **Document usage** - Document how to use the system

---

**Last Updated:** Based on current codebase analysis  
**Estimated Effort:** 58-76 hours (~1.5-2 weeks)  
**Priority:** High (enables secure multi-user access)

