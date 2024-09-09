package user_detail_accounts

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// Transfers 用户详情 - 平台转账
func (ths *UserDetailAccounts) Transfers(c *gin.Context) {
	cond := builder.NewCond()
	status := c.DefaultQuery("status", "")
	if len(status) > 0 {
		cond = cond.And(builder.Eq{"status": status})
	}
	transferType := c.DefaultQuery("transfer_type", "")
	if len(transferType) > 0 {
		cond = cond.And(builder.Eq{"transfer_type": transferType})
	}
	in := c.DefaultQuery("transfer_in_account", "")
	if len(in) > 0 {
		cond = cond.And(builder.Eq{"transfer_in_account": in})
	}
	out := c.DefaultQuery("transfer_out_account", "")
	if len(out) > 0 {
		cond = cond.And(builder.Eq{"transfer_out_account": out})
	}
	create := c.Query("created")
	if create != "" { //对时间进行处理
		areas := strings.Split(create, " - ")
		startAt, _ := time.Parse("2006-01-02 15:04:05", areas[0]+" 00:00:00")
		endAt, _ := time.Parse("2006-01-02 15:04:05", areas[1]+" 23:59:59")
		//cond = cond.And(builder.Gte{"created": start_at.Unix()}).And(builder.Lte{"created": end_at.Unix()})
		if endAt != startAt {
			cond = cond.And(builder.Gte{"created": startAt.UnixMicro()}).And(builder.Lte{"created": endAt.UnixMicro()})
		} else {
			cond = cond.And(builder.Gte{"created": startAt.UnixMicro()}).And(builder.Lte{"created": (endAt.UnixMicro() + tools.SecondToMicro(86400))})
		}
	}
	// 过滤金额为0的记录
	cond = cond.And(builder.Neq{"money": 0})

	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	rows := []models.UserTransfer{}
	cond = cond.And(builder.Eq{"user_id": id})
	limit, offset := request.GetOffsets(c)
	total, err := dbSession.Table("user_transfers").OrderBy("created desc").Where(cond).Limit(limit, offset).FindAndCount(&rows)
	if err != nil {
		log.Err(err.Error())
		return
	}
	res, err := dbSession.QueryString("SELECT name, ename FROM game_venues WHERE pid = 0")
	if err != nil {
		log.Err(err.Error())
		return
	}

	cRes := []models.TotalFloat{}
	cSql := "select sum(money) as total from user_transfers where user_id=" + id
	if err := dbSession.SQL(cSql).Find(&cRes); err != nil {
		log.Err(err.Error())
	}
	viewData := response.ViewData{
		"rows":            rows,
		"total":           total,
		"id":              id,
		"venue":           res,
		"transfers_money": cRes[0].Total,
	}
	if create == "" {
		response.Render(c, "users/detail_transfers.html", viewData)
		return
	}

	response.Render(c, "users/_detail_transfers.html", viewData)
}
