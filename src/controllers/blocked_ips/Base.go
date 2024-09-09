package blocked_ips

import "sports-admin/controllers/base_controller"

type BlockedIps struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
}
