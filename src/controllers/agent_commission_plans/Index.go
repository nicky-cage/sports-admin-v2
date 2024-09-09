package agent_commission_plans

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissionPlans) Index(c *gin.Context) {
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select agent_commission,created,creat_admin, " +
		"(select min(rate) from agent_commission_plans b where b.agent_commission=a.agent_commission and user_id=0 ) as min, (select max(rate) from agent_commission_plans b where b.agent_commission=a.agent_commission and user_id=0) as max " +
		"from agent_commission_plans a " +
		"group by agent_commission "
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		return
	}
	base_controller.SetLoginAdmin(c)
	response.Render(c, "agents/_commission_plans.html", pongo2.Context{"rows": res})
}
