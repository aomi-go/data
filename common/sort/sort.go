package sort

const (
	ASC = iota
	DESC
)

type Sort struct {
	Sort string `form:"sort" json:"sort" describe:"排序"`
}
