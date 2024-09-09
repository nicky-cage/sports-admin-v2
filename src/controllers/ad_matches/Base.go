package ad_matches

import "sports-admin/controllers/base_controller"

// AdMatches 轮播图
type AdMatches struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
	*base_controller.ActionState
}
