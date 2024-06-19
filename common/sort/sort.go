package sort

import "strings"

type Direction string

const (
	ASC              = Direction("asc")
	DESC             = Direction("desc")
	DefaultDirection = ASC
)

type Order struct {
	Property  string    `form:"property" json:"property" describe:"排序字段"`
	Direction Direction `form:"direction" json:"direction" describe:"排序方向"`
}

func NewSortByStr(sortStr string) Sort {
	return Sort{
		Sort: sortStr,
	}
}

func NewSortBy(direction Direction, properties ...string) Sort {
	orders := make([]Order, len(properties))
	for i, property := range properties {
		orders[i] = Order{
			Property:  property,
			Direction: direction,
		}
	}
	return NewSort(orders...)
}

func NewSort(orders ...Order) Sort {
	return Sort{
		orders: orders,
	}
}

type Sort struct {
	Sort string `form:"sort" json:"sort" describe:"排序"`

	orders []Order
}

func (s Sort) GetOrders() []Order {
	if nil != s.orders {
		return s.orders
	}
	if "" == s.Sort {
		return make([]Order, 0)
	}
	parts := strings.Split(s.Sort, ",")
	if len(parts) != 2 {
		return make([]Order, 0)
	}
	field := parts[0]
	direction := parts[1]
	if "" == field {
		return make([]Order, 0)
	}
	// 如果 direction 等于空 或者 不等于 asc、desc 则返回默认值
	if "" == direction || (strings.ToLower(direction) != string(ASC) && strings.ToLower(direction) != string(DESC)) {
		direction = string(DefaultDirection)
	}

	var orders []Order
	orders = append(orders, Order{Property: field, Direction: Direction(direction)})
	return orders
}
