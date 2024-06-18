package sort

const (
	ASC = iota
	DESC
)

type Sort struct {
	Sort string `json:"sort" describe:"排序"`
}
