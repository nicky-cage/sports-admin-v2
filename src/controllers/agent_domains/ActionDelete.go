package agent_domains

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionDelete = &base_controller.ActionDelete{
	Model: models.AgentDomains,
}
