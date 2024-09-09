package agent_logs

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

var ActionEsList = &base_controller.ActionEsList{
	ModelEs: models.EsLoginLogsies,
	Row: func() interface{} {
		return &models.EsLoginLogs{}
	},
	QueryCond: map[string]interface{}{
		"username": "=",
		"login_ip": "=",
	},
	GetQueryCond: func(c *gin.Context) *elastic.BoolQuery { //默认只查当月
		cond := elastic.NewBoolQuery()
		var start_at int64
		var end_at int64
		if value, exists := c.GetQuery("created"); exists {
			areas := strings.Split(value, " - ")
			start_at = tools.GetTimeStampByString(areas[0] + " 0:00:00")
			end_at = tools.GetTimeStampByString(areas[1] + " 0:00:00")
			if start_at != end_at {
				cond.Filter(elastic.NewRangeQuery("created_at").Gte(start_at).Lte(end_at))
			} else {
				cond.Filter(elastic.NewRangeQuery("created_at").Gte(start_at))
			}
		}

		return cond
	},
	ViewFile: "agent_logs/_logins.html",
}
