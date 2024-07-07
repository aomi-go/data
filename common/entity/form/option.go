package form

// Option 选项
type Option struct {
	// Label 选项标签
	Label string `json:"label" bson:"label" description:"选项标签" example:"选项1"`
	// Value 选项值
	Value string `json:"value" bson:"value" description:"选项值" example:"1"`
}
