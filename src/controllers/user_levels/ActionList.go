package user_levels

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.UserLevels,
	ViewFile: "user_levels/index.html",
	Rows: func() interface{} {
		return &[]models.UserLevel{}
	},
}
