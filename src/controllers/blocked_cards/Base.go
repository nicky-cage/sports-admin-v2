package blocked_cards

import "sports-admin/controllers/base_controller"

// BlockedCards 黑名单 - 银行卡号
type BlockedCards struct {
	*base_controller.ActionList
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionDelete
}
