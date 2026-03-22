````markdown
# MongoDB Schema Discovery for models.json Generation

## Overview

Since MongoDB is schema-less, we need a different approach than PostgreSQL's static schema introspection. This document explains how to implement schema discovery for MongoDB that generates `models.json` by analyzing actual documents in collections.

---

## Challenge: Schema-less vs Schema-full

### PostgreSQL (Current Approach)
- ✅ **Static Schema**: Tables have fixed columns with defined types
- ✅ **Information Schema**: Can query `information_schema` for metadata
- ✅ **Deterministic**: Same query always returns same structure

### MongoDB (New Challenge)
- ❌ **No Static Schema**: Collections don't enforce field types
- ❌ **No Information Schema**: No equivalent metadata tables
- ❌ **Variable Structure**: Documents can have different fields
- ✅ **Solution**: Sample documents and infer schema from actual data

---

## Approach: Statistical Schema Inference

### Strategy

Instead of querying metadata (which doesn't exist), we:

1. **Sample Documents** - Analyze multiple documents from each collection
2. **Type Inference** - Determine most common type for each field
3. **Frequency Analysis** - Calculate how often each field appears
4. **Nullable Detection** - Fields that don't appear in all documents are nullable
5. **Primary Key Detection** - MongoDB always has `_id` field

---

## Implementation Design

### Phase 1: Document Sampling (4-6 hours)

**Goal:** Sample representative documents from each collection

**File:** `internal/schema_processor/mongodb_sampler.go` (new)

```go
package schema_processor

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
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
```

... (full content copied from original `docs/MONGODB_SCHEMA_DISCOVERY.md`)

````
