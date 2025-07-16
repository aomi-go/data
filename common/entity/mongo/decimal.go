package mongo

import (
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

type DecimalEncoder struct {
}

func (e *DecimalEncoder) EncodeValue(context bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	if !value.IsValid() || value.Type() != reflect.TypeOf(decimal.Decimal{}) {
		return bsoncodec.ValueEncoderError{
			Name:     "DecimalEncoder",
			Types:    []reflect.Type{reflect.TypeOf(decimal.Decimal{})},
			Received: value,
		}
	}

	dec := value.Interface().(decimal.Decimal)
	d128, err := primitive.ParseDecimal128(dec.String())
	if err != nil {
		return err
	}

	return writer.WriteDecimal128(d128)
}

type DecimalDecoder struct {
}

func (d *DecimalDecoder) DecodeValue(context bsoncodec.DecodeContext, reader bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != reflect.Struct || val.Type() != reflect.TypeOf(decimal.Decimal{}) {
		return bsoncodec.ValueDecoderError{
			Name:     "DecimalDecoder",
			Kinds:    []reflect.Kind{reflect.Struct},
			Received: val,
		}
	}

	var strVal string
	var err error

	switch reader.Type() {
	case bson.TypeDecimal128:
		d128, err := reader.ReadDecimal128()
		if err != nil {
			return err
		}
		strVal = d128.String()
	case bson.TypeString:
		s, err := reader.ReadString()
		if err != nil {
			return err
		}
		strVal = strings.TrimSpace(s)
	default:
		return fmt.Errorf("received invalid type for decimal: %s", reader.Type())
	}

	dec, err := decimal.NewFromString(strVal)
	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(dec))
	return nil
}
