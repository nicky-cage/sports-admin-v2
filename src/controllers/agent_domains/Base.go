package agent_domains

import "sports-admin/controllers/base_controller"

type AgentDomains struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
	*base_controller.ActionState
}
