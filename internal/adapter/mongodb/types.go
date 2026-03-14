package mongodb

// MongoQuery represents a MongoDB operation to be executed
// Operation can be one of: find, aggregate, insert, update, delete
// Filter is a bson.M representing query filter
// Update is a bson.M representing update document
// Document is a bson.M for insert operations
// Pipeline is a mongo.Pipeline for aggregations
// Options are find options

type MongoQuery struct {
	Collection string
	Operation  string
	Filter     interface{}
	Pipeline   interface{}
	Update     interface{}
	Document   interface{}
	Options    interface{}
}
