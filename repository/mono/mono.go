package mono

import (
	"context"
	"github.com/aomi-go/data/repository"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository[Entity interface{}] interface {
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
}
