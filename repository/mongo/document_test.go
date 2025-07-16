package mongo

import (
	"context"
	"fmt"
	"github.com/aomi-go/data/common/page"
	"github.com/aomi-go/data/common/sort"
	"github.com/aomi-go/data/mongo/mongoxcodec"
	"github.com/aomi-go/data/mongo/mongoxentity"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

type User struct {
	ID mongoxentity.StrObjectId `bson:"_id,omitempty"`

	UserIdTest mongoxentity.StrObjectId `bson:"user_id,omitempty"`
	Name       string                   `bson:"name"`
}

func TestSave(t *testing.T) {
	client := getClient()

	repository := NewDocumentRepositoryWithEntity[User](client.Database("crypto"), User{})

	//id, _ := primitive.ObjectIDFromHex("68773f19dcfdef2276d06ad6")
	save, err := repository.Save(context.TODO(), &User{ID: "68773f19dcfdef2276d06ad6", UserIdTest: "68773f19dcfdef2276d06ad6", Name: "testxx888"})
	if nil != err {
		t.Errorf("Save() error = %v", err)
		return
	}
	fmt.Println(save.ID)

	var users []*User

	for i := 0; i < 10; i++ {
		users = append(users, &User{Name: fmt.Sprintf("bbbbb-%d", i)})
	}
	newUsers, err := repository.SaveMany(context.TODO(), users)

	fmt.Println(newUsers)
}

func TestQueryPage(t *testing.T) {
	client := getClient()

	repository := NewDocumentRepositoryWithEntity[User](client.Database("crypto"), User{})

	var filter map[string]interface{}
	page, err := repository.QueryWithPage(context.TODO(), filter, page.NewPageableWithSort(0, 2, sort.NewSortBy(sort.DESC, "_id")))
	if nil != err {
		return
	}

	fmt.Println(page)

	users, err := repository.QueryWithSort(context.TODO(), filter, nil)
	if nil != err {
		return
	}

	fmt.Println(users)
}

func createRegistry() *bsoncodec.Registry {

	registry := bson.NewRegistry()

	// 注册自定义 decimal.Decimal 编解码器
	decimalType := reflect.TypeOf(decimal.Decimal{})
	registry.RegisterTypeEncoder(decimalType, &mongoxcodec.DecimalEncoder{})
	registry.RegisterTypeDecoder(decimalType, &mongoxcodec.DecimalDecoder{})

	strIdType := reflect.TypeOf(mongoxentity.StrObjectId(""))
	registry.RegisterTypeEncoder(strIdType, &mongoxcodec.StrObjectIdEncoder{})
	registry.RegisterTypeDecoder(strIdType, &mongoxcodec.StrObjectIdDecoder{})

	return registry
}

func getClient() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().
		SetRegistry(createRegistry()).
		ApplyURI("mongodb://xxx:xxx@127.0.0.1:27017/xxx"))
	if nil != err {
		panic(err)
	}
	return client
}
