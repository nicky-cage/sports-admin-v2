package user_detail_commissions

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

func (ths *UserDetailCommissions) Accounts(c *gin.Context) {
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select agent_balance from accounts where user_id=" + id
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		return
	}

	c.Set("_admin", base_controller.GetLoginAdmin(c))
	response.Render(c, "users/commissions.html", response.ViewData{
		"r":  res[0]["agent_balance"],
		"id": id,
	})
}
