package ad_sponsors

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionList = &base_controller.ActionList{
	Model:    models.AdSponsors,
	ViewFile: "ads/sponsors.html",
	OrderBy: func(*gin.Context) string {
		return "id DESC"
	},
	Rows: func() interface{} {
		return &[]models.AdSponsor{}
	},
	QueryCond: map[string]interface{}{
		"title": "=",
		"state": "=",
	},
	GetQueryCond: func(c *gin.Context) builder.Cond {
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
