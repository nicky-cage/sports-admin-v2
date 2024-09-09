package controllers

import (
	"sports-common/consts"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// UserAccessLogs 日志列表
var UserAccessLogs = struct {
	*ActionList
	*ActionUpdate
	*ActionSave
	*ActionState
}{
	ActionList: &ActionList{
		Model: models.UserAccessLogs,
		Rows: func() interface{} {
			return &[]models.UserAccessLogInfo{}
		},
		OrderBy: func(c *gin.Context) string {
			return "user_access_logs.id DESC"
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			var cond builder.Cond = builder.NewCond()
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

			cond = cond.And(builder.Gte{"user_access_logs.created": timeStart})
			cond = cond.And(builder.Lte{"user_access_logs.created": timeEnd})
			return cond
		},
		ViewFile: "user_access_logs/list.html",
		QueryCond: map[string]interface{}{
			"ip":        "%",
			"user_name": "%",
			"method":    "=",
		},
		ProcessRow: func(c *gin.Context, rows interface{}) {
			rs := rows.(*[]models.UserAccessLogInfo)
			for k, r := range *rs {
				// 修改注册域名
				rArr := strings.Split(r.RegisterDomain, "/")
				(*rs)[k].RegisterDomain = rArr[2]

				// 修改注册域名类型
				if registerType, err := strconv.Atoi(r.RegisterType); err == nil {
					if typeName, exists := consts.UserLoginTypes[int8(registerType)]; exists {
						(*rs)[k].RegisterType = typeName
					}
				}
			}
		},
	},
}
