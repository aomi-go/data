package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type AbstractEntity struct {
	Id string `bson:"_id,omitempty"`
}

func (a *AbstractEntity) GetId() interface{} {
	return a.Id
}

func (a *AbstractEntity) SetId(id interface{}) {
	switch i := id.(type) {
	case primitive.ObjectID:
		a.Id = i.Hex()
	case string:
		a.Id = i
	}
}
