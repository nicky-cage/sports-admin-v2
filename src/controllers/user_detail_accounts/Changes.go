package user_detail_accounts

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *UserDetailAccounts) Changes(c *gin.Context) {
	cond := builder.NewCond()
	status := c.DefaultQuery("type_id", "")
	if len(status) > 0 {
		cond = cond.And(builder.Eq{"multiple_type": status})
	}
	create := c.Query("created")
	if create != "" {
		//对时间进行处理
		areas := strings.Split(create, " - ")
		startAt := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
		endAt := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
		cond = cond.And(builder.Gte{"updated": startAt}).And(builder.Lte{"updated": endAt})
	}
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	rows := []models.Transaction{}
	cond = cond.And(builder.Eq{"user_id": id})
	limit, offset := request.GetOffsets(c)
	total, err := dbSession.Table("transactions").Where(cond).OrderBy("created DESC").Limit(limit, offset).FindAndCount(&rows)
	if err != nil {
		log.Err(err.Error())
		return
	}
	sql := "select sum(amount) as money from transactions where user_id=" + id
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
	}
	amount, _ := strconv.ParseFloat(res[0]["money"], 64)
	viewData := response.ViewData{
		"rows":         rows,
		"total":        total,
		"id":           id,
		"amount_money": amount,
	}
	if create == "" {
		response.Render(c, "users/detail_changes.html", viewData)
	} else {
		response.Render(c, "users/_detail_changes.html", viewData)
	}
}
