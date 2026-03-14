package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"udv/internal/schema_processor"

	_ "github.com/lib/pq"
)

func main() {
	// Define command-line flags
	dbType := flag.String("type", "postgres", "Database type: postgres or mongodb")
	databaseURL := flag.String("db", "", "PostgreSQL connection string (or use DATABASE_URL env var)")
	mongodbURI := flag.String("mongodb-uri", "", "MongoDB connection URI (or use MONGODB_URI env var)")
	mongodbDB := flag.String("mongodb-db", "", "MongoDB database name (or use MONGODB_DATABASE env var)")
	outputPath := flag.String("output", "configs/models.json", "Output path for generated models.json")
	tableNamesStr := flag.String("tables", "", "Comma-separated list of table names to process (default: all tables)")
	collectionNamesStr := flag.String("collections", "", "Comma-separated list of collection names to process (MongoDB only)")
	sampleSize := flag.Int("sample-size", 100, "Number of documents to sample per collection (MongoDB only)")
	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		printHelp()
		os.Exit(0)
	}

	switch *dbType {
	case "mongodb":
		generateMongoDBModels(*mongodbURI, *mongodbDB, *collectionNamesStr, *sampleSize, *outputPath)
	case "postgres", "":
		generatePostgresModels(*databaseURL, *tableNamesStr, *outputPath)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unsupported database type: %s\n", *dbType)
		os.Exit(1)
	}
}

func generatePostgresModels(dbURL, tableNamesStr, outputPath string) {
	// Get database URL from flag or environment variable
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	if dbURL == "" {
		fmt.Println("Error: Database URL not provided")
		fmt.Println("Usage: generate-models -type postgres -db 'postgresql://...' [-output path/to/models.json]")
		fmt.Println("\nOr set DATABASE_URL environment variable")
		os.Exit(1)
	}

	// Connect to database
	log.Println("Connecting to PostgreSQL database...")
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
	log.Println("✓ Connected to PostgreSQL database")

	// Create schema processor
	processor := schema_processor.NewSchemaProcessor(db)

	// Parse table names if provided
	var tableNames []string
	if tableNamesStr != "" {
		tableNames = strings.Split(tableNamesStr, ",")
		for i := range tableNames {
			tableNames[i] = strings.TrimSpace(tableNames[i])
		}
		log.Printf("Processing tables: %v", tableNames)
	}

	// Generate and save models
	log.Println("Introspecting database schema...")
	err = processor.GenerateAndSaveModels(outputPath, tableNames)
	if err != nil {
		log.Fatalf("Failed to generate models: %v", err)
	}

	fmt.Printf("\n✓ Models generated successfully at: %s\n", outputPath)
}

func generateMongoDBModels(mongoURI, mongoDBName, collectionNamesStr string, sampleSize int, outputPath string) {
	// Get MongoDB URI from flag or environment
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGODB_URI")
	}
	if mongoURI == "" {
		fmt.Println("Error: MongoDB URI not provided")
		fmt.Println("Usage: generate-models -type mongodb -mongodb-uri 'mongodb://...' -mongodb-db dbname")
		fmt.Println("\nOr set MONGODB_URI and MONGODB_DATABASE environment variables")
		os.Exit(1)
	}

	// Get MongoDB database name from flag or environment
	if mongoDBName == "" {
		mongoDBName = os.Getenv("MONGODB_DATABASE")
	}
	if mongoDBName == "" {
		fmt.Println("Error: MongoDB database name not provided")
		fmt.Println("Usage: generate-models -type mongodb -mongodb-db dbname")
		fmt.Println("\nOr set MONGODB_DATABASE environment variable")
		os.Exit(1)
	}

	// Parse collection names if provided
	var collectionNames []string
	if collectionNamesStr != "" {
		collectionNames = strings.Split(collectionNamesStr, ",")
		for i := range collectionNames {
			collectionNames[i] = strings.TrimSpace(collectionNames[i])
		}
		log.Printf("Processing collections: %v", collectionNames)
	}

	log.Println("Connecting to MongoDB...")
	processor, err := schema_processor.NewMongoDBProcessor(mongoURI, mongoDBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer processor.Close()

	log.Printf("✓ Connected to MongoDB, sampling %d documents per collection\n", sampleSize)
	log.Println("Introspecting MongoDB schema...")

	err = processor.GenerateAndSaveModels(outputPath, collectionNames, sampleSize)
	if err != nil {
		log.Fatalf("Failed to generate models: %v", err)
	}

	fmt.Printf("\n✓ Models generated successfully at: %s\n", outputPath)
}

func printHelp() {
	fmt.Print(`
Universal Data Viewer - Schema Processor CLI

USAGE:
  generate-models [flags]

FLAGS:
  -type string
    	Database type: postgres or mongodb
    	Default: postgres

PostgreSQL FLAGS:
  -db string
    	PostgreSQL connection string
    	If not provided, uses DATABASE_URL environment variable
    	Example: postgresql://user:password@localhost:5432/dbname

MongoDB FLAGS:
  -mongodb-uri string
    	MongoDB connection URI
    	If not provided, uses MONGODB_URI environment variable
    	Example: mongodb://localhost:27017

  -mongodb-db string
    	MongoDB database name
    	If not provided, uses MONGODB_DATABASE environment variable

  -sample-size int
    	Number of documents to sample per collection
    	Default: 100

COMMON FLAGS:
  -output string
    	Output path for generated models.json
    	Default: configs/models.json

  -tables string
    	Comma-separated list of table names to process (PostgreSQL only)
    	Default: all tables in database

  -collections string
    	Comma-separated list of collection names to process (MongoDB only)
    	Default: all collections in database

  -help
    	Show this help message

EXAMPLES - PostgreSQL:
  # Using environment variable
  export DATABASE_URL="postgresql://postgres:password@localhost:5432/mydb"
  generate-models -type postgres

  # Specify database URL directly
  generate-models -type postgres -db "postgresql://user:pass@host:5432/db"

  # Custom output path
  generate-models -type postgres -output /custom/path/models.json

EXAMPLES - MongoDB:
  # Using environment variables
  export MONGODB_URI="mongodb://localhost:27017"
  export MONGODB_DATABASE="mydb"
  generate-models -type mongodb

  # Specify MongoDB connection
  generate-models -type mongodb -mongodb-uri "mongodb://localhost:27017" -mongodb-db mydb

  # Specific collections with larger sample size
  generate-models -type mongodb \
    -mongodb-uri "mongodb://localhost:27017" \
    -mongodb-db mydb \
    -collections "users,orders,products" \
    -sample-size 200

ENVIRONMENT VARIABLES:
  DATABASE_URL
    	PostgreSQL connection string (PostgreSQL mode)

  MONGODB_URI
    	MongoDB connection URI (MongoDB mode)

  MONGODB_DATABASE
    	MongoDB database name (MongoDB mode)

DESCRIPTION:
  Introspects a database and automatically generates a models.json configuration
  file that can be used with UDV.

  For PostgreSQL:
  - Discovers all tables in the public schema
  - Maps PostgreSQL data types to UDV types
  - Detects primary keys
  - Identifies nullable columns
  - Generates properly formatted models.json

  For MongoDB:
  - Samples documents from collections
  - Infers schema from document analysis
  - Detects field types by frequency
  - Identifies nullable fields
  - Generates properly formatted models.json
`)
}
