package entity

type Entity interface {
	GetId() interface{}
	GetStrId() string
	SetId(id interface{})
}
