package blocked_ips

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model: models.BlockedIps,
	Row: func() interface{} {
		return &models.BlockedIp{}
	},
	ViewFile: "blacklists/ips_edit.html",
}
