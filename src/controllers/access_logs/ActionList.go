package access_logs

import (
	"regexp"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionList = &base_controller.ActionList{
	Model: models.AccessLogs,
	Rows: func() interface{} {
		return &[]models.AccessLog{}
	},
	OrderBy: func(c *gin.Context) string {
		return "id DESC"
	},
	ViewFile: "access_logs/list.html",
	QueryCond: map[string]interface{}{
		"menu_name":  "%",
		"ip":         "%",
		"admin_name": "%",
		"method":     "=",
	},
	ProcessRow: func(c *gin.Context, rows interface{}) {
		rs := rows.(*[]models.AccessLog)
		platform := request.GetPlatform(c)

		pat := "\\D(1[0-9]{5})\\D"
		reg, _ := regexp.Compile(pat)

		for k, r := range *rs {
			results := reg.FindStringSubmatch(r.Data)
			if len(results) == 0 {
				continue
			}
			for rk, rv := range results {
				if rk == 0 {
					continue
				}
				userIdStr := rv
				userId, _ := strconv.Atoi(userIdStr)
				userName := models.Users.GetUserNameById(platform, userId)
				if userName != "" {
					(*rs)[k].UserId = userId
					(*rs)[k].UserName = userName
					break
				}
			}
		}
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
		cond = cond.And(builder.Gte{"created": tools.SecondToMicro(timeStart)})
		cond = cond.And(builder.Lte{"created": tools.SecondToMicro(timeEnd)})

		userName := c.DefaultQuery("user_name", "")
		if userName != "" {
			platform := request.GetPlatform(c)
			userId := models.Users.GetIdByName(platform, userName)
			cond = cond.And(builder.Like{"data", strconv.Itoa(userId)})
		}

		return cond
	},
}
