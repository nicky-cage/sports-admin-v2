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

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// UserDepositCoins 存款管理-存款列表
var UserDepositCoins = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		if value, exists := c.GetQuery("created"); exists {
			areas := strings.Split(value, " - ")
			startAt, endAt := tools.GetMicroTimeStampByString(areas[0]), tools.GetMicroTimeStampByString(areas[1])
			cond = cond.And(builder.Gte{"user_deposits.created": startAt}).And(builder.Lte{"user_deposits.created": endAt})
		}
		cond = cond.And(builder.Eq{"user_deposits.type": 4}) // 4: 虚拟币
		request.QueryCondEq(c, &cond, map[string]string{
			"status":          "user_deposits.status",
			"vip":             "users.vip",
			"channel_type":    "user_deposits.channel_type",
			"account_by_name": "user_deposits.account_by_name",
		})
		request.QueryCondLike(c, &cond, map[string]string{
			"deposit_name": "user_deposits.deposit_name",
			"username":     "user_deposits.username",
			"bill_no":      "user_deposits.order_no",
		})
		cond = cond.And(builder.Lte{"user_deposits.status": 1})
		userDeposits := make([]UserDepositsStruct, 0)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		dbSession.Table("user_deposits").
			Select("user_deposits.*, users.vip, users.label").
			Join("LEFT OUTER", "users", "user_deposits.user_id = users.id").
			Where(cond).
			OrderBy("user_deposits.id DESC")
		if request.IsExportExcel(c) { // 如果只是导出数据
			err := dbSession.Find(&userDeposits)
			if err != nil {
				response.Err(c, "获取数据有误: "+err.Error())
				return
			}
			response.Result(c, userDeposits)
			return
		}

		limit, offset := request.GetOffsets(c)
		queryTotal, err := dbSession.Limit(limit, offset).FindAndCount(&userDeposits)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取存款列表错误")
			return
		}
		// 计算本次查询的所有订单数
		pageInfo := struct {
			Total           int     `json:"total"`
			Money           float64 `json:"money"`
			Success         int     `json:"success"`
			Discount        float64 `json:"discount"`
			SuccessMoney    float64 `json:"success_money"`
			SuccessDiscount float64 `json:"success_discount"`
		}{}
		pageInfo.Total = len(userDeposits)
		for _, v := range userDeposits {
			pageInfo.Money += v.Money
			pageInfo.Discount += v.Discount
			if v.Status == 2 { // 成功
				pageInfo.Success += 1
				pageInfo.Money += v.Money
				pageInfo.SuccessDiscount += v.Discount
			}
		}

		// 查询总计 -
		queryRow := struct {
			Money    float64 `json:"money"`
			Discount float64 `json:"discount"`
		}{}
		queryNoStatus, _ := dbSession.Table("user_deposits").Where(cond).Sums(&queryRow, "money", "discount") // 默认条件
		cond = cond.And(builder.Eq{"status": 2})                                                              // 加上 成功状态
		queryStatus, _ := dbSession.Table("user_deposits").Where(cond).Sums(&queryRow, "money", "discount")
		queryCount := struct {
			Id int `json:"id" xorm:"id"`
		}{}
		querySuccess, _ := dbSession.Table("user_deposits").Where(cond).Count(&queryCount)
		queryInfo := struct {
			Total           int     `json:"total"`
			Money           float64 `json:"money"`
			Discount        float64 `json:"discount"`
			Success         int     `json:"success"`
			SuccessMoney    float64 `json:"success_money"`
			SuccessDiscount float64 `json:"success_discount"`
		}{
			Total:           int(queryTotal),   // 总记录数
			Money:           queryNoStatus[0],  // 总存款数
			Success:         int(querySuccess), // 成功记录数
			Discount:        queryNoStatus[1],  //
			SuccessMoney:    queryStatus[0],    // 成功存款数
			SuccessDiscount: queryStatus[1],    // 成功存款优惠
		}

		// 计算所有查询统计 - 总计存款 /总计成功存款 - 总计存款笔数/总计成功存款笔数 - 总计上分
		totalInfo := struct {
			Total           int
			TotalMoney      float64
			Discount        float64
			Success         int
			SuccessMoney    float64
			SuccessDiscount float64
		}{}
		type TotalRow struct {
			Kind     string  `json:"kind"`
			Total    int     `json:"total"`
			Money    float64 `json:"money"`
			Discount float64 `json:"discount"`
		}
		totalType := "4"
		totalArr := []string{
			fmt.Sprintf("(SELECT 'to' AS Kind, COUNT(*) AS total, SUM(money) AS money, SUM(discount) AS discount "+
				"FROM user_deposits WHERE type IN (%s))", totalType), // 总计存款
			fmt.Sprintf("(SELECT 'ts' AS Kind, COUNT(*) AS total, SUM(money) AS money, SUM(discount) AS discount "+
				"FROM user_deposits WHERE status = 2 AND type IN (%s))", totalType), // 总计成功存款
		}
		totalRows := []TotalRow{}
		if err := dbSession.SQL(strings.Join(totalArr, " UNION")).Find(&totalRows); err == nil {
			for _, r := range totalRows {
				if r.Kind == "to" { // 总计
					totalInfo.Total += r.Total
					totalInfo.TotalMoney += r.Money
					totalInfo.Discount += r.Discount
				} else if r.Kind == "ts" { // 成功
					totalInfo.Success += r.Total
					totalInfo.SuccessMoney += r.Money
					totalInfo.SuccessDiscount += r.Discount
				}
			}
		}

		viewData := pongo2.Context{
			"rows":                     userDeposits,
			"total":                    queryTotal,
			"page_total":               pageInfo.Total,                                              // page
			"page_money":               pageInfo.Money,                                              // page
			"page_success":             pageInfo.Success,                                            // page
			"page_success_money":       pageInfo.SuccessMoney,                                       // page
			"page_success_discount":    pageInfo.SuccessDiscount,                                    // page
			"page_success_rate":        float64(pageInfo.Success) / float64(pageInfo.Total) * 100.0, //
			"page_success_money_rate":  pageInfo.SuccessMoney / pageInfo.Money * 100.0,              //
			"page_discount":            pageInfo.Discount,
			"query_total":              queryInfo.Total, // query
			"query_money":              queryInfo.Money, // query
			"query_success":            queryInfo.Success,
			"query_success_money":      queryInfo.SuccessMoney, // query
			"query_success_discount":   queryStatus[1],         // query
			"query_success_rate":       float64(queryInfo.Success) / float64(queryInfo.Total) * 100.0,
			"query_success_money_rate": pageInfo.SuccessMoney / pageInfo.Money * 100.0,
			"query_discount":           pageInfo.Discount,
			"total_record":             totalInfo.Total,                                               // total
			"total_money":              totalInfo.TotalMoney,                                          // total
			"total_success_record":     totalInfo.Success,                                             // total - success
			"total_success_money":      totalInfo.SuccessMoney,                                        // total - success
			"total_success_discount":   totalInfo.SuccessDiscount,                                     // total - success
			"total_success_rate":       float64(totalInfo.Success) / float64(totalInfo.Total) * 100.0, // total - rate
			"total_success_money_rate": totalInfo.SuccessMoney / totalInfo.TotalMoney * 100.0,         // total
			"total_discount":           totalInfo.Discount,
			"vipLevels":                caches.UserLevels.All(platform),
			"depositVirtuals":          models.DepositVirtuals.GetVirtuals(platform),
		}

		viewFile := func() string {
			if request.IsAjax(c) {
				return "user_deposit_coins/_user_deposits.html"
			}
			return "user_deposit_coins/user_deposits.html"
		}()
		SetLoginAdmin(c)
		response.Render(c, viewFile, viewData)
	},
}
