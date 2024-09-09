package admin_roles

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionDetail = &base_controller.ActionDetail{
	Model: models.AdminRoles,
	Row: func() interface{} {
		return &models.AdminRole{}
	},
}
