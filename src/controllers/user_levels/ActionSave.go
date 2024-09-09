package user_levels

import (
	"github.com/gin-gonic/gin"
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionSave = &base_controller.ActionSave{
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		if (*m)["venue_all_chose"].(string) == "1" && (*m)["code_list"].(string) == "" {
			(*m)["code_list"] = " "
			(*m)["venues"] = " "
		}
		return nil
	},
	Model: models.UserLevels,
}
