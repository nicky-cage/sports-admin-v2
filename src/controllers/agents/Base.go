package agents

import (
	"sports-admin/controllers/base_controller"
)

// Agents 代理
type Agents struct {
	*base_controller.ActionList
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionCreate
}
