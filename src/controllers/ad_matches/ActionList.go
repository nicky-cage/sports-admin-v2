package ad_matches

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionList = &base_controller.ActionList{
	Model:    models.AdMatches,
	ViewFile: "ads/matches.html",
	OrderBy: func(*gin.Context) string {
		return "id DESC"
	},
	Rows: func() interface{} {
		return &[]models.AdMatche{}
	},
	QueryCond: map[string]interface{}{
		"title":       "=",
		"device_type": "=",
		"state":       "=",
	},
	GetQueryCond: func(c *gin.Context) builder.Cond { //默认只查当月
		cond := builder.NewCond()
		if value, exists := c.GetQuery("created"); exists {
			areas := strings.Split(value, " - ")
			startAt := tools.GetTimeStampByString(areas[0])
			endAt := tools.GetTimeStampByString(areas[1])
			cond = cond.And(builder.Gte{"time_start": startAt}).And(builder.Lte{"time_end": endAt})
		}
		return cond
	},
}
