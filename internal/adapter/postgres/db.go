package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Database wraps a PostgreSQL connection pool
type Database struct {
	db *sql.DB
}

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

// Query executes a parameterized query and returns rows
func (d *Database) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	return d.db.Query(sql, args...)
}

// QueryRow executes a query that returns a single row
func (d *Database) QueryRow(sql string, args ...interface{}) *sql.Row {
	return d.db.QueryRow(sql, args...)
}

// Exec executes a query that doesn't return rows
func (d *Database) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return d.db.Exec(sql, args...)
}

// ExecuteAndFetchRows executes a query and returns results as []map[string]interface{}
func (d *Database) ExecuteAndFetchRows(sql string, args ...interface{}) ([]map[string]interface{}, error) {
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
