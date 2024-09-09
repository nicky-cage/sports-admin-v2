package blocked_mails

import "sports-admin/controllers/base_controller"

// BlockedMails 黑名单 - 电子邮件
type BlockedMails struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
}
