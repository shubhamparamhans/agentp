# MongoDB Schema Discovery and Models Generation

Automatically discover MongoDB database schemas and generate `models.json` configuration files through intelligent sampling and type inference. This process extends the Data Modelling Processor to support MongoDB databases, similar to PostgreSQL introspection.

## Overview

The MongoDB schema discovery process:
- Connects to MongoDB using URI and database name
- Lists collections or processes specified collections
- Samples configurable number of documents from each collection
- Statistically infers field types based on actual document content
- Detects nullable fields by analyzing field presence frequency
- Generates a valid `models.json` configuration file automatically

## How It Works

### 1. **Connection & Collection Discovery**
```
MongoDB URI + Database Name → Connected Client → Collection List
```

### 2. **Document Sampling**
- Retrieves sample documents from each collection
- Default sample size: 100 documents per collection
- Configurable via `-sample-size` flag
- Handles sparse or large collections efficiently

### 3. **Type Inference**
Analyzes BSON document structures to determine field types:
- **string** - Text values
- **integer** - Whole numbers (int32, int64)
- **decimal** - Floating-point numbers
- **boolean** - True/false values
- **timestamp** - Date/time objects
- **uuid** - ObjectID fields
- **array** - Collections of values
- **object** - Nested documents

### 4. **Nullability Detection**
```
Field Present Count / Total Documents
├─ > 90% present  → NOT nullable
└─ ≤ 90% present  → NULLABLE
```

### 5. **Primary Key Assignment**
- `_id` field automatically marked as primary key
- Preserved from MongoDB ObjectID

## Getting Started

### Prerequisites

Ensure the `generate-models` tool is built:

```bash
cd /Users/shubhamparamhans/Workspace/udv
go build ./cmd/generate-models
```

### Basic Usage

#### Using Environment Variables

```bash
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="your_database_name"

./generate-models -type mongodb
```

#### Using Command-Line Flags

```bash
./generate-models \
  -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "your_database_name"
```

## Command-Line Flags

### Required Flags
| Flag | Description | Env Variable |
|------|-------------|--------------|
| `-type mongodb` | Specify MongoDB as database type | - |
| `-mongodb-uri` | MongoDB connection URI | `MONGODB_URI` |
| `-mongodb-db` | MongoDB database name | `MONGODB_DATABASE` |

### Optional Flags
| Flag | Description | Default |
|------|-------------|---------|
| `-collections` | Comma-separated collection names to sample | All collections |
| `-sample-size` | Number of documents to sample per collection | 100 |
| `-output` | Output path for generated models.json | `configs/models.json` |
| `-help` | Display help message | - |

## Usage Examples

### Example 1: Basic Generation (All Collections)

```bash
./generate-models \
  -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "ecommerce"
```

**Output:** `configs/models.json` with models for all collections

### Example 2: Specific Collections with Custom Sample Size

```bash
./generate-models \
  -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "ecommerce" \
  -collections "users,orders,products" \
  -sample-size 200 \
  -output models/ecommerce_models.json
```

**Output:** Models for `users`, `orders`, and `products` based on 200 samples each

### Example 3: Using Environment Variables

```bash
export MONGODB_URI="mongodb+srv://user:password@cluster.mongodb.net/?retryWrites=true"
export MONGODB_DATABASE="production_db"

./generate-models -type mongodb -collections "users,sessions"
```

### Example 4: Local Development

```bash
# Start MongoDB in Docker
docker run -d -p 27017:27017 mongo:latest

# Generate models
./generate-models \
  -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "dev_db" \
  -sample-size 50
```

## Generated Models.json Format

The tool generates a `models.json` file with this structure:

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
          "name": "profile",
          "type": "object",
          "nullable": true
        }
      ]
    }
  ]
}
```

## Integration with Your Application

### Step 1: Generate Models
```bash
./generate-models -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "myapp"
```

### Step 2: Load Models
The application automatically loads `configs/models.json` on startup

### Step 3: Use MongoDB Adapter
Set environment variable for runtime database selection:
```bash
export DB_TYPE=mongodb
export MONGODB_URI="mongodb://localhost:27017"
export MONGODB_DATABASE="myapp"

./server  # Server will use MongoDB
```

## Features & Capabilities

✅ **Automatic Type Detection** - Infers types from actual data  
✅ **Nullability Analysis** - Determines required vs optional fields  
✅ **Sparse Data Support** - Handles documents with different field sets  
✅ **Large Dataset Handling** - Efficient sampling from large collections  
✅ **Multiple Collections** - Process entire database or specific collections  
✅ **Primary Key Detection** - Automatically identifies `_id` as primary key  
✅ **Nested Objects** - Supports object and array types  

## Advanced Configuration

### High-Volume Collections
For collections with millions of documents:
```bash
./generate-models -type mongodb \
  -mongodb-uri "mongodb://localhost:27017" \
  -mongodb-db "analytics" \
  -collections "events,logs" \
  -sample-size 500  # Larger sample for better type detection
```

### Production Deployments
```bash
# Use Atlas connection string
./generate-models \
  -type mongodb \
  -mongodb-uri "mongodb+srv://user:pass@cluster.mongodb.net/?retryWrites=true" \
  -mongodb-db "production" \
  -output configs/production_models.json
```

### Incremental Updates
```bash
# Process only new collections
./generate-models \
  -type mongodb \
  -mongodb-db "mydb" \
  -collections "new_collection1,new_collection2"
```

## Type Mapping Reference

| BSON Type | Inferred Type | Example |
|-----------|--------------|---------|
| String | string | "John Doe" |
| Int32/Int64 | integer | 42, 1000000 |
| Double/Decimal128 | decimal | 3.14, 99.99 |
| Boolean | boolean | true, false |
| Date | timestamp | ISODate("2024-01-15") |
| ObjectId | uuid | ObjectId("...") |
| Array | array | [1, 2, 3] |
| Object/Document | object | { "nested": "value" } |

## Troubleshooting

### Connection Issues
```bash
# Test MongoDB connection
mongosh "mongodb://localhost:27017"

# Verify with URI
mongosh "mongodb+srv://user:pass@cluster.mongodb.net"
```

### Collection Not Found
- Verify collection name (case-sensitive)
- Check database name is correct
- Ensure MongoDB user has read permissions

### Type Detection Issues
- Increase sample size for better accuracy: `-sample-size 500`
- Check for mixed-type fields in documents
- Review generated `models.json` for corrections

## Related Documentation

- [MONGODB_IMPLEMENTATION_PLAN.md](MONGODB_IMPLEMENTATION_PLAN.md) - Architecture and design details
- [MONGODB_SCHEMA_DISCOVERY.md](MONGODB_SCHEMA_DISCOVERY.md) - Deep dive into schema discovery algorithm
- [MONGODB_SUPPORT_ANALYSIS.md](MONGODB_SUPPORT_ANALYSIS.md) - Feature analysis and capabilities
- [DATA_MODELLING_PROCESSOR.md](DATA_MODELLING_PROCESSOR.md) - General data modelling guide

## Notes

- **Schema-Less Support**: This process allows model generation for MongoDB's schema-less nature using actual document data
- **Statistical Inference**: Type detection is based on sampling - ensure adequate sample size for accuracy
- **Primary Key Convention**: `_id` field is always preserved as primary key
- **Backward Compatible**: Works alongside existing PostgreSQL schema introspection
- **Reproducible**: Same MongoDB database produces consistent models.json files
