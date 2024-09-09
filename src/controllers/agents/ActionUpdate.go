package agents

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model: models.Users,
	Row: func() interface{} {
		return &models.User{}
	},
	ViewFile: "agents/agents_update.html",
}
