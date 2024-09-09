package blocked_devices

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.BlockedDevices,
	ViewFile: "blacklists/devices.html",
	Rows: func() interface{} {
		return &[]models.BlockedDevice{}
	},
	QueryCond: map[string]interface{}{
		"device_no":  "%",
		"user_names": "%",
	},
}
