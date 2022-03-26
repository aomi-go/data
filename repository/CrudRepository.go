package repository

type CrudRepository interface {
	Save(entity interface{}) interface{}

	FindById(id interface{}) interface{}

	ExistsById(id interface{}) bool

	DeleteById(id interface{}) bool
}
