package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type GraphLookup struct {
	From             string `bson:"from"`
	StartWith        string `bson:"startWith"`
	ConnectFromField string `bson:"connectFromField"`
	ConnectToField   string `bson:"connectToField"`
	MaxDepth         int    `bson:"maxDepth,omitempty"`
	DepthField       string `bson:"depthField,omitempty"`
	As               string `bson:"as"`
}

type Lookup struct {
	From         string         `bson:"from,omitempty"`
	As           string         `bson:"as,omitempty"`
	Let          bson.M         `bson:"let,omitempty"`
	LocalField   string         `bson:"localField,omitempty"`
	ForeignField string         `bson:"foreignField,omitempty"`
	Pipeline     mongo.Pipeline `bson:"pipeline,omitempty"`
}

type GeoNear struct {
	DistanceField      string  `bson:"distanceField,omitempty"`
	DistanceMultiplier float64 `bson:"distanceMultiplier,omitempty"`
	IncludeLocs        string  `bson:"includeLocs,omitempty"`
	Key                string  `bson:"key,omitempty"`
	MaxDistance        float64 `bson:"maxDistance,omitempty"`
	MinDistance        float64 `bson:"minDistance,omitempty"`
	Spherical          bool    `bson:"spherical,omitempty"`
	Near               Geo     `bson:"near"`
	Query              bson.M  `bson:"query,omitempty"`
}

type Geo struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type QueryFilter interface {
	AddFields(val map[string]interface{}) *filter
	Bucket(groupBy string, boundaries []any, def any, output map[string]interface{}) *filter
	BucketAuto(groupBy string, buckets int32, output map[string]interface{}, granularity string) *filter
	Count(field string) *filter
	Densify(field string, bounds []time.Time, step int, unit string, partitionByFields ...string) *filter
	Documents(val ...map[string]interface{}) *filter
	Facet(val map[string]mongo.Pipeline) *filter
	Fill(sortBy map[string]interface{}, output map[string]interface{}, partition ...string) *filter
	GeoNear(geo GeoNear) *filter
	GraphLookup(gl GraphLookup) *filter
	Group(id any, fields map[string]interface{}) *filter
	Limit(val int64) *filter
	Lookup(val Lookup) *filter
	Match(val map[string]interface{}) *filter
	Merge(db string, coll string, onMatch any, notMatch string, let map[string]interface{}, on ...string) *filter
	Out(db string, coll string) *filter
	Project(val map[string]interface{}) *filter
	ReplaceRoot(val string) *filter
	Sample(size int64) *filter
	Set(val map[string]interface{}) *filter
	Skip(val int64) *filter
	Sort(val map[string]interface{}) *filter
	UnionWith(coll string, pipeline map[string]interface{}) *filter
	Unset(val ...string) *filter
	Unwind(path string, preserveNullAndEmptyArrays bool) *filter
	Use() mongo.Pipeline
	Concat(filt QueryFilter) *filter
}
