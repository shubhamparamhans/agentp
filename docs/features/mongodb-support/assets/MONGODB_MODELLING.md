````markdown
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

... (full content copied from original `docs/MONGODB_MODELLING.md`)

````
