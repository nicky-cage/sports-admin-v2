package agent_users

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionDetail = &base_controller.ActionDetail{
	Model: models.Users,
	Row: func() interface{} {
		return &models.User{}
	},
}
