package blocked_devices

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.BlockedDevices,
	ViewFile: "blacklists/devices_edit.html",
}
