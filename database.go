package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Database interface {
	Collection(name string) collection
}

type DB struct {
	db     *mongo.Database
	client *mongo.Client
}

func (d *DB) Collection(name string) collection {
	return collection{d.db.Collection(name), context.TODO()}
}

type Tx struct {
	db   *mongo.Database
	sess mongo.Session
	ctx  context.Context
}

func (tx *Tx) Collection(name string) collection {
	return collection{tx.db.Collection(name), tx.ctx}
}

func (tx *Tx) Context() context.Context {
	return tx.ctx
}

func (d *DB) Transaction(ctx context.Context) (*Tx, error) {
	sess, _ := d.client.StartSession(option.Session())
	err := sess.StartTransaction(option.Transaction().SetWriteConcern(&writeconcern.WriteConcern{W: 1}).SetReadConcern(readconcern.Majority()))
	return &Tx{
		db:   d.db,
		sess: sess,
		ctx:  mongo.NewSessionContext(ctx, sess),
	}, err
}

func (tx *Tx) Rollback() error {
	err := tx.sess.AbortTransaction(tx.ctx)
	tx.sess.EndSession(tx.ctx)
	return err
}

func (tx *Tx) Commit() error {
	err := tx.sess.CommitTransaction(tx.ctx)
	tx.sess.EndSession(tx.ctx)
	return err
}
