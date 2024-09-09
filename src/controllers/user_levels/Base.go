package user_levels

import "sports-admin/controllers/base_controller"

type UserLevels struct {
	*base_controller.ActionList
	*base_controller.ActionUpdate
	*base_controller.ActionSave
}
