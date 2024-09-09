package controllers

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// UserWithdrawHrsStruct 代币提现历史记录
type UserWithdrawHrsStruct struct {
	models.UserWithdraw `xorm:"extends"`
	Label               string `json:"label"`
	Vip                 int32
}

// UserWithdrawHrSumStruct 提现历史记录统计
type UserWithdrawHrSumStruct struct {
	Money float64
}

// UserWithdrawHrs 提款管理-历史记录
var UserWithdrawHrs = struct {
	List func(*gin.Context)
	*ActionExport
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		startAt, endAt := func() (int64, int64) { // 申请时间
			value, exists := c.GetQuery("created")
			if !exists {
				currentTime := time.Now().Unix()
				startAt := currentTime - currentTime%86400
				return tools.SecondToMicro(startAt), tools.SecondToMicro(startAt + 86400)
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
		cond = cond.And(builder.Eq{"user_withdraws.wallet_id": 0}) // 钱包id == 0
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
		}
		viewFile := request.GetViewFile(c, "user_withdraws/%shistory_records.html")
		response.Render(c, viewFile, viewData)
	},
	ActionExport: &ActionExport{
		Columns: []ExportHeader{
			{"序号", "id"},
			{"订单编号", "bill_no"},
			{"会员编号", "user_id"},
			{"会员名称", "username"},
			{"会员等级", "vip"},
			{"订单金额", "money"},
			{"手续费", "transaction_fee"},
			{"行政费", "withdraw_cost"},
			{"实需出款", "actual_money"},
			{"实已出款", "actual_withdraw_money"},
			{"银行卡信息", "bank_info"},
			{"申请时间", "created"},
			{"代付", "payout_third"},
			{"风控审核时间", "risk_process_at"},
			{"风控审核人", "risk_admin"},
			{"完成时间", "finance_process_at"},
			{"出款人", "finance_admin"},
			{"出款卡号/商户", "card_number"},
			{"状态", "status"},
		},
		GetSQL: func(s string, context *gin.Context) string {
			sql := "SELECT uw.*, u.vip " +
				"FROM user_withdraws AS uw " +
				"LEFT JOIN users AS u ON uw.user_id = u.id " +
				"ORDER BY uw.id DESC"
			return sql
		},
		ProcessRow: func(m *map[string]interface{}, c *gin.Context) {
			needWithdrawMoney, _ := strconv.ParseFloat(fmt.Sprintf("%v", (*m)["money"]), 64)
			withdrawCost, _ := strconv.ParseFloat(fmt.Sprintf("%v", (*m)["withdraw_cost"]), 64)
			actualMoney := needWithdrawMoney - withdrawCost
			(*m)["money"] = fmt.Sprintf("%.2f", (*m)["money"].(float64))
			(*m)["transaction_fee"] = fmt.Sprintf("%.2f", (*m)["transaction_fee"].(float64))
			(*m)["withdraw_cost"] = fmt.Sprintf("%.2f", (*m)["withdraw_cost"].(float64))
			(*m)["actual_money"] = fmt.Sprintf("%.2f", actualMoney)
			if (*m)["status"] != "1" && (*m)["status"] != "2" {
				needWithdrawMoney, _ := strconv.ParseFloat(fmt.Sprintf("%v", (*m)["money"]), 64)
				withdrawCost, _ := strconv.ParseFloat(fmt.Sprintf("%v", (*m)["withdraw_cost"]), 64)
				actualMoney := needWithdrawMoney - withdrawCost
				(*m)["actual_withdraw_money"] = fmt.Sprintf("%.2f", actualMoney)
			} else {
				(*m)["actual_withdraw_money"] = "0.00"
			}
			(*m)["bank_info"] = fmt.Sprintf("%v/%v/%v/%v", (*m)["bank_realname"], (*m)["bank_name"], (*m)["bank_card"], (*m)["bank_address"])
			created := (*m)["created"].(float64)
			(*m)["created"] = base_controller.FieldToDateTime(fmt.Sprintf("%d", int(created)))
			if (*m)["business_type"] == "0" {
				(*m)["payout-third"] = "否"
			} else if (*m)["business_tpe"] == "1" {
				(*m)["payout_third"] = "是"
			} else {
				(*m)["payout_third"] = "未知"
			}
			(*m)["status"] = func(val string) string {
				if val == "1" {
					return "处理中"
				} else if val == "2" {
					return "成功"
				} else {
					return "失败"
				}
			}(fmt.Sprintf("%v", (*m)["status"]))
		},
	},
}
