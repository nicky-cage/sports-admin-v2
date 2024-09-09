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

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *UserDetailAccounts) Dividends(c *gin.Context) {
	limit, offset := request.GetOffsets(c)
	id := c.Query("id")
	cond := builder.NewCond()
	if val, err := strconv.Atoi(c.DefaultQuery("type", "")); err == nil && val > 0 {
		cond = cond.And(builder.Eq{"type": val})
	}
	if val, err := strconv.Atoi(c.DefaultQuery("state", "")); err == nil && val >= 0 {
		cond = cond.And(builder.Eq{"state": val})
	}
	create := c.Query("created")
	if create != "" { //对时间进行处理
		areas := strings.Split(create, " - ")
		startAt := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
		endAt := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
		cond = cond.And(builder.Gte{"updated": startAt}).And(builder.Lte{"updated": endAt})
	}
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	rows := []models.UserDividend{}
	cond = cond.And(builder.Eq{"user_id": id})
	total, err := dbSession.Table("user_dividends").Where(cond).OrderBy("created DESC").Limit(limit, offset).FindAndCount(&rows)
	if err != nil {
		log.Err(err.Error())
		return
	}
	tRow := struct {
		TotalApply float64 `json:"total_apply"`
		TotalPass  float64 `json:"total_pass"`
		TotalDeny  float64 `json:"total_deny"`
	}{}
	sql := "SELECT SUM(IF(state = 1, money, 0)) AS total_apply, SUM(IF(state = 2, money, 0)) AS total_pass, SUM(IF(state = 3, money, 0)) AS total_deny " +
		"FROM user_dividends WHERE user_id = " + id
	_, err = dbSession.SQL(sql).Get(&tRow)
	if err != nil {
		log.Err(err.Error())
	}
	//红利 类型,
	type dividendTotal struct {
		Apply float64 `json:"apply"`
		Pass  float64 `json:"pass"`
		Deny  float64 `json:"deny"`
	}
	pageTotal := func() dividendTotal {
		totalAcc := dividendTotal{}
		for _, r := range rows {
			if r.State == 1 {
				totalAcc.Apply += r.Money
			} else if r.State == 2 {
				totalAcc.Pass += r.Money
			} else if r.State == 3 {
				totalAcc.Deny += r.Money
			}
		}
		return totalAcc
	}()
	data := pongo2.Context{
		"rows":   rows,
		"total":  total,
		"id":     id,
		"t":      tRow,
		"tTotal": tRow.TotalApply + tRow.TotalDeny + tRow.TotalPass,
		"p":      pageTotal,
		"pTotal": pageTotal.Apply + pageTotal.Pass + pageTotal.Deny,
	}
	if create == "" {
		response.Render(c, "users/detail_dividends.html", data)
	} else {
		response.Render(c, "users/_detail_dividends.html", data)
	}
}
