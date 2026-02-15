package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"udv/internal/schema_processor"

	_ "github.com/lib/pq"
)

func main() {
	// Define command-line flags
	databaseURL := flag.String("db", "", "PostgreSQL connection string (or use DATABASE_URL env var)")
	outputPath := flag.String("output", "configs/models.json", "Output path for generated models.json")
	tableNamesStr := flag.String("tables", "", "Comma-separated list of table names to process (default: all tables)")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	// Get database URL from flag or environment variable
	dbURL := *databaseURL
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	if dbURL == "" {
		fmt.Println("Error: Database URL not provided")
		fmt.Println("Usage: generate-models -db 'postgresql://...' [-output path/to/models.json] [-tables table1,table2]")
		fmt.Println("\nOr set DATABASE_URL environment variable")
		os.Exit(1)
	}

	// Connect to database
	log.Println("Connecting to database...")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✓ Connected to database")

	// Create schema processor
	processor := schema_processor.NewSchemaProcessor(db)

	// Parse table names if provided
	var tableNames []string
	if *tableNamesStr != "" {
		// Parse comma-separated table names
		// TODO: implement parsing if needed
		log.Println("Note: Specific table selection not yet implemented, processing all tables")
	}

	// Generate and save models
	log.Println("Introspecting database schema...")
	err = processor.GenerateAndSaveModels(*outputPath, tableNames)
	if err != nil {
		log.Fatalf("Failed to generate models: %v", err)
	}

	fmt.Printf("\n✓ Models generated successfully at: %s\n", *outputPath)
}

func printHelp() {
	fmt.Print(`
Universal Data Viewer - Schema Processor CLI

USAGE:
  generate-models [flags]

FLAGS:
  -db string
    	PostgreSQL connection string
    	If not provided, uses DATABASE_URL environment variable
    	Example: postgresql://user:password@localhost:5432/dbname

  -output string
    	Output path for generated models.json
    	Default: configs/models.json

  -tables string
    	Comma-separated list of table names to process
    	Default: all tables in database

  -help
    	Show this help message

EXAMPLES:
  # Using environment variable
  export DATABASE_URL="postgresql://postgres:password@localhost:5432/mydb"
  generate-models

  # Specify database URL directly
  generate-models -db "postgresql://user:pass@host:5432/db"

  # Custom output path
  generate-models -output /custom/path/models.json

ENVIRONMENT VARIABLES:
  DATABASE_URL
    	PostgreSQL connection string (alternative to -db flag)

DESCRIPTION:
  Introspects a PostgreSQL database and automatically generates a
  models.json configuration file that can be used with UDV.

  The processor:
  - Discovers all tables in the public schema
  - Maps PostgreSQL data types to UDV types
  - Detects primary keys
  - Identifies nullable columns
  - Generates properly formatted models.json

SUPPORTED PostgreSQL TYPES:
  - integer (int, int4, int8, serial, bigserial)
  - string (text, varchar, character)
  - decimal (numeric, money, float)
  - boolean
  - timestamp (with/without time zone)
  - uuid
  - json/jsonb
  - binary (bytea, bit)
`)
}
