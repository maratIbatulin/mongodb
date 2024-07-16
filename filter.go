package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
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
	Let          map[string]any `bson:"let,omitempty"`
	LocalField   string         `bson:"localField,omitempty"`
	ForeignField string         `bson:"foreignField,omitempty"`
	Pipeline     *filter        `bson:"pipeline,omitempty"`
}

type GeoNear struct {
	DistanceField      string         `bson:"distanceField,omitempty"`
	DistanceMultiplier float64        `bson:"distanceMultiplier,omitempty"`
	IncludeLocs        string         `bson:"includeLocs,omitempty"`
	Key                string         `bson:"key,omitempty"`
	MaxDistance        float64        `bson:"maxDistance,omitempty"`
	MinDistance        float64        `bson:"minDistance,omitempty"`
	Spherical          bool           `bson:"spherical,omitempty"`
	Near               Geo            `bson:"near"`
	Query              map[string]any `bson:"query,omitempty"`
}

type Geo struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type filter []bson.D

func Filter() *filter {
	return &filter{}
}

// AddFields спецификация дополнительных полей, которые должны быть получены в результате запроса
func (f *filter) AddFields(val map[string]interface{}) *filter {
	*f = append(*f, bson.D{{"$addFields", val}})
	return f
}

// Bucket распределение всех записей по нескольким категориям
func (f *filter) Bucket(groupBy string, boundaries []any, def any, output map[string]interface{}) *filter {
	*f = append(*f, bson.D{{"$bucket", bson.M{
		"groupBy":    groupBy,
		"boundaries": boundaries,
		"default":    def,
		"output":     output,
	}}})
	return f
}

// BucketAuto распределение всех записей по конкретно заданному числу корзин
func (f *filter) BucketAuto(groupBy string, buckets int32, output map[string]interface{}, granularity string) *filter {
	*f = append(*f, bson.D{{"$bucketAuto", map[string]any{
		"groupBy":     groupBy,
		"buckets":     buckets,
		"output":      output,
		"granularity": granularity,
	}}})
	return f
}

// Count считает количество записей и результат подсчёта выводит в выбранное поле
func (f *filter) Count(field string) *filter {
	*f = append(*f, bson.D{{"$count", field}})
	return f
}

// Densify создаёт дополнительные записи для заполнения промежутков времени и значения в это время
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
	*f = append(*f, bson.D{{"$densify", val}})
	return f
}

// Documents записи, которые не будут созданы в базе данных, однако пройдут процесс фильтрации
func (f *filter) Documents(val ...map[string]interface{}) *filter {
	*f = append(*f, bson.D{{"$documents", val}})
	return f
}

// Facet изменяет значение поля/полей при их получении на установленное в фильтре
func (f *filter) Facet(val map[string]mongo.Pipeline) *filter {
	*f = append(*f, bson.D{{"$facet", val}})
	return f
}

// Fill заполнение полей без значений в записях
func (f *filter) Fill(sortBy map[string]interface{}, output map[string]interface{}, partition ...string) *filter {
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
	*f = append(*f, bson.D{{"$fill", val}})
	return f
}

// GeoNear изменяет значение поля/полей при их получении на установленное в фильтре
func (f *filter) GeoNear(geo GeoNear) *filter {
	*f = append(*f, bson.D{{"$geoNear", geo}})
	return f
}

// GraphLookup позволяет построить граф связи двух полей с возможностью настройки глубины
func (f *filter) GraphLookup(gl GraphLookup) *filter {
	*f = append(*f, bson.D{{"$graphLookup", gl}})
	return f
}

// Group объединение записей согласно некоторому набору параметров
func (f *filter) Group(id any, fields map[string]interface{}) *filter {
	fields["_id"] = id
	*f = append(*f, bson.D{{"$group", fields}})
	return f
}

// Limit ограничение количества записей, которые будут отданы
func (f *filter) Limit(val int64) *filter {
	*f = append(*f, bson.D{{"$limit", val}})
	return f
}

// Lookup присоединение таблиц к результату запроса
func (f *filter) Lookup(val Lookup) *filter {
	*f = append(*f, bson.D{{"$lookup", val}})
	return f
}

// Match спецификация для поиска записей согласно заданному фильтру
func (f *filter) Match(val map[string]interface{}) *filter {
	*f = append(*f, bson.D{{"$match", val}})
	return f
}

// Merge объединение записей из одной коллекции с записями из другой коллекции
func (f *filter) Merge(db string, coll string, onMatch any, notMatch string, let map[string]interface{}, on ...string) *filter {
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
	*f = append(*f, bson.D{{"$merge", val}})
	return f
}

// Out сохраняет результат запроса в отдельную коллекцию mongo.(ЕСЛИ КОЛЛЕКЦИЯ УЖЕ СУЩЕсТВУЕТ, ТО ОН ЕЁ ЗАМЕНИТ)
func (f *filter) Out(db string, coll string) *filter {
	*f = append(*f, bson.D{{"$out", bson.M{
		"db":   db,
		"coll": coll,
	}}})
	return f
}

// Project спецификация полей, которые будут отображены
func (f *filter) Project(val map[string]interface{}) *filter {
	*f = append(*f, bson.D{{"$project", val}})
	return f
}

// ReplaceRoot смена корня записи документа
func (f *filter) ReplaceRoot(val string) *filter {
	*f = append(*f, bson.D{{"$replaceRoot", map[string]any{"newRoot": val}}})
	return f
}

// Sample выборка случайных n записей из базы
func (f *filter) Sample(size int64) *filter {
	*f = append(*f, bson.D{{"$sample", map[string]any{
		"size": size,
	}}})
	return f
}

// Set изменяет значение поля/полей при их получении на установленное в фильтре
func (f *filter) Set(val map[string]interface{}) *filter {
	*f = append(*f, bson.D{{"$set", val}})
	return f
}

// Skip пропуск первых n записей, которые нашлись согласно фильтру
func (f *filter) Skip(val int64) *filter {
	*f = append(*f, bson.D{{"$skip", val}})
	return f
}

// Sort сортировка записей по полям в порядке убывания/возрастания
func (f *filter) Sort(val map[string]interface{}) *filter {
	*f = append(*f, bson.D{{"$sort", val}})
	return f
}

// UnionWith объединение записей из одной коллекции в другую с возможностью использования pipeline при слиянии
func (f *filter) UnionWith(coll string, pipeline map[string]interface{}) *filter {
	val := map[string]any{
		"coll": coll,
	}
	if pipeline != nil && reflect.DeepEqual(pipeline, bson.M{}) {
		val["pipeline"] = pipeline
	}
	*f = append(*f, bson.D{{"$unionWith", val}})
	return f
}

// Unset удаление из вывода значений поля или полей
func (f *filter) Unset(val ...string) *filter {
	*f = append(*f, bson.D{{"$unset", val}})
	return f
}

// Unwind разделение поля выводящего массив на несколько результатов
func (f *filter) Unwind(path string, preserveNullAndEmptyArrays bool) *filter {
	*f = append(*f, bson.D{{"$unwind", map[string]any{"path": path, "preserveNullAndEmptyArrays": preserveNullAndEmptyArrays}}})
	return f
}

// Use convert filter to mongo.Pipeline
func (f *filter) Use() mongo.Pipeline {
	return []bson.D(*f)
}

// Concat объединение двух фильтров в 1
func (f *filter) Concat(filt *filter) *filter {
	for _, val := range *filt {
		*f = append(*f, val)
	}
	return f
}
