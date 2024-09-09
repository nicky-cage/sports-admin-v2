package blocked_phones

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.BlockedPhones,
	ViewFile: "blacklists/phones.html",
	Rows: func() interface{} {
		return &[]models.BlockedPhone{}
	},
	QueryCond: map[string]interface{}{
		"phone":      "%",
		"user_names": "%",
	},
}
