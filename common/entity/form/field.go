package form

// Field 表单字段
type Field struct {
	// Name 字段名
	DataIndex string `json:"dataIndex" bson:"dataIndex"`
	// ValueTYpe 值类型
	ValueType string `json:"valueType" bson:"valueType"`
	// Label 字段标签
	Title string `json:"title" bson:"title" description:"字段标签" example:"名称"`
	// Placeholder 字段占位符
	Placeholder string `json:"placeholder" bson:"placeholder" description:"字段占位符" example:"请输入名称"`
	// Required 是否必填
	Required bool `json:"required" bson:"required" description:"是否必填"`
	// Options 选项
	Options []*Option `json:"options" bson:"options" description:"选项"`
	// Error 错误信息
	Error string `json:"error" bson:"error" description:"错误信息"`
	// Help 帮助信息
	Help         string `json:"help" bson:"help" description:"帮助信息"`
	InitialValue any    `json:"defaultValue" bson:"defaultValue" description:"默认值"`
}
