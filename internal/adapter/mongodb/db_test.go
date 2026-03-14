package mongodb

import (
	"fmt"
	"testing"

	"udv/internal/adapter"
)

// MockMongoDB provides a mock MongoDB connection for testing
// In a real integration test, you would use testcontainers or a test MongoDB instance
type MockMongoDB struct {
	connected  bool
	documents  map[string][]map[string]interface{}
	insertedID string
	rowsAff    int64
}

func NewMockMongoDB() *MockMongoDB {
	return &MockMongoDB{
		documents: make(map[string][]map[string]interface{}),
	}
}

func (m *MockMongoDB) Connect(uri, dbName string) error {
	if uri == "" {
		return fmt.Errorf("empty uri provided")
	}
	if dbName == "" {
		return fmt.Errorf("empty dbName provided")
	}
	m.connected = true
	return nil
}

func (m *MockMongoDB) Close() error {
	if !m.connected {
		return fmt.Errorf("not connected")
	}
	m.connected = false
	return nil
}

func (m *MockMongoDB) Ping() error {
	if !m.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (m *MockMongoDB) ExecuteQuery(query interface{}, args ...interface{}) ([]map[string]interface{}, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	collection := mongoQuery.Collection
	if mongoQuery.Operation == "find" {
		if docs, exists := m.documents[collection]; exists {
			return docs, nil
		}
		return []map[string]interface{}{}, nil
	}

	return nil, fmt.Errorf("unsupported operation: %s", mongoQuery.Operation)
}

func (m *MockMongoDB) Exec(query interface{}, args ...interface{}) (adapter.ExecResult, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}

	mongoQuery, ok := query.(*MongoQuery)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}

	switch mongoQuery.Operation {
	case "insert":
		m.insertedID = "new_id"
		m.rowsAff = 1
		return &ExecInsertResult{InsertedID: m.insertedID}, nil
	case "update":
		m.rowsAff = 1
		return &ExecUpdateResult{ModifiedCount: 1}, nil
	case "delete":
		m.rowsAff = 1
		return &ExecUpdateResult{ModifiedCount: 1}, nil
	}

	return nil, fmt.Errorf("unsupported operation: %s", mongoQuery.Operation)
}

func TestConnect_Success(t *testing.T) {
	mock := NewMockMongoDB()

	err := mock.Connect("mongodb://localhost:27017", "testdb")
	if err != nil {
		t.Errorf("Connect failed: %v", err)
	}

	if !mock.connected {
		t.Error("Expected connected to be true")
	}
}

func TestConnect_EmptyURI(t *testing.T) {
	mock := NewMockMongoDB()

	err := mock.Connect("", "testdb")
	if err == nil {
		t.Error("Expected error for empty URI")
	}
}

func TestConnect_EmptyDB(t *testing.T) {
	mock := NewMockMongoDB()

	err := mock.Connect("mongodb://localhost:27017", "")
	if err == nil {
		t.Error("Expected error for empty database name")
	}
}

func TestClose_Success(t *testing.T) {
	mock := NewMockMongoDB()

	mock.connected = true

	err := mock.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	if mock.connected {
		t.Error("Expected connected to be false")
	}
}

func TestClose_NotConnected(t *testing.T) {
	mock := NewMockMongoDB()

	err := mock.Close()
	if err == nil {
		t.Error("Expected error when closing non-connected database")
	}
}

func TestPing_Success(t *testing.T) {
	mock := NewMockMongoDB()

	mock.connected = true

	err := mock.Ping()
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}
}

func TestPing_NotConnected(t *testing.T) {
	mock := NewMockMongoDB()

	err := mock.Ping()
	if err == nil {
		t.Error("Expected error when pinging non-connected database")
	}
}

func TestExecuteQuery_Find_Success(t *testing.T) {
	mock := NewMockMongoDB()

	mock.connected = true
	mock.documents["users"] = []map[string]interface{}{
		{"_id": "1", "name": "John"},
		{"_id": "2", "name": "Jane"},
	}

	mongoQuery := &MongoQuery{
		Collection: "users",
		Operation:  "find",
	}

	results, err := mock.ExecuteQuery(mongoQuery, nil)
	if err != nil {
		t.Errorf("ExecuteQuery failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results[0]["name"] != "John" {
		t.Errorf("Expected name 'John', got '%v'", results[0]["name"])
	}
}

func TestExecuteQuery_NotConnected(t *testing.T) {
	mock := NewMockMongoDB()

	mongoQuery := &MongoQuery{
		Collection: "users",
		Operation:  "find",
	}

	_, err := mock.ExecuteQuery(mongoQuery, nil)
	if err == nil {
		t.Error("Expected error when executing query on non-connected database")
	}
}

func TestExecuteQuery_InvalidQueryType(t *testing.T) {
	mock := NewMockMongoDB()

	mock.Connect("mongodb://localhost:27017", "testdb")

	_, err := mock.ExecuteQuery("invalid")
	if err == nil {
		t.Error("Expected error for invalid query type")
	}
}

func TestExec_Insert_Success(t *testing.T) {
	mock := NewMockMongoDB()

	mock.connected = true

	mongoQuery := &MongoQuery{
		Collection: "users",
		Operation:  "insert",
		Document:   map[string]interface{}{"name": "John"},
	}

	result, err := mock.Exec(mongoQuery, nil)
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected failed: %v", err)
	}

	if rowsAff != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAff)
	}
}

func TestExec_Update_Success(t *testing.T) {
	mock := NewMockMongoDB()

	mock.connected = true

	mongoQuery := &MongoQuery{
		Collection: "users",
		Operation:  "update",
		Filter:     map[string]interface{}{"_id": "1"},
		Update:     map[string]interface{}{"$set": map[string]interface{}{"name": "Jane"}},
	}

	result, err := mock.Exec(mongoQuery, nil)
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected failed: %v", err)
	}

	if rowsAff != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAff)
	}
}

func TestExec_Delete_Success(t *testing.T) {
	mock := NewMockMongoDB()

	mock.connected = true

	mongoQuery := &MongoQuery{
		Collection: "users",
		Operation:  "delete",
		Filter:     map[string]interface{}{"_id": "1"},
	}

	result, err := mock.Exec(mongoQuery, nil)
	if err != nil {
		t.Errorf("Exec failed: %v", err)
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected failed: %v", err)
	}

	if rowsAff != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAff)
	}
}

func TestExec_NotConnected(t *testing.T) {
	mock := NewMockMongoDB()

	mongoQuery := &MongoQuery{
		Collection: "users",
		Operation:  "insert",
	}

	_, err := mock.Exec(mongoQuery, nil)
	if err == nil {
		t.Error("Expected error when executing on non-connected database")
	}
}

func TestExec_UnsupportedOperation(t *testing.T) {
	mock := NewMockMongoDB()

	mock.connected = true

	mongoQuery := &MongoQuery{
		Collection: "users",
		Operation:  "unknown",
	}

	_, err := mock.Exec(mongoQuery)
	if err == nil {
		t.Error("Expected error for unsupported operation")
	}
}

func TestExecInsertResult_RowsAffected(t *testing.T) {
	result := &ExecInsertResult{InsertedID: "123"}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		t.Errorf("RowsAffected failed: %v", err)
	}

	if rowsAff != 1 {
		t.Errorf("Expected 1 row affected, got %d", rowsAff)
	}
}

func TestExecUpdateResult_RowsAffected(t *testing.T) {
	tests := []struct {
		name     string
		modified int64
		expected int64
	}{
		{"single row", 1, 1},
		{"multiple rows", 5, 5},
		{"no rows", 0, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := &ExecUpdateResult{ModifiedCount: test.modified}

			rowsAff, err := result.RowsAffected()
			if err != nil {
				t.Errorf("RowsAffected failed: %v", err)
			}

			if rowsAff != test.expected {
				t.Errorf("Expected %d rows affected, got %d", test.expected, rowsAff)
			}
		})
	}
}

// Test interface implementations
func TestDatabaseInterfaceImplementation(t *testing.T) {
	mock := NewMockMongoDB()
	var _ adapter.Database = mock
}

func TestExecResultInterfaceImplementation(t *testing.T) {
	insertResult := &ExecInsertResult{InsertedID: "123"}
	var _ adapter.ExecResult = insertResult

	updateResult := &ExecUpdateResult{ModifiedCount: 1}
	var _ adapter.ExecResult = updateResult
}
