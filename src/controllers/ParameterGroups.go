package controllers

import (
	models "sports-models"
)

// ParameterGroups 参数分组
var ParameterGroups = struct {
	*ActionList
	*ActionCreate
	*ActionSave
}{
	ActionList: &ActionList{
		Model:    models.ParameterGroups,
		ViewFile: "parameter_groups/list.html",
		Rows: func() interface{} {
			return &[]models.ParameterGroup{}
		},
		QueryCond: map[string]interface{}{
			"title": "%",
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.ParameterGroups,
		ViewFile: "parameter_groups/add.html",
	},
	ActionSave: &ActionSave{
		Model: models.ParameterGroups,
	},
}
