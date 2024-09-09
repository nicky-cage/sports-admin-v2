package user_tags

import "sports-admin/controllers/base_controller"

// UserTags 坐员标签
type UserTags struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
}
