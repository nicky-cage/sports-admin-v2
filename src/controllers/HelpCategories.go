package controllers

import (
	"sports-admin/caches"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// HelpCategories 帮助分类
var HelpCategories = struct {
	*ActionCreate
	*ActionSave
}{
	ActionCreate: &ActionCreate{
		Model:    models.HelpCategories,
		ViewFile: "help_categories/edit.html",
	},
	ActionSave: &ActionSave{
		Model: models.HelpCategories,
		SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
			platform := request.GetPlatform(c)
			caches.HelpCategories.Load(platform)
		},
	},
}
