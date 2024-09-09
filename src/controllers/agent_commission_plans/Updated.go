package agent_commission_plans

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissionPlans) Updated(c *gin.Context) {
	plan := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select id,level,level_id,active_num,negative_profit,rate,agent_commission,type from agent_commission_plans where agent_commission like '%" + plan + "%' and user_id=0 order by level_id"
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		return
	}
	if res[0]["type"] == "1" {
		response.Render(c, "agents/commission_plan_v.html", pongo2.Context{"rows": res, "plan": res[0]["agent_commission"]})

	} else {
		temp, _ := strconv.ParseFloat(res[0]["rate"], 64)
		rate := tools.ToFixed(temp*100, 0)
		response.Render(c, "agents/commission_special_update.html", pongo2.Context{"rows": res[0], "plan": res[0]["agent_commission"], "rate": rate})
	}
}
