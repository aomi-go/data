package form

type Type string

const (
	// Text 文本
	Text Type = "text"
	// Number 数字
	Number Type = "number"
	// Radio 单选
	Radio Type = "radio"
	// Checkbox 多选
	Checkbox Type = "checkbox"
	JSON     Type = "json"
)

// Field 表单字段
type Field struct {
	// Name 字段名
	Name string `json:"name" bson:"name" description:"字段名" example:"name"`
	// Type 字段类型
	Type Type `json:"type" bson:"type" description:"字段类型" example:"text"`
	// Label 字段标签
	Label string `json:"label" bson:"label" description:"字段标签" example:"名称"`
	// Placeholder 字段占位符
	Placeholder string `json:"placeholder" bson:"placeholder" description:"字段占位符" example:"请输入名称"`
	// Required 是否必填
	Required bool `json:"required" bson:"required" description:"是否必填"`
	// Options 选项
	Options []*Option `json:"options" bson:"options" description:"选项"`
	// Error 错误信息
	Error string `json:"error" bson:"error" description:"错误信息"`
	// Help 帮助信息
	Help string `json:"help" bson:"help" description:"帮助信息"`
}
