package agent_commissions

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.AgentCommissionLogs,
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		dataTemp, _ := strconv.ParseFloat((*m)["commission_adjust"].(string), 64)
		(*m)["commission_adjust"] = tools.ToFixed(dataTemp, 0)
		return nil
	},
	//SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
	//	_, _ = models.AgentCommissionAdjusts.Create(*m)
	//},
}
