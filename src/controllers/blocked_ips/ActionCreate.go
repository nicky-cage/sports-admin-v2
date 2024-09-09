package blocked_ips

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.BlockedIps,
	ViewFile: "blacklists/ips_edit.html",
}
