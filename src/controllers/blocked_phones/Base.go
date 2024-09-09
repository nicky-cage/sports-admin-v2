package blocked_phones

import "sports-admin/controllers/base_controller"

// BlockedPhones 黑名单 - 手机号码
type BlockedPhones struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
}
