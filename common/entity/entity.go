package entity

type Entity interface {
	GetId() interface{}
	GetStrId() string
	IdIsNil() bool
	SetId(id interface{})
}
