package agent_users

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model:    models.Users,
	ViewFile: "agents/users_update.html",
	Row: func() interface{} {
		return &models.User{}
	},
	ExtendData: func(c *gin.Context) pongo2.Context {
		var top string
		id := c.Query("id")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select top_name from users where id= " + id
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
		}
		if len(res) > 0 {
			top = res[0]["top_name"]
		}
		return pongo2.Context{"top_name": top}
	},
}
