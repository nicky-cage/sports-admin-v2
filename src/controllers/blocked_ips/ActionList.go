package blocked_ips

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.BlockedIps,
	ViewFile: "blacklists/ips.html",
	Rows: func() interface{} {
		return &[]models.BlockedIp{}
	},
	QueryCond: map[string]interface{}{
		"ip":         "%",
		"user_names": "%",
	},
}
