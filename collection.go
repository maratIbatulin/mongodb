package mongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
)

type collection struct {
	coll *mongo.Collection
	ctx  context.Context
}

type Collection interface {
	Aggregate(ctx context.Context, filter *filter, opts ...*option.AggregateOptions) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter D, opts ...*option.FindOneOptions) *mongo.SingleResult
	FindOneAndUpdate(ctx context.Context, filter D, update D, opts ...*option.FindOneAndUpdateOptions) *mongo.SingleResult
	InsertOne(ctx context.Context, body any, opts ...*option.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(ctx context.Context, body []any, opts ...*option.InsertManyOptions) (*mongo.InsertManyResult, error)
	UpdateOne(ctx context.Context, filter D, update D, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(ctx context.Context, filter D, update D, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter D, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(ctx context.Context, filter D, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter D, opts ...*option.CountOptions) (int64, error)
	Watch(ctx context.Context, filter *filter, opts ...*option.ChangeStreamOptions) (*mongo.ChangeStream, error)
	Collection() *mongo.Collection
}

func (c collection) Aggregate(ctx context.Context, filter *filter, opts ...*option.AggregateOptions) (*mongo.Cursor, error) {
	if ctx == nil {
		ctx = c.ctx
	}
	cursor, err := c.coll.Aggregate(ctx, filter.Use(), opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return cursor, nil
	}
	return cursor, err
}

func (c collection) FindOne(ctx context.Context, filter D, opts ...*option.FindOneOptions) *mongo.SingleResult {
	if ctx == nil {
		ctx = c.ctx
	}
	return c.coll.FindOne(ctx, filter, opts...)
}

func (c collection) FindOneAndUpdate(ctx context.Context, filter D, update D, opts ...*option.FindOneAndUpdateOptions) *mongo.SingleResult {
	bn := bson.D{}
	for _, e := range update {
		bn = append(bn, e)
	}
	if ctx == nil {
		ctx = c.ctx
	}
	return c.coll.FindOneAndUpdate(ctx, filter, bn, opts...)
}

func (c collection) InsertOne(ctx context.Context, body any, opts ...*option.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if ctx == nil {
		ctx = c.ctx
	}
	insertedId, err := c.coll.InsertOne(ctx, body, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return insertedId, nil
	}
	return insertedId, err
}

func (c collection) InsertMany(ctx context.Context, body []any, opts ...*option.InsertManyOptions) (*mongo.InsertManyResult, error) {
	if ctx == nil {
		ctx = c.ctx
	}
	insertedId, err := c.coll.InsertMany(ctx, body, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return insertedId, nil
	}

	return insertedId, err
}

func (c collection) UpdateOne(ctx context.Context, filter D, update D, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error) {
	bn := bson.D{}
	for _, e := range update {
		bn = append(bn, e)
	}
	if ctx == nil {
		ctx = c.ctx
	}
	upd, err := c.coll.UpdateOne(ctx, filter, bn, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return upd, nil
	}
	return upd, err
}

func (c collection) UpdateMany(ctx context.Context, filter D, update D, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error) {
	bn := bson.D{}
	for _, e := range update {
		bn = append(bn, e)
	}
	if ctx == nil {
		ctx = c.ctx
	}
	upd, err := c.coll.UpdateMany(ctx, filter, bn, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return upd, nil
	}
	return upd, err
}

func (c collection) DeleteOne(ctx context.Context, filter D, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error) {
	if ctx == nil {
		ctx = c.ctx
	}
	del, err := c.coll.DeleteOne(ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return del, nil
	}
	return del, err
}

func (c collection) DeleteMany(ctx context.Context, filter D, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error) {
	if ctx == nil {
		ctx = c.ctx
	}
	del, err := c.coll.DeleteMany(ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return del, nil
	}
	return del, err
}

func (c collection) CountDocuments(ctx context.Context, filter D, opts ...*option.CountOptions) (int64, error) {
	if ctx == nil {
		ctx = c.ctx
	}
	count, err := c.coll.CountDocuments(ctx, filter, opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return count, nil
	}
	return count, err
}

func (c collection) Watch(ctx context.Context, filter *filter, opts ...*option.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	if ctx == nil {
		ctx = c.ctx
	}
	stream, err := c.coll.Watch(ctx, filter.Use(), opts...)
	if errors.Is(err, mongo.ErrUnacknowledgedWrite) {
		return stream, nil
	}
	return stream, err
}

func (c collection) Collection() *mongo.Collection {
	return c.coll
}
