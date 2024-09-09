package agent_domains

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model: models.AgentDomains,
	Row: func() interface{} {
		return &models.AgentDomain{}
	},
	ViewFile: "agent_domains/edit.html",
}
