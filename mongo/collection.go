package mongo

import (
	"context"
	"github.com/maratIbatulin/mongodb/collection"
	. "github.com/maratIbatulin/mongodb/filter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	//Aggregate find request to mongo with big bool of filters
	Aggregate(ctx context.Context, filter QueryFilter, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	//FindOne filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	//Find filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	//InsertOne default id type is primitive.ObjectID but you can insert int/string id's
	//body can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	InsertOne(ctx context.Context, body interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	//InsertMany default id type is primitive.ObjectID but you can insert int/string id's
	//body can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	InsertMany(ctx context.Context, body []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	//UpdateOne filter and update can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//UpdateMany can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//DeleteOne filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//DeleteMany filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//CountDocuments filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	//Collection return mongo.collection to use CRUD ignoring package
	Collection() *mongo.Collection
	//Transaction start with goroutine to use transaction context in different functions if needed
	Transaction(sc chan collection.Transaction, acknowledged bool)
}
