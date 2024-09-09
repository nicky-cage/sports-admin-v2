package user_tags

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.UserTags,
	ViewFile: "user_tags/list.html",
	Rows: func() interface{} {
		return &[]models.UserTag{}
	},
	QueryCond: map[string]interface{}{
		"name":   "%",
		"remark": "%",
	},
}
