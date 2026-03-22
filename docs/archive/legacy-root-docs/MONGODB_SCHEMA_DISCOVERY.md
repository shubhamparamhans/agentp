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

### Phase 2: Type Inference (8-12 hours)

**Goal:** Analyze documents and infer field types

**File:** `internal/schema_processor/mongodb_infer.go` (new)

```go
package schema_processor

import (
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "reflect"
    "strings"
)

// FieldStats tracks type occurrences for a field
type FieldStats struct {
    TypeCounts    map[FieldType]int  // How many times each type appears
    TotalCount    int                // Total documents analyzed
    NullCount     int                // How many times field is missing/null
    SampleValues  []interface{}      // Sample values for debugging
}

// CollectionSchema represents inferred schema for a collection
type CollectionSchema struct {
    CollectionName string
    Fields         map[string]*FieldStats
    DocumentCount  int
}

// InferSchema analyzes documents and infers field types
func InferSchema(documents []bson.M) *CollectionSchema {
    schema := &CollectionSchema{
        Fields: make(map[string]*FieldStats),
    }
    
    schema.DocumentCount = len(documents)
    
    // Analyze each document
    for _, doc := range documents {
        analyzeDocument(doc, "", schema.Fields)
    }
    
    return schema
}

// analyzeDocument recursively analyzes a document
func analyzeDocument(doc bson.M, prefix string, fields map[string]*FieldStats) {
    for key, value := range doc {
        fieldPath := key
        if prefix != "" {
            fieldPath = prefix + "." + key
        }
        
        // Initialize field stats if not exists
        if fields[fieldPath] == nil {
            fields[fieldPath] = &FieldStats{
                TypeCounts:   make(map[FieldType]int),
                SampleValues: make([]interface{}, 0, 5),
            }
        }
        
        stats := fields[fieldPath]
        stats.TotalCount++
        
        // Handle null/missing values
        if value == nil {
            stats.NullCount++
            continue
        }
        
        // Infer type from value
        fieldType := inferTypeFromValue(value)
        stats.TypeCounts[fieldType]++
        
        // Store sample value (limit to 5)
        if len(stats.SampleValues) < 5 {
            stats.SampleValues = append(stats.SampleValues, value)
        }
        
        // Handle nested documents
        if nestedDoc, ok := value.(bson.M); ok {
            analyzeDocument(nestedDoc, fieldPath, fields)
        }
        
        // Handle arrays
        if arr, ok := value.(bson.A); ok && len(arr) > 0 {
            // Analyze first element to infer array element type
            if len(arr) > 0 {
                elemType := inferTypeFromValue(arr[0])
                arrayFieldPath := fieldPath + "[]"
                if fields[arrayFieldPath] == nil {
                    fields[arrayFieldPath] = &FieldStats{
                        TypeCounts: make(map[FieldType]int),
                    }
                }
                fields[arrayFieldPath].TypeCounts[elemType]++
                fields[arrayFieldPath].TotalCount++
            }
        }
    }
}

// inferTypeFromValue determines FieldType from a Go value
func inferTypeFromValue(value interface{}) FieldType {
    if value == nil {
        return TypeString // Default for null
    }
    
    switch v := value.(type) {
    case bool:
        return TypeBoolean
    
    case int, int32, int64:
        return TypeInteger
    
    case float32, float64:
        return TypeDecimal
    
    case string:
        // Check if it's a UUID format
        if isUUID(v) {
            return TypeUUID
        }
        // Check if it's a date string
        if isDateString(v) {
            return TypeTimestamp
        }
        return TypeString
    
    case primitive.ObjectID:
        return TypeUUID // MongoDB ObjectID as UUID
    
    case primitive.DateTime, primitive.Timestamp:
        return TypeTimestamp
    
    case bson.M, map[string]interface{}:
        return TypeJSON // Nested document as JSON
    
    case bson.A, []interface{}:
        return TypeJSON // Array as JSON
    
    case primitive.Binary:
        return TypeBinary
    
    default:
        // Use reflection for unknown types
        rt := reflect.TypeOf(value)
        switch rt.Kind() {
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            return TypeInteger
        case reflect.Float32, reflect.Float64:
            return TypeDecimal
        case reflect.Bool:
            return TypeBoolean
        case reflect.String:
            return TypeString
        default:
            return TypeString // Default fallback
        }
    }
}

// Helper functions
func isUUID(s string) bool {
    // UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
    parts := strings.Split(s, "-")
    return len(parts) == 5 && len(s) == 36
}

func isDateString(s string) bool {
    // Check common date formats
    dateFormats := []string{
        "2006-01-02",
        "2006-01-02T15:04:05Z",
        "2006-01-02T15:04:05.000Z",
        "2006-01-02 15:04:05",
    }
    for _, format := range dateFormats {
        if _, err := time.Parse(format, s); err == nil {
            return true
        }
    }
    return false
}
```

### Phase 3: Type Resolution (6-8 hours)

**Goal:** Determine final type for each field based on statistics

**File:** `internal/schema_processor/mongodb_resolver.go` (new)

```go
package schema_processor

// ResolveFieldType determines the most likely type for a field
func ResolveFieldType(stats *FieldStats) (FieldType, bool) {
    if stats == nil || len(stats.TypeCounts) == 0 {
        return TypeString, true // Default to string, nullable
    }
    
    // Find most common type
    maxCount := 0
    var mostCommonType FieldType
    
    for fieldType, count := range stats.TypeCounts {
        if count > maxCount {
            maxCount = count
            mostCommonType = fieldType
        }
    }
    
    // Calculate nullability
    // Field is nullable if it's missing in more than 10% of documents
    nullabilityThreshold := float64(stats.TotalCount) * 0.1
    isNullable := float64(stats.NullCount) > nullabilityThreshold
    
    // If type is ambiguous (multiple types with similar counts), prefer more specific types
    if maxCount < stats.TotalCount/2 {
        // Less than 50% of documents have this type - might be ambiguous
        // Prefer more specific types
        if stats.TypeCounts[TypeInteger] > 0 && stats.TypeCounts[TypeString] > 0 {
            // If we see both integer and string, check if strings are numeric
            if areAllNumericStrings(stats.SampleValues) {
                return TypeInteger, isNullable
            }
        }
    }
    
    return mostCommonType, isNullable
}

// areAllNumericStrings checks if string values are actually numbers
func areAllNumericStrings(values []interface{}) bool {
    for _, v := range values {
        if str, ok := v.(string); ok {
            if _, err := strconv.Atoi(str); err != nil {
                return false
            }
        }
    }
    return true
}

// GenerateModelFromSchema converts inferred schema to Model
func GenerateModelFromSchema(collectionName string, schema *CollectionSchema) Model {
    var fields []Field
    
    for fieldPath, stats := range schema.Fields {
        // Skip array element types (handled separately if needed)
        if strings.HasSuffix(fieldPath, "[]") {
            continue
        }
        
        fieldType, isNullable := ResolveFieldType(stats)
        
        field := Field{
            Name:     fieldPath,
            Type:     fieldType,
            Nullable: isNullable,
        }
        
        fields = append(fields, field)
    }
    
    // Determine primary key (always _id in MongoDB)
    primaryKey := "_id"
    
    // Check if _id exists, if not use first field
    hasID := false
    for _, field := range fields {
        if field.Name == "_id" {
            hasID = true
            break
        }
    }
    
    if !hasID {
        // Add _id field if missing (shouldn't happen, but safety check)
        fields = append([]Field{{
            Name:     "_id",
            Type:     TypeUUID, // ObjectID as UUID
            Nullable: false,
        }}, fields...)
    }
    
    return Model{
        Name:       collectionName,
        Table:      collectionName, // Collection name
        PrimaryKey: primaryKey,
        Fields:     fields,
    }
}
```

### Phase 4: MongoDB Schema Processor (6-8 hours)

**Goal:** Main processor that orchestrates sampling, inference, and generation

**File:** `internal/schema_processor/mongodb_processor.go` (new)

```go
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
        
        log.Printf("✓ Generated model for %s with %d fields", collectionName, len(model.Fields))
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
```

### Phase 5: CLI Tool Update (4-6 hours)

**Goal:** Update generate-models CLI to support MongoDB

**File:** `cmd/generate-models/main.go` (modify)

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    
    "udv/internal/schema_processor"
)

func main() {
    // Database type flag
    dbType := flag.String("type", "postgres", "Database type: postgres or mongodb")
    
    // PostgreSQL flags
    postgresURL := flag.String("db", "", "PostgreSQL connection string")
    outputPath := flag.String("output", "configs/models.json", "Output path")
    tableNamesStr := flag.String("tables", "", "Comma-separated table names")
    
    // MongoDB flags
    mongodbURI := flag.String("mongodb-uri", "", "MongoDB connection URI")
    mongodbDB := flag.String("mongodb-db", "", "MongoDB database name")
    collectionNamesStr := flag.String("collections", "", "Comma-separated collection names")
    sampleSize := flag.Int("sample-size", 100, "Number of documents to sample (MongoDB only)")
    
    help := flag.Bool("help", false, "Show help")
    flag.Parse()
    
    if *help {
        printHelp()
        os.Exit(0)
    }
    
    switch *dbType {
    case "mongodb":
        generateMongoDBModels(*mongodbURI, *mongodbDB, *collectionNamesStr, *sampleSize, *outputPath)
    case "postgres", "":
        generatePostgreSQLModels(*postgresURL, *tableNamesStr, *outputPath)
    default:
        log.Fatalf("Unsupported database type: %s", *dbType)
    }
}

func generateMongoDBModels(uri, dbName, collectionsStr string, sampleSize int, outputPath string) {
    // Get URI from flag or environment
    if uri == "" {
        uri = os.Getenv("MONGODB_URI")
    }
    if uri == "" {
        log.Fatal("MongoDB URI not provided. Use -mongodb-uri or set MONGODB_URI env var")
    }
    
    // Get database name from flag or environment
    if dbName == "" {
        dbName = os.Getenv("MONGODB_DATABASE")
    }
    if dbName == "" {
        log.Fatal("MongoDB database name not provided. Use -mongodb-db or set MONGODB_DATABASE env var")
    }
    
    // Parse collection names
    var collectionNames []string
    if collectionsStr != "" {
        collectionNames = strings.Split(collectionsStr, ",")
        for i := range collectionNames {
            collectionNames[i] = strings.TrimSpace(collectionNames[i])
        }
    }
    
    log.Println("Connecting to MongoDB...")
    processor, err := schema_processor.NewMongoDBProcessor(uri, dbName)
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer processor.Close()
    
    log.Printf("Sampling %d documents per collection...", sampleSize)
    log.Println("Introspecting MongoDB schema...")
    
    err = processor.GenerateAndSaveModels(outputPath, collectionNames, sampleSize)
    if err != nil {
        log.Fatalf("Failed to generate models: %v", err)
    }
    
    fmt.Printf("\n✓ Models generated successfully at: %s\n", outputPath)
}

func generatePostgreSQLModels(dbURL, tablesStr, outputPath string) {
    // Existing PostgreSQL code...
}
```

---

## Configuration Options

### Sample Size Strategy

**Default:** 100 documents per collection

**Rationale:**
- Too few (10-20): May miss fields that appear rarely
- Too many (1000+): Slower, but more accurate
- 100: Good balance between speed and accuracy

**Configurable via CLI:**
```bash
generate-models -type mongodb -sample-size 200
```

### Type Resolution Strategy

**Most Common Type Wins:**
- If 80% of documents have `integer`, 20% have `string` → Type: `integer`
- If 50% have `integer`, 50% have `string` → Type: `string` (fallback)

**Nullability Detection:**
- Field appears in < 90% of documents → `nullable: true`
- Field appears in ≥ 90% of documents → `nullable: false`

---

## Usage Examples

### Basic Usage
```bash
# Using environment variables
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="mydb"
generate-models -type mongodb

# Using flags
generate-models -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "mydb" \
  -output configs/models.json
```

### Specific Collections
```bash
generate-models -type mongodb \
  -mongodb-uri "mongodb://..." \
  -mongodb-db "mydb" \
  -collections "users,orders,products" \
  -sample-size 200
```

### Custom Sample Size
```bash
# Sample 500 documents for better accuracy
generate-models -type mongodb \
  -mongodb-uri "mongodb://..." \
  -mongodb-db "mydb" \
  -sample-size 500
```

---

## Generated models.json Example

```json
{
  "models": [
    {
      "name": "users",
      "table": "users",
      "primaryKey": "_id",
      "fields": [
        {
          "name": "_id",
          "type": "uuid",
          "nullable": false
        },
        {
          "name": "name",
          "type": "string",
          "nullable": false
        },
        {
          "name": "email",
          "type": "string",
          "nullable": false
        },
        {
          "name": "age",
          "type": "integer",
          "nullable": true
        },
        {
          "name": "metadata",
          "type": "json",
          "nullable": true
        },
        {
          "name": "created_at",
          "type": "timestamp",
          "nullable": true
        }
      ]
    }
  ]
}
```

---

## Limitations & Considerations

### 1. Schema Evolution
- **Issue:** MongoDB documents can change structure over time
- **Solution:** Re-run schema discovery periodically or on schema changes
- **Recommendation:** Add `--force` flag to overwrite existing models.json

### 2. Sparse Fields
- **Issue:** Fields that appear in < 10% of documents might be missed
- **Solution:** Increase sample size or use `--min-frequency` flag
- **Future:** Add field frequency threshold option

### 3. Type Ambiguity
- **Issue:** Same field might have different types in different documents
- **Solution:** Use most common type, log warnings for ambiguity
- **Future:** Support union types or "any" type

### 4. Nested Documents
- **Issue:** Deeply nested structures can create long field paths
- **Solution:** Flatten with dot notation (e.g., `address.city`)
- **Future:** Option to keep nested structure as JSON

### 5. Arrays
- **Issue:** Arrays can contain mixed types
- **Solution:** Analyze first element, treat as JSON if mixed
- **Future:** Support array element type specification

---

## Testing Strategy

### Unit Tests
```go
func TestInferTypeFromValue(t *testing.T) {
    tests := []struct {
        value    interface{}
        expected FieldType
    }{
        {true, TypeBoolean},
        {42, TypeInteger},
        {3.14, TypeDecimal},
        {"hello", TypeString},
        {primitive.NewObjectID(), TypeUUID},
    }
    
    for _, tt := range tests {
        result := inferTypeFromValue(tt.value)
        assert.Equal(t, tt.expected, result)
    }
}

func TestInferSchema(t *testing.T) {
    documents := []bson.M{
        {"name": "John", "age": 30, "active": true},
        {"name": "Jane", "age": 25, "active": false},
        {"name": "Bob", "age": 35},
    }
    
    schema := InferSchema(documents)
    
    assert.NotNil(t, schema.Fields["name"])
    assert.Equal(t, TypeString, schema.Fields["name"].TypeCounts[TypeString])
    assert.True(t, schema.Fields["active"].Nullable) // Missing in one doc
}
```

### Integration Tests
- Test with real MongoDB instance
- Test with various document structures
- Test edge cases (empty collections, single document, etc.)

---

## Implementation Effort

| Phase | Task | Hours |
|-------|------|-------|
| **Phase 1** | Document Sampling | 4-6 |
| **Phase 2** | Type Inference | 8-12 |
| **Phase 3** | Type Resolution | 6-8 |
| **Phase 4** | MongoDB Processor | 6-8 |
| **Phase 5** | CLI Updates | 4-6 |
| **Testing** | Unit & Integration | 6-8 |
| **TOTAL** | | **34-48 hours** |

**Estimated Timeline:** 1-1.5 weeks

---

## Comparison: PostgreSQL vs MongoDB

| Aspect | PostgreSQL | MongoDB |
|--------|------------|---------|
| **Schema Source** | information_schema | Document sampling |
| **Accuracy** | 100% (static schema) | ~95% (statistical inference) |
| **Speed** | Fast (metadata query) | Slower (document sampling) |
| **Schema Changes** | Manual update needed | Re-run discovery |
| **Nullability** | From column definition | From field frequency |
| **Type Detection** | From column type | From value analysis |

---

## Future Enhancements

### 1. Incremental Updates
- Detect schema changes and update models.json incrementally
- Only re-analyze changed collections

### 2. Schema Validation
- Use generated models.json to validate incoming documents
- Warn about type mismatches

### 3. Schema Versioning
- Track schema evolution over time
- Support multiple schema versions

### 4. Custom Type Mappings
- Allow user-defined type mappings
- Support custom field type inference rules

### 5. Field Frequency Reporting
- Report how often each field appears
- Help identify optional vs required fields

---

## Success Criteria

✅ **MongoDB Schema Discovery Complete When:**
1. Can connect to MongoDB database
2. Can sample documents from collections
3. Can infer field types from samples
4. Can detect nullability from field frequency
5. Can generate valid models.json
6. CLI tool supports MongoDB option
7. Tests pass for all inference logic
8. Documentation updated

---

## Next Steps

1. **Review this design** - Confirm approach
2. **Implement Phase 1** - Document sampling
3. **Implement Phase 2** - Type inference
4. **Implement Phase 3** - Type resolution
5. **Implement Phase 4** - Main processor
6. **Update CLI** - Add MongoDB support
7. **Test thoroughly** - Unit and integration tests
8. **Document usage** - Update quick start guide

---

**Last Updated:** Based on current PostgreSQL schema processor analysis  
**Estimated Effort:** 34-48 hours (~1-1.5 weeks)  
**Priority:** High (enables MongoDB support)

