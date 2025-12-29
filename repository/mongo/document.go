package mongo

import (
	"context"
	"reflect"
	"regexp"
	"strings"

	"github.com/aomi-go/data/common/page"
	"github.com/aomi-go/data/common/sort"
	"github.com/aomi-go/data/mongo/mongoxentity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDocumentRepositoryWithEntity creates a new DocumentRepository.
func NewDocumentRepositoryWithEntity[E interface{}](db *mongo.Database, emptyEntity E, collectionOpts ...*options.CollectionOptions) *DocumentRepository[E] {
	collectionName := GetCollectionName(&emptyEntity)
	return NewDocumentRepository[E](db, collectionName, collectionOpts...)
}

// NewDocumentRepository creates a new DocumentRepository.
func NewDocumentRepository[E interface{}](db *mongo.Database, collectionName string, collectionOpts ...*options.CollectionOptions) *DocumentRepository[E] {
	return &DocumentRepository[E]{
		db:             db,
		collection:     db.Collection(collectionName, collectionOpts...),
		collectionName: collectionName,
		IDFieldName:    "ID",
	}
}

type DocumentRepository[Entity interface{}] struct {
	db             *mongo.Database
	collection     *mongo.Collection
	collectionName string
	IDFieldName    string
}

func (d *DocumentRepository[Entity]) Save(ctx context.Context, entity *Entity) (*Entity, error) {

	idFieldValue, idFieldOk := d.getIdFieldValue(entity)
	idOk := false
	var id primitive.ObjectID
	if idFieldOk {
		id, idOk = d.ToObjectIdWithCheck(idFieldValue.Interface())
	}

	if idOk && !id.IsZero() {
		filter := bson.M{"_id": id}
		opts := options.Replace().SetUpsert(true) // This option will create a new document if no document matches the filter

		_, err := d.collection.ReplaceOne(ctx, filter, entity, opts)
		if err != nil {
			return nil, err
		}
		return entity, nil
	} else {
		r, err := d.collection.InsertOne(ctx, entity)
		if nil == err && idFieldOk {
			d.setIdFieldValue(idFieldValue, r.InsertedID.(primitive.ObjectID))
		}
		return entity, err
	}
}

func (d *DocumentRepository[Entity]) FindAll(ctx context.Context) ([]*Entity, error) {
	vs, err := d.Find(ctx, bson.M{})
	if err := toErr(err); nil != err {
		return nil, err
	}
	return vs, nil
}
func (d *DocumentRepository[Entity]) FindAllById(ctx context.Context, ids ...interface{}) ([]*Entity, error) {
	if len(ids) == 0 {
		return []*Entity{}, nil
	}

	var objectIds []primitive.ObjectID
	for _, id := range ids {
		if oid, ok := d.ToObjectIdWithCheck(id); ok && !oid.IsZero() {
			objectIds = append(objectIds, oid)
		}
	}

	if len(objectIds) == 0 {
		return []*Entity{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	return d.Find(ctx, filter)
}

func (d *DocumentRepository[Entity]) FindById(ctx context.Context, id interface{}) (*Entity, error) {
	var result Entity
	err := d.collection.FindOne(ctx, map[string]interface{}{"_id": d.ToObjectId(id)}).Decode(&result)
	if err := toErr(err); nil != err {
		return nil, err
	}
	return &result, err
}

func (d *DocumentRepository[Entity]) ExistsById(ctx context.Context, id interface{}) (bool, error) {
	return d.Exist(ctx, map[string]interface{}{"_id": d.ToObjectId(id)})
}
func (d *DocumentRepository[Entity]) DeleteById(ctx context.Context, id interface{}) (bool, error) {
	r, err := d.collection.DeleteOne(ctx, map[string]interface{}{"_id": d.ToObjectId(id)})
	if nil == err {
		return r.DeletedCount > 0, nil
	}
	return false, err
}

func (d *DocumentRepository[Entity]) SaveMany(ctx context.Context, entities []*Entity) ([]*Entity, error) {
	var models []mongo.WriteModel

	for _, entity := range entities {
		idFieldValue, idFieldOk := d.getIdFieldValue(entity)
		idOk := false
		var id primitive.ObjectID
		if idFieldOk {
			id, idOk = d.ToObjectIdWithCheck(idFieldValue.Interface())
		}

		if idOk {
			// 如果实体的 ID 不为空，则表示这是一个现有实体，需要更新
			filter := bson.M{"_id": id}
			model := mongo.NewReplaceOneModel().
				SetFilter(filter).
				SetReplacement(entity).
				SetUpsert(true)
			models = append(models, model)
		} else {
			// 如果实体的 ID 为空，则表示这是一个新实体，需要插入
			if idFieldOk {
				d.setIdFieldValue(idFieldValue, primitive.NewObjectID())
			}
			model := mongo.NewInsertOneModel().SetDocument(entity)
			models = append(models, model)
		}
	}

	// 执行批量写操作
	opts := options.BulkWrite()
	_, err := d.collection.BulkWrite(ctx, models, opts)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (d *DocumentRepository[Entity]) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*Entity, error) {
	cursor, err := d.collection.Find(ctx, filter, opts...)
	if err := toErr(err); nil != err {
		return nil, err
	}

	var result = make([]*Entity, 0)

	for cursor.Next(ctx) {
		var item Entity
		err := cursor.Decode(&item)
		if nil != err {
			return nil, err
		}
		result = append(result, &item)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (d *DocumentRepository[Entity]) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*Entity, error) {
	var result Entity
	err := d.collection.FindOne(ctx, filter, opts...).Decode(&result)
	if err := toErr(err); nil != err {
		return nil, err
	}
	return &result, err
}

func (d *DocumentRepository[Entity]) FindOneAndModify(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*Entity, error) {
	var result Entity
	err := d.collection.FindOneAndUpdate(ctx, filter, update, opts...).Decode(&result)
	if err := toErr(err); nil != err {
		return nil, err
	}
	return &result, err
}

func (d *DocumentRepository[Entity]) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return d.collection.CountDocuments(ctx, filter, opts...)
}

func (d *DocumentRepository[Entity]) Exist(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (bool, error) {
	if count, err := d.Count(ctx, filter, opts...); nil == err {
		return count > 0, nil
	} else {
		return false, err
	}
}

func (d *DocumentRepository[Entity]) Delete(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	r, e := d.collection.DeleteMany(ctx, filter, opts...)
	if nil != e {
		return 0, e
	}
	return r.DeletedCount, nil
}

// QueryWithPage 分页查询
func (d *DocumentRepository[Entity]) QueryWithPage(ctx context.Context, filter interface{}, pageable *page.Pageable) (*page.Page[Entity], error) {
	if nil == pageable {
		pageable = page.NewDefaultPageable()
	}
	total, err := d.Count(ctx, filter)
	if nil != err {
		return nil, err
	}

	if total == 0 {
		return page.NewPage[Entity](make([]*Entity, 0), 0, pageable), nil
	}

	pageOpts := options.Find().SetSkip(pageable.GetOffset()).SetLimit(int64(pageable.GetSize()))
	sortOpts := GetSortOpts(pageable.Sort)

	entities, err := d.Find(ctx, filter, pageOpts, sortOpts)
	if nil != err {
		return nil, err
	}

	return page.NewPage[Entity](entities, total, pageable), nil
}

// QueryWithSort 排序查询
func (d *DocumentRepository[Entity]) QueryWithSort(ctx context.Context, filter interface{}, sort *sort.Sort) ([]*Entity, error) {
	var opts *options.FindOptions
	if nil != sort {
		opts = GetSortOpts(*sort)
	}
	return d.Find(ctx, filter, opts)
}
func (d *DocumentRepository[Entity]) FindWithCursor(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	cursor, err := d.collection.Find(ctx, filter, opts...)
	if err := toErr(err); nil != err {
		return nil, err
	}
	return cursor, nil
}

func (d *DocumentRepository[Entity]) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	r, err := d.collection.UpdateMany(ctx, filter, update, opts...)
	if nil != err {
		return 0, err
	}
	return r.ModifiedCount, nil
}
func (d *DocumentRepository[Entity]) UpdateMany(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (int64, error) {
	r, err := d.collection.UpdateMany(ctx, filter, update, opts...)
	if nil != err {
		return 0, err
	}
	return r.ModifiedCount, nil
}

func (d *DocumentRepository[Entity]) convertEntitiesToInterface(entities []*Entity) []interface{} {
	result := make([]interface{}, len(entities))
	for i, e := range entities {
		result[i] = e
	}
	return result
}

func (d *DocumentRepository[Entity]) getIdFieldValue(doc *Entity) (reflect.Value, bool) {
	v := reflect.ValueOf(doc)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return reflect.Value{}, false
	}
	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	// 查找 ID 字段
	idField := elem.FieldByName("ID")
	//if !idField.IsValid() || idField.Type() != reflect.TypeOf(primitive.ObjectID{}) {
	if !idField.IsValid() {
		return reflect.Value{}, false
	}
	return idField, true
}

func (d *DocumentRepository[Entity]) setIdFieldValue(idField reflect.Value, id primitive.ObjectID) {
	if idField.Type() == reflect.TypeOf(primitive.ObjectID{}) {
		idField.Set(reflect.ValueOf(id))
	} else if idField.Kind() == reflect.String {
		if idField.Type().ConvertibleTo(reflect.TypeOf("")) {
			converted := reflect.ValueOf(id.Hex()).Convert(idField.Type())
			idField.Set(converted)
		}
	}
}

func (d *DocumentRepository[Entity]) ToObjectIdWithCheck(id interface{}) (primitive.ObjectID, bool) {
	if nil == id {
		return primitive.NilObjectID, false
	}
	switch v := id.(type) {
	case primitive.ObjectID:
		return v, true
	case string:
		if oid, err := primitive.ObjectIDFromHex(v); nil == err {
			return oid, true
		}
	case BaseObjectId:
		tmp := v.ObjectId()
		return tmp, !tmp.IsZero()
	default:
		// 使用反射处理包装类型
		val := reflect.ValueOf(id)
		if val.Kind() == reflect.String {
			strVal := val.Convert(reflect.TypeOf("")).Interface().(string)
			if oid, err := primitive.ObjectIDFromHex(strVal); err == nil {
				return oid, true
			}
		}
	}
	return primitive.NilObjectID, false
}
func (d *DocumentRepository[Entity]) ToObjectId(id interface{}) primitive.ObjectID {
	v, _ := d.ToObjectIdWithCheck(id)
	return v
}

func (d *DocumentRepository[Entity]) GetCollection() *mongo.Collection {
	return d.collection
}

// GetCollectionName returns the collection name for the given entity.
// 判断 emptyEntity 是否实现了 mongoxentity.EntityDocument 接口，如果是，则调用其 CollectionName 方法获取集合名称。(同时支持，值和指针两种方式)
// 如果不是，则使用反射获取结构体名称，并转换为 snake_case 格式。
func GetCollectionName(emptyEntity any) string {
	// 优先判断是否实现 mongoxentity.EntityDocument 接口（值类型和指针类型都支持）
	if v, ok := emptyEntity.(mongoxentity.EntityDocument); ok {
		return v.CollectionName()
	}

	entityType := reflect.TypeOf(emptyEntity)
	if entityType.Kind() == reflect.Ptr {
		entityType = entityType.Elem() // 获取指针指向的类型
	}
	structName := entityType.Name()
	return toSnakeCase(structName)
}

// toSnakeCase converts a CamelCase string to snake_case.
func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
