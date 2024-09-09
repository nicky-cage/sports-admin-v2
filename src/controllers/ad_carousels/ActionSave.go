package ad_carousels

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.AdCarousels,
	SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
		var startAt int64
		var endAt int64
		if value, exists := (*data)["created"].(string); !exists {
			current_time := time.Now().Unix()
			startAt = current_time - current_time%86400
			endAt = startAt + 86400
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0])
			endAt = tools.GetTimeStampByString(areas[1])
		}
		(*data)["time_start"] = startAt
		(*data)["time_end"] = endAt
		delete((*data), "created")
		return nil
	},
	CreateAfter: func(c *gin.Context, data *map[string]interface{}) {
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		redis.Del(consts.BannerList + "0" + "1")
		redis.Del(consts.BannerList + "0" + "2")
		redis.Del(consts.BannerList + "1" + "1")
		redis.Del(consts.BannerList + "1" + "2")
		redis.Del(consts.BannerList + "2" + "1")
		redis.Del(consts.BannerList + "2" + "2")
	},
	UpdateAfter: func(c *gin.Context, data *map[string]interface{}) {
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		redis.Del(consts.BannerList + "0" + "1")
		redis.Del(consts.BannerList + "0" + "2")
		redis.Del(consts.BannerList + "1" + "1")
		redis.Del(consts.BannerList + "1" + "2")
		redis.Del(consts.BannerList + "2" + "1")
		redis.Del(consts.BannerList + "2" + "2")
	},
}
