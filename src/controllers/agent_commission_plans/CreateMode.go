package agent_commission_plans

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissionPlans) CreateMode(c *gin.Context) {
	mode := c.Query("mode")
	if mode == "1" { //普通
		response.Render(c, "agents/commission_normal.html", pongo2.Context{})
		return
	}
	//占成
	response.Render(c, "agents/commission_special.html", pongo2.Context{})
}
