package admin_roles

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.AdminRoles,
	ViewFile: "admin_roles/list.html",
	Rows: func() interface{} {
		return &[]models.AdminRole{}
	},
}
