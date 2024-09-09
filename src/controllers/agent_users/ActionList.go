package agent_users

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.Users,
	ViewFile: "agents/users.html",
	QueryCond: map[string]interface{}{
		"agent_code": "=",
		"username":   "%",
		"status":     "=",
	},
	Rows: func() interface{} {
		return &[]models.User{}
	},
}
