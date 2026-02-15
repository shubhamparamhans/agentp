package schema_processor

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestInferSchema_SingleDocument(t *testing.T) {
	docs := []bson.M{
		{
			"_id":   primitive.NewObjectID(),
			"name":  "John",
			"age":   30,
			"email": "john@example.com",
		},
	}

	schema := InferSchema(docs)

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	if schema.DocumentCount != 1 {
		t.Errorf("Expected 1 document, got %d", schema.DocumentCount)
	}

	if len(schema.Fields) == 0 {
		t.Error("Expected non-empty fields")
	}

	// Check for expected fields
	expectedFields := []string{"_id", "name", "age", "email"}
	for _, field := range expectedFields {
		if _, exists := schema.Fields[field]; !exists {
			t.Errorf("Expected field %q not found", field)
		}
	}
}

func TestInferSchema_MultipleDocuments(t *testing.T) {
	docs := []bson.M{
		{
			"_id":   primitive.NewObjectID(),
			"name":  "John",
			"age":   30,
			"email": "john@example.com",
		},
		{
			"_id":   primitive.NewObjectID(),
			"name":  "Jane",
			"age":   25,
			"email": "jane@example.com",
		},
		{
			"_id":   primitive.NewObjectID(),
			"name":  "Bob",
			"age":   35,
			"email": "bob@example.com",
		},
	}

	schema := InferSchema(docs)

	if schema.DocumentCount != 3 {
		t.Errorf("Expected 3 documents, got %d", schema.DocumentCount)
	}

	if len(schema.Fields) == 0 {
		t.Error("Expected non-empty fields")
	}
}

func TestInferSchema_EmptyDocuments(t *testing.T) {
	docs := []bson.M{}

	schema := InferSchema(docs)

	if schema == nil {
		t.Fatal("Expected non-nil schema")
	}

	if schema.DocumentCount != 0 {
		t.Errorf("Expected 0 documents, got %d", schema.DocumentCount)
	}

	if len(schema.Fields) != 0 {
		t.Errorf("Expected 0 fields for empty documents, got %d", len(schema.Fields))
	}
}

func TestInferSchema_SparseFields(t *testing.T) {
	docs := []bson.M{
		{
			"_id":   primitive.NewObjectID(),
			"name":  "John",
			"email": "john@example.com",
		},
		{
			"_id":   primitive.NewObjectID(),
			"name":  "Jane",
			"phone": "123-456",
		},
		{
			"_id":   primitive.NewObjectID(),
			"age":   30,
			"email": "someone@example.com",
		},
	}

	schema := InferSchema(docs)

	if len(schema.Fields) == 0 {
		t.Error("Expected non-empty fields for sparse data")
	}

	// All fields should be tracked
	expectedFields := []string{"_id", "name", "email", "phone", "age"}
	for _, field := range expectedFields {
		if _, exists := schema.Fields[field]; !exists {
			t.Errorf("Expected field %q not found", field)
		}
	}
}

func TestInferSchema_FieldStats(t *testing.T) {
	docs := []bson.M{
		{"_id": primitive.NewObjectID(), "value": "string"},
		{"_id": primitive.NewObjectID(), "value": "another"},
		{"_id": primitive.NewObjectID(), "value": "more"},
	}

	schema := InferSchema(docs)

	if _, exists := schema.Fields["value"]; !exists {
		t.Fatal("Expected 'value' field")
	}

	stats := schema.Fields["value"]
	if stats == nil {
		t.Fatal("Expected non-nil FieldStats")
	}

	if stats.TotalCount != 3 {
		t.Errorf("Expected TotalCount 3, got %d", stats.TotalCount)
	}
}

func TestGenerateModelFromSchema_BasicFields(t *testing.T) {
	schema := &CollectionSchema{
		DocumentCount: 100,
		Fields: map[string]*FieldStats{
			"_id": {
				TypeCounts: map[FieldType]int{TypeUUID: 100},
				TotalCount: 100,
				NullCount:  0,
			},
			"name": {
				TypeCounts: map[FieldType]int{TypeString: 100},
				TotalCount: 100,
				NullCount:  0,
			},
			"age": {
				TypeCounts: map[FieldType]int{TypeInteger: 100},
				TotalCount: 100,
				NullCount:  0,
			},
		},
	}

	model := GenerateModelFromSchema("users", schema)

	if model.Name != "users" {
		t.Errorf("Expected name 'users', got %q", model.Name)
	}

	if model.Table != "users" {
		t.Errorf("Expected table 'users', got %q", model.Table)
	}

	if model.PrimaryKey != "_id" {
		t.Errorf("Expected primary key '_id', got %q", model.PrimaryKey)
	}

	if len(model.Fields) == 0 {
		t.Error("Expected non-empty fields")
	}
}

func TestGenerateModelFromSchema_Nullability(t *testing.T) {
	schema := &CollectionSchema{
		DocumentCount: 100,
		Fields: map[string]*FieldStats{
			"_id": {
				TypeCounts: map[FieldType]int{TypeUUID: 100},
				TotalCount: 100,
				NullCount:  0, // 0% missing = not nullable
			},
			"optional_field": {
				TypeCounts: map[FieldType]int{TypeString: 50},
				TotalCount: 100,
				NullCount:  50, // 50% missing = nullable
			},
		},
	}

	model := GenerateModelFromSchema("test", schema)

	// Find fields and check nullability
	var idField, optField *Field
	for i := range model.Fields {
		if model.Fields[i].Name == "_id" {
			idField = &model.Fields[i]
		}
		if model.Fields[i].Name == "optional_field" {
			optField = &model.Fields[i]
		}
	}

	if idField == nil {
		t.Fatal("Expected '_id' field")
	}

	if optField == nil {
		t.Fatal("Expected 'optional_field' field")
	}

	if idField.Nullable {
		t.Error("_id field should not be nullable")
	}

	if !optField.Nullable {
		t.Error("optional_field should be nullable")
	}
}

func TestGenerateModelFromSchema_FieldTypes(t *testing.T) {
	schema := &CollectionSchema{
		DocumentCount: 100,
		Fields: map[string]*FieldStats{
			"string_field": {
				TypeCounts: map[FieldType]int{TypeString: 100},
				TotalCount: 100,
				NullCount:  0,
			},
			"int_field": {
				TypeCounts: map[FieldType]int{TypeInteger: 100},
				TotalCount: 100,
				NullCount:  0,
			},
			"decimal_field": {
				TypeCounts: map[FieldType]int{TypeDecimal: 100},
				TotalCount: 100,
				NullCount:  0,
			},
			"bool_field": {
				TypeCounts: map[FieldType]int{TypeBoolean: 100},
				TotalCount: 100,
				NullCount:  0,
			},
		},
	}

	model := GenerateModelFromSchema("test", schema)

	expectedTypes := map[string]FieldType{
		"string_field":  TypeString,
		"int_field":     TypeInteger,
		"decimal_field": TypeDecimal,
		"bool_field":    TypeBoolean,
	}

	for _, field := range model.Fields {
		expectedType, exists := expectedTypes[field.Name]
		if !exists {
			continue
		}

		if field.Type != expectedType {
			t.Errorf("Field %q: expected type %q, got %q", field.Name, expectedType, field.Type)
		}
	}
}

func TestResolveFieldType_SingleType(t *testing.T) {
	stats := &FieldStats{
		TypeCounts: map[FieldType]int{
			TypeString: 100,
		},
		TotalCount: 100,
		NullCount:  0,
	}

	fieldType, isNullable := ResolveFieldType(stats)

	if fieldType != TypeString {
		t.Errorf("Expected TypeString, got %s", fieldType)
	}

	if isNullable {
		t.Error("Expected not nullable")
	}
}

func TestResolveFieldType_WithNulls(t *testing.T) {
	stats := &FieldStats{
		TypeCounts: map[FieldType]int{
			TypeString: 80,
		},
		TotalCount: 100,
		NullCount:  20, // 20% missing
	}

	fieldType, isNullable := ResolveFieldType(stats)

	if fieldType != TypeString {
		t.Errorf("Expected TypeString, got %s", fieldType)
	}

	if !isNullable {
		t.Error("Expected nullable")
	}
}

func TestResolveFieldType_NilStats(t *testing.T) {
	fieldType, isNullable := ResolveFieldType(nil)

	if fieldType != TypeString {
		t.Errorf("Expected TypeString default, got %s", fieldType)
	}

	if !isNullable {
		t.Error("Expected nullable for nil stats")
	}
}

func TestResolveFieldType_MixedTypes(t *testing.T) {
	stats := &FieldStats{
		TypeCounts: map[FieldType]int{
			TypeString:  50,
			TypeInteger: 50,
		},
		TotalCount: 100,
		NullCount:  0,
	}

	fieldType, _ := ResolveFieldType(stats)

	// Should pick one of the types present
	if fieldType != TypeString && fieldType != TypeInteger {
		t.Errorf("Expected TypeString or TypeInteger, got %s", fieldType)
	}
}

func TestAnalyzeDocument_SimpleDocument(t *testing.T) {
	doc := bson.M{
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}

	fields := make(map[string]*FieldStats)
	analyzeDocument(doc, "", fields)

	expectedFields := []string{"name", "age", "email"}
	for _, field := range expectedFields {
		if _, exists := fields[field]; !exists {
			t.Errorf("Expected field %q not found", field)
		}
	}
}

func TestAnalyzeDocument_MultipleDocuments(t *testing.T) {
	schema := &CollectionSchema{
		Fields: make(map[string]*FieldStats),
	}

	doc1 := bson.M{"name": "John", "age": 30}
	doc2 := bson.M{"name": "Jane", "age": 25}

	analyzeDocument(doc1, "", schema.Fields)
	analyzeDocument(doc2, "", schema.Fields)

	// Both documents analyzed
	if len(schema.Fields) == 0 {
		t.Error("Expected non-empty fields")
	}

	for _, field := range []string{"name", "age"} {
		stats, exists := schema.Fields[field]
		if !exists {
			t.Errorf("Expected field %q", field)
			continue
		}

		if stats.TotalCount != 2 {
			t.Errorf("Field %q: expected TotalCount 2, got %d", field, stats.TotalCount)
		}
	}
}

func TestAnalyzeDocument_NestedFields(t *testing.T) {
	doc := bson.M{
		"user": bson.M{
			"name": "John",
			"age":  30,
		},
	}

	fields := make(map[string]*FieldStats)
	analyzeDocument(doc, "", fields)

	// Should track nested fields
	if len(fields) == 0 {
		t.Error("Expected fields to be analyzed")
	}
}
