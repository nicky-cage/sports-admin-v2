package user_detail_accounts

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *UserDetailAccounts) Withdraws(c *gin.Context) {
	var part string
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	rows := []models.UserWithdraw{}
	cond := builder.NewCond()
	cond = cond.And(builder.Eq{"user_id": id}).And(builder.Eq{"type": 1})
	state := c.DefaultQuery("status", "")
	if len(state) > 0 {
		cond = cond.And(builder.Eq{"status": state})
		part = part + " and status=" + state
	} else {
		part = part + " and status=2 "
	}
	create := c.Query("created")
	if create != "" { //对时间进行处理
		areas := strings.Split(create, " - ")
		startAt, _ := time.Parse("2006-01-02 15:04:05", areas[0]+" 00:00:00")
		endAt, _ := time.Parse("2006-01-02 15:04:05", areas[1]+" 23:59:59")
		cond = cond.And(builder.Gte{"updated": startAt.UnixMicro()}).And(builder.Lte{"updated": endAt.UnixMicro()})
		part = part + " and updated >= " + strconv.Itoa(int(startAt.UnixMicro())) + " and updated < " + strconv.Itoa(int(endAt.UnixMicro()))
	}
	limit, offset := request.GetOffsets(c)
	total, err := dbSession.Table("user_withdraws").Where(cond).
		OrderBy("finance_process_at DESC").Limit(limit, offset).FindAndCount(&rows)
	if err != nil {
		log.Err(err.Error())
		return
	}
	res := []models.TotalFloat{}
	sql := "select sum(money) as total from user_withdraws where user_id = " + id
	if err := dbSession.SQL(sql + part).Find(&res); err != nil {
		log.Err(err.Error())
	}
	for i := 0; i < len(rows); i++ {
		if rows[i].BankName == "其他银行" {
			sqlCard := "select bank_name from user_cards where card_number= " + rows[i].BankCard
			dataCard, _ := dbSession.QueryString(sqlCard)
			if dataCard[0]["bank_name"] != "" {
				rows[i].BankName = dataCard[0]["bank_name"]
			}
		}
	}
	viewData := response.ViewData{
		"id":              id,
		"rows":            rows,
		"total":           total,
		"withdraws_money": res[0].Total,
	}
	if create == "" {
		response.Render(c, "users/detail_withdraws.html", viewData)
	} else {
		response.Render(c, "users/_detail_withdraws.html", viewData)
	}
}
