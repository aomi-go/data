package mongo

import (
	"context"
	"fmt"
	cmongo "github.com/aomi-go/data/common/entity/mongo"
	"github.com/aomi-go/data/common/page"
	"github.com/aomi-go/data/common/sort"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

type User struct {
	*cmongo.AbstractEntity `bson:",inline"`

	Name string `bson:"name"`
}

func TestSave(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://admin:admin@127.0.0.1:27017/?authSource=admin"))
	if nil != err {
		return
	}

	repository := NewDocumentRepositoryWithEntity[User](client.Database("crypto"), User{})

	save, err := repository.Save(context.TODO(), &User{Name: "test"})
	if nil != err {
		return
	}
	fmt.Println(save.ID)

	var users []*User

	for i := 0; i < 10; i++ {
		users = append(users, &User{Name: fmt.Sprintf("test-%d", i)})
	}
	newUsers, err := repository.SaveMany(context.TODO(), users)

	fmt.Println(newUsers)
}

func TestQueryPage(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://admin:admin@127.0.0.1:27017/?authSource=admin"))
	if nil != err {
		return
	}

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
