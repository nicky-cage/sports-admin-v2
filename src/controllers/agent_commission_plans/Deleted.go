package agent_commission_plans

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissionPlans) Deleted(c *gin.Context) {
	plan := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	ssql := "select user_id from agent_commission_plans where agent_commission ='" + plan + "' GROUP BY user_id "
	res, serr := dbSession.QueryString(ssql)
	if serr != nil {
		log.Err(serr.Error())
		return
	}
	if res[0]["user_id"] != "0" {
		response.Err(c, "方案已被使用，不能删除")
		return
	}
	sql := "delete from agent_commission_plans where agent_commission like '%" + plan + "%'"
	_, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		return
	}
	response.Ok(c)
}
