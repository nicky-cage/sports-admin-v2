package blocked_phones

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.BlockedPhones,
	ViewFile: "blacklists/phones_edit.html",
}
