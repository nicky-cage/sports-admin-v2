package ads

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
	Model: models.AdLanuches,
	SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
		var startAt int64
		var endAt int64
		if value, exists := (*data)["created"].(string); !exists {
			currentTime := time.Now().Unix()
			startAt = currentTime - currentTime%86400
			endAt = startAt + 86400
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0])
			endAt = tools.GetTimeStampByString(areas[1])
		}
		(*data)["time_start"] = startAt
		(*data)["time_end"] = endAt
		delete((*data), "created")

		//imageAndroid := ""
		//for i := 0; i < 10; i++ {
		//	tempI := strconv.Itoa(i)
		//	v, ok := (*data)["image_android["+tempI+"]"].(string)
		//	if ok {
		//		imageAndroid += v + ","
		//	}
		//}
		//(*data)["image_android"] = strings.TrimRight(imageAndroid, ",")
		//imageIos := ""
		//for i := 0; i < 10; i++ {
		//	tempI := strconv.Itoa(i)
		//	v, ok := (*data)["image_ios["+tempI+"]"].(string)
		//	if ok {
		//		imageIos += v + ","
		//	}
		//}
		//(*data)["image_ios"] = strings.TrimRight(imageIos, ",")
		//imageIosx := ""
		//for i := 0; i < 10; i++ {
		//	tempI := strconv.Itoa(i)
		//	v, ok := (*data)["image_iosx["+tempI+"]"].(string)
		//	if ok {
		//		imageIosx += v + ","
		//	}
		//}
		//(*data)["image_iosx"] = strings.TrimRight(imageIosx, ",")
		return nil
	},
	SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
		platform := request.GetPlatform(c)
		platformType := (*m)["platform_type"].(string)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		redis.Del(consts.AdLaunchList + platformType)
	},
}
