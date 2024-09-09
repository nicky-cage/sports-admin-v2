package user_level_changes

import (
	"github.com/gin-gonic/gin"
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.UserVipLogs,
	ViewFile: "user_changes/_user_level_changes.html",
	QueryCond: map[string]interface{}{
		"username":    "%",
		"admin":       "%",
		"adjust_type": "=",
		"created":     "[_]|timestamp",
	},
	Rows: func() interface{} {
		return &[]models.UserVipLog{}
	},
	OrderBy: func(c *gin.Context) string {
		return "created DESC"
	},
}
