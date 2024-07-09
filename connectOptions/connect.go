package connect

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

// Options create default connect options
func Options() Opt {
	opt := &option{}
	return opt
}

func (o *option) Default(appName string) *option {
	timeout := 10 * time.Second
	zLevel := 1
	*o = option{
		appName:                &appName,
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
		zlibLevel:              &zLevel,
		zstdLevel:              &zLevel,
	}
	return o
}

// Auth add auth credential
func (o *option) Auth(auth Auth) *option {
	o.auth = &options.Credential{AuthSource: auth.DB, Username: auth.UserName, Password: auth.Password}
	return o
}

// Hosts set hosts of mongod instances
func (o *option) Hosts(h []string) *option {
	o.hosts = h
	return o
}

// Pools change number of min and max pool of connections
func (o *option) Pools(min, max uint64) *option {
	o.minPoolSize = min
	o.maxPoolSize = max
	return o
}

// Replica set replication name
func (o *option) Replica(name string) *option {
	o.replicaSet = &name
	return o
}

// RetryWrites retry write operations
func (o *option) RetryWrites(rw bool) *option {
	o.retryWrites = rw
	return o
}

// RetryReads retry read operations
func (o *option) RetryReads(rr bool) *option {
	o.retryReads = rr
	return o
}

// ServerTimeout time to find available server to execute operation
func (o *option) ServerTimeout(dur time.Duration) *option {
	o.serverSelectionTimeout = &dur
	return o
}

// Timeout time to execute operation
func (o *option) Timeout(dur time.Duration) *option {
	o.timeout = &dur
	return o
}

// Acknowledged is all mongo request will wait answer from mongod instance that operation was success or error
func (o *option) Acknowledged(acknow bool) *option {
	switch acknow {
	case false:
		o.writeConcern = writeconcern.Unacknowledged()
	case true:
		o.writeConcern = writeconcern.W1()
	}
	return o
}

func (o *option) ReadConcern(level int) *option {
	switch level {
	case 1:
		o.readConcern = readconcern.Available()
	case 2:
		o.readConcern = readconcern.Majority()
	}
	return o
}

func (o *option) Connect(database string) (*mongo.Database, *mongo.Client, error) {
	clOps := &options.ClientOptions{
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
		ZlibLevel:    o.zlibLevel,
		ZstdLevel:    o.zstdLevel,
	}

	client, err := mongo.Connect(context.TODO(), clOps)
	if err != nil {
		return nil, nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, nil, err
	}

	db := client.Database(database)

	return db, client, err
}
