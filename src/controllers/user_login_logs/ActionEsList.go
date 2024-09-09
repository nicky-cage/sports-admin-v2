package user_login_logs

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

var ActionEsList = &base_controller.ActionEsList{
	ModelEs:  models.EsLoginLogsies,
	ViewFile: "user_login_logs/index.html",
	/*		Rows: func() interface{} {
			return &[]models.EsLoginLogs{}
		},*/
	Row: func() interface{} {
		return &models.EsLoginLogs{}
	},
	QueryCond: map[string]interface{}{
		"username":   "=",
		"device_no":  "=",
		"top_name":   "=",
		"login_type": "=",
		"login_ip":   "=",
		"last_ip":    "=",
	},
	GetQueryCond: func(c *gin.Context) *elastic.BoolQuery { //默认只查当月
		cond := elastic.NewBoolQuery()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); exists {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0] + " 0:00:00")
			endAt = tools.GetTimeStampByString(areas[1] + " 0:00:00")
			if startAt != endAt {
				cond.Filter(elastic.NewRangeQuery("created_at").Gte(startAt).Lte(endAt))
			} else {
				cond.Filter(elastic.NewRangeQuery("created_at").Gte(startAt))
			}
		}

		return cond
	},
}
