package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Database defines the interface for database operations
type Database interface {
	// Collection returns a collection instance for the given name
	Collection(name string) collection
	// Transaction starts a new transaction and returns a transaction object
	Transaction(ctx context.Context) (*Tx, error)
	// Ping verifies a connection to the database is still alive
	Ping(ctx context.Context, timeout time.Duration) error
	// Disconnect closes the connection to the database
	Disconnect(ctx context.Context) error
}

// DB represents a MongoDB database connection
type DB struct {
	db     *mongo.Database
	client *mongo.Client
}

// Collection returns a collection instance for the given name
func (d *DB) Collection(name string) collection {
	return collection{d.db.Collection(name), context.TODO()}
}

// Ping verifies a connection to the database is still alive
// It accepts a timeout duration and returns an error if the connection cannot be established within that time
func (d *DB) Ping(ctx context.Context, timeout time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}

	pingCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return d.client.Ping(pingCtx, nil)
}

// Disconnect closes the connection to the database
// It should be called when the application is shutting down to release resources
func (d *DB) Disconnect(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return d.client.Disconnect(ctx)
}

// Tx represents a MongoDB transaction
type Tx struct {
	db   *mongo.Database
	sess mongo.Session
	ctx  context.Context
}

// Collection returns a collection instance within the transaction context
func (tx *Tx) Collection(name string) collection {
	return collection{tx.db.Collection(name), tx.ctx}
}

// Context returns the transaction context
func (tx *Tx) Context() context.Context {
	return tx.ctx
}

// Transaction starts a new MongoDB transaction
// It returns a transaction object that can be used to perform operations within the transaction
func (d *DB) Transaction(ctx context.Context) (*Tx, error) {
	sess, _ := d.client.StartSession(option.Session())
	err := sess.StartTransaction(option.Transaction().SetWriteConcern(&writeconcern.WriteConcern{W: 1}).SetReadConcern(readconcern.Majority()))
	return &Tx{
		db:   d.db,
		sess: sess,
		ctx:  mongo.NewSessionContext(ctx, sess),
	}, err
}

// Rollback aborts the transaction and ends the session
func (tx *Tx) Rollback() error {
	err := tx.sess.AbortTransaction(tx.ctx)
	tx.sess.EndSession(tx.ctx)
	return err
}

// Commit commits the transaction and ends the session
func (tx *Tx) Commit() error {
	err := tx.sess.CommitTransaction(tx.ctx)
	tx.sess.EndSession(tx.ctx)
	return err
}
