package controllers

import (
	models "sports-models"
)

// Platforms 平台列表
var Platforms = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionState
}{
	ActionList: &ActionList{
		Model:    models.Platforms,
		ViewFile: "platforms/list.html",
		Rows: func() interface{} {
			return &[]models.Platform{}
		},
		QueryCond: map[string]interface{}{
			"name":   "%",
			"remark": "%",
			"code":   "%",
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Platforms,
		ViewFile: "platforms/edit.html",
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Platforms,
		ViewFile: "platforms/edit.html",
		Row: func() interface{} {
			return &models.Platform{}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Platforms,
	},
	ActionState: &ActionState{
		Model: models.Platforms,
		Field: "status",
	},
}
