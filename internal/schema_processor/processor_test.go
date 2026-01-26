package schema_processor

import (
	"testing"
)

// TestMapPostgreSQLTypeToJSON tests the data type mapping
func TestMapPostgreSQLTypeToJSON(t *testing.T) {
	tests := []struct {
		pgType   string
		expected FieldType
	}{
		// Integer types
		{"integer", TypeInteger},
		{"int", TypeInteger},
		{"int4", TypeInteger},
		{"smallint", TypeInteger},
		{"bigint", TypeInteger},
		{"serial", TypeInteger},
		{"bigserial", TypeInteger},
		{"integer[]", TypeInteger},

		// String types
		{"text", TypeString},
		{"character varying", TypeString},
		{"varchar", TypeString},
		{"character", TypeString},
		{"char(50)", TypeString},
		{"varchar(255)", TypeString},

		// Decimal types
		{"numeric", TypeDecimal},
		{"decimal", TypeDecimal},
		{"money", TypeDecimal},
		{"double precision", TypeDecimal},
		{"float8", TypeDecimal},
		{"real", TypeDecimal},
		{"float4", TypeDecimal},

		// Boolean
		{"boolean", TypeBoolean},
		{"bool", TypeBoolean},

		// Timestamp types
		{"timestamp", TypeTimestamp},
		{"timestamp without time zone", TypeTimestamp},
		{"timestamp with time zone", TypeTimestamp},
		{"timestamptz", TypeTimestamp},
		{"date", TypeTimestamp},
		{"time", TypeTimestamp},

		// JSON types
		{"json", TypeJSON},
		{"jsonb", TypeJSON},

		// UUID
		{"uuid", TypeUUID},

		// Binary types
		{"bytea", TypeBinary},
		{"bit", TypeBinary},

		// Unknown type (should default to string)
		{"unknown_type", TypeString},
		{"custom_type", TypeString},
	}

	for _, test := range tests {
		t.Run(test.pgType, func(t *testing.T) {
			result := mapPostgreSQLTypeToJSON(test.pgType)
			if result != test.expected {
				t.Errorf("mapPostgreSQLTypeToJSON(%q) = %v, want %v", test.pgType, result, test.expected)
			}
		})
	}
}

// TestTypeValues ensures FieldType constants are strings
func TestFieldTypeValues(t *testing.T) {
	tests := []struct {
		ft       FieldType
		expected string
	}{
		{TypeInteger, "integer"},
		{TypeString, "string"},
		{TypeDecimal, "decimal"},
		{TypeBoolean, "boolean"},
		{TypeTimestamp, "timestamp"},
		{TypeJSON, "json"},
		{TypeUUID, "uuid"},
		{TypeBinary, "binary"},
	}

	for _, test := range tests {
		if string(test.ft) != test.expected {
			t.Errorf("FieldType value mismatch: got %q, want %q", test.ft, test.expected)
		}
	}
}
