package mongoxentity

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewStrObjectId() StrObjectId {
	return StrObjectIdFromObjectId(primitive.NewObjectID())
}

func StrObjectIdFromAny(id interface{}) StrObjectId {
	if id == nil {
		return ""
	}

	switch v := id.(type) {
	case StrObjectId:
		return v
	case primitive.ObjectID:
		return StrObjectIdFromObjectId(v)
	case string:
		return StrObjectIdFromStringZero(v)
	default:
		return ""
	}
}

func StrObjectIdFromStringZero(id string) StrObjectId {
	primitiveId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		primitiveId = primitive.NilObjectID
	}

	return StrObjectIdFromObjectId(primitiveId)
}

func StrObjectIdFromString(id string) (StrObjectId, error) {
	primitiveId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	return StrObjectIdFromObjectId(primitiveId), nil
}

func StrObjectIdFromObjectId(id primitive.ObjectID) StrObjectId {
	return StrObjectId(id.Hex())
}

// StrObjectId 必须添加omitempty才能让MongoDB自动生成ID bson:"_id,omitempty"
type StrObjectId string

func (id StrObjectId) String() string {
	return string(id)
}

func (id StrObjectId) ObjectId() primitive.ObjectID {
	if r, e := primitive.ObjectIDFromHex(id.String()); e == nil {
		return r
	}
	return primitive.NilObjectID
}

func (id StrObjectId) IsZero() bool {
	return id == "" || id.ObjectId().IsZero()
}
