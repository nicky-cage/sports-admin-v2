package agent_commissions

import (
	"sports-admin/controllers/base_controller"
)

type AgentCommissions struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionSave
}
