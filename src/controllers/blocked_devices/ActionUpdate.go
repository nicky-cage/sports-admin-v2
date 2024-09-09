package blocked_devices

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model: models.BlockedDevices,
	Row: func() interface{} {
		return &models.BlockedDevice{}
	},
	ViewFile: "blacklists/devices_edit.html",
}
