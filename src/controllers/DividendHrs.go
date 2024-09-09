package controllers

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// DividendHrs 红利-历史列表
var DividendHrs = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		var start_at int64
		var end_at int64
		if value, exists := c.GetQuery("created"); !exists {
			current_time := time.Now().Unix()
			start_at = current_time - current_time%86400
			end_at = start_at + 86400
		} else {
			areas := strings.Split(value, " - ")
			start_at = tools.GetTimeStampByString(areas[0] + " 00:00:00")
			end_at = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
		}
		cond = cond.And(builder.Gte{"created": tools.SecondToMicro(start_at)}).And(builder.Lte{"created": tools.SecondToMicro(end_at)})
		var uStartAt int64
		var uEndAt int64
		if value, exists := c.GetQuery("updated"); !exists {
			current_time := time.Now().Unix()
			uStartAt = current_time - current_time%86400
			uEndAt = start_at + 86400
		} else {
			areas := strings.Split(value, " - ")
			uStartAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
			uEndAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
		}
		cond = cond.And(builder.Gte{"updated": tools.SecondToMicro(uStartAt)}).And(builder.Lte{"updated": tools.SecondToMicro(uEndAt)})
		cond = cond.And(builder.Neq{"state": 1})
		username := c.DefaultQuery("username", "")
		billNo := c.DefaultQuery("bill_no", "")
		applicant := c.DefaultQuery("applicant", "")
		sType := c.DefaultQuery("type", "")
		moneyType := c.DefaultQuery("money_type", "")
		flowLimit := c.DefaultQuery("flow_limit", "")
		reviewer := c.DefaultQuery("reviewer", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"username": username})
		}
		if len(applicant) > 0 {
			cond = cond.And(builder.Eq{"applicant": applicant})
		}
		if len(sType) > 0 {
			cond = cond.And(builder.Eq{"type": sType})
		}
		if len(billNo) > 0 {
			cond = cond.And(builder.Eq{"bill_no": billNo})
		}
		if len(moneyType) > 0 {
			cond = cond.And(builder.Eq{"money_type": moneyType})
		}
		if len(flowLimit) > 0 {
			cond = cond.And(builder.Eq{"flow_limit": flowLimit})
		}
		if len(reviewer) > 0 {
			cond = cond.And(builder.Eq{"reviewer": reviewer})
		}
		limit, offset := request.GetOffsets(c)
		userDividends := make([]models.UserDividend, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("user_dividends").Where(cond).OrderBy("id DESC").Limit(limit, offset).FindAndCount(&userDividends)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ss := new(SumStruct)
		sumTotal, _ := engine.Table("user_dividends").Where(cond).Sum(ss, "money")
		viewData := pongo2.Context{
			"rows":      userDividends,
			"total":     total,
			"sum_money": sumTotal,
		}
		viewFile := "dividend_managements/records.html"
		if request.IsAjax(c) {
			viewFile = "dividend_managements/_records.html"
		}
		response.Render(c, viewFile, viewData)
	},
}
