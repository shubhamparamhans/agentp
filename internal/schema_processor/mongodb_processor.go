package schema_processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBProcessor struct {
	sampler *MongoDBSampler
	ctx     context.Context
}

func NewMongoDBProcessor(uri string, dbName string) (*MongoDBProcessor, error) {
	ctx := context.Background()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	sampler := NewMongoDBSampler(client, dbName)

	return &MongoDBProcessor{
		sampler: sampler,
		ctx:     ctx,
	}, nil
}

// GenerateModels generates models from MongoDB collections
func (mp *MongoDBProcessor) GenerateModels(collectionNames []string, sampleSize int) ([]Model, error) {
	var collections []string
	var err error

	// Get all collections if none specified
	if len(collectionNames) == 0 {
		collections, err = mp.sampler.GetAllCollections()
		if err != nil {
			return nil, fmt.Errorf("failed to get collections: %w", err)
		}
	} else {
		collections = collectionNames
	}

	if len(collections) == 0 {
		return nil, fmt.Errorf("no collections found")
	}

	log.Printf("Found %d collections", len(collections))

	var models []Model

	for _, collectionName := range collections {
		log.Printf("Analyzing collection: %s", collectionName)

		// Sample documents
		documents, err := mp.sampler.SampleDocuments(collectionName, sampleSize)
		if err != nil {
			log.Printf("Warning: Failed to sample documents from %s: %v", collectionName, err)
			continue
		}

		if len(documents) == 0 {
			log.Printf("Warning: Collection %s is empty, skipping", collectionName)
			continue
		}

		// Infer schema from samples
		schema := InferSchema(documents)

		// Generate model
		model := GenerateModelFromSchema(collectionName, schema)

		models = append(models, model)

		log.Printf("7 Generated model for %s with %d fields", collectionName, len(model.Fields))
	}

	return models, nil
}

// GenerateAndSaveModels generates models and saves to file
func (mp *MongoDBProcessor) GenerateAndSaveModels(
	outputPath string,
	collectionNames []string,
	sampleSize int,
) error {
	models, err := mp.GenerateModels(collectionNames, sampleSize)
	if err != nil {
		return fmt.Errorf("failed to generate models: %w", err)
	}

	if len(models) == 0 {
		return fmt.Errorf("no valid models generated")
	}

	// Create config
	config := ModelConfig{
		Models: models,
	}

	// Marshal to JSON
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

	return nil
}

func (mp *MongoDBProcessor) Close() error {
	return mp.sampler.client.Disconnect(mp.ctx)
}
