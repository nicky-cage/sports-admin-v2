package user_codes

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionList = &base_controller.ActionList{
	RequireParameters: true,
	Model:             models.UserCodes,
	ViewFile:          "user_codes/list.html",
	QueryCond: map[string]interface{}{
		"username": "%",
	},
	GetQueryCond: func(c *gin.Context) builder.Cond {
		cond := builder.NewCond()
		if _, exists := c.GetQuery("username"); !exists {
			cond = cond.And(builder.Eq{"username": "nihao"})
		}
		return cond
	},
	Rows: func() interface{} {
		return &[]models.UserCode{}
	},
	OrderBy: func(C *gin.Context) string {
		return "created DESC"
	},
}
