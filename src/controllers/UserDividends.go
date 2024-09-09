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

// UserDevidendsStruct 红利获取
type UserDevidendsStruct struct {
	models.UserDividend `xorm:"extends"`
	Vip                 int32
}

// UserDevidends 用户红利记录
var UserDevidends = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) { //默 认首页
		cond := builder.NewCond()
		if value, exists := c.GetQuery("created"); !exists {
			currentTime := time.Now().Unix()
			startAt := tools.SecondToMicro(currentTime - currentTime%86400 - 8*3600)
			endAt := startAt + tools.SecondToMicro(86400)
			cond = cond.And(builder.Gte{"user_dividends.created": startAt}).And(builder.Lte{"user_dividends.created": endAt})
		} else {
			areas := strings.Split(value, " - ")
			startAt := tools.GetMicroTimeStampByString(areas[0])
			endAt := tools.GetMicroTimeStampByString(areas[1])
			cond = cond.And(builder.Gte{"user_dividends.created": startAt}).And(builder.Lte{"user_dividends.created": endAt})
		}
		request.QueryCond(c, &cond, map[string]map[string]string{
			"%": {
				"username": "user_dividends.username",
				"bill_no":  "user_dividends.bill_no",
				"reviewer": "user_dividends.reviewer",
			},
			"=": {
				"type":       "user_dividends.type",
				"state":      "user_dividends.state",
				"flow_limit": "user_dividends.flow_limit",
			},
		})
		limit, offset := request.GetOffsets(c)
		userDevidends := make([]UserDevidendsStruct, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("user_dividends").
			Join("LEFT OUTER", "users", "user_dividends.user_id = users.id").
			Where(cond).
			OrderBy("user_dividends.id DESC").
			Limit(limit, offset).
			FindAndCount(&userDevidends)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ArtificialMoneySum := new(OperationsDividendMoneyStruct)
		TotalArtificialMoneySum, _ := engine.Table("user_dividends").Where(cond).Sum(ArtificialMoneySum, "money")
		viewData := pongo2.Context{
			"rows":            userDevidends,
			"total":           total,
			"total_dividends": TotalArtificialMoneySum,
		}
		response.Render(c, "user_changes/_user_dividends.html", viewData)
	},
}
