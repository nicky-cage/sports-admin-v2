package agents

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
	Model: models.Agents,
	QueryCond: map[string]interface{}{
		"agent_code": "=",
		"username":   "%",
		"is_ma":      "=",
	},
	ViewFile: "agents/index.html",
	Rows: func() interface{} {
		return &[]models.Agent{}
	},
	GetQueryCond: func(c *gin.Context) builder.Cond {
		cond := builder.NewCond()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); !exists {
			currentTime := time.Now().Unix()
			startAt = currentTime - currentTime%86400
			endAt = startAt + 86400
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
			endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
		}

		cond.And(builder.Gte{"created": tools.SecondToMicro(startAt)}).And(builder.Lt{"created": tools.SecondToMicro(endAt)})
		return cond
	},
}
