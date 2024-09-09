package controllers

import (
	"sports-admin/caches"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// Banks 银行列表
var Banks = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionState
}{
	ActionList: &ActionList{
		Model:    models.Banks,
		ViewFile: "banks/_list.html",
		OrderBy: func(*gin.Context) string {
			return "sort DESC"
		},
		QueryCond: map[string]interface{}{
			"name": "%",
			"code": "%",
		},
		Rows: func() interface{} {
			return &[]models.Bank{}
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Banks,
		ViewFile: "banks/edit.html",
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Banks,
		ViewFile: "banks/edit.html",
		Row: func() interface{} {
			return &models.Bank{}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Banks,
		SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
			platform := request.GetPlatform(c)
			caches.Banks.Load(platform)
		},
	},
	ActionState: &ActionState{
		Model: models.Banks,
		StateAfter: func(c *gin.Context) {
			platform := request.GetPlatform(c)
			caches.Banks.Load(platform)
		},
	},
}
