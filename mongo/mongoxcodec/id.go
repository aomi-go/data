package mongoxcodec

import (
	"fmt"
	"github.com/aomi-go/data/mongo/mongoxentity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

type StrObjectIdEncoder struct {
}

func (e *StrObjectIdEncoder) EncodeValue(context bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	if !value.IsValid() || value.Type() != reflect.TypeOf(mongoxentity.StrObjectId("")) {
		return bsoncodec.ValueEncoderError{
			Name:     "StrObjectIdEncoder",
			Types:    []reflect.Type{reflect.TypeOf(mongoxentity.StrObjectId(""))},
			Received: value,
		}
	}

	dec := value.Interface().(mongoxentity.StrObjectId)
	id := dec.ObjectId()
	if id.IsZero() {
		return nil
	}
	return writer.WriteObjectID(id)
}

type StrObjectIdDecoder struct {
}

func (d *StrObjectIdDecoder) DecodeValue(context bsoncodec.DecodeContext, reader bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != reflect.String {
		return bsoncodec.ValueDecoderError{
			Name:     "StrObjectIdDecoder",
			Kinds:    []reflect.Kind{reflect.String},
			Received: val,
		}
	}

	var strVal string

	switch reader.Type() {
	case bson.TypeObjectID:
		id, err := reader.ReadObjectID()
		if err != nil {
			return err
		}
		strVal = id.Hex()
	case bson.TypeString:
		s, err := reader.ReadString()
		if err != nil {
			return err
		}
		id, err := primitive.ObjectIDFromHex(strings.TrimSpace(s))
		if err != nil {
			return err
		}
		strVal = id.Hex()
	default:
		return fmt.Errorf("received invalid type for object id: %s", reader.Type())
	}

	val.Set(reflect.ValueOf(mongoxentity.StrObjectId(strVal)))
	return nil
}
