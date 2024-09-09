package controllers

import (
	"errors"
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/pgsql"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var getFloat = func(num string) float64 {
	if val, err := strconv.ParseFloat(num, 64); err == nil {
		return val
	}
	return 0
}

var getInt = func(num string) int {
	if val, err := strconv.Atoi(num); err == nil {
		return val
	}
	return 0
}

var getUserAuditIDs = func(c *gin.Context) (int, int, error) { // 获取编号相关信息
	IDStr := c.DefaultQuery("id", "0")
	userIDStr := c.DefaultQuery("user_id", "0")
	ID, IDErr := strconv.Atoi(IDStr)
	userID, userIDErr := strconv.Atoi(userIDStr)
	if IDErr != nil || userIDErr != nil || ID <= 0 || userID <= 0 {
		return 0, 0, errors.New("稽核记录编号错误")
	}

	return ID, userID, nil
}

var UserAudits = struct {
	*ActionList
	Detail   func(*gin.Context)
	DetailV2 func(*gin.Context)
	DetailV3 func(*gin.Context)
	BetsV2   func(*gin.Context)
	Update   func(*gin.Context)
	Sync     func(*gin.Context)
	Delete   func(*gin.Context)
}{
	ActionList: &ActionList{
		RequireParameters: true,
		Model:             models.UserAudits,
		ViewFile:          "user_audits/list.html",
		QueryCond: map[string]interface{}{
			"username": "%",
		},
		Rows: func() interface{} {
			return &[]models.UserAudit{}
		},
		OrderBy: func(C *gin.Context) string {
			return "vip DESC"
		},
		ProcessRow: func(c *gin.Context, res interface{}) {
			rs := res.(*[]models.UserAudit)
			platform := request.GetPlatform(c)
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			for k, r := range *rs {
				// rows, err := dbSession.QueryString("CALL USER_AUDIT_INFO(?)", r.Id)
				rows := models.UserAudits.GetAuditInfo(dbSession, int(r.Id))
				if len(rows) >= 1 {
					row := rows[0]
					(*rs)[k].TotalDeposit = getFloat(row["total_deposit"])
					(*rs)[k].TotalWithdraw = getFloat(row["total_withdraw"])
					(*rs)[k].TotalActivity = getFloat(row["total_activity"])
					(*rs)[k].TotalActivityApply = getFloat(row["total_activity_apply"])
					(*rs)[k].TotalUserReset = getFloat(row["total_user_reset"])
					(*rs)[k].TotalAccountSet = getFloat(row["total_account_set"])
					(*rs)[k].TotalDividend = getFloat(row["total_dividend"])
					(*rs)[k].DepositFirstMoney = getFloat(row["deposit_first_money"])
					(*rs)[k].DepositLastMoney = getFloat(row["deposit_last_money"])
					(*rs)[k].WithdrawFirstMoney = getFloat(row["withdraw_first_money"])
					(*rs)[k].WithdrawLastMoney = getFloat(row["withdraw_last_money"])
					(*rs)[k].DepositCount = getInt(row["deposit_count"])
					(*rs)[k].WithdrawCount = getInt(row["withdraw_count"])
					(*rs)[k].WithdrawFirstTime = uint32(getFloat(row["withdraw_first_time"]))
					(*rs)[k].WithdrawLastTime = uint32(getInt(row["withdraw_last_time"]))
					(*rs)[k].DepositFirstTime = uint32(getFloat(row["deposit_first_time"]))
					(*rs)[k].DepositLastTime = uint32(getInt(row["deposit_last_time"]))
				}
			}
		},
	},
	Detail: func(c *gin.Context) {
		userIDStr := c.DefaultQuery("id", "0")
		if userIDStr == "0" {
			response.Err(c, "缺少用户编号")
			return
		}
		userID, err := strconv.Atoi(userIDStr)
		if err != nil || userID <= 0 {
			response.Err(c, err.Error())
			return
		}
		user := models.User{}
		platform := request.GetPlatform(c)
		exists, err := models.Users.FindById(platform, userID, &user)
		if err != nil || !exists {
			response.Err(c, "用户信息查找失败")
			return
		}

		// 刷新最新的稽核信息
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		debugStart := tools.TimeDebugBegin("开始计算稽核")
		if err := AuditTrigger.Audit(dbSession, platform, userID); err != nil {
			response.Err(c, err.Error())
			return
		}
		_ = tools.TimeDebugAt(debugStart, "结束计算稽核")

		var firstDepositAfterLastWithdraw int64 = 0
		info, err := models.UserAudits.Detail(platform, userID, &firstDepositAfterLastWithdraw)
		if err != nil {
			response.Err(c, err.Error())
			return
		}
		totalWithdrawCount := 0 // 总提款次数
		totalDepositCount := 0  // 总存款次数
		depositLastTime := "-"
		withdrawLastTime := "-"
		totalDeposit := 0.0
		totalWithdraw := 0.0

		//rows, err := dbSession.QueryString("CALL USER_AUDIT_INFO(?)", userID)
		rows := models.UserAudits.GetAuditInfo(dbSession, userID, true)
		if err == nil && len(rows) > 0 {
			row := rows[0]
			totalWithdrawCount = getInt(row["withdraw_count"])
			totalDepositCount = getInt(row["deposit_count"])
			totalDeposit = getFloat(row["total_deposit"])
			totalWithdraw = getFloat(row["total_withdraw"])

			if v := getInt(row["withdraw_last_time"]); v > 0 {
				withdrawLastTime = time.Unix(int64(v), 0).Format("2006-01-02 15:04:05")
			}
			if v := getInt(row["deposit_last_time"]); v > 0 {
				depositLastTime = time.Unix(int64(v), 0).Format("2006-01-02 15:04:05")
			}
		}
		// 警告, 此用户没有任何存款信息, 但是账上有钱
		isWarn := false
		if len(info.Rows) > 0 && totalDeposit == 0.0 && totalDepositCount == 0 {
			isWarn = true
		}

		viewData := ViewData{
			"userID":                   userID,
			"isWarn":                   isWarn,
			"user":                     user,
			"lastWithdraw":             "-",
			"info":                     info,
			"manageFeeRate":            info.ManageFeeRate,       // 行政费率
			"rows":                     info.Rows,                // 稽核记录
			"totalRows":                info.TotalRows,           // 记录总数
			"totalFlowNeed":            info.TotalFlowNeed,       // 所需流水
			"totalFlowCurrent":         info.TotalFlowCurrent,    // 当前流水
			"totalFlowLeft":            info.TotalFlowLeft,       // 剩余流水
			"totalMoney":               info.TotalMoney,          // 存款总计
			"totalDiscount":            info.TotalDiscount,       // 优惠总计
			"totalManageFee":           info.TotalManageFee,      // 行政费用总计
			"totalDiscountDeduct":      info.TotalDiscountDeduct, // 总计优惠扣除
			"totalWithdrawCount":       totalWithdrawCount,
			"totalDepositCount":        totalDepositCount,
			"withdrawLastTime":         withdrawLastTime,
			"depositLastTime":          depositLastTime,
			"totalDeposit":             totalDeposit,
			"totalWithdraw":            totalWithdraw,
			"totalPro":                 totalDeposit - totalWithdraw, // 存提差
			"depositAfterLastWithdraw": firstDepositAfterLastWithdraw,
		}
		if info.LastWithdraw > 0 {
			viewData["lastWithdraw"] = time.Unix(info.LastWithdraw, 0).Format("2006-01-02 15:04:05")
		}
		response.Render(c, "user_audits/detail.html", viewData)
	},
	DetailV2: func(c *gin.Context) {
		userID, err := strconv.Atoi(c.DefaultQuery("id", "0"))
		if err != nil || userID <= 0 {
			response.Err(c, err.Error())
			return
		}
		user := models.User{}
		platform := request.GetPlatform(c)
		exists, err := models.Users.FindById(platform, userID, &user)
		if err != nil || !exists {
			response.Err(c, "用户信息查找失败")
			return
		}

		// 刷新最新的稽核信息
		dbSession := common.Mysql(platform, true)
		defer dbSession.Close()
		if err := AuditTrigger.Audit(dbSession, platform, userID); err != nil {
			response.Err(c, err.Error())
			return
		}

		info, err := models.UserAudits.DetailV2(platform, userID)
		if err != nil {
			response.Err(c, err.Error())
			return
		}

		totalWithdrawCount := 0 // 总提款次数
		totalDepositCount := 0  // 总存款次数
		depositLastTime := "-"
		withdrawLastTime := "-"
		totalDeposit := 0.0
		totalWithdraw := 0.0
		rows := models.UserAudits.GetAuditInfo(dbSession, userID, true)
		if err == nil && len(rows) > 0 {
			row := rows[0]
			totalWithdrawCount = getInt(row["withdraw_count"]) // 总计提款数量
			totalDepositCount = getInt(row["deposit_count"])   // 总计存款数量
			totalDeposit = getFloat(row["total_deposit"])      // 总计存款金额
			totalWithdraw = getFloat(row["total_withdraw"])    // 总计提款金额
			if v := getInt(row["withdraw_last_time"]); v > 0 {
				withdrawLastTime = time.Unix(int64(v), 0).Format("2006-01-02 15:04:05") // 最后提款时间
			}
			if v := getInt(row["deposit_last_time"]); v > 0 {
				depositLastTime = time.Unix(int64(v), 0).Format("2006-01-02 15:04:05") // 最后存款时间
			}
		}
		var firstDepositAfterLastWithdraw int64 = 0 // 最后一次成功出款之后的第一次存款时间
		isWarn := false                             // 警告, 此用户没有任何存款信息, 但是账上有钱 ****
		if len(info.Rows) > 0 && totalDeposit == 0.0 && totalDepositCount == 0 {
			isWarn = true // 是否警告信息
		}

		viewData := ViewData{
			"userID":                   userID,
			"isWarn":                   isWarn,
			"user":                     user,
			"lastWithdraw":             "-",
			"info":                     info,
			"manageFeeRate":            info.ManageFeeRate,       // 行政费率
			"rows":                     info.Rows,                // 稽核记录
			"totalRows":                info.TotalRows,           // 记录总数
			"totalFlowNeed":            info.TotalFlowNeed,       // 所需流水
			"totalFlowCurrent":         info.TotalFlowCurrent,    // 当前流水
			"totalFlowLeft":            info.TotalFlowLeft,       // 剩余流水
			"totalMoney":               info.TotalMoney,          // 存款总计
			"totalDiscount":            info.TotalDiscount,       // 优惠总计
			"totalManageFee":           info.TotalManageFee,      // 行政费用总计
			"totalDiscountDeduct":      info.TotalDiscountDeduct, // 总计优惠扣除
			"totalWithdrawCount":       totalWithdrawCount,
			"totalDepositCount":        totalDepositCount,
			"withdrawLastTime":         withdrawLastTime,
			"depositLastTime":          depositLastTime,
			"totalDeposit":             totalDeposit,
			"totalWithdraw":            totalWithdraw,
			"totalPro":                 totalDeposit - totalWithdraw, // 存提差
			"depositAfterLastWithdraw": firstDepositAfterLastWithdraw,
		}
		if info.LastWithdraw > 0 {
			viewData["lastWithdraw"] = time.Unix(info.LastWithdraw, 0).Format("2006-01-02 15:04:05")
		}
		response.Render(c, "user_audits/detail_v2.html", viewData)
	},
	DetailV3: func(c *gin.Context) {
		userID, err := strconv.Atoi(c.DefaultQuery("id", "0"))
		if err != nil || userID <= 0 {
			fmt.Println("获取稽核详情错误原因:", err)
			response.Err(c, err.Error())
			return
		}
		user := models.User{}
		platform := request.GetPlatform(c)
		exists, err := models.Users.FindById(platform, userID, &user)
		if err != nil || !exists {
			response.Err(c, "用户信息查找失败")
			return
		}

		// 刷新最新的稽核信息
		dbSession := common.Mysql(platform, true)
		defer dbSession.Close()
		if err := AuditTrigger.Audit(dbSession, platform, userID); err != nil {
			fmt.Println("err1:", err)
			response.Err(c, err.Error())
			return
		}

		info, err := models.UserAudits.DetailV3(platform, userID)
		if err != nil {
			fmt.Println("err2:", err)
			response.Err(c, err.Error())
			return
		}

		totalWithdrawCount := 0 // 总提款次数
		totalDepositCount := 0  // 总存款次数
		depositLastTime := "-"
		withdrawLastTime := "-"
		totalDeposit := 0.0
		totalWithdraw := 0.0
		rows := models.UserAudits.GetAuditInfo(dbSession, userID, true)
		if err == nil && len(rows) > 0 {
			row := rows[0]
			totalWithdrawCount = getInt(row["withdraw_count"]) // 总计提款数量
			totalDepositCount = getInt(row["deposit_count"])   // 总计存款数量
			totalDeposit = getFloat(row["total_deposit"])      // 总计存款金额
			totalWithdraw = getFloat(row["total_withdraw"])    // 总计提款金额
			if v := getInt(row["withdraw_last_time"]); v > 0 {
				withdrawLastTime = time.Unix(int64(v), 0).Format("2006-01-02 15:04:05") // 最后提款时间
			}
			if v := getInt(row["deposit_last_time"]); v > 0 {
				depositLastTime = time.Unix(int64(v), 0).Format("2006-01-02 15:04:05") // 最后存款时间
			}
		}
		var firstDepositAfterLastWithdraw int64 = 0 // 最后一次成功出款之后的第一次存款时间
		isWarn := false                             // 警告, 此用户没有任何存款信息, 但是账上有钱 ****
		if len(info.Rows) > 0 && totalDeposit == 0.0 && totalDepositCount == 0 {
			isWarn = true // 是否警告信息
		}

		viewData := ViewData{
			"userID":                   userID,
			"isWarn":                   isWarn,
			"user":                     user,
			"lastWithdraw":             "-",
			"info":                     info,
			"manageFeeRate":            info.ManageFeeRate,           // 行政费率
			"rows":                     info.Rows,                    // 稽核记录
			"totalRows":                info.TotalRows,               // 记录总数
			"totalFlowNeed":            info.TotalFlowNeed,           // 所需流水
			"totalFlowCurrent":         info.TotalFlowCurrent,        // 当前流水
			"totalFlowLeft":            info.TotalFlowLeft,           // 剩余流水
			"totalMoney":               info.TotalMoney,              // 存款总计
			"totalDiscount":            info.TotalDiscount,           // 优惠总计
			"totalManageFee":           info.TotalManageFee,          // 行政费用总计
			"totalDiscountDeduct":      info.TotalDiscountDeduct,     // 总计优惠扣除
			"totalWithdrawCount":       totalWithdrawCount,           //
			"totalDepositCount":        totalDepositCount,            //
			"withdrawLastTime":         withdrawLastTime,             //
			"depositLastTime":          depositLastTime,              //
			"totalDeposit":             totalDeposit,                 //
			"totalWithdraw":            totalWithdraw,                // 总计提款
			"totalPro":                 totalDeposit - totalWithdraw, // 存提差
			"depositAfterLastWithdraw": firstDepositAfterLastWithdraw,
			"currentTime":              tools.Now(), // 当前时间
		}
		if info.LastWithdraw > 0 {
			viewData["lastWithdraw"] = time.Unix(info.LastWithdraw, 0).Format("2006-01-02 15:04:05")
		}
		response.Render(c, "user_audits/detail_v2.html", viewData)
	},
	Update: func(c *gin.Context) {
		ID, userID, err := getUserAuditIDs(c)
		if err != nil {
			response.Err(c, err.Error())
			return
		}

		// 获取提交数据信息
		postedData := request.GetPostedData(c)
		getFloatVal := func(key string) float64 {
			v, exists := postedData[key]
			if !exists {
				return 0.0
			}
			val, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
			if err != nil {
				return 0.0
			}
			return val
		}
		flowNeedOld := getFloatVal("flow_old")           // 所需流水旧值
		flowNeed := getFloatVal("flow")                  // 所需流水
		discount := getFloatVal("discount")              // 优惠原额
		discountDeduct := getFloatVal("discount_deduct") // 优惠扣除
		if flowNeed < 0.0 || flowNeedOld < 0.0 || discount < 0.0 || discountDeduct < 0.0 {
			response.Err(c, "获取流水相关信息失败: 值不正确")
			return
		}
		if flowNeed == flowNeedOld && discount == discountDeduct {
			response.Err(c, "稽核相关信息没有变更")
			return
		}

		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		row := models.UserAuditRecord{}
		exists, err := dbSession.SQL("SELECT * FROM user_audits WHERE id = ? AND user_id = ?", ID, userID).Get(&row)
		if err != nil || !exists {
			response.Err(c, "查询稽核记录失败: "+err.Error())
			return
		}

		if row.Deleted > 0.0 {
			response.Err(c, "稽核记录状态异常")
			return
		}

		sql := fmt.Sprintf("UPDATE user_audits SET %s WHERE id = %d AND user_id = %d", func() string {
			tArr := []string{}
			if flowNeed != flowNeedOld {
				if flowNeed == 0.0 {
					tArr = append(tArr, "audit_multiple = 0") // 修改倍数为0
				}
				tArr = append(tArr, fmt.Sprintf("flow_need = %.2f", flowNeed)) // 修改为新值
			}
			if discount != discountDeduct {
				tArr = append(tArr, fmt.Sprintf("discount_deduct = %.2f", discountDeduct))
				tArr = append(tArr, fmt.Sprintf("discount_deducted  = %d", tools.Now()))
			}
			return strings.Join(tArr, ", ")
		}(), ID, userID)
		result, err := dbSession.Exec(sql)
		if err != nil {
			response.Err(c, "修改稽核记录状态失败: "+err.Error())
			return
		}

		affectedCount, err := result.RowsAffected()
		if err != nil || affectedCount == 0 {
			response.Err(c, "修改稽核记录失败: 无效操作")
			return
		}

		response.Ok(c)

	},
	Delete: func(c *gin.Context) {
		ID, userID, err := getUserAuditIDs(c)
		if err != nil {
			response.Err(c, err.Error())
			return
		}

		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		row := models.UserAuditRecord{}
		exists, err := dbSession.SQL("SELECT * FROM user_audits WHERE id = ? AND user_id = ?", ID, userID).Get(&row)
		if err != nil || !exists {
			response.Err(c, "查询稽核记录失败: "+err.Error())
			return
		}

		if row.Deleted > 0 {
			response.Err(c, "稽核记录状态异常")
			return
		}

		sql := fmt.Sprintf("UPDATE user_audits SET deleted = %d WHERE id = %d AND user_id = %d", tools.Now(), ID, userID)
		result, err := dbSession.Exec(sql)
		if err != nil {
			response.Err(c, "修改稽核记录状态失败: "+err.Error())
			return
		}

		affectedCount, err := result.RowsAffected()
		if err != nil || affectedCount == 0 {
			response.Err(c, "修改稽核记录失败: 无效操作")
			return
		}

		response.Ok(c)
	},
	BetsV2: func(c *gin.Context) {
		timeStart, _ := strconv.Atoi(c.DefaultQuery("time_start", "0"))
		timeEnd, _ := strconv.Atoi(c.DefaultQuery("time_end", "0"))
		userID, _ := strconv.Atoi(c.DefaultQuery("user_id", "0"))
		platform := request.GetPlatform(c)
		pConn := pgsql.GetConnForReading(platform)
		if pConn == nil {
			fmt.Println("获取PG连接失败")
			response.Err(c, "获取PG连接失败")
			return
		}
		defer pConn.Close()
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize := 15

		rows := []models.WagerRecord{}
		tableName := "wager_records"
		if currMonth := time.Unix(int64(timeStart), 0).Format("200601"); currMonth == time.Unix(int64(timeEnd), 0).Format("200601") {
			tableName += "_" + currMonth
		}
		wSQL := fmt.Sprintf("user_id = %d AND created_at >= %d AND created_at < %d ", userID, timeStart, timeEnd)
		sql := fmt.Sprintf("SELECT * FROM %s WHERE "+wSQL+"ORDER BY created_at DESC LIMIT %d OFFSET %d", tableName, pageSize, (page-1)*pageSize)
		_, err := pConn.Query(&rows, sql)
		total := 0
		if err != nil {
			fmt.Println("获取记录信息出错:", err)
		}
		totalInfo := struct {
			Total int `json:"total"`
		}{}
		if _, err = pConn.QueryOne(&totalInfo, fmt.Sprintf("SELECT COUNT(*) AS total FROM %s WHERE "+wSQL, tableName)); err == nil {
			total = totalInfo.Total
		}

		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		gameVenus := make([]models.GameVenue, 0)
		if err := dbSession.Table("game_venues").Where("pid > 0 AND venue_type > 0").Find(&gameVenus); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "无法查找到场馆信息: "+err.Error())
			return
		}
		temp := make(map[string]string)
		for _, v := range gameVenus {
			temp[v.Code+"_"+strconv.Itoa(int(v.VenueType))] = v.Name
		}

		newRows := make([]AllWagerRecord, len(rows))
		for i, v := range rows { // 修正 IM 类型错误历史遗留问题
			if v.GameCode == "IM" {
				rows[i].GameType = consts.SportsType
				rows[i].GameCodeType = rows[i].GameCode + "-" + strconv.Itoa(rows[i].GameType)
				v.GameType = consts.SportsType
				v.GameCodeType = v.GameCode + "-" + strconv.Itoa(v.GameType)
			}
			newRows[i].WagerRecord = v
			newRows[i].GameName = temp[v.GameCode+"_"+strconv.Itoa(v.GameType)]
		}

		queryTotal := BetTotalInfo{} // 本次查询总计
		_, err = pConn.QueryOne(&queryTotal, fmt.Sprintf("SELECT "+
			"SUM(bet_money) AS bet_money, "+
			"SUM(valid_money) AS valid_money, "+
			"SUM(net_money) AS net_money, "+
			"SUM(rebate_money) AS rebate_money FROM %s WHERE "+wSQL, tableName))
		if err != nil {
			fmt.Println("获取统计信息出错:", err)
		}
		pageTotal := BetTotalInfo{}
		for _, r := range rows {
			pageTotal.NetMoney += r.NetMoney
			pageTotal.BetMoney += r.BetMoney
			pageTotal.ValidMoney += r.ValidMoney
			pageTotal.RebateMoney += r.RebateMoney
		}

		viewData := ViewData{
			"rows":        newRows,
			"rows_total":  len(newRows),
			"total":       total,
			"page_total":  pageTotal,
			"query_total": queryTotal,
			"game_venue":  caches.GameVenues.All(platform),
			"user_id":     userID,
			"time_start":  timeStart,
			"time_end":    timeEnd,
		}
		SetLoginAdmin(c)
		viewFile := "user_audits/bets_v2.html"
		if request.IsAjax(c) {
			viewFile = "user_audits/_bets_v2.html"
		}
		response.Render(c, viewFile, viewData)
	},
	Sync: func(c *gin.Context) {
		billNo := c.DefaultQuery("bill", "")
		if billNo == "" {
			response.Err(c, "")
			return
		}
		userID, err := strconv.Atoi(c.DefaultQuery("id", "0"))
		if err != nil {
			response.Err(c, err.Error())
			return
		}

		platform := request.GetPlatform(c)
		mConn := common.Mysql(platform)
		defer mConn.Close()
		sql := fmt.Sprintf("SELECT * FROM user_withdraws WHERE user_id = %d ANd bill_no = '%s' LIMIT 1", userID, billNo)
		userWithdraw := models.UserWithdraw{}
		if exists, err := mConn.SQL(sql).Get(&userWithdraw); err != nil {
			response.Err(c, err.Error())
			return
		} else if !exists {
			response.Err(c, "提款记录不存在")
			return
		}

		info, err := models.UserAudits.DetailV2(platform, userID)
		if err != nil {
			response.Err(c, err.Error())
			return
		}
		userWithdraw.FlowTotal = info.TotalFlowNeed
		userWithdraw.FlowCurrent = info.TotalFlowCurrent
		userWithdraw.WithdrawCost = info.TotalManageFee
		if userWithdraw.WithdrawCost == 0.0 { // 如果没有行政费
			userWithdraw.ActualMoney = userWithdraw.Money // 则实际出款等于申请出款
		}

		_, err = mConn.ID(userWithdraw.Id).Update(&userWithdraw)
		if err != nil { // 如果不为空
			response.Err(c, err.Error())
			return
		}

		response.Result(c, map[string]string{
			"withdraw_cost": fmt.Sprintf("%.2f", userWithdraw.WithdrawCost),                    // 行政费用
			"actual_money":  fmt.Sprintf("%.2f", userWithdraw.Money-userWithdraw.WithdrawCost), // 实需出款
		})
	},
}
