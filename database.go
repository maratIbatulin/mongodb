package mongo

import (
	"context"
	"fmt"
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

func (d *DB) Transaction() *Tx {
	sess, _ := d.client.StartSession(option.Session())
	err := sess.StartTransaction(option.Transaction().SetWriteConcern(&writeconcern.WriteConcern{W: 1}).SetReadConcern(readconcern.Snapshot()))
	fmt.Println(err)
	return &Tx{
		db:   d.db,
		sess: sess,
		ctx:  mongo.NewSessionContext(context.TODO(), sess),
	}
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
