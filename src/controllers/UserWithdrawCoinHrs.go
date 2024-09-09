package controllers

import (
	"fmt"
	"sports-admin/caches"
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

// UserWithdrawCoinHrsStruct 代币提现历史记录
type UserWithdrawCoinHrsStruct struct {
	models.UserWithdraw `xorm:"extends"`
	Label               string `json:"label"`
	Vip                 int32
}

// UserWithdrawCoinHrSumStruct 提现历史记录统计
type UserWithdrawCoinHrSumStruct struct {
	Money float64
}

// UserWithdrawCoinHrs 提款管理-历史记录
var UserWithdrawCoinHrs = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		startAt, endAt := func() (int64, int64) { // 申请时间
			value, exists := c.GetQuery("created")
			if !exists {
				currentTime := time.Now().Unix()
				startAt := currentTime - currentTime%86400
				return startAt, startAt + 86400
			}
			areas := strings.Split(value, " - ")
			return tools.GetMicroTimeStampByString(areas[0]), tools.GetMicroTimeStampByString(areas[1])
		}()
		cond = cond.And(builder.Gte{"user_withdraws.created": startAt}).And(builder.Lte{"user_withdraws.created": endAt})
		uStartAt, uEndAt := func() (int64, int64) { // 完成时间
			value, exists := c.GetQuery("updated")
			if !exists {
				currentTime := time.Now().Unix()
				startAt := currentTime - currentTime%86400
				return tools.SecondToMicro(startAt), tools.SecondToMicro(startAt + 86400)
			}
			areas := strings.Split(value, " - ")
			return tools.GetMicroTimeStampByString(areas[0]), tools.GetMicroTimeStampByString(areas[1])
		}()
		cond = cond.And(builder.Gte{"user_withdraws.updated": uStartAt}).And(builder.Lte{"user_withdraws.updated": uEndAt})
		cond = cond.And(builder.Neq{"user_withdraws.status": 1}).And(builder.Eq{"user_withdraws.type": 1})
		if min, ok := c.GetQuery("money_min"); ok {
			cond = cond.And(builder.Gte{"user_withdraws.money": min})
		}
		if max, ok := c.GetQuery("money_max"); ok {
			cond = cond.And(builder.Lte{"user_withdraws.money": max})
		}
		if username := c.DefaultQuery("username", ""); username != "" {
			cond = cond.And(builder.Eq{"user_withdraws.username": username})
		}
		request.QueryCondEq(c, &cond, map[string]string{
			"username":       "user_withdraws.username",
			"bill_no":        "user_withdraws.bill_no",
			"risk_admin":     "user_withdraws.risk_admin",
			"finance_admin":  "user_withdraws.finance_admin",
			"status":         "user_withdraws.status",
			"business_type":  "user_withdraws.payment_method",
			"payment_method": "user_withdraws.business_type",
		})
		cond = cond.And(builder.Gt{"user_withdraws.wallet_id": 0}) // 钱包id > 0
		limit, offset := request.GetOffsets(c)
		userWithdrawsHrs := make([]UserWithdrawHrsStruct, 0)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		dbSession.Table("user_withdraws").
			Join("LEFT OUTER", "users", "user_withdraws.user_id = users.id").
			Where(cond).OrderBy("user_withdraws.id DESC")
		if request.IsExportExcel(c) { // 如果是导出数据
			err := dbSession.Find(&userWithdrawsHrs)
			if err != nil {
				response.Err(c, "获取提款数据失败")
				return
			}
			response.Result(c, userWithdrawsHrs)
			return
		}

		total, err := dbSession.Limit(limit, offset).FindAndCount(&userWithdrawsHrs)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}

		// 统计数据结构
		type TotalInfo struct {
			Total    int     `json:"total"`    // 订单笔数
			Order    float64 `json:"order"`    // 订单金额
			Fee      float64 `json:"fee"`      // 手续费
			Cost     float64 `json:"cost"`     // 行政费用
			Real     float64 `json:"real"`     // 实需出款
			Withdraw float64 `json:"withdraw"` // 实已出款
			Rate     float64 `json:"rate"`     // 成功率
		}
		totalPage := TotalInfo{}
		totalQuery := TotalInfo{}
		totalAll := TotalInfo{}
		for _, v := range userWithdrawsHrs {
			totalPage.Total += 1
			totalPage.Order += v.Money
			if v.Status == 2 {
				totalPage.Fee += v.TransactionFee
				totalPage.Real += v.ActualMoney
				totalPage.Cost += v.WithdrawCost
				totalPage.Withdraw += v.ActualMoney
			}
		}
		sqlCond, condParams, _ := builder.ToSQL(cond)
		sqlQueryTotal := "SELECT COUNT(1) AS total, SUM(user_withdraws.money) AS `order`, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.withdraw_cost, 0)) AS cost, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.transaction_fee, 0)) AS fee, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.actual_money, 0)) AS `real`, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.actual_money, 0)) AS withdraw " +
			"FROM user_withdraws LEFT JOIN users ON user_withdraws.user_id = users.id WHERE " + sqlCond
		if _, err := dbSession.Where(cond).SQL(sqlQueryTotal, condParams...).Get(&totalQuery); err != nil {
			fmt.Println("查询统计出错:", err)
		}
		sqlTotalAll := "SELECT COUNT(1) AS total, SUM(user_withdraws.money) AS `order`, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.withdraw_cost, 0)) AS cost, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.transaction_fee, 0)) AS fee, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.actual_money, 0)) AS `real`, " +
			"SUM(IF(user_withdraws.status = 2, user_withdraws.actual_money, 0)) AS withdraw " +
			"FROM user_withdraws LEFT JOIN users ON user_withdraws.user_id = users.id "
		if _, err := dbSession.SQL(sqlTotalAll).Get(&totalAll); err != nil {
			fmt.Println("总计统计出错:", err)
		}

		totalPage.Rate = tools.Fixed((totalPage.Withdraw/totalPage.Order)*100, 2)
		totalQuery.Rate = tools.Fixed((totalQuery.Withdraw/totalQuery.Order)*100, 2)
		totalAll.Rate = tools.Fixed((totalAll.Withdraw/totalAll.Order)*100, 2)
		viewData := pongo2.Context{
			"rows":        userWithdrawsHrs,
			"total":       total,
			"total_page":  totalPage,
			"total_query": totalQuery,
			"total_all":   totalAll,
			"vipLevels":   caches.UserLevels.All(platform),
		}
		viewFile := request.GetViewFile(c, "user_withdraw_coins/%shistory_records.html")
		response.Render(c, viewFile, viewData)
	},
}
