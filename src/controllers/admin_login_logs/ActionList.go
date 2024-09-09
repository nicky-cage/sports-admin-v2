package admin_login_logs

import (
	"sports-admin/controllers"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var actionList = &controllers.ActionList{
	Model: models.AdminLoginLogs,
	Rows: func() interface{} {
		return &[]models.AdminLoginLog{}
	},
	OrderBy: func(c *gin.Context) string {
		return "id DESC"
	},
	ViewFile: "admin_login_logs/_list.html",
	QueryCond: map[string]interface{}{
		"ip":         "%",
		"admin_name": "%",
	},
	GetQueryCond: func(c *gin.Context) builder.Cond {
		var cond = builder.NewCond()
		var timeStart int64 = 0
		var timeEnd int64 = 0
		sArr := strings.Split(c.DefaultQuery("created", ""), " - ")
		if len(sArr) < 2 {
			timeStart = tools.GetTodayBegin()
			timeEnd = tools.GetTodayEnd()
		} else {
			timeStart = tools.GetDayBegin(sArr[0])
			timeEnd = tools.GetDayEnd(sArr[1])
		}

		cond = cond.And(builder.Gte{"created": tools.SecondToMicro(timeStart)})
		cond = cond.And(builder.Lte{"created": tools.SecondToMicro(timeEnd)})
		return cond
	},
}
