package adapter

import "udv/internal/planner"

// Database represents a generic database connection abstraction
type Database interface {
	// Connection management
	Close() error
	Ping() error

	// Query execution
	ExecuteQuery(query interface{}, args ...interface{}) ([]map[string]interface{}, error)
	Exec(query interface{}, args ...interface{}) (ExecResult, error)
}

// ExecResult wraps the result of an exec operation
type ExecResult interface {
	RowsAffected() (int64, error)
}

// QueryBuilder converts a QueryPlan into a database-specific query
type QueryBuilder interface {
	BuildQuery(plan *planner.QueryPlan) (query interface{}, args []interface{}, err error)
}
