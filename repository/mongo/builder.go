package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueryBuilder struct {
	filter bson.M
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{filter: bson.M{}}
}
func QueryBuilderOf(filter bson.M) *QueryBuilder {
	return &QueryBuilder{filter: filter}
}

// Is 构建等于查询条件
func (b *QueryBuilder) Is(field string, value interface{}) *QueryBuilder {
	b.filter[field] = value
	return b
}

// Ne 构建不等于查询��件
func (b *QueryBuilder) Ne(field string, value interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$ne": value}
	return b
}

// Gt 构建大于查询条件
func (b *QueryBuilder) Gt(field string, value interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$gt": value}
	return b
}

// Lt 构建小于查询条件
func (b *QueryBuilder) Lt(field string, value interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$lt": value}
	return b
}

// Gte 构建大于等于查询条件
func (b *QueryBuilder) Gte(field string, value interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$gte": value}
	return b
}

// Lte 构建小于等于查询条件
func (b *QueryBuilder) Lte(field string, value interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$lte": value}
	return b
}

// Like 构建模糊查询条件
func (b *QueryBuilder) Like(field string, value string) *QueryBuilder {
	b.filter[field] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + value + ".*", Options: "i"}}
	return b
}

// LeftLike 构建左模糊查询条件
func (b *QueryBuilder) LeftLike(field string, value string) *QueryBuilder {
	b.filter[field] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + value, Options: "i"}}
	return b
}

// RightLike 构建右模糊查询条件
func (b *QueryBuilder) RightLike(field string, value string) *QueryBuilder {
	b.filter[field] = bson.M{"$regex": primitive.Regex{Pattern: value + ".*", Options: "i"}}
	return b
}

// Between 构建区间查询条件
func (b *QueryBuilder) Between(field string, min interface{}, max interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$gte": min, "$lte": max}
	return b
}

// Exists 构建存在查询条件
func (b *QueryBuilder) Exists(field string, exists bool) *QueryBuilder {
	b.filter[field] = bson.M{"$exists": exists}
	return b
}

func (b *QueryBuilder) In(field string, values ...interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$in": values}
	return b
}

func (b *QueryBuilder) NotIn(field string, values ...interface{}) *QueryBuilder {
	b.filter[field] = bson.M{"$nin": values}
	return b
}

func (b *QueryBuilder) Where(field string, value interface{}) *QueryBuilder {
	b.filter[field] = bson.M{field: value}
	return b
}

// And 构建AND查询条件
func (b *QueryBuilder) And(conditions ...interface{}) *QueryBuilder {
	b.filter["$and"] = conditions
	return b
}

// Or 构建OR查询条件
func (b *QueryBuilder) Or(conditions ...interface{}) *QueryBuilder {
	b.filter["$or"] = conditions
	return b
}

func (b *QueryBuilder) Build() interface{} {
	return b.filter
}
