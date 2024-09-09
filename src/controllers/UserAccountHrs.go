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

// UserAccountHrs 上下分-历史列表
var UserAccountHrs = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		timeStart, timeEnd := func() (int64, int64) {
			if value, exists := c.GetQuery("created"); !exists {
				currentDayTime := time.Now().Format("2006-01-02")
				return tools.GetMicroTimeStampByString(currentDayTime + " 00:00:00"), tools.GetMicroTimeStampByString(currentDayTime + " 23:59:59")
			} else {
				areas := strings.Split(value, " - ")
				return tools.GetMicroTimeStampByString(areas[0] + " 00:00:00"), tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
			}
		}()
		cond = cond.And(builder.Gte{"user_account_sets.updated": timeStart}).And(builder.Lte{"user_account_sets.updated": timeEnd})
		cond = cond.And(builder.Neq{"user_account_sets.status": 1})
		if min, ok := c.GetQuery("money_min"); ok {
			cond = cond.And(builder.Gte{"user_account_sets.money": min})
		}
		if max, ok := c.GetQuery("money_max"); ok {
			cond = cond.And(builder.Lte{"user_account_sets.money": max})
		}
		request.QueryCondEq(c, &cond, map[string]string{
			"type":   "user_account_sets.type",
			"status": "user_account_sets.status",
		})
		request.QueryCondLike(c, &cond, map[string]string{
			"username": "user_account_sets.username",
			"audit":    "user_account_sets.audit",
		})
		limit, offset := request.GetOffsets(c)
		userAccountSets := make([]models.UserAccountSet, 0)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		dbSession.Table("user_account_sets").Where(cond).OrderBy("user_account_sets.id DESC")
		if exportExcel := c.DefaultQuery("export_excel", ""); exportExcel != "" { // 如果是导出excel
			dbSession.Find(&userAccountSets)
			response.Result(c, userAccountSets)
			return
		}
		total, err := dbSession.Limit(limit, offset).FindAndCount(&userAccountSets)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}

		ss := new(SumStruct)
		sumTotal, _ := dbSession.Table("user_account_sets").Where(cond).Sum(ss, "money")
		viewData := pongo2.Context{
			"rows":      userAccountSets,
			"total":     total,
			"sum_money": sumTotal,
		}
		viewFile := request.GetViewFile(c, "user_account_sets/%shistory_records.html")
		response.Render(c, viewFile, viewData)
	},
}
