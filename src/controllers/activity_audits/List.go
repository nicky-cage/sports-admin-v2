package activity_audits

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

func (ths *ActivityAudits) List(c *gin.Context) {
	cond := builder.NewCond()
	var startAt int64
	var endAt int64
	if value, exists := c.GetQuery("created"); !exists {
		current_time := time.Now().Unix()
		startAt = current_time - current_time%86400
		endAt = startAt + 86400
	} else {
		areas := strings.Split(value, " - ")
		startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
		endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
	}
	cond = cond.And(builder.Gte{"created": tools.SecondToMicro(startAt)}).And(builder.Lte{"created": tools.SecondToMicro(endAt)})
	cond = cond.And(builder.Eq{"state": 1})
	username := c.DefaultQuery("username", "")
	billNo := c.DefaultQuery("bill_no", "")
	applicant := c.DefaultQuery("applicant", "")
	sType := c.DefaultQuery("type", "")
	moneyType := c.DefaultQuery("money_type", "")
	flowLimit := c.DefaultQuery("flow_limit", "")
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
	limit, offset := request.GetOffsets(c)
	userDividends := make([]models.UserDividend, 0)
	engine := common.Mysql(request.GetPlatform(c))
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
	viewFile := "dividend_managements/applies.html"
	if request.IsAjax(c) {
		viewFile = "dividend_managements/_applies.html"
	}
	response.Render(c, viewFile, viewData)
}
