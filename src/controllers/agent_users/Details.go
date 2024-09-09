package agent_users

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentUsers) Details(c *gin.Context) { //默认首页
	id, exists := c.GetQuery("id")
	if !exists {
		return
	}
	var userRow []models.User
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	err := dbSession.Table("users").Where("id=?", id).Find(&userRow)
	if err != nil {
		log.Err(err.Error())
		return
	}
	var usercRow []models.Account
	cerr := dbSession.Table("users").Where("user_id=?", id).Find(&usercRow)
	if cerr != nil {
		log.Err(cerr.Error())
		return
	}
	viewData := pongo2.Context{"r": userRow[0], "account": usercRow[0]}
	response.Render(c, "agents/users_details.html", viewData)
}
