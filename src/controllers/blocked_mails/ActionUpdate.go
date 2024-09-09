package blocked_mails

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model: models.BlockedMails,
	Row: func() interface{} {
		return &models.BlockedMail{}
	},
	ViewFile: "blacklists/mails_edit.html",
}
