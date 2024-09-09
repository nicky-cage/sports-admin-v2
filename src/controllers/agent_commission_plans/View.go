package agent_commission_plans

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissionPlans) View(c *gin.Context) {
	response.Render(c, "agents/commission_v.html", pongo2.Context{})
}
