package ads

import "sports-admin/controllers/base_controller"

// Ads app启动
type Ads struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
	*base_controller.ActionState
}
