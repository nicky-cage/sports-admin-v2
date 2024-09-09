package ad_carousels

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionState = &base_controller.ActionState{
	Model: models.AdCarousels,
	StateAfter: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		redis.Del(consts.BannerList + "01")
		redis.Del(consts.BannerList + "11")
		redis.Del(consts.BannerList + "12")
	},
}
