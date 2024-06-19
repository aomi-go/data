package mongo

import (
	"context"
	"github.com/aomi-go/data/common/page"
	"github.com/aomi-go/data/common/sort"
	"github.com/aomi-go/data/repository"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository[Entity interface{}] interface {
	repository.CrudRepository[Entity]
	SaveMany(ctx context.Context, entities []*Entity) ([]*Entity, error)
	// Find 根据条件查询数据
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]*Entity, error)

	// FindOne 查找单挑数据
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*Entity, error)

	FindOneAndModify(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*Entity, error)

	// Count 统计数据
	Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)

	// Exist 存在
	Exist(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (bool, error)

	// QueryWithPage 分页查询
	QueryWithPage(ctx context.Context, filter interface{}, pageable *page.Pageable) (*page.Page[Entity], error)

	// QueryWithSort 排序查询
	QueryWithSort(ctx context.Context, filter interface{}, sort *sort.Sort) ([]*Entity, error)
}
