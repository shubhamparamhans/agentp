package mongodb

import (
	"fmt"
	"strings"

	"udv/internal/dsl"
	"udv/internal/planner"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryBuilder struct{}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) BuildQuery(plan *planner.QueryPlan) (interface{}, []interface{}, error) {
	switch plan.Operation {
	case dsl.OpSelect:
		mq, err := qb.buildFindQuery(plan)
		return mq, nil, err
	case dsl.OpCreate:
		mq, err := qb.buildInsert(plan)
		return mq, nil, err
	case dsl.OpUpdate:
		mq, err := qb.buildUpdate(plan)
		return mq, nil, err
	case dsl.OpDelete:
		mq, err := qb.buildDelete(plan)
		return mq, nil, err
	default:
		return nil, nil, fmt.Errorf("unsupported operation: %s", plan.Operation)
	}
}

func (qb *QueryBuilder) buildFindQuery(plan *planner.QueryPlan) (*MongoQuery, error) {
	filter, err := qb.buildFilterFromExpr(plan.Filters)
	if err != nil {
		return nil, err
	}

	opt := options.Find()
	if plan.Pagination.Limit > 0 {
		opt.SetLimit(int64(plan.Pagination.Limit))
	}
	if plan.Pagination.Offset > 0 {
		opt.SetSkip(int64(plan.Pagination.Offset))
	}

	// Build sort
	if len(plan.Sort) > 0 {
		sortDoc := bson.D{}
		for _, s := range plan.Sort {
			direction := 1
			if strings.ToLower(s.Direction) == "desc" {
				direction = -1
			}
			fieldName := s.Column.ColumnName
			if s.Column.TableAlias != "" {
				fieldName = s.Column.TableAlias + "." + fieldName
			}
			sortDoc = append(sortDoc, bson.E{Key: fieldName, Value: direction})
		}
		opt.SetSort(sortDoc)
	}

	return &MongoQuery{
		Collection: plan.RootModel.Table,
		Operation:  "find",
		Filter:     filter,
		Options:    opt,
	}, nil
}

func (qb *QueryBuilder) buildFilterFromExpr(expr planner.FilterExpr) (bson.M, error) {
	if expr == nil {
		return bson.M{}, nil
	}

	switch f := expr.(type) {
	case *planner.ComparisonFilterIR:
		return qb.buildComparisonFilter(f)
	case *planner.LogicalFilterIR:
		return qb.buildLogicalFilter(f)
	default:
		return nil, fmt.Errorf("unsupported filter type: %T", expr)
	}
}

func (qb *QueryBuilder) buildComparisonFilter(f *planner.ComparisonFilterIR) (bson.M, error) {
	filter := make(bson.M)
	fieldName := f.Left.ColumnName

	mongoOp, mongoVal, err := qb.convertOperator(string(f.Operator), f.Value.Value)
	if err != nil {
		return nil, err
	}

	if mongoOp == "$eq" {
		filter[fieldName] = mongoVal
	} else {
		filter[fieldName] = bson.M{mongoOp: mongoVal}
	}

	return filter, nil
}

func (qb *QueryBuilder) buildLogicalFilter(f *planner.LogicalFilterIR) (bson.M, error) {
	filter := make(bson.M)

	switch strings.ToUpper(f.Op) {
	case "AND":
		andClauses := make([]bson.M, 0)
		for _, node := range f.Nodes {
			childFilter, err := qb.buildFilterFromExpr(node)
			if err != nil {
				return nil, err
			}
			andClauses = append(andClauses, childFilter)
		}
		filter["$and"] = andClauses
	case "OR":
		orClauses := make([]bson.M, 0)
		for _, node := range f.Nodes {
			childFilter, err := qb.buildFilterFromExpr(node)
			if err != nil {
				return nil, err
			}
			orClauses = append(orClauses, childFilter)
		}
		filter["$or"] = orClauses
	default:
		return nil, fmt.Errorf("unsupported logical operator: %s", f.Op)
	}

	return filter, nil
}

func (qb *QueryBuilder) convertOperator(op string, value interface{}) (string, interface{}, error) {
	switch op {
	case "=", "eq":
		return "$eq", value, nil
	case "!=", "ne":
		return "$ne", value, nil
	case ">", "gt":
		return "$gt", value, nil
	case ">=", "gte":
		return "$gte", value, nil
	case "<", "lt":
		return "$lt", value, nil
	case "<=", "lte":
		return "$lte", value, nil
	case "in":
		return "$in", value, nil
	case "not_in", "nin":
		return "$nin", value, nil
	case "like", "contains":
		strVal, ok := value.(string)
		if !ok {
			return "", nil, fmt.Errorf("like/contains operator requires string value")
		}
		return "$regex", strVal, nil
	case "is_null":
		return "$exists", false, nil
	case "not_null":
		return "$exists", true, nil
	default:
		return "", nil, fmt.Errorf("unsupported operator: %s", op)
	}
}

func (qb *QueryBuilder) buildInsert(plan *planner.QueryPlan) (*MongoQuery, error) {
	if len(plan.Data) == 0 {
		return nil, fmt.Errorf("insert data required")
	}

	doc := bson.M(plan.Data)

	return &MongoQuery{
		Collection: plan.RootModel.Table,
		Operation:  "insert",
		Document:   doc,
	}, nil
}

func (qb *QueryBuilder) buildUpdate(plan *planner.QueryPlan) (*MongoQuery, error) {
	filter, err := qb.buildFilterFromExpr(plan.Filters)
	if err != nil {
		return nil, err
	}

	updateDoc := bson.M{"$set": plan.Data}

	return &MongoQuery{
		Collection: plan.RootModel.Table,
		Operation:  "update",
		Filter:     filter,
		Update:     updateDoc,
	}, nil
}

func (qb *QueryBuilder) buildDelete(plan *planner.QueryPlan) (*MongoQuery, error) {
	filter, err := qb.buildFilterFromExpr(plan.Filters)
	if err != nil {
		return nil, err
	}

	return &MongoQuery{
		Collection: plan.RootModel.Table,
		Operation:  "delete",
		Filter:     filter,
	}, nil
}
