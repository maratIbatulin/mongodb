package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
)

// collection represents a MongoDB collection wrapper with embedded context
type collection struct {
	coll *mongo.Collection
	ctx  context.Context
}

// Collection defines an interface for MongoDB collection operations
type Collection interface {
	// Aggregate executes an aggregation pipeline on the collection
	Aggregate(ctx context.Context, filter *filter, opts ...*option.AggregateOptions) (*mongo.Cursor, error)

	// FindOne returns a single document that matches the filter
	FindOne(ctx context.Context, filter D, opts ...*option.FindOneOptions) *mongo.SingleResult

	// FindOneAndUpdate finds a single document and updates it, returning the original
	FindOneAndUpdate(ctx context.Context, filter D, update D, opts ...*option.FindOneAndUpdateOptions) *mongo.SingleResult

	// InsertOne inserts a single document into the collection
	InsertOne(ctx context.Context, body any, opts ...*option.InsertOneOptions) (*mongo.InsertOneResult, error)

	// InsertMany inserts multiple documents into the collection
	InsertMany(ctx context.Context, body []any, opts ...*option.InsertManyOptions) (*mongo.InsertManyResult, error)

	// UpdateOne updates a single document matching the filter
	UpdateOne(ctx context.Context, filter D, update D, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error)

	// UpdateMany updates multiple documents matching the filter
	UpdateMany(ctx context.Context, filter D, update D, opts ...*option.UpdateOptions) (*mongo.UpdateResult, error)

	// DeleteOne deletes a single document matching the filter
	DeleteOne(ctx context.Context, filter D, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error)

	// DeleteMany deletes multiple documents matching the filter
	DeleteMany(ctx context.Context, filter D, opts ...*option.DeleteOptions) (*mongo.DeleteResult, error)

	// CountDocuments returns the count of documents matching the filter
	CountDocuments(ctx context.Context, filter D, opts ...*option.CountOptions) (int64, error)

	// Watch returns a change stream for watching changes to the collection
	Watch(ctx context.Context, filter *filter, opts ...*option.ChangeStreamOptions) (*mongo.ChangeStream, error)

	// Collection returns the underlying MongoDB collection
	Collection() *mongo.Collection
}

// Aggregate executes an aggregation pipeline on the collection
// If context is nil, uses the collection's default context
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

// FindOne returns a single document that matches the filter
// If context is nil, uses the collection's default context
func (c collection) FindOne(ctx context.Context, filter D, opts ...*option.FindOneOptions) *mongo.SingleResult {
	if ctx == nil {
		ctx = c.ctx
	}
	return c.coll.FindOne(ctx, filter, opts...)
}

// FindOneAndUpdate finds a single document and updates it, returning the original
// If context is nil, uses the collection's default context
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

// InsertOne inserts a single document into the collection
// If context is nil, uses the collection's default context
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

// InsertMany inserts multiple documents into the collection
// If context is nil, uses the collection's default context
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

// UpdateOne updates a single document matching the filter
// If context is nil, uses the collection's default context
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

// UpdateMany updates multiple documents matching the filter
// If context is nil, uses the collection's default context
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

// DeleteOne deletes a single document matching the filter
// If context is nil, uses the collection's default context
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

// DeleteMany deletes multiple documents matching the filter
// If context is nil, uses the collection's default context
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

// CountDocuments returns the count of documents matching the filter
// If context is nil, uses the collection's default context
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

// Watch returns a change stream for watching changes to the collection
// If context is nil, uses the collection's default context
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

// Collection returns the underlying MongoDB collection
func (c collection) Collection() *mongo.Collection {
	return c.coll
}
