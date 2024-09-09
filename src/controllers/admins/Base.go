package admins

import (
	"sports-admin/controllers/base_controller"
)

// Admins 系统账号管理
type Admins struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionState
	*base_controller.ActionDelete
}
