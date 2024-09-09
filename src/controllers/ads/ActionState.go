package ads

import (
	"sports-admin/controllers/base_controller"
	"sports-common/consts"
	"sports-common/redis"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionState = &base_controller.ActionState{
	Model: models.AdLanuches,
	StateAfter: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		redis.DeleteMatchedKeys(platform, consts.AdLaunchList+"*")
	},
}
