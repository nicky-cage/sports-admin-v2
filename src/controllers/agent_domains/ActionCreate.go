package agent_domains

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.AgentDomains,
	ViewFile: "agent_domains/edit.html",
}
