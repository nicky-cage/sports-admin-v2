package activities

import (
	"encoding/json"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// Turntable 转盘
func (*Activities) Turntable(c *gin.Context) { // 活动条件判断依据
	playCon := map[int]string{
		1: "充值金额",
		2: "会员有效投注",
		3: "体育有效投注",
		4: "电竞有效投注",
		5: "电子有效投注",
		6: "真人有效投注",
		7: "棋牌有效投注",
		8: "彩票有效投注",
		9: "捕鱼有效投注",
	}
	play, _ := json.Marshal(playCon)

	// 获取转盘活动配置
	platform := request.GetPlatform(c)
	row := models.GetTurntableConfig(platform)

	turntableTimeStart := tools.GetDateTimeByTimeStamp(int64(row.TurntableStart))
	turntableTimeEnd := tools.GetDateTimeByTimeStamp(int64(row.TurntableEnd))
	showTimeStart := tools.GetDateTimeByTimeStamp(int64(row.ShowTimeStart))
	showTimeEnd := tools.GetDateTimeByTimeStamp(int64(row.ShowTimeEnd))
	viewData := response.ViewData{
		"playConJson":     string(play),
		"turntable":       row,
		"playList":        playCon,
		"turntable_start": turntableTimeStart,
		"turntable_end":   turntableTimeEnd,
		"show_time_start": showTimeStart,
		"show_time_end":   showTimeEnd,
	}
	response.Render(c, "activities/_activity_turntable.html", viewData)
}
