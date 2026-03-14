package schema_processor

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBSampler struct {
	client   *mongo.Client
	database *mongo.Database
	ctx      context.Context
}

func NewMongoDBSampler(client *mongo.Client, dbName string) *MongoDBSampler {
	return &MongoDBSampler{
		client:   client,
		database: client.Database(dbName),
		ctx:      context.Background(),
	}
}

// SampleDocuments samples N documents from a collection
func (s *MongoDBSampler) SampleDocuments(collectionName string, sampleSize int) ([]bson.M, error) {
	collection := s.database.Collection(collectionName)

	// Use aggregation with $sample for random sampling
	pipeline := []bson.M{
		{"$sample": bson.M{"size": sampleSize}},
	}

	cursor, err := collection.Aggregate(s.ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	var documents []bson.M
	if err = cursor.All(s.ctx, &documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GetAllCollections lists all collections in the database
func (s *MongoDBSampler) GetAllCollections() ([]string, error) {
	collections, err := s.database.ListCollectionNames(s.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return collections, nil
}
