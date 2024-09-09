package user_detail_accounts

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *UserDetailAccounts) Commissions(c *gin.Context) {
	cond := builder.NewCond()
	create := c.Query("created")
	if create != "" {
		//对时间进行处理
		areas := strings.Split(create, " - ")
		startAt, _ := time.Parse("2006-01-02 15:04:05", areas[0]+" 00:00:00")
		endAt, _ := time.Parse("2006-01-02 15:04:05", areas[1]+" 23:59:59")
		if endAt != startAt {
			cond = cond.And(builder.Gte{"user_rebate_records.created": startAt.Unix()}).And(builder.Lte{"user_rebate_records.created": endAt.Unix()})
		} else {
			cond = cond.And(builder.Gte{"user_rebate_records.created": startAt.Unix()}).And(builder.Lte{"user_rebate_records.created": endAt.Unix() + 86400})
		}
	}
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	rows := []RebateAndVip{}
	cond = cond.And(builder.Eq{"user_id": id})
	limit, offset := request.GetOffsets(c)
	total, err := dbSession.Table("user_rebate_records").
		Join("LEFT OUTER", "users", "user_rebate_records.user_id = users.id").
		Where(cond).
		OrderBy("user_rebate_records.created DESC").
		Limit(limit, offset).
		FindAndCount(&rows)
	if err != nil {
		log.Err(err.Error())
		return
	}
	var rCount SumAgent
	countNum, rerr := dbSession.Table("user_rebate_records").Where(cond).Sum(rCount, "money")
	if rerr != nil {
		log.Err(rerr.Error())
	}
	viewData := response.ViewData{
		"rows":    rows,
		"total":   total,
		"id":      id,
		"account": countNum,
	}
	if create == "" {
		response.Render(c, "users/detail_commissions.html", viewData)
	} else {
		response.Render(c, "users/_detail_commissions.html", viewData)
	}
}
