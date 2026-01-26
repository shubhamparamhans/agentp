package schema_processor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// FieldType represents the JSON type for a field
type FieldType string

const (
	TypeInteger   FieldType = "integer"
	TypeString    FieldType = "string"
	TypeDecimal   FieldType = "decimal"
	TypeBoolean   FieldType = "boolean"
	TypeTimestamp FieldType = "timestamp"
	TypeJSON      FieldType = "json"
	TypeUUID      FieldType = "uuid"
	TypeBinary    FieldType = "binary"
)

// Field represents a table column in the JSON config
type Field struct {
	Name     string `json:"name"`
	Type     FieldType `json:"type"`
	Nullable bool   `json:"nullable"`
}

// Model represents a database table in the JSON config
type Model struct {
	Name       string  `json:"name"`
	Table      string  `json:"table"`
	PrimaryKey string  `json:"primaryKey"`
	Fields     []Field `json:"fields"`
}

// ModelConfig represents the complete models.json structure
type ModelConfig struct {
	Models []Model `json:"models"`
}

// ColumnInfo holds PostgreSQL column metadata
type ColumnInfo struct {
	ColumnName    string
	DataType      string
	IsNullable    bool
	ColumnDefault *string
	OrdinalPos    int
}

// TableInfo holds PostgreSQL table metadata
type TableInfo struct {
	TableName string
	Columns   []ColumnInfo
	PrimaryKey string
}

// SchemaProcessor handles database schema introspection
type SchemaProcessor struct {
	db *sql.DB
}

// NewSchemaProcessor creates a new schema processor
func NewSchemaProcessor(db *sql.DB) *SchemaProcessor {
	return &SchemaProcessor{db: db}
}

// mapPostgreSQLTypeToJSON maps PostgreSQL data types to JSON model types
func mapPostgreSQLTypeToJSON(pgType string) FieldType {
	pgType = strings.ToLower(pgType)
	pgType = strings.TrimSpace(pgType)

	// Handle array types (e.g., "integer[]" → "integer")
	pgType = strings.TrimSuffix(pgType, "[]")

	// Handle types with parameters (e.g., "character varying" → "character varying")
	basePGType := strings.Split(pgType, "(")[0]
	basePGType = strings.TrimSpace(basePGType)

	switch basePGType {
	case "integer", "int", "int4", "smallint", "int2", "bigint", "int8", "serial", "serial4", "bigserial", "serial8":
		return TypeInteger
	case "text", "character varying", "varchar", "character", "char", "name":
		return TypeString
	case "numeric", "decimal", "money", "double precision", "float8", "real", "float4":
		return TypeDecimal
	case "boolean", "bool":
		return TypeBoolean
	case "timestamp", "timestamp without time zone", "timestamp with time zone", "timestamptz", "timestamp at time zone", "date", "time", "time without time zone", "time with time zone", "timetz":
		return TypeTimestamp
	case "json", "jsonb":
		return TypeJSON
	case "uuid":
		return TypeUUID
	case "bytea", "bit", "bit varying", "varbit":
		return TypeBinary
	default:
		// Default to string for unknown types
		log.Printf("Warning: Unknown PostgreSQL type '%s', defaulting to string", pgType)
		return TypeString
	}
}

// GetTableColumns fetches column information for a table
func (sp *SchemaProcessor) GetTableColumns(tableName string) ([]ColumnInfo, error) {
	query := `
		SELECT
			column_name,
			data_type,
			is_nullable,
			column_default,
			ordinal_position
		FROM
			information_schema.columns
		WHERE
			table_name = $1 AND table_schema = 'public'
		ORDER BY
			ordinal_position ASC
	`

	rows, err := sp.db.Query(query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns: %w", err)
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		var isNullable string
		err := rows.Scan(&col.ColumnName, &col.DataType, &isNullable, &col.ColumnDefault, &col.OrdinalPos)
		if err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}
		col.IsNullable = isNullable == "YES"
		columns = append(columns, col)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating columns: %w", err)
	}

	return columns, nil
}

// GetPrimaryKey fetches the primary key for a table
func (sp *SchemaProcessor) GetPrimaryKey(tableName string) (string, error) {
	query := `
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		JOIN pg_class t ON t.oid = i.indrelid
		WHERE t.relname = $1 AND i.indisprimary
		LIMIT 1
	`

	var pkName string
	err := sp.db.QueryRow(query, tableName).Scan(&pkName)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("failed to query primary key: %w", err)
	}
	if err == sql.ErrNoRows {
		// Fallback: assume "id" if no primary key found
		return "id", nil
	}

	return pkName, nil
}

// GetAllTables fetches all table names from the database
func (sp *SchemaProcessor) GetAllTables() ([]string, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
		ORDER BY table_name ASC
	`

	rows, err := sp.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}

	return tables, nil
}

// GenerateModels creates Model objects from database schema
func (sp *SchemaProcessor) GenerateModels(tableNames []string) ([]Model, error) {
	var models []Model

	for _, tableName := range tableNames {
		// Get columns
		columns, err := sp.GetTableColumns(tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get columns for table %s: %w", tableName, err)
		}

		if len(columns) == 0 {
			log.Printf("Warning: Table '%s' has no columns, skipping", tableName)
			continue
		}

		// Get primary key
		pkName, err := sp.GetPrimaryKey(tableName)
		if err != nil {
			return nil, fmt.Errorf("failed to get primary key for table %s: %w", tableName, err)
		}

		// Convert columns to fields
		var fields []Field
		for _, col := range columns {
			field := Field{
				Name:     col.ColumnName,
				Type:     mapPostgreSQLTypeToJSON(col.DataType),
				Nullable: col.IsNullable,
			}
			fields = append(fields, field)
		}

		// Create model
		model := Model{
			Name:       tableName,
			Table:      tableName,
			PrimaryKey: pkName,
			Fields:     fields,
		}

		models = append(models, model)
	}

	return models, nil
}

// GenerateAndSaveModels generates models from database and saves to file
func (sp *SchemaProcessor) GenerateAndSaveModels(outputPath string, tableNames []string) error {
	var tables []string
	var err error

	// If no tables specified, get all tables
	if len(tableNames) == 0 {
		tables, err = sp.GetAllTables()
		if err != nil {
			return fmt.Errorf("failed to get all tables: %w", err)
		}

		if len(tables) == 0 {
			return fmt.Errorf("no tables found in database")
		}

		log.Printf("Found %d tables in database", len(tables))
	} else {
		tables = tableNames
	}

	// Generate models
	models, err := sp.GenerateModels(tables)
	if err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	if len(models) == 0 {
		return fmt.Errorf("no valid models generated from database")
	}

	// Create config
	config := ModelConfig{
		Models: models,
	}

	// Marshal to JSON with pretty printing
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal models to JSON: %w", err)
	}

	// Write to file
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write models file: %w", err)
	}

	log.Printf("Successfully generated models.json with %d models at %s", len(models), outputPath)
	fmt.Printf("Generated models for: %v\n", tables)

	return nil
}
