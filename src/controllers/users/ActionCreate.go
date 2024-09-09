package users

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.Users,
	ViewFile: "users/create.html",
	ExtendData: func(c *gin.Context) pongo2.Context {
		cond := builder.NewCond()
		//需要获取总代
		var list []models.User
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		//总代的标准
		cond = cond.And(builder.Eq{"top_id": 0}).And(builder.Eq{"is_agent": 1})
		err := dbSession.Table("users").Where(cond).Find(&list)
		if err != nil {
			log.Err(err.Error())
		}
		return pongo2.Context{"rows": list}
	},
}
