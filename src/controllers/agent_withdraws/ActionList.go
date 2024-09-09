package agent_withdraws

import (
	"sports-admin/controllers/base_controller"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionList = &base_controller.ActionList{
	Model:    models.UserWithdraws,
	ViewFile: "agents/withdraws.html",
	Rows: func() interface{} {
		return &[]models.UserWithdraw{}
	},
	QueryCond: map[string]interface{}{
		"bill_no":    "=",
		"username":   "%",
		"agent_type": "=",
		"user_id":    "=",
	},
	GetQueryCond: func(c *gin.Context) builder.Cond { //默认只查当月
		cond := builder.NewCond()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); !exists {
			currentTime := time.Now().Unix()
			startAt = currentTime - currentTime%86400 - 3600*8
			endAt = startAt + 86400
			cond = cond.And(builder.Gte{"created": tools.SecondToMicro(startAt)}).And(builder.Lt{"created": tools.SecondToMicro(endAt)})
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
			endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
			cond = cond.And(builder.Gte{"created": tools.SecondToMicro(startAt)}).And(builder.Lt{"created": tools.SecondToMicro(endAt)})
		}

		cond = cond.And(builder.Eq{"status": 1}).And(builder.Eq{"type": 2}).And(builder.Eq{"risk_admin": ""})
		return cond
	},
}
