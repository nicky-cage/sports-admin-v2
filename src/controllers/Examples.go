package controllers

import (
	models "sports-models"
)

// Examples 配置信息
var Examples = struct {
	*ActionList
}{
	ActionList: &ActionList{
		Model:    models.Examples,
		ViewFile: "examples/list.html",
		Rows: func() interface{} {
			return &[]models.Example{}
		},
	},
}
