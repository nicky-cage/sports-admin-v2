package agent_withdraw_records

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionList = &base_controller.ActionList{
	Model: models.UserWithdraws,
	Rows: func() interface{} {
		return &[]models.UserWithdraw{}
	},
	QueryCond: map[string]interface{}{
		"bill_no":       "=",
		"status":        "=",
		"type":          "=",
		"risk_admin":    "%",
		"username":      "%",
		"business_type": "=",
	},
	GetQueryCond: func(c *gin.Context) builder.Cond { //默认只查当月
		cond := builder.NewCond()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); exists {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
			endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
			if areas[0] != areas[1] {
				cond = cond.And(builder.Gte{"created": tools.SecondToMicro(startAt)}).And(builder.Lt{"created": tools.SecondToMicro(endAt)})
			}
		}

		if values, exists := c.GetQuery("risk_process_at"); exists {
			areas := strings.Split(values, " - ")
			if areas[0] != areas[1] {
				cond = cond.And(builder.Gte{"risk_process_at": areas[0]}).And(builder.Lt{"risk_process_at": areas[1]})
			}
		}
		cond = cond.And(builder.Eq{"type": 2}).And(builder.Neq{"risk_admin": ""})
		return cond
	},
	ViewFile: "agents/_withdraws_records.html",
	OrderBy: func(c *gin.Context) string {
		return "created DESC"
	},
}
