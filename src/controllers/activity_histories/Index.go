package activity_histories

import (
	"sports-admin/controllers/base_controller"
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

func (ths *ActivityHrs) List(c *gin.Context) {
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
	cond = cond.And(builder.Gte{"created": tools.SecondToMicro(startAt)}).And(builder.Lte{"created": tools.SecondToMicro(endAt)})
	var startu_at int64
	var endu_at int64
	if value, exists := c.GetQuery("updated"); !exists {
		currentTime := time.Now().Unix()
		startu_at = currentTime - currentTime%86400
		endu_at = startAt + 86400
	} else {
		areas := strings.Split(value, " - ")
		startu_at = tools.GetTimeStampByString(areas[0] + " 00:00:00")
		endu_at = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
	}
	cond = cond.And(builder.Gte{"updated": tools.SecondToMicro(startu_at)}).And(builder.Lte{"updated": tools.SecondToMicro(endu_at)})
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
	ss := new(base_controller.SumStruct)
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
}
