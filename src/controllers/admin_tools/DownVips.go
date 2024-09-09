package admin_tools

import (
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// DownVips 降级vip
func DownVips(c *gin.Context) {
	platform := request.GetPlatform(c)
	username := "" // 默认用户名称

	currentTime := tools.NowTime()
	dateStart := tools.FirstDayOfLastMon(currentTime) // 上月第1天
	dateEnd := tools.LastDayOfLastMon(currentTime)    // 上月最后一天
	if userName := c.DefaultQuery("username", ""); userName != "" {
		username = userName
	}
	// 开始日期和结束日期必须是同一年同一月
	if !func() bool {
		sArr := strings.Split(dateStart, "-")
		eArr := strings.Split(dateEnd, "-")
		if len(sArr) != 3 || len(eArr) != 3 {
			return false
		}
		startTime := tools.Unix(tools.GetTimeStampByDate(dateStart) + 86400*2)
		startDate := tools.FirstDayOfMon(startTime)
		endDate := tools.LastDayOfMon(startTime)
		if dateStart != startDate || dateEnd != endDate {
			return false
		}
		return true
	}() {
		response.Err(c, "开始日期和结束日期必须在同一个月, 且格式正确")
		return
	}

	// 分页
	page := func() int {
		pageStr := c.DefaultQuery("page", "1")
		if v, err := strconv.Atoi(pageStr); err == nil && v > 1 {
			return v
		}
		return 1
	}()

	total := 0 // 总计记录数量
	rows := models.GetShouldDownUsers(dateStart, dateEnd, platform, page, &total, username)

	viewFile := "admin_tools/down_vips.html"
	if request.IsAjax(c) {
		viewFile = "admin_tools/_down_vips.html"
	}
	response.Render(c, viewFile, response.ViewData{
		"total":      total,
		"rows":       rows,
		"rows_count": len(rows),
	})
}
