package main

import (
	"fmt"
	"net/http"
	"os"

	"udv/internal/config"
    "udv/internal/api"
	"udv/internal/adapter/postgres"
	"udv/internal/schema"
)

func main() {
	// Load configuration
	configPath := "configs/models.json"
	if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Log loaded models
	fmt.Printf("Loaded %d model(s):\n", len(cfg.Models))
	for _, model := range cfg.Models {
		fmt.Printf("  - %s (table: %s, primaryKey: %s)\n", model.Name, model.Table, model.PrimaryKey)
	}

	// Initialize schema registry
	registry := schema.NewRegistry()
	if err := registry.LoadFromConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize schema registry: %v\n", err)
		os.Exit(1)
	}

	// Log registry initialization
	fmt.Printf("Schema registry initialized with %d model(s)\n", len(registry.ListModels()))

	// Initialize database connection (optional)
	var db *postgres.Database
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		var err error
		db, err = postgres.Connect(dbURL)
		if err != nil {
			fmt.Printf("Warning: Could not connect to database: %v\n", err)
			fmt.Println("Running in SQL-generation-only mode")
		} else {
			defer db.Close()
			fmt.Println("Database connection established")
		}
	} else {
		fmt.Println("DATABASE_URL not set, running in SQL-generation-only mode")
	}

	// Health check endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	// Register API routes
	apiSrv := api.New(registry, db)
	apiSrv.RegisterRoutes(mux)


	// CORS middleware
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		mux.ServeHTTP(w, r)
	})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
