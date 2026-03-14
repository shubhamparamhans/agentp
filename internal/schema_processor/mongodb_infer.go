package schema_processor

import (
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FieldStats tracks type occurrences for a field
// FieldType is assumed to be defined elsewhere

type FieldStats struct {
	TypeCounts   map[FieldType]int // How many times each type appears
	TotalCount   int               // Total documents analyzed
	NullCount    int               // How many times field is missing/null
	SampleValues []interface{}     // Sample values for debugging
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
