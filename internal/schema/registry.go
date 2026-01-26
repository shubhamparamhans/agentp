package schema

import (
	"fmt"
	"sync"

	"udv/internal/config"
)

// RelationType represents the type of relationship between models
type RelationType string

const (
	OneToOne   RelationType = "one_to_one"
	OneToMany  RelationType = "one_to_many"
	ManyToOne  RelationType = "many_to_one"
	ManyToMany RelationType = "many_to_many"
)

// Field represents a model field
type Field struct {
	Name          string
	Type          string
	Nullable      bool
	Filterable    bool
	Groupable     bool
	Aggregatable  bool
}

// Relation represents a relationship to another model
type Relation struct {
	Type          RelationType
	TargetModel   string // Name of the related model
	ForeignKey    string // Local field name
	ReferenceKey  string // Field in target model
}

// Model represents a data model with its fields and relationships
type Model struct {
	Name        string
	Table       string
	PrimaryKey  string
	Fields      map[string]*Field
	Relations   map[string]*Relation
	FieldOrder  []string // Preserve field order
}

// Registry is the in-memory schema registry
type Registry struct {
	mu     sync.RWMutex
	models map[string]*Model
}

// NewRegistry creates a new empty registry
func NewRegistry() *Registry {
	return &Registry{
		models: make(map[string]*Model),
	}
}

// LoadFromConfig populates the registry from a config
func (r *Registry) LoadFromConfig(cfg *config.Config) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// First pass: create all models
	for _, cfgModel := range cfg.Models {
		model := &Model{
			Name:       cfgModel.Name,
			Table:      cfgModel.Table,
			PrimaryKey: cfgModel.PrimaryKey,
			Fields:     make(map[string]*Field),
			Relations:  make(map[string]*Relation),
			FieldOrder: []string{},
		}

		// Add fields with sensible defaults
		for _, cfgField := range cfgModel.Fields {
			field := &Field{
				Name:          cfgField.Name,
				Type:          cfgField.Type,
				Nullable:      cfgField.Nullable,
				Filterable:    true,  // Default: fields are filterable
				Groupable:     true,  // Default: fields are groupable
				Aggregatable:  true,  // All fields are aggregatable; validateAggregateForType validates function-type compatibility
			}

			model.Fields[cfgField.Name] = field
			model.FieldOrder = append(model.FieldOrder, cfgField.Name)
		}

		r.models[cfgModel.Name] = model
	}

	return nil
}

// GetModel returns a model by name (read-only)
func (r *Registry) GetModel(name string) *Model {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.models[name]
}

// GetField returns a field from a model
func (r *Registry) GetField(modelName, fieldName string) (*Field, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, exists := r.models[modelName]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", modelName)
	}

	field, exists := model.Fields[fieldName]
	if !exists {
		return nil, fmt.Errorf("field not found: %s.%s", modelName, fieldName)
	}

	return field, nil
}

// ModelExists checks if a model exists in the registry
func (r *Registry) ModelExists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.models[name]
	return exists
}

// FieldExists checks if a field exists in a model
func (r *Registry) FieldExists(modelName, fieldName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, exists := r.models[modelName]
	if !exists {
		return false
	}

	_, exists = model.Fields[fieldName]
	return exists
}

// ListModels returns all model names in the registry
func (r *Registry) ListModels() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	models := make([]string, 0, len(r.models))
	for name := range r.models {
		models = append(models, name)
	}
	return models
}

// GetModelFields returns all fields of a model in their defined order
func (r *Registry) GetModelFields(modelName string) ([]*Field, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, exists := r.models[modelName]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", modelName)
	}

	fields := make([]*Field, len(model.FieldOrder))
	for i, fieldName := range model.FieldOrder {
		fields[i] = model.Fields[fieldName]
	}
	return fields, nil
}
