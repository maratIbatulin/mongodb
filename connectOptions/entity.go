package connect

import (
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type option struct {
	appName                *string
	auth                   *options.Credential
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
	zlibLevel              *int
	zstdLevel              *int
}

type Auth struct {
	DB       string
	UserName string
	Password string
}

type Opt interface {
	Default(appName string) *option
	Auth(auth Auth) *option
	Hosts(h []string) *option
	Pools(min, max uint64) *option
	Replica(name string) *option
	RetryWrites(rw bool) *option
	RetryReads(rr bool) *option
	ServerTimeout(dur time.Duration) *option
	Timeout(dur time.Duration) *option
	Acknowledged(acknow bool) *option
	Connect(db string) (*mongo.Database, *mongo.Client, error)
}
