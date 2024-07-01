package mongo

import (
	. "github.com/maratIbatulin/mongodb/filter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection interface {
	//Aggregate find request to mongo with big bool of filters
	Aggregate(filter QueryFilter, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
	//FindOne filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
	//Find filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error)
	//InsertOne default id type is primitive.ObjectID but you can insert int/string id's
	//body can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	InsertOne(body interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	//InsertMany default id type is primitive.ObjectID but you can insert int/string id's
	//body can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	InsertMany(body []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)
	//UpdateOne filter and update can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//UpdateMany can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	//DeleteOne filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//DeleteMany filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)
	//CountDocuments filter can be struct/map[string]interface{}/[]map[string]interface(struct/bson.M/bson.D)
	CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error)
	//Collection return mongo.collection to use CRUD ignoring package
	Collection() *mongo.Collection
}

type Database interface {
	Collection(name string) Collection
}
