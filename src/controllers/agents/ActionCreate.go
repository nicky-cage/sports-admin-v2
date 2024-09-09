package agents

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.Users,
	ViewFile: "agents/agents_add.html",
	ExtendData: func(c *gin.Context) pongo2.Context {
		sql := "select agent_commission from agent_commission_plans group by agent_commission"
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return nil
		}
		return pongo2.Context{"rows": res}
	},
}
