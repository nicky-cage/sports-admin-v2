package ad_sponsors

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.AdSponsors,
	SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
		var startAt int64
		var endAt int64
		if value, exists := (*data)["created"].(string); !exists {
			currentDayTime := time.Now().Format("2006-01-02")
			startAt = tools.GetTimeStampByString(currentDayTime + " 00:00:00")
			endAt = tools.GetTimeStampByString(currentDayTime + " 23:59:59")
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
}
