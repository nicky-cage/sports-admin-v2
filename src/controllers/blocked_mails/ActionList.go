package blocked_mails

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.BlockedMails,
	ViewFile: "blacklists/mails.html",
	Rows: func() interface{} {
		return &[]models.BlockedMail{}
	},
	QueryCond: map[string]interface{}{
		"mail":       "%",
		"user_names": "%",
	},
}
