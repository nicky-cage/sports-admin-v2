package admin_roles

import (
	"sports-admin/controllers/base_controller"
)

// AdminRoles 角色管理
type AdminRoles struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDetail
	*base_controller.ActionDelete
}
