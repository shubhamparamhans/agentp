package mongodb

import (
	"context"
	"fmt"

	"udv/internal/adapter"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database represents a MongoDB database connection with its context.
type Database struct {
	ctx      context.Context
	database *mongo.Database
	client   *mongo.Client
}

// Compile-time assertion that Database implements adapter.Database interface
var _ adapter.Database = (*Database)(nil)

// Connect creates a new MongoDB client and connects to the given URI and database name.
func Connect(uri string, databaseName string) (*Database, error) {
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(databaseName)

	return &Database{
		ctx:      ctx,
		database: db,
		client:   client,
	}, nil
}

// Close disconnects the MongoDB client.
func (d *Database) Close() error {
	return d.client.Disconnect(d.ctx)
}

// Ping checks the connection to the MongoDB server.
func (d *Database) Ping() error {
	return d.client.Ping(d.ctx, nil)
}

// ExecResult is an interface representing the result of an exec operation
// like Insert, Update or Delete.
type ExecResult interface {
	RowsAffected() (int64, error)
}

// ExecInsertResult holds the result of an insert operation.
type ExecInsertResult struct {
	InsertedID interface{}
}

// RowsAffected for ExecInsertResult returns 1 since a single document is inserted.
func (r *ExecInsertResult) RowsAffected() (int64, error) {
	return 1, nil
}

// ExecUpdateResult holds the result of an update or delete operation.
type ExecUpdateResult struct {
	ModifiedCount int64
}

// RowsAffected for ExecUpdateResult returns the number of documents modified.
func (r *ExecUpdateResult) RowsAffected() (int64, error) {
	return r.ModifiedCount, nil
}

// Ensure our types implement adapter.ExecResult
var (
	_ adapter.ExecResult = (*ExecInsertResult)(nil)
	_ adapter.ExecResult = (*ExecUpdateResult)(nil)
)

// ExecuteQuery executes a read operation like find or aggregate and returns the results.
func (d *Database) ExecuteQuery(query interface{}, args ...interface{}) ([]map[string]interface{}, error) {
	mq, ok := query.(*MongoQuery)
	if !ok {
		return nil, fmt.Errorf("ExecuteQuery: invalid query type %T", query)
	}

	coll := d.database.Collection(mq.Collection)

	switch mq.Operation {
	case "find":
		cursor, err := coll.Find(d.ctx, mq.Filter, mq.Options.(*options.FindOptions))
		if err != nil {
			return nil, err
		}
		defer cursor.Close(d.ctx)

		var results []map[string]interface{}
		err = cursor.All(d.ctx, &results)
		if err != nil {
			return nil, err
		}
		return results, nil

	case "aggregate":
		cursor, err := coll.Aggregate(d.ctx, mq.Pipeline)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(d.ctx)

		var results []map[string]interface{}
		err = cursor.All(d.ctx, &results)
		if err != nil {
			return nil, err
		}
		return results, nil

	default:
		return nil, fmt.Errorf("ExecuteQuery: unsupported operation %s", mq.Operation)
	}
}

// Exec executes insert, update or delete operations and returns the result.
func (d *Database) Exec(query interface{}, args ...interface{}) (adapter.ExecResult, error) {
	mq, ok := query.(*MongoQuery)
	if !ok {
		return nil, fmt.Errorf("Exec: invalid query type %T", query)
	}

	coll := d.database.Collection(mq.Collection)

	switch mq.Operation {
	case "insert":
		insertResult, err := coll.InsertOne(d.ctx, mq.Document)
		if err != nil {
			return nil, err
		}
		return &ExecInsertResult{InsertedID: insertResult.InsertedID}, nil

	case "update":
		res, err := coll.UpdateMany(d.ctx, mq.Filter, mq.Update)
		if err != nil {
			return nil, err
		}
		return &ExecUpdateResult{ModifiedCount: res.ModifiedCount}, nil

	case "delete":
		res, err := coll.DeleteMany(d.ctx, mq.Filter)
		if err != nil {
			return nil, err
		}
		return &ExecUpdateResult{ModifiedCount: res.DeletedCount}, nil

	default:
		return nil, fmt.Errorf("Exec: unsupported operation %s", mq.Operation)
	}
}
