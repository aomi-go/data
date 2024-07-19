package repository

import "context"

type CrudRepository[Entity interface{}] interface {
	Save(ctx context.Context, entity *Entity) (*Entity, error)

	FindAll(ctx context.Context) ([]*Entity, error)

	FindById(ctx context.Context, id interface{}) (*Entity, error)

	ExistsById(ctx context.Context, id interface{}) (bool, error)

	DeleteById(ctx context.Context, id interface{}) (bool, error)
}
