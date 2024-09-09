package agent_commission_plans

import (
	"github.com/gin-gonic/gin"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
)

func (ths *AgentCommissionPlans) Lists(c *gin.Context) {
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select * from agent_commission_plans"
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		response.Result(c, "获取佣金方案错误")
		return
	}
	response.Result(c, res)
}
