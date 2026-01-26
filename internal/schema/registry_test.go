package schema

import (
	"testing"

	"udv/internal/config"
)

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Errorf("NewRegistry() returned nil")
	}
	if len(reg.models) != 0 {
		t.Errorf("NewRegistry() should create empty registry, got %d models", len(reg.models))
	}
}

func TestLoadFromConfig(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "users",
				Table:      "users",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: false},
					{Name: "email", Type: "string", Nullable: false},
					{Name: "age", Type: "integer", Nullable: true},
					{Name: "created_at", Type: "timestamp", Nullable: false},
				},
			},
			{
				Name:       "orders",
				Table:      "orders",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "user_id", Type: "integer", Nullable: false},
					{Name: "total", Type: "decimal", Nullable: false},
					{Name: "created_at", Type: "timestamp", Nullable: false},
				},
			},
		},
	}

	reg := NewRegistry()
	err := reg.LoadFromConfig(cfg)
	if err != nil {
		t.Errorf("LoadFromConfig() error = %v, want nil", err)
	}

	// Check models loaded
	if !reg.ModelExists("users") {
		t.Errorf("LoadFromConfig() did not load users model")
	}
	if !reg.ModelExists("orders") {
		t.Errorf("LoadFromConfig() did not load orders model")
	}

	// Check model structure
	usersModel := reg.GetModel("users")
	if usersModel == nil {
		t.Errorf("GetModel(users) returned nil")
	}
	if usersModel.Name != "users" {
		t.Errorf("GetModel(users).Name = %s, want users", usersModel.Name)
	}
	if usersModel.Table != "users" {
		t.Errorf("GetModel(users).Table = %s, want users", usersModel.Table)
	}
	if usersModel.PrimaryKey != "id" {
		t.Errorf("GetModel(users).PrimaryKey = %s, want id", usersModel.PrimaryKey)
	}

	// Check fields
	if len(usersModel.Fields) != 5 {
		t.Errorf("users model has %d fields, want 5", len(usersModel.Fields))
	}
	if len(usersModel.FieldOrder) != 5 {
		t.Errorf("users model FieldOrder has %d items, want 5", len(usersModel.FieldOrder))
	}
}

func TestGetField(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "users",
				Table:      "users",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: false},
					{Name: "age", Type: "integer", Nullable: true},
				},
			},
		},
	}

	reg := NewRegistry()
	reg.LoadFromConfig(cfg)

	tests := []struct {
		modelName string
		fieldName string
		wantErr   bool
		checkFn   func(*Field) bool
	}{
		{
			modelName: "users",
			fieldName: "id",
			wantErr:   false,
			checkFn: func(f *Field) bool {
				return f.Name == "id" && f.Type == "integer" && !f.Nullable
			},
		},
		{
			modelName: "users",
			fieldName: "name",
			wantErr:   false,
			checkFn: func(f *Field) bool {
				return f.Name == "name" && f.Type == "string" && !f.Nullable && f.Filterable && f.Groupable
			},
		},
		{
			modelName: "users",
			fieldName: "age",
			wantErr:   false,
			checkFn: func(f *Field) bool {
				return f.Name == "age" && f.Type == "integer" && f.Nullable && f.Aggregatable
			},
		},
		{
			modelName: "users",
			fieldName: "nonexistent",
			wantErr:   true,
		},
		{
			modelName: "nonexistent",
			fieldName: "id",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		field, err := reg.GetField(tt.modelName, tt.fieldName)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetField(%s, %s) error = %v, wantErr %v", tt.modelName, tt.fieldName, err, tt.wantErr)
		}
		if !tt.wantErr && !tt.checkFn(field) {
			t.Errorf("GetField(%s, %s) returned unexpected field: %+v", tt.modelName, tt.fieldName, field)
		}
	}
}

func TestFieldExists(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "users",
				Table:      "users",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: false},
				},
			},
		},
	}

	reg := NewRegistry()
	reg.LoadFromConfig(cfg)

	tests := []struct {
		modelName string
		fieldName string
		want      bool
	}{
		{"users", "id", true},
		{"users", "name", true},
		{"users", "nonexistent", false},
		{"nonexistent", "id", false},
	}

	for _, tt := range tests {
		got := reg.FieldExists(tt.modelName, tt.fieldName)
		if got != tt.want {
			t.Errorf("FieldExists(%s, %s) = %v, want %v", tt.modelName, tt.fieldName, got, tt.want)
		}
	}
}

func TestListModels(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "users",
				Table:      "users",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
				},
			},
			{
				Name:       "orders",
				Table:      "orders",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
				},
			},
			{
				Name:       "products",
				Table:      "products",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
				},
			},
		},
	}

	reg := NewRegistry()
	reg.LoadFromConfig(cfg)

	models := reg.ListModels()
	if len(models) != 3 {
		t.Errorf("ListModels() returned %d models, want 3", len(models))
	}

	modelMap := make(map[string]bool)
	for _, m := range models {
		modelMap[m] = true
	}

	expectedModels := []string{"users", "orders", "products"}
	for _, expected := range expectedModels {
		if !modelMap[expected] {
			t.Errorf("ListModels() missing model %s", expected)
		}
	}
}

func TestGetModelFields(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "users",
				Table:      "users",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "name", Type: "string", Nullable: false},
					{Name: "email", Type: "string", Nullable: false},
				},
			},
		},
	}

	reg := NewRegistry()
	reg.LoadFromConfig(cfg)

	fields, err := reg.GetModelFields("users")
	if err != nil {
		t.Errorf("GetModelFields(users) error = %v, want nil", err)
	}

	if len(fields) != 3 {
		t.Errorf("GetModelFields(users) returned %d fields, want 3", len(fields))
	}

	// Check order is preserved
	expectedOrder := []string{"id", "name", "email"}
	for i, field := range fields {
		if field.Name != expectedOrder[i] {
			t.Errorf("GetModelFields(users)[%d].Name = %s, want %s", i, field.Name, expectedOrder[i])
		}
	}

	// Test nonexistent model
	_, err = reg.GetModelFields("nonexistent")
	if err == nil {
		t.Errorf("GetModelFields(nonexistent) error = nil, want error")
	}
}

func TestAggregatableFields(t *testing.T) {
	cfg := &config.Config{
		Models: []config.Model{
			{
				Name:       "stats",
				Table:      "stats",
				PrimaryKey: "id",
				Fields: []config.Field{
					{Name: "id", Type: "integer", Nullable: false},
					{Name: "count", Type: "integer", Nullable: false},
					{Name: "amount", Type: "decimal", Nullable: false},
					{Name: "price", Type: "float", Nullable: false},
					{Name: "name", Type: "string", Nullable: false},
					{Name: "is_active", Type: "boolean", Nullable: false},
					{Name: "created_at", Type: "timestamp", Nullable: false},
				},
			},
		},
	}

	reg := NewRegistry()
	reg.LoadFromConfig(cfg)

	tests := []struct {
		fieldName    string
		wantAgg      bool
		wantFilter   bool
		wantGroupable bool
	}{
		{"id", true, true, true},          // integer - aggregatable
		{"count", true, true, true},       // integer - aggregatable
		{"amount", true, true, true},      // decimal - aggregatable
		{"price", true, true, true},       // float - aggregatable
		{"name", true, true, true},        // string - aggregatable (COUNT/MIN/MAX work on all types)
		{"is_active", true, true, true},   // boolean - aggregatable (COUNT/MIN/MAX work on all types)
		{"created_at", true, true, true},  // timestamp - aggregatable (COUNT/MIN/MAX work on all types)
	}

	for _, tt := range tests {
		field, _ := reg.GetField("stats", tt.fieldName)
		if field.Aggregatable != tt.wantAgg {
			t.Errorf("Field %s Aggregatable = %v, want %v", tt.fieldName, field.Aggregatable, tt.wantAgg)
		}
		if field.Filterable != tt.wantFilter {
			t.Errorf("Field %s Filterable = %v, want %v", tt.fieldName, field.Filterable, tt.wantFilter)
		}
		if field.Groupable != tt.wantGroupable {
			t.Errorf("Field %s Groupable = %v, want %v", tt.fieldName, field.Groupable, tt.wantGroupable)
		}
	}
}
