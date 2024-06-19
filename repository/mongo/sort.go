package mongo

import (
	sort2 "github.com/aomi-go/data/common/sort"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetSortOpts 获取排序选项
// @param sortStr 排序字符串: "id,desc"
func GetSortOpts(s sort2.Sort) *options.FindOptions {
	var sort = bson.D{}
	for _, order := range s.GetOrders() {
		value := 1
		if order.Direction == sort2.DESC {
			value = -1
		}
		sort = append(sort, bson.E{Key: order.Property, Value: value})
	}

	return options.Find().SetSort(sort)
}
