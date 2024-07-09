package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	. "mongodb/filter"
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

type Transaction struct {
	Context mongo.SessionContext
	Error   error
}

type collection struct {
	coll *mongo.Collection
	ctx  context.Context
}

func new(db *mongo.Database, cl string, ctx context.Context) collection {
	return collection{db.Collection(cl), ctx}
}

func (c collection) Aggregate(filter QueryFilter, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := c.coll.Aggregate(c.ctx, filter.Use(), opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return cursor, nil
	}
	return cursor, err
}

// FindOne find function with limit 1
func (c collection) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.coll.FindOne(c.ctx, filter, opts...)
}

func (c collection) Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := c.coll.Find(c.ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return cursor, nil
	}
	return cursor, err
}

// InsertOne default id type is primitive.ObjectID but you can insert int/string id's
func (c collection) InsertOne(body interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	insertedId, err := c.coll.InsertOne(c.ctx, body, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return insertedId, nil
	}
	return insertedId, err
}

// InsertMany default id type is primitive.ObjectID but you can insert int/string id's
func (c collection) InsertMany(body []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	insertedId, err := c.coll.InsertMany(c.ctx, body, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return insertedId, nil
	}

	return insertedId, err
}

func (c collection) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	upd, err := c.coll.UpdateOne(c.ctx, filter, update, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return upd, nil
	}
	return upd, err
}

func (c collection) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	upd, err := c.coll.UpdateMany(c.ctx, filter, update, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return upd, nil
	}
	return upd, err
}

func (c collection) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	del, err := c.coll.DeleteOne(c.ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return del, nil
	}
	return del, err
}

func (c collection) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	del, err := c.coll.DeleteMany(c.ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return del, nil
	}
	return del, err
}

func (c collection) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := c.coll.CountDocuments(c.ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return count, nil
	}
	return count, err
}

func (c collection) Collection() *mongo.Collection {
	return c.coll
}
