package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type AbstractEntity struct {
	ID string `bson:"_id,omitempty"`
}

func (a *AbstractEntity) GetId() interface{} {
	return a.ID
}

func (a *AbstractEntity) SetId(id interface{}) {
	switch i := id.(type) {
	case primitive.ObjectID:
		a.ID = i.Hex()
	case string:
		a.ID = i
	}
}
