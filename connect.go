package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type options struct {
	appName                *string
	auth                   *option.Credential
	hosts                  []string
	maxPoolSize            uint64
	minPoolSize            uint64
	poolMonitor            *event.PoolMonitor
	monitor                *event.CommandMonitor
	readConcern            *readconcern.ReadConcern
	replicaSet             *string
	retryReads             bool
	retryWrites            bool
	serverSelectionTimeout *time.Duration
	timeout                *time.Duration
	writeConcern           *writeconcern.WriteConcern
}

var timeout time.Duration = 10 * time.Second

func Options(app string) options {
	return options{
		appName:                &app,
		auth:                   nil,
		hosts:                  nil,
		maxPoolSize:            100,
		minPoolSize:            10,
		poolMonitor:            nil,
		monitor:                nil,
		readConcern:            readconcern.Available(),
		replicaSet:             nil,
		retryReads:             false,
		retryWrites:            false,
		serverSelectionTimeout: &timeout,
		timeout:                &timeout,
		writeConcern:           writeconcern.Unacknowledged(),
	}
}

// Auth add auth credential
func (o *options) Auth(db, username, password string) *options {
	o.auth = &option.Credential{AuthSource: db, Username: username, Password: password}
	return o
}

// Hosts set hosts of mongod instances
func (o *options) Hosts(h []string) *options {
	o.hosts = h
	return o
}

// Pools change number of min and max pool of connections
func (o *options) Pools(min, max uint64) *options {
	o.minPoolSize = min
	o.maxPoolSize = max
	return o
}

// Replica set replication name
func (o *options) Replica(name string) *options {
	o.replicaSet = &name
	return o
}

// RetryWrites retry write operations
func (o *options) RetryWrites(rw bool) *options {
	o.retryWrites = rw
	return o
}

// RetryReads retry read operations
func (o *options) RetryReads(rr bool) *options {
	o.retryReads = rr
	return o
}

// ServerTimeout time to find available server to execute operation
func (o *options) ServerTimeout(dur time.Duration) *options {
	o.serverSelectionTimeout = &dur
	return o
}

// Timeout time to execute operation
func (o *options) Timeout(dur time.Duration) *options {
	o.timeout = &dur
	return o
}

// Acknowledged is all mongo request will wait answer from mongod instance that operation was success or error
func (o *options) Acknowledged(acknow bool) *options {
	switch acknow {
	case false:
		o.writeConcern = writeconcern.Unacknowledged()
	case true:
		o.writeConcern = writeconcern.W1()
	}
	return o
}

func (o *options) ReadConcern(level int) *options {
	switch level {
	case 1:
		o.readConcern = readconcern.Available()
	case 2:
		o.readConcern = readconcern.Majority()
	}
	return o
}

func (o *options) Connect(database string) (DB, error) {
	clOps := &option.ClientOptions{
		AppName:      o.appName,
		Auth:         o.auth,
		Hosts:        o.hosts,
		MaxPoolSize:  &o.maxPoolSize,
		MinPoolSize:  &o.minPoolSize,
		ReadConcern:  o.readConcern,
		ReplicaSet:   o.replicaSet,
		RetryReads:   &o.retryReads,
		RetryWrites:  &o.retryWrites,
		Timeout:      o.timeout,
		WriteConcern: o.writeConcern,
	}

	db := DB{
		db:     nil,
		client: nil,
	}
	var err error

	db.client, err = mongo.Connect(context.TODO(), clOps)
	if err != nil {
		return db, err
	}

	if err = db.client.Ping(context.TODO(), nil); err != nil {
		return db, err
	}

	db.db = db.client.Database(database)

	return db, err
}
