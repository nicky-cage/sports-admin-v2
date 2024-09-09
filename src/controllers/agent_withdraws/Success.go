package agent_withdraws

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentWithdraws) Success(c *gin.Context) { //成功
	idStr, exists := c.GetQuery("id")
	if !exists || idStr == "" {
		response.Err(c, "无法获取id信息!\n")
		return
	}
	sql := "select a.*,b.vip,b.agent_type from user_withdraws a left join users b on a.user_id=b.id where a.id=" + idStr
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	data, err := dbSession.QueryString(sql)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	viewData := pongo2.Context{"r": data[0]}
	response.Render(c, "agent_withdraws/success.html", viewData)
}
