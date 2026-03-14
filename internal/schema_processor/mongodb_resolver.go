package schema_processor

import (
	"strconv"
	"strings"
)

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
