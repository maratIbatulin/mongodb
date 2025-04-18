package mongo

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GraphLookup represents the $graphLookup MongoDB aggregation stage configuration
// Used to perform recursive searches through a collection
type GraphLookup struct {
	From             string `bson:"from"`                 // The collection to use for the graph lookup
	StartWith        string `bson:"startWith"`            // Starting value for the recursive search
	ConnectFromField string `bson:"connectFromField"`     // Field in the from collection to connect from
	ConnectToField   string `bson:"connectToField"`       // Field in the from collection to connect to
	MaxDepth         int    `bson:"maxDepth,omitempty"`   // Optional maximum depth for the recursive search
	DepthField       string `bson:"depthField,omitempty"` // Optional field to add to each output document showing the number of connections needed
	As               string `bson:"as"`                   // Output array field name
}

// SetWindowFields represents the $setWindowFields MongoDB aggregation stage configuration
// Used for window operations processing data based on specified range or window
type SetWindowFields struct {
	PartitionBy string `bson:"partitionBy"` // Field to partition documents by
	SortBy      D      `bson:"sortBy"`      // Sort specification for the window
	Output      D      `bson:"output"`      // Output field specification
}

// Lookup represents the $lookup MongoDB aggregation stage configuration
// Used to perform a left outer join to another collection
type Lookup struct {
	From         string         `bson:"from,omitempty"`         // Collection to join with
	As           string         `bson:"as,omitempty"`           // Output array field name
	Let          map[string]any `bson:"let,omitempty"`          // Optional variables to use in the pipeline stage
	LocalField   string         `bson:"localField,omitempty"`   // Field from the input documents
	ForeignField string         `bson:"foreignField,omitempty"` // Field from the documents of the "from" collection
	Pipeline     *filter        `bson:"pipeline,omitempty"`     // Pipeline to run on the joined collection
}

// Facet is a map of field names to arrays of pipeline stages
// Used for multi-faceted aggregation
type Facet map[string][]D

// GeoNear represents the $geoNear MongoDB aggregation stage configuration
// Used to find documents near a specified geospatial point
type GeoNear struct {
	DistanceField      string         `bson:"distanceField,omitempty"`      // Output field that contains the calculated distance
	DistanceMultiplier float64        `bson:"distanceMultiplier,omitempty"` // Optional multiplier for the calculated distance
	IncludeLocs        string         `bson:"includeLocs,omitempty"`        // Optional output field for the location used to calculate the distance
	Key                string         `bson:"key,omitempty"`                // Optional index to use for the query
	MaxDistance        float64        `bson:"maxDistance,omitempty"`        // Optional maximum distance
	MinDistance        float64        `bson:"minDistance,omitempty"`        // Optional minimum distance
	Spherical          bool           `bson:"spherical,omitempty"`          // Whether to use spherical geometry
	Near               Geo            `bson:"near"`                         // The geospatial point to find documents near
	Query              map[string]any `bson:"query,omitempty"`              // Optional additional query conditions
}

// Geo represents a geospatial point for MongoDB geospatial queries
type Geo struct {
	Type        string    `bson:"type"`        // Type of the geospatial object (e.g., "Point")
	Coordinates []float64 `bson:"coordinates"` // Coordinates of the geospatial point
}

// filter represents a MongoDB aggregation pipeline
type filter []D

// Filter creates a new empty MongoDB aggregation pipeline
func Filter() *filter {
	return &filter{}
}

// AddFields adds new fields to documents
// Incorporates the functionality of $addFields stage in MongoDB
func (f *filter) AddFields(val D) *filter {
	*f = append(*f, D{{"$addFields", val}})
	return f
}

// Bucket categorizes incoming documents into groups, called buckets
// Corresponds to $bucket stage in MongoDB aggregation pipeline
func (f *filter) Bucket(groupBy string, boundaries []any, def any, output D) *filter {
	*f = append(*f, D{{"$bucket", bson.M{
		"groupBy":    groupBy,
		"boundaries": boundaries,
		"default":    def,
		"output":     output,
	}}})
	return f
}

// BucketAuto categorizes incoming documents into a specific number of groups
// Corresponds to $bucketAuto stage in MongoDB aggregation pipeline
func (f *filter) BucketAuto(groupBy string, buckets int, output D, granularity string) *filter {
	*f = append(*f, D{{"$bucketAuto", map[string]any{
		"groupBy":     groupBy,
		"buckets":     buckets,
		"output":      output,
		"granularity": granularity,
	}}})
	return f
}

// Count counts the number of documents in the pipeline stage
// Corresponds to $count stage in MongoDB aggregation pipeline
func (f *filter) Count(field string) *filter {
	*f = append(*f, D{{"$count", field}})
	return f
}

// Densify creates additional documents to fill in gaps in time ranges
// Corresponds to $densify stage in MongoDB aggregation pipeline
func (f *filter) Densify(field string, bounds []time.Time, step int, unit string, partitionByFields ...string) *filter {
	val := map[string]any{
		"field": field,
		"range": map[string]any{
			"bounds": bounds,
			"step":   step,
			"unit":   unit,
		},
	}
	if len(partitionByFields) > 0 {
		val["partitionByFields"] = partitionByFields
	}
	*f = append(*f, D{{"$densify", val}})
	return f
}

// Documents defines new documents to be passed to the next stage of the pipeline
// Corresponds to $documents stage in MongoDB aggregation pipeline
func (f *filter) Documents(val ...D) *filter {
	*f = append(*f, D{{"$documents", val}})
	return f
}

// Facet processes multiple aggregation pipelines in parallel
// Corresponds to $facet stage in MongoDB aggregation pipeline
func (f *filter) Facet(val Facet) *filter {
	*f = append(*f, D{{"$facet", val}})
	return f
}

// Fill populates missing values in documents
// Corresponds to $fill stage in MongoDB aggregation pipeline
func (f *filter) Fill(sortBy D, output D, partition ...string) *filter {
	val := map[string]any{
		"sortBy": sortBy,
		"output": output,
	}
	if len(partition) == 1 {
		val["partitionBy"] = partition[0]
	}
	if len(partition) > 1 {
		val["partitionByFields"] = partition
	}
	*f = append(*f, D{{"$fill", val}})
	return f
}

// GeoNear returns documents in order of nearest to farthest from a specified point
// Corresponds to $geoNear stage in MongoDB aggregation pipeline
func (f *filter) GeoNear(geo GeoNear) *filter {
	*f = append(*f, D{{"$geoNear", geo}})
	return f
}

// GraphLookup performs a recursive search on a collection
// Corresponds to $graphLookup stage in MongoDB aggregation pipeline
func (f *filter) GraphLookup(gl GraphLookup) *filter {
	*f = append(*f, D{{"$graphLookup", gl}})
	return f
}

// Group groups documents by a specified expression
// Corresponds to $group stage in MongoDB aggregation pipeline
func (f *filter) Group(id any, fields D) *filter {
	fields = append(fields, primitive.E{Key: "_id", Value: id})
	*f = append(*f, D{{"$group", fields}})
	return f
}

// Limit limits the number of documents passed to the next stage
// Corresponds to $limit stage in MongoDB aggregation pipeline
func (f *filter) Limit(val int) *filter {
	*f = append(*f, D{{"$limit", val}})
	return f
}

// Lookup performs a left outer join to another collection
// Corresponds to $lookup stage in MongoDB aggregation pipeline
func (f *filter) Lookup(val Lookup) *filter {
	*f = append(*f, D{{"$lookup", val}})
	return f
}

// Match filters documents to pass only those that match the specified condition
// Corresponds to $match stage in MongoDB aggregation pipeline
func (f *filter) Match(val D) *filter {
	*f = append(*f, D{{"$match", val}})
	return f
}

// Merge writes the results of the aggregation pipeline to a specified collection
// Corresponds to $merge stage in MongoDB aggregation pipeline
func (f *filter) Merge(db string, coll string, onMatch any, notMatch string, let D, on ...string) *filter {
	val := map[string]any{
		"into": map[string]any{
			"db":   db,
			"coll": coll,
		},
		"whenMatched":    onMatch,
		"whenNOtMatched": notMatch,
	}
	if len(on) > 0 {
		val["on"] = on
	}
	if (let != nil) && (!reflect.DeepEqual(let, bson.M{})) {
		val["let"] = let
	}
	*f = append(*f, D{{"$merge", val}})
	return f
}

// Out writes the results of the aggregation pipeline to a specified collection
// Corresponds to $out stage in MongoDB aggregation pipeline
// WARNING: If the collection already exists, it will be replaced
func (f *filter) Out(db string, coll string) *filter {
	*f = append(*f, D{{"$out", bson.M{
		"db":   db,
		"coll": coll,
	}}})
	return f
}

// Project reshapes the documents in the pipeline
// Corresponds to $project stage in MongoDB aggregation pipeline
func (f *filter) Project(val D) *filter {
	*f = append(*f, D{{"$project", val}})
	return f
}

// ReplaceRoot replaces the document with the specified embedded document
// Corresponds to $replaceRoot stage in MongoDB aggregation pipeline
func (f *filter) ReplaceRoot(val string) *filter {
	*f = append(*f, D{{"$replaceRoot", map[string]any{"newRoot": val}}})
	return f
}

// Sample randomly selects the specified number of documents
// Corresponds to $sample stage in MongoDB aggregation pipeline
func (f *filter) Sample(size int) *filter {
	*f = append(*f, D{{"$sample", map[string]any{
		"size": size,
	}}})
	return f
}

// Set adds new fields to documents
// Corresponds to $set stage in MongoDB aggregation pipeline (an alias for $addFields)
func (f *filter) Set(val D) *filter {
	*f = append(*f, D{{"$set", val}})
	return f
}

// SetWindowField performs operations on a window/range of documents
// Corresponds to $setWindowFields stage in MongoDB aggregation pipeline
func (f *filter) SetWindowField(val SetWindowFields) *filter {
	*f = append(*f, D{{"$setWindowFields", val}})
	return f
}

// Skip skips the specified number of documents
// Corresponds to $skip stage in MongoDB aggregation pipeline
func (f *filter) Skip(val int) *filter {
	*f = append(*f, D{{"$skip", val}})
	return f
}

// Sort reorders the documents based on the specified sort keys
// Corresponds to $sort stage in MongoDB aggregation pipeline
func (f *filter) Sort(val D) *filter {
	*f = append(*f, D{{"$sort", val}})
	return f
}

// UnionWith combines the pipeline with documents from another collection
// Corresponds to $unionWith stage in MongoDB aggregation pipeline
func (f *filter) UnionWith(coll string, pipeline D) *filter {
	val := map[string]any{
		"coll": coll,
	}
	if pipeline != nil && reflect.DeepEqual(pipeline, bson.M{}) {
		val["pipeline"] = pipeline
	}
	*f = append(*f, D{{"$unionWith", val}})
	return f
}

// Unset removes specified fields from documents
// Corresponds to $unset stage in MongoDB aggregation pipeline
func (f *filter) Unset(val ...string) *filter {
	*f = append(*f, D{{"$unset", val}})
	return f
}

// Unwind deconstructs an array field to output a document for each element
// Corresponds to $unwind stage in MongoDB aggregation pipeline
func (f *filter) Unwind(path string, preserveNullAndEmptyArrays bool) *filter {
	*f = append(*f, D{{"$unwind", map[string]any{"path": path, "preserveNullAndEmptyArrays": preserveNullAndEmptyArrays}}})
	return f
}

// Use converts the filter to a MongoDB pipeline format that can be used in aggregation operations
func (f *filter) Use() []D {
	return *f
}

// Concat combines two filters into one by appending all stages from the second filter
func (f *filter) Concat(filt *filter) *filter {
	*f = append(*f, *filt...)
	return f
}
