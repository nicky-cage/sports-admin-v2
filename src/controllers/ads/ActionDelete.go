package ads

import (
	"sports-admin/controllers/base_controller"
	"sports-common/consts"
	"sports-common/redis"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionDelete = &base_controller.ActionDelete{
	Model: models.AdLanuches,
	DeleteAfter: func(c *gin.Context, data interface{}) {
		platform := request.GetPlatform(c)
		redis.DeleteMatchedKeys(platform, consts.AdLaunchList+"*")
	},
}
