package blocked_phones

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.BlockedPhones,
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		admin := base_controller.GetLoginAdmin(c)
		(*m)["admin_id"] = admin.Id
		(*m)["admin_name"] = admin.Name
		return nil
	},
}
