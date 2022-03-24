package repository

type CrudRepository[T interface{}, ID interface{}] interface {
	Save(entity *T)

	FindById(id ID) *T

	ExistsById(id ID) bool

	Count() uint

	DeleteById(id ID)
}
