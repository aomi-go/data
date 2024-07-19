package mongo

import (
	"errors"
	"github.com/aomi-go/data/common"
	"go.mongodb.org/mongo-driver/mongo"
)

func toErr(err error) error {
	if nil != err {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return common.ErrNoResult
		}
	}
	return err
}
