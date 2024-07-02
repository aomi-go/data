package mongo

import (
	"context"
	ce "github.com/aomi-go/data/common/entity"
	cmongo "github.com/aomi-go/data/common/entity/mongo"
	"github.com/aomi-go/data/common/page"
	"github.com/aomi-go/data/common/sort"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

// NewDocumentRepositoryWithEntity creates a new DocumentRepository.
func NewDocumentRepositoryWithEntity[E ce.Entity](db *mongo.Database, emptyEntity any, collectionOpts ...*options.CollectionOptions) *DocumentRepository[E] {
	entityType := reflect.TypeOf(emptyEntity)
	structName := entityType.Name()
	collectionName := toCamelCase(structName)
	return NewDocumentRepository[E](db, collectionName, collectionOpts...)
}

// NewDocumentRepository creates a new DocumentRepository.
func NewDocumentRepository[E ce.Entity](db *mongo.Database, collectionName string, collectionOpts ...*options.CollectionOptions) *DocumentRepository[E] {
	return &DocumentRepository[E]{
		db:                     db,
		collection:             db.Collection(collectionName, collectionOpts...),
		collectionName:         collectionName,
		AutoConvertId2ObjectId: true,
	}
}

type DocumentRepository[Entity ce.Entity] struct {
	db                     *mongo.Database
	collection             *mongo.Collection
	collectionName         string
	AutoConvertId2ObjectId bool
}

func (d *DocumentRepository[Entity]) Save(ctx context.Context, entity *Entity) (*Entity, error) {
	initializeEntity(entity)
	if (*entity).IdIsNil() {
		(*entity).SetId(primitive.NewObjectID())
	}
	filter := bson.M{"_id": d.GetId((*entity).GetId())}
	opts := options.Replace().SetUpsert(true) // This option will create a new document if no document matches the filter

	result, err := d.collection.ReplaceOne(ctx, filter, entity, opts)
	if err != nil {
		return nil, err
	}

	if result.UpsertedID != nil {
		(*entity).SetId(result.UpsertedID)
	}

	return entity, nil
}

func (d *DocumentRepository[Entity]) FindById(ctx context.Context, id interface{}) (*Entity, error) {
	var result Entity
	err := d.collection.FindOne(ctx, map[string]interface{}{"_id": d.GetId(id)}).Decode(&result)
	return &result, err
}

func (d *DocumentRepository[Entity]) ExistsById(ctx context.Context, id interface{}) (bool, error) {
	return d.Exist(ctx, map[string]interface{}{"_id": d.GetId(id)})
}
func (d *DocumentRepository[Entity]) DeleteById(ctx context.Context, id interface{}) (bool, error) {
	r, err := d.collection.DeleteOne(ctx, map[string]interface{}{"_id": d.GetId(id)})
	if nil == err {
		return r.DeletedCount > 0, nil
	}
	return false, err
}

func (d *DocumentRepository[Entity]) SaveMany(ctx context.Context, entities []*Entity) ([]*Entity, error) {
	var models []mongo.WriteModel

	for _, entity := range entities {
		initializeEntity(entity)

		// 如果实体的 ID 为空，则表示这是一个新实体，需要插入
		if (*entity).IdIsNil() {
			(*entity).SetId(primitive.NewObjectID())
			model := mongo.NewInsertOneModel().SetDocument(entity)
			models = append(models, model)
		} else {
			// 如果实体的 ID 不为空，则表示这是一个现有实体，需要更新
			filter := bson.M{"_id": d.GetId((*entity).GetId())}
			model := mongo.NewReplaceOneModel().
				SetFilter(filter).
				SetReplacement(entity).
				SetUpsert(true)
			models = append(models, model)
		}
	}

	// 执行批量写操作
	opts := options.BulkWrite().SetOrdered(false)
	_, err := d.collection.BulkWrite(ctx, models, opts)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

func (d *DocumentRepository[Entity]) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*Entity, error) {
	cursor, err := d.collection.Find(ctx, filter, opts...)
	if nil != err {
		return nil, err
	}

	var result []*Entity

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
	return &result, err
}

func (d *DocumentRepository[Entity]) FindOneAndModify(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*Entity, error) {
	var result Entity
	err := d.collection.FindOneAndUpdate(ctx, filter, update, opts...).Decode(&result)
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
		return page.NewPage[Entity](nil, 0, pageable), nil
	}

	pageOpts := options.Find().SetSkip(pageable.GetOffset()).SetLimit(int64(pageable.Size))
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

func (d *DocumentRepository[Entity]) convertEntitiesToInterface(entities []*Entity) []interface{} {
	result := make([]interface{}, len(entities))
	for i, e := range entities {
		result[i] = e
	}
	return result
}

func (d *DocumentRepository[Entity]) GetId(id interface{}) interface{} {
	if !d.AutoConvertId2ObjectId {
		return id
	}
	if idStr, ok := id.(string); ok {
		if oid, err := primitive.ObjectIDFromHex(idStr); nil == err {
			return oid
		}
	}
	return id
}

// InitializeEntity initializes the embedded AbstractEntity pointer
func initializeEntity(entity any) {
	entityValue := reflect.ValueOf(entity).Elem()

	for i := 0; i < entityValue.NumField(); i++ {
		field := entityValue.Field(i)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			fieldType := field.Type()
			if fieldType.Elem() == reflect.TypeOf(cmongo.AbstractEntity{}) {
				field.Set(reflect.New(fieldType.Elem()))
				return
			}
		}
	}

}

func toCamelCase(name string) string {
	// 分割字符串
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	// 首字母小写
	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = strings.ToTitle(part)
		}
	}

	return strings.Join(parts, "")
}
