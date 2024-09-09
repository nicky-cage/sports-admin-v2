package admins

import (
	"sports-common/request"
	models "sports-models"

	"sports-admin/caches"
	"sports-admin/controllers/base_controller"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// ActionList 数据列表
var ActionList = &base_controller.ActionList{
	Model:    models.Admins,
	ViewFile: "admins/list.html",
	Rows: func() interface{} {
		return &[]models.Admin{}
	},
	ExtendData: func(c *gin.Context) pongo2.Context {
		platform := request.GetPlatform(c)
		return pongo2.Context{
			"admin_roles": caches.AdminRoles.All(platform),
		}
	},
	ProcessRow: func(c *gin.Context, rows interface{}) {
		platform := request.GetPlatform(c)
		rs := rows.(*[]models.Admin)
		for k, v := range *rs {
			isOnline := models.LoginAdmins.IsOnline(platform, int(v.Id))
			(*rs)[k].IsOnline = isOnline
		}
	},
	GetQueryCond: func(c *gin.Context) builder.Cond {
		cond := builder.NewCond()
		cond = cond.And(builder.Neq{"name": "admin"}) // 不显示 admin 用户 - 谁都不可以使用
		return cond
	},
	QueryCond: map[string]interface{}{
		"name":    "%",
		"role_id": "=",
		"state":   "=",
	},
}
