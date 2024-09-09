package blocked_mails

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.BlockedMails,
	ViewFile: "blacklists/mails_edit.html",
}
