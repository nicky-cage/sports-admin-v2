package agents

import (
	"fmt"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"strconv"
)

func (ths *Agents) Updated(c *gin.Context) {
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select * from users where id =%s"
	sqll := fmt.Sprintf(sql, id)
	res, err := dbSession.QueryString(sqll)

	if err != nil {
		log.Err(err.Error())
		return
	}
	//所有
	commissionSql := "select agent_commission,type from  agent_commission_plans where user_id=0 group by agent_commission"
	comRes, err := dbSession.QueryString(commissionSql)
	if err != nil {
		log.Err(err.Error())
	}
	var rates float64
	var stype string

	agentCommission := "select * from agent_commission_plans where agent_commission='" + res[0]["agent_commission"] + "'"
	agentCommissionRes, _ := dbSession.QueryString(agentCommission)
	//单独
	stype = "1"
	if len(agentCommissionRes) > 0 && agentCommissionRes[0]["type"] == "2" {
		userRate := "select rate from  agent_commission_plans where user_id=? and agent_commission=? "
		rate, err := dbSession.QueryString(userRate, res[0]["id"], res[0]["agent_commission"])
		if err != nil {
			log.Err(err.Error())
		}
		if len(rate) > 0 {
			temp, _ := strconv.ParseFloat(rate[0]["rate"], 64)
			rates = temp * 100
		}
		stype = "2"
	}

	response.Render(c, "agents/agents_update.html", pongo2.Context{"r": res[0], "agent_commission": comRes, "rate": rates, "types": stype})
}
