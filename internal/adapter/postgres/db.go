package postgres

import (
	"database/sql"
	"fmt"

	"udv/internal/adapter"

	_ "github.com/lib/pq"
)

// Database wraps a PostgreSQL connection pool
type Database struct {
	db *sql.DB
}

// Compile-time assertion that Database implements adapter.Database interface
var _ adapter.Database = (*Database)(nil)

// Connect opens a connection to a PostgreSQL database using a DSN
func Connect(dsn string) (*Database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Database{db: db}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// Ping checks the connection to the database
func (d *Database) Ping() error {
	return d.db.Ping()
}

// Query executes a parameterized query and returns rows
func (d *Database) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(sql, args...)
}

// QueryRow executes a query that returns a single row
func (d *Database) QueryRow(sql string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(sql, args...)
}

// Exec executes a query that doesn't return rows
func (d *Database) Exec(query interface{}, args ...interface{}) (adapter.ExecResult, error) {
	sql, ok := query.(string)
	if !ok {
		return nil, fmt.Errorf("expected query to be string, got %T", query)
	}

	result, err := d.db.Exec(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("exec failed: %w", err)
	}

	return &PostgresExecResult{result: result}, nil
}

// ExecuteQuery executes a query and returns results as []map[string]interface{}
func (d *Database) ExecuteQuery(query interface{}, args ...interface{}) ([]map[string]interface{}, error) {
	sql, ok := query.(string)
	if !ok {
		return nil, fmt.Errorf("expected query to be string, got %T", query)
	}

	rows, err := d.db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Fetch all rows
	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// Scan the row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert to map
		entry := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			// Convert []byte to string for better JSON serialization
			if b, ok := val.([]byte); ok {
				entry[col] = string(b)
			} else {
				entry[col] = val
			}
		}
		results = append(results, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// ExecuteAndFetchRows is kept for backward compatibility
func (d *Database) ExecuteAndFetchRows(sql string, args ...interface{}) ([]map[string]interface{}, error) {
	return d.ExecuteQuery(sql, args...)
}

// PostgresExecResult wraps sql.Result to implement adapter.ExecResult
type PostgresExecResult struct {
	result sql.Result
}

// RowsAffected returns the number of rows affected by the operation
func (r *PostgresExecResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}
