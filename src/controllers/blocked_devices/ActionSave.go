package blocked_devices

import (
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.BlockedDevices,
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		admin := base_controller.GetLoginAdmin(c)
		(*m)["admin_id"] = admin.Id
		(*m)["admin_name"] = admin.Name
		return nil
	},
	SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
		platform := request.GetPlatform(c)
		if v, exists := (*m)["disabled_all"]; exists && v == "1" { // 如果需要设定为全部禁用
			if v, exists := (*m)["usernames"]; exists && v != "" { // disable all user names
				_ = base_controller.BlockedUserNames(platform, v.(string)) // TODO: 记录错误信息
			}
		}
	},
}
