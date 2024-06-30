package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type AbstractEntity struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`
}

func (a *AbstractEntity) GetId() interface{} {
	return a.Id
}

func (a *AbstractEntity) SetId(id interface{}) {
	switch i := id.(type) {
	case primitive.ObjectID:
		a.Id = i
	case string:
		oid, err := primitive.ObjectIDFromHex(i)
		if err == nil {
			a.Id = oid
		}
	}
}
