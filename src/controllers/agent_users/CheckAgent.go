package agent_users

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

func (ths *AgentUsers) CheckAgent(c *gin.Context) {
	name := c.Query("username")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select id ,is_agent from users where username='" + name + "'"
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		response.Err(c, "代理账号不正确")
	}
	response.Result(c, res[0])
}
