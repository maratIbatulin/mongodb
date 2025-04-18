// Package mongo provides a MongoDB client implementation with configuration options.
package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	option "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// options represents MongoDB connection configuration parameters.
type options struct {
	appName                *string                    // Application name identifier
	auth                   *option.Credential         // Authentication credentials
	hosts                  []string                   // MongoDB server addresses
	maxPoolSize            uint64                     // Maximum number of connections in the pool
	minPoolSize            uint64                     // Minimum number of connections in the pool
	poolMonitor            *event.PoolMonitor         // Pool event monitor
	monitor                *event.CommandMonitor      // Command execution monitor
	readConcern            *readconcern.ReadConcern   // Read concern level
	replicaSet             *string                    // Replica set name
	retryReads             bool                       // Whether to retry read operations
	retryWrites            bool                       // Whether to retry write operations
	serverSelectionTimeout *time.Duration             // Timeout for server selection
	timeout                *time.Duration             // Operation timeout
	writeConcern           *writeconcern.WriteConcern // Write concern level
}

// Default timeout duration for MongoDB operations.
var timeout time.Duration = 10 * time.Second

// Options creates a new options instance with default values and the specified application name.
// It provides fluent interface for MongoDB connection configuration.
func Options(app string) *options {
	return &options{
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

// Auth adds authentication credentials to the MongoDB connection options.
// Parameters:
//   - db: The authentication database name
//   - username: The username for authentication
//   - password: The password for authentication
//
// Returns the options instance for method chaining.
func (o *options) Auth(db, username, password string) *options {
	o.auth = &option.Credential{AuthSource: db, Username: username, Password: password}
	return o
}

// Hosts sets the MongoDB server addresses to connect to.
// Parameter:
//   - h: Slice of host addresses in the format "host:port"
//
// Returns the options instance for method chaining.
func (o *options) Hosts(h []string) *options {
	o.hosts = h
	return o
}

// Pools configures the connection pool size limits.
// Parameters:
//   - min: Minimum number of connections to maintain in the pool
//   - max: Maximum number of connections allowed in the pool
//
// Returns the options instance for method chaining.
func (o *options) Pools(min, max uint64) *options {
	o.minPoolSize = min
	o.maxPoolSize = max
	return o
}

// Replica sets the replica set name for the MongoDB deployment.
// Parameter:
//   - name: The name of the replica set
//
// Returns the options instance for method chaining.
func (o *options) Replica(name string) *options {
	o.replicaSet = &name
	return o
}

// RetryWrites configures whether write operations should be retried on certain errors.
// Parameter:
//   - rw: Boolean indicating if writes should be retried
//
// Returns the options instance for method chaining.
func (o *options) RetryWrites(rw bool) *options {
	o.retryWrites = rw
	return o
}

// RetryReads configures whether read operations should be retried on certain errors.
// Parameter:
//   - rr: Boolean indicating if reads should be retried
//
// Returns the options instance for method chaining.
func (o *options) RetryReads(rr bool) *options {
	o.retryReads = rr
	return o
}

// ServerTimeout sets the timeout for server selection operations.
// Parameter:
//   - dur: Duration for server selection timeout
//
// Returns the options instance for method chaining.
func (o *options) ServerTimeout(dur time.Duration) *options {
	o.serverSelectionTimeout = &dur
	return o
}

// Timeout sets the timeout for MongoDB operations.
// Parameter:
//   - dur: Duration for operation timeout
//
// Returns the options instance for method chaining.
func (o *options) Timeout(dur time.Duration) *options {
	o.timeout = &dur
	return o
}

// Acknowledged configures whether write operations require acknowledgment from the server.
// Parameter:
//   - acknow: If true, requires server acknowledgment (W1); if false, uses unacknowledged writes
//
// Returns the options instance for method chaining.
func (o *options) Acknowledged(acknow bool) *options {
	switch acknow {
	case false:
		o.writeConcern = writeconcern.Unacknowledged()
	case true:
		o.writeConcern = writeconcern.W1()
	}
	return o
}

// ReadConcern sets the read concern level for MongoDB read operations.
// Parameter:
//   - level: Integer representing the read concern level (1=Available, 2=Majority)
//
// Returns the options instance for method chaining.
func (o *options) ReadConcern(level int) *options {
	switch level {
	case 1:
		o.readConcern = readconcern.Available()
	case 2:
		o.readConcern = readconcern.Majority()
	}
	return o
}

// Connect establishes a connection to the MongoDB database using the configured options.
// Parameter:
//   - database: The name of the database to connect to
//
// Returns:
//   - DB: A database connection object
//   - error: Any error encountered during connection
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
