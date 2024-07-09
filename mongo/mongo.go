package mongo

import (
	"context"
	"fmt"
	connect "github.com/maratIbatulin/mongodb/connectOptions"
	"github.com/maratIbatulin/mongodb/filter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type DB struct {
	db     *mongo.Database
	client *mongo.Client
}

type Tx struct {
	db   *mongo.Database
	sess mongo.Session
	ctx  context.Context
}

func (d *DB) Transaction() *Tx {
	sess, _ := d.client.StartSession(options.Session())
	err := sess.StartTransaction(options.Transaction().SetWriteConcern(&writeconcern.WriteConcern{W: 1}).SetReadConcern(readconcern.Snapshot()))
	fmt.Println(err)
	return &Tx{
		db:   d.db,
		sess: sess,
		ctx:  mongo.NewSessionContext(context.TODO(), sess),
	}
}

func ConnectDB(connect connect.Opt, dbName string) (*DB, error) {
	db, client, err := connect.Connect(dbName)
	if err != nil {
		return nil, err
	}

	return &DB{db: db, client: client}, nil
}

func Connection() connect.Opt {
	return connect.Options()
}

func Filter() filter.QueryFilter {
	return filter.New()
}

func (d *DB) Collection(name string) Collection {
	return new(d.db, name, context.TODO())
}

func (tx *Tx) Collection(name string) Collection {
	return new(tx.db, name, tx.ctx)
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
