package blocked_phones

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model: models.BlockedPhones,
	Row: func() interface{} {
		return &models.BlockedPhone{}
	},
	ViewFile: "blacklists/phones_edit.html",
}
