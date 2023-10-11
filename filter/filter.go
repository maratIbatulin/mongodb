package filter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"
)

type filter struct {
	Filter []bson.D
}

func New() *filter {
	return &filter{}
}

// AddFields спецификация дополнительных полей, которые должны быть получены в результате запроса
func (f *filter) AddFields(val bson.M) *filter {
	f.Filter = append(f.Filter, bson.D{{"$addFields", val}})
	return f
}

// Bucket распределение всех записей по нескольким категориям
func (f *filter) Bucket(groupBy string, boundaries []int, def string, output bson.M) *filter {
	f.Filter = append(f.Filter, bson.D{{"$bucket", bson.M{
		"groupBy":    groupBy,
		"boundaries": boundaries,
		"default":    def,
		"output":     output,
	}}})
	return f
}

// BucketAuto распределение всех записей по конкретно заданному числу корзин
func (f *filter) BucketAuto(groupBy string, buckets int32, output bson.M, granularity string) *filter {
	f.Filter = append(f.Filter, bson.D{{"$bucketAuto", bson.M{
		"groupBy":     groupBy,
		"buckets":     buckets,
		"output":      output,
		"granularity": granularity,
	}}})
	return f
}

// Count считает количество записей и результат подсчёта выводит в выбранное поле
func (f *filter) Count(field string) *filter {
	f.Filter = append(f.Filter, bson.D{{"$count", field}})
	return f
}

// Densify создаёт дополнительные записи для заполнения промежутков времени и значения в это время
func (f *filter) Densify(field string, bounds []time.Time, step int, unit string, partitionByFields ...string) *filter {
	val := bson.M{
		"field": field,
		"range": bson.M{
			"bounds": bounds,
			"step":   step,
			"unit":   unit,
		},
	}
	if len(partitionByFields) > 0 {
		val["partitionByFields"] = partitionByFields
	}
	f.Filter = append(f.Filter, bson.D{{"$densify", val}})
	return f
}

// Documents записи, которые не будут созданы в базе данных, однако пройдут процесс фильтрации
func (f *filter) Documents(val ...bson.M) *filter {
	f.Filter = append(f.Filter, bson.D{{"$documents", val}})
	return f
}

// Facet изменяет значение поля/полей при их получении на установленное в фильтре
func (f *filter) Facet(val map[string]mongo.Pipeline) *filter {
	f.Filter = append(f.Filter, bson.D{{"$facet", val}})
	return f
}

// Fill заполнение полей без значений в записях
func (f *filter) Fill(sortBy bson.M, output bson.M, partition ...string) *filter {
	val := bson.M{
		"sortBy": sortBy,
		"output": output,
	}
	if len(partition) == 1 {
		val["partitionBy"] = partition[0]
	}
	if len(partition) > 1 {
		val["partitionByFields"] = partition
	}
	f.Filter = append(f.Filter, bson.D{{"$fill", val}})
	return f
}

// GeoNear изменяет значение поля/полей при их получении на установленное в фильтре
func (f *filter) GeoNear(geo GeoNear) *filter {
	f.Filter = append(f.Filter, bson.D{{"$geoNear", geo}})
	return f
}

// GraphLookup позволяет построить граф связи двух полей с возможностью настройки глубины
func (f *filter) GraphLookup(gl GraphLookup) *filter {
	f.Filter = append(f.Filter, bson.D{{"$graphLookup", gl}})
	return f
}

// Group объединение записей согласно некоторому набору параметров
func (f *filter) Group(id any, fields bson.M) *filter {
	fields["_id"] = id
	f.Filter = append(f.Filter, bson.D{{"$group", fields}})
	return f
}

// Limit ограничение количества записей, которые будут отданы
func (f *filter) Limit(val int64) *filter {
	f.Filter = append(f.Filter, bson.D{{"$limit", val}})
	return f
}

// Lookup присоединение таблиц к результату запроса
func (f *filter) Lookup(val Lookup) *filter {
	f.Filter = append(f.Filter, bson.D{{"$lookup", val}})
	return f
}

// Match спецификация для поиска записей согласно заданному фильтру
func (f *filter) Match(val bson.M) *filter {
	f.Filter = append(f.Filter, bson.D{{"$match", val}})
	return f
}

// Merge объединение записей из одной коллекции с записями из другой коллекции
func (f *filter) Merge(db string, coll string, onMatch any, notMatch string, let bson.M, on ...string) *filter {
	val := bson.M{
		"into": bson.M{
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
	f.Filter = append(f.Filter, bson.D{{"$merge", val}})
	return f
}

// Out сохраняет результат запроса в отдельную коллекцию mongo.(ЕСЛИ КОЛЛЕКЦИЯ УЖЕ СУЩЕсТВУЕТ, ТО ОН ЕЁ ЗАМЕНИТ)
func (f *filter) Out(db string, coll string) *filter {
	f.Filter = append(f.Filter, bson.D{{"$out", bson.M{
		"db":   db,
		"coll": coll,
	}}})
	return f
}

// Project спецификация полей, которые будут отображены
func (f *filter) Project(val bson.M) *filter {
	f.Filter = append(f.Filter, bson.D{{"$project", val}})
	return f
}

// Sample выборка случайных n записей из базы
func (f *filter) Sample(size int64) *filter {
	f.Filter = append(f.Filter, bson.D{{"$sample", bson.M{
		"size": size,
	}}})
	return f
}

// Set изменяет значение поля/полей при их получении на установленное в фильтре
func (f *filter) Set(val bson.M) *filter {
	f.Filter = append(f.Filter, bson.D{{"$set", val}})
	return f
}

// Skip пропуск первых n записей, которые нашлись согласно фильтру
func (f *filter) Skip(val int64) *filter {
	f.Filter = append(f.Filter, bson.D{{"$skip", val}})
	return f
}

// Sort сортировка записей по полям в порядке убывания/возрастания
func (f *filter) Sort(val bson.M) *filter {
	f.Filter = append(f.Filter, bson.D{{"$sort", val}})
	return f
}

// UnionWith объединение записей из одной коллекции в другую с возможностью использования pipeline при слиянии
func (f *filter) UnionWith(coll string, pipeline bson.M) *filter {
	val := bson.M{
		"coll": coll,
	}
	if pipeline != nil && reflect.DeepEqual(pipeline, bson.M{}) {
		val["pipeline"] = pipeline
	}
	f.Filter = append(f.Filter, bson.D{{"$unionWith", val}})
	return f
}

// Unset удаление из вывода значений поля или полей
func (f *filter) Unset(val ...string) *filter {
	f.Filter = append(f.Filter, bson.D{{"$unset", val}})
	return f
}

// Unwind разделение поля выводящего массив на несколько результатов
func (f *filter) Unwind(path string, preserveNullAndEmptyArrays bool) *filter {
	f.Filter = append(f.Filter, bson.D{{"$unwind", bson.M{"path": path, "preserveNullAndEmptyArrays": preserveNullAndEmptyArrays}}})
	return f
}

// Use convert filter to mongo.Pipeline
func (f *filter) Use() mongo.Pipeline {
	return f.Filter
}
