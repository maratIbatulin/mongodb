package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
)

type collection struct {
	coll *mongo.Collection
	ctx  context.Context
}

type Collection interface {
	Aggregate(filter *filter, opts ...*option.AggregateOptions) (*mongo.Cursor, error)
	FindOne(filter map[string]any, opts ...*option.FindOneOptions) *mongo.SingleResult
	FindOneAndUpdate(filter map[string]any, update any, opts ...*option.FindOneAndUpdateOptions) *mongo.SingleResult
	InsertOne(body any, opts ...*option.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(body []any, opts ...*option.InsertManyOptions) (*mongo.InsertManyResult, error)
	UpdateOne(filter map[string]any, update any, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(filter map[string]any, update any, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(filter map[string]any, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter map[string]any, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(filter map[string]any, opts ...*option.CountOptions) (int64, error)
	Collection() *mongo.Collection
}

func (c collection) Aggregate(filter *filter, opts ...*option.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := c.coll.Aggregate(c.ctx, filter.Use(), opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return cursor, nil
	}
	return cursor, err
}

func (c collection) FindOne(filter map[string]any, opts ...*option.FindOneOptions) *mongo.SingleResult {
	return c.coll.FindOne(c.ctx, filter, opts...)
}

func (c collection) FindOneAndUpdate(filter map[string]any, update any, opts ...*option.FindOneAndUpdateOptions) *mongo.SingleResult {
	return c.coll.FindOneAndUpdate(c.ctx, filter, map[string]any{"$set": update}, opts...)
}

func (c collection) InsertOne(body any, opts ...*option.InsertOneOptions) (*mongo.InsertOneResult, error) {
	insertedId, err := c.coll.InsertOne(c.ctx, body, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return insertedId, nil
	}
	return insertedId, err
}

func (c collection) InsertMany(body []any, opts ...*option.InsertManyOptions) (*mongo.InsertManyResult, error) {
	insertedId, err := c.coll.InsertMany(c.ctx, body, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return insertedId, nil
	}

	return insertedId, err
}

func (c collection) UpdateOne(filter map[string]any, update any, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error) {
	upd, err := c.coll.UpdateOne(c.ctx, filter, map[string]any{"$set": update}, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return upd, nil
	}
	return upd, err
}

func (c collection) UpdateMany(filter map[string]any, update any, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error) {
	upd, err := c.coll.UpdateMany(c.ctx, filter, map[string]any{"$set": update}, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return upd, nil
	}
	return upd, err
}

func (c collection) DeleteOne(filter map[string]any, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error) {
	del, err := c.coll.DeleteOne(c.ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return del, nil
	}
	return del, err
}

func (c collection) DeleteMany(filter map[string]any, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error) {
	del, err := c.coll.DeleteMany(c.ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return del, nil
	}
	return del, err
}

func (c collection) CountDocuments(filter map[string]any, opts ...*option.CountOptions) (int64, error) {
	count, err := c.coll.CountDocuments(c.ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return count, nil
	}
	return count, err
}

func (c collection) Collection() *mongo.Collection {
	return c.coll
}
