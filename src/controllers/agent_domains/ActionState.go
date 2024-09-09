package agent_domains

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionState = &base_controller.ActionState{
	Model: models.AgentDomains,
	StateAfter: func(c *gin.Context) { //修改状态后处理
	},
}
