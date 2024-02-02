package collection

import (
	"context"
	"github.com/maratIbatulin/mongodb/filter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Transaction struct {
	Context mongo.SessionContext
	Error   error
}

type collection struct {
	coll   *mongo.Collection
	client *mongo.Client
}

func New(db *mongo.Database, client *mongo.Client, cl string) collection {
	return collection{db.Collection(cl), client}
}

// Transaction start with goroutine to use transaction context in different functions if needed
func (c collection) Transaction(sc chan Transaction, acknowledged bool) {
	wc := writeconcern.New(writeconcern.WMajority(), writeconcern.J(acknowledged))
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	sess, err := c.client.StartSession()
	if err != nil {
		sc <- Transaction{Error: err}
		return
	}
	defer sess.EndSession(context.Background())
	_ = mongo.WithSession(context.Background(), sess, func(sessCtx mongo.SessionContext) error {
		err = sess.StartTransaction(txnOpts)
		if err != nil {
			sc <- Transaction{Error: err}
			return err
		}
		sc <- Transaction{Context: sessCtx, Error: nil}
		result := <-sc
		if result.Error != nil {
			_ = sessCtx.AbortTransaction(sessCtx)
		} else {
			err = sessCtx.CommitTransaction(sessCtx)
			if err != nil {
				_ = sessCtx.AbortTransaction(sessCtx)
				sc <- Transaction{Error: err}
			} else {
				sc <- Transaction{Error: nil}
			}
		}
		return nil
	})
}

func (c collection) Aggregate(ctx context.Context, filter filter.QueryFilter, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	cursor, err := c.coll.Aggregate(ctx, filter.Use(), opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return cursor, nil
	}
	return cursor, err
}

// FindOne find function with limit 1
func (c collection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.coll.FindOne(ctx, filter, opts...)
}

func (c collection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := c.coll.Find(ctx, filter, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return cursor, nil
	}
	return cursor, err
}

// InsertOne default id type is primitive.ObjectID but you can insert int/string id's
func (c collection) InsertOne(ctx context.Context, body interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	insertedId, err := c.coll.InsertOne(ctx, body, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return insertedId.InsertedID, nil
	}
	return insertedId.InsertedID, err
}

// InsertMany default id type is primitive.ObjectID but you can insert int/string id's
func (c collection) InsertMany(ctx context.Context, body []interface{}, opts ...*options.InsertManyOptions) ([]interface{}, error) {
	insertedId, err := c.coll.InsertMany(ctx, body, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return insertedId.InsertedIDs, nil
	}

	return insertedId.InsertedIDs, err
}

func (c collection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	upd, err := c.coll.UpdateOne(ctx, filter, update, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return upd, nil
	}
	return upd, err
}

func (c collection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	upd, err := c.coll.UpdateMany(ctx, filter, update, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return upd, nil
	}
	return upd, err
}

func (c collection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	del, err := c.coll.DeleteOne(ctx, filter, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return del, nil
	}
	return del, err
}

func (c collection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	del, err := c.coll.DeleteMany(ctx, filter, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return del, nil
	}
	return del, err
}

func (c collection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	count, err := c.coll.CountDocuments(ctx, filter, opts...)
	if err == mongo.ErrUnacknowledgedWrite {
		return count, nil
	}
	return count, err
}

func (c collection) Collection() *mongo.Collection {
	return c.coll
}
