package user_notes

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.UserNotes,
	SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
		delete(*data, "id")
		(*data)["id"] = 0
		admin := base_controller.GetLoginAdmin(c)
		(*data)["admin_id"] = admin.Id
		(*data)["admin_name"] = admin.Name
		return nil
	},
}
