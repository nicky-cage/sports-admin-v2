package controllers

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
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

// UserWithdrawsStruct 提款
type UserWithdrawsStruct struct {
	models.UserWithdraw `xorm:"extends"`
	Label               string `jsong:"label"`
	Vip                 int32
}

// UserWithdrawSumStruct 提款统计
type UserWithdrawSumStruct struct {
	Money float64
}

// UserWithdraws 提款管理
var UserWithdraws = struct {
	List       func(*gin.Context)
	Success    func(*gin.Context)
	SuccessDo  func(*gin.Context)
	Failure    func(*gin.Context)
	FailureDo  func(*gin.Context)
	GetStatus  func(*gin.Context)
	Notify     func(*gin.Context) // 三方出款的异步回调通知 - 从支付平台过来, 不是直接从第三方平台过来
	SaveConfig func(*gin.Context)
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); !exists {
			currentTime := time.Now().Unix()
			startAt = tools.SecondToMicro(currentTime - currentTime%86400)
			endAt = startAt + tools.SecondToMicro(86400)
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetMicroTimeStampByString(areas[0])
			endAt = tools.GetMicroTimeStampByString(areas[1])
		}
		cond = cond.And(builder.Gte{"user_withdraws.created": startAt}).And(builder.Lte{"user_withdraws.created": endAt})
		cond = cond.And(builder.Eq{"process_step": 3}).And(builder.Eq{"user_withdraws.type": 1}).And(builder.Eq{"user_withdraws.status": 1})
		if min, ok := c.GetQuery("money_min"); ok {
			cond = cond.And(builder.Gte{"user_withdraws.money": min})
		}
		if max, ok := c.GetQuery("money_max"); ok {
			cond = cond.And(builder.Lte{"user_withdraws.money": max})
		}
		username := c.DefaultQuery("username", "")
		billNo := c.DefaultQuery("bill_no", "")
		riskAdmin := c.DefaultQuery("risk_admin", "")
		status := c.DefaultQuery("status", "")
		vipLevel := c.DefaultQuery("vip", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.username": username})
		}
		if len(billNo) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.bill_no": billNo})
		}
		if len(riskAdmin) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.risk_admin": riskAdmin})
		}
		if len(status) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.status": status})
		}
		if vipLevel != "" {
			cond = cond.And(builder.Eq{"users.vip": vipLevel})
		}
		cond = cond.And(builder.Eq{"user_withdraws.wallet_id": 0}) // 钱包id == 0
		limit, offset := request.GetOffsets(c)
		userWithdraws := make([]UserWithdrawsStruct, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()

		//出款方式
		depositCards := make([]models.DepositCard, 0)
		if err := engine.Table("deposit_cards").Find(&depositCards); err != nil {
			log.Logger.Error(err.Error())
		}
		depositCardsNew := make([]string, 0)
		for _, v := range depositCards {
			depositCardsNew = append(depositCardsNew, v.Byname)
		}
		depositCardsNew = append(depositCardsNew, "shipu_daifu")
		total, err := engine.Table("user_withdraws").
			Join("LEFT OUTER", "users", "user_withdraws.user_id = users.id").Where(cond).
			OrderBy("user_withdraws.id DESC").
			Limit(limit, offset).FindAndCount(&userWithdraws)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ss := new(UserWithdrawSumStruct)
		totalApply, _ := engine.Table("user_withdraws").Where(cond).Sum(ss, "money")
		totalCost, _ := engine.Table("user_withdraws").Where(cond).Sum(ss, "withdraw_cost")
		subtotalApply := 0.0
		subtotalCost := 0.0
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		for i := 0; i < len(userWithdraws); i++ {
			subtotalApply += userWithdraws[i].Money
			subtotalCost += userWithdraws[i].WithdrawCost
			// if userWithdraws[i].BankName == "其他银行" {
			if userWithdraws[i].BankName != "" {
				cSQL := "SELECT uc.bank_name, p.name AS province_name, c.name AS city_name, d.name AS district_name " +
					"FROM user_cards AS uc, provinces AS p, cities AS c, districts AS d " +
					"WHERE uc.card_number = '" + userWithdraws[i].BankCard + "' AND uc.province_id = p.id AND uc.city_id = c.id AND uc.district_id = d.id "
				if dataCard, err := dbSession.QueryString(cSQL); err == nil && len(dataCard) >= 1 {
					c := dataCard[0]
					bankText := c["province_name"] + "|" + c["city_name"] + "|" + c["district_name"] + "|"
					if c["bank_name"] != "" && userWithdraws[i].BankName == "其他银行" {
						bankText += "|" + c["bank_name"]
					} else {
						bankText += "|" + userWithdraws[i].BankName
					}
					userWithdraws[i].BankName = bankText
				}
			}
			//} else {
			//	cSQL := "SELECT uc.bank_name, p.name AS province_name, c.name AS city_name, d.name AS district_name " +
			//		"FROM user_cards AS uc, provinces AS p, cities AS c, districts AS d " +
			//		"WHERE uc.card_number = '" + userWithdraws[i].BankCard + "' AND uc.province_id = p.id AND uc.city_id = c.id AND c.district_id = d.id "
			//	fmt.Println("SQL: ", cSQL)
			//	if dataCard, err := dbSession.QueryString(cSQL); err == nil && len(dataCard) >= 1 {
			//		c := dataCard[0]
			//		if c["bank_name"] != "" {
			//			userWithdraws[i].BankName = c["province_name"] + "|" + c["city_name"] + "|" + c["district_name"] + "|" + c["bank_name"]
			//		}
			//	}
			//}
		}
		viewData := pongo2.Context{
			"rows":                  userWithdraws,
			"payment_method":        depositCardsNew,
			"total":                 total,
			"subtotal_apply":        subtotalApply,
			"subtotal_cost":         subtotalCost,
			"subtotal_actual":       subtotalApply - subtotalCost,
			"total_apply":           totalApply,
			"total_cost":            totalCost,
			"total_actual":          totalApply - totalCost,
			"vipLevels":             caches.UserLevels.All(platform),
			"rate":                  6.43, //  tools.GetExchangeRate(),
			"min_withdraw":          models.Parameters.GetValueByFloat(platform, "min_withdraw", 100.0, "单次最小提款金额"),
			"max_withdraw":          models.Parameters.GetValueByFloat(platform, "max_withdraw", 49999.0, "单次最大提款金额"),
			"max_withdraw_day":      models.Parameters.GetValueByFloat(platform, "max_withdraw_day", 999999.0, "单日最大提款金额"),
			"withdraw_usdt_min":     models.Parameters.GetValueByFloat(platform, "withdraw_usdt_min", 20.0, "单次最小USDT提款金额"),
			"withdraw_usdt_max":     models.Parameters.GetValueByFloat(platform, "withdraw_usdt_max", 5000.0, "单次最大USDT提款金额"),
			"withdraw_usdt_max_day": models.Parameters.GetValueByFloat(platform, "withdraw_usdt_max_day", 10000.0, "单日最大USDT提款金额"),
			"withdraw_auto_rate":    models.Parameters.GetValueByInt(platform, "withdraw_auto_rate", 2, "提款自动获取汇率"),
			"withdraw_rate_float":   models.Parameters.GetValueByFloat(platform, "withdraw_rate_float", 0.02, "提款浮动汇率"),
			"withdraw_fixed_rate":   models.Parameters.GetValueByFloat(platform, "withdraw_fixed_rate", 1.0, "提款固定汇率"),
		}

		base_controller.SetLoginAdmin(c)
		viewFile := "user_withdraws/user_withdraws.html"
		if request.IsAjax(c) {
			viewFile = "user_withdraws/_user_withdraws.html"
		}
		response.Render(c, viewFile, viewData)
	},
	Success: func(c *gin.Context) { //成功
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		sql := "select a.*,b.vip,b.label from user_withdraws a left join users b on a.user_id = b.id where a.id = " + idStr
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		data, err := dbSession.QueryString(sql)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		for _, v := range data {
			if v["bank_name"] == "其他银行" {
				sql_card := "select bank_name from user_cards where card_number= " + v["bank_card"]
				data_card, _ := dbSession.QueryString(sql_card)
				if data_card[0]["bank_name"] != "" {
					v["bank_name"] = data_card[0]["bank_name"]
				}
			}
		}
		//虚拟币
		var coin float64
		temp, _ := strconv.ParseFloat(data[0]["money"], 64)
		temp1, _ := strconv.ParseFloat(data[0]["withdraw_cost"], 64)
		var actualMoney float64 = 0.0
		actualMoney = temp - temp1
		if data[0]["wallet_id"] != "0" { // 如果代币提款
			var num float64
			rate := tools.GetExchangeRate()
			if models.Parameters.GetValueByInt(platform, "withdraw_auto_rate", 0) == 1 { // 表示不是自动获取汇率
				num = models.Parameters.GetValueByFloat(platform, "withdraw_fixed_rate", 0)
			} else { // 表示自动获取汇率
				temp := models.Parameters.GetValueByFloat(platform, "withdraw_rate_float", 0.02, "提款浮动汇率")
				num = (rate*100 + temp*100) / 100
			}
			coin = actualMoney / num
		}
		viewData := ViewData{
			"r":            data[0],
			"payouts":      models.PaymentThirds.GetPayouts(platform),
			"coin":         tools.ToFixed(coin, 2),
			"actual_money": actualMoney,
		}
		response.Render(c, "user_withdraws/success.html", viewData)
	},
	SuccessDo: func(c *gin.Context) { //成功
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)

		//防止多人同时更改
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		rKey := postedData["bill_no"].(string) + "_" + idStr + "_" + postedData["username"].(string)
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		if num > 1 {
			response.Err(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)

		// 查找出款信息
		userWithdrawInfo := &models.UserWithdraw{}
		if exists, err := models.UserWithdraws.FindById(platform, id, userWithdrawInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户提现订单失败")
			return
		}
		if userWithdrawInfo.Status != 1 { // 只有为1才会被操作
			response.Err(c, "该订单已经被处理")
			return
		}
		accountInfo := &models.Account{} // 查找用户账户信息
		if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": userWithdrawInfo.UserId})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户账户信息失败")
			return
		}
		userInfo := &models.User{} // 查找用户信息
		if exists, err := models.Users.FindById(platform, int(userWithdrawInfo.UserId), userInfo); !exists || err != nil {
			response.Err(c, "查找用户信息失败")
			return
		}
		administrator := GetLoginAdmin(c)
		uMap := map[string]interface{}{
			"status":             2, // 默认是成功
			"finance_process_at": time.Now().Format("2006-01-02 15:04:05"),
			"finance_admin":      administrator.Name,
			"remark":             postedData["remark"],
		}
		isThirdPayout := false                       // 是否是第三方出款
		thirdPayoutCode := ""                        // 三方代付代码
		if postedData["pay_method"] == "bank_card" { //银行卡出款
			//bankCardNumStr := postedData["bank_card_num"].(string)
			//bankCardNum, err := strconv.Atoi(bankCardNumStr)
			//if err != nil { // 检测卡数量是否正确
			//	response.Err(c, "")
			//	return
			//}
			paymentMethodTemp := ""
			moneyActual := 0.0
			transactionFee := 0.0
			//rep, repErr := regexp.Compile(`^\d{12,22}$`)
			//if repErr != nil {
			//	response.Err(c, "正则表达式的格式错误")
			//	return
			//}
			//for i := 0; i < bankCardNum; i++ {
			//	//bankCard := postedData["de_bank_card["+strconv.Itoa(i)+"]"].(string)
			//	//if matched := rep.MatchString(bankCard); !matched {
			//	//	response.Err(c, "卡号格式错误")
			//	//	return
			//	//}
			//	//paymentMethodTemp += bankCard + ","
			//	deMoney, err := strconv.Atoi(postedData["de_money["+strconv.Itoa(i)+"]"].(string))
			//	if err != nil { // 检测是否有数字错误
			//		response.Err(c, "出款金额有误")
			//		return
			//	}
			//	moneyActual += float64(deMoney)
			//	deTransactionFee, err := strconv.Atoi(postedData["de_transaction_fee["+strconv.Itoa(i)+"]"].(string))
			//	if err != nil {
			//		response.Err(c, "手续费格式有误")
			//		return
			//	}
			//	if deTransactionFee > deMoney {
			//		response.Err(c, "手续费必须小于出款金额")
			//		return
			//	}
			//	transactionFee += float64(deTransactionFee)
			//}
			fMoney, err := strconv.ParseFloat(postedData["money"].(string), 64)
			if err != nil {
				response.Err(c, "出款金额格式有误")
				return
			}
			if moneyActual <= 0.0 {
				moneyActual = userWithdrawInfo.Money - userWithdrawInfo.WithdrawCost // 实际出款 = 提款金额 - 行政费用
			}
			if moneyActual > fMoney {
				response.Err(c, "出款金额不能大于提款金额")
				return
			}
			uMap["payment_method"] = strings.TrimRight(paymentMethodTemp, ",") // 允许为空
			uMap["card_number"] = strings.TrimRight(paymentMethodTemp, ",")    // 允许为空
			uMap["transaction_fee"] = transactionFee                           // 手续费用
			uMap["actual_money"] = moneyActual                                 //实际出款金额
			uMap["updated"] = tools.NowMicro()                                 // 最后时间
		} else { //代付出款-暂时只有一种
			if v, exists := postedData["payout"]; exists {
				isThirdPayout = true
				thirdPayoutCode = v.(string)
				uMap["business_type"] = 1 //shipu代付
				uMap["payment_method"] = thirdPayoutCode
				uMap["actual_money"] = postedData["money"].(string) //实际出款金额
			} else {
				response.Err(c, "未选择任何一种代付方式")
				return
			}
		}
		//事务操作
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		if _, err := session.Table("user_withdraws").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "更新失败")
			return
		}
		creatTime, _ := strconv.Atoi(postedData["created"].(string))
		consuming := time.Now().Unix() - int64(creatTime/1000000)
		imap := map[string]interface{}{
			"bill_no":   postedData["bill_no"],
			"type":      1,
			"operating": "财务出款结束",
			"result":    "成功",
			"operator":  administrator.Name,
			"consuming": consuming,
			"remark":    postedData["remark"],
			"created":   tools.NowMicro(),
		}
		if _, err := session.Table("finance_logs").Insert(imap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "更新失败")
			return
		}
		transAction := &models.Transaction{}
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"description":   "会员提款-成功",
			"administrator": administrator.Name,
			"admin_user_id": administrator.Id,
			"serial_number": userWithdrawInfo.BillNo,
		}
		transType := consts.TransTypeRechargeWithdrawLess
		if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, userWithdrawInfo.Money, extraMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, err.Error())
			return
		}

		if isThirdPayout { // 此处是第三方代付相关操作 - 向支付平台发起支付
			if thirdPayoutCode == "" {
				_ = session.Rollback()
				response.Err(c, "缺少三方代付编码")
				return
			}
			err = models.PaymentThirds.AutoPayout(platform, thirdPayoutCode, userWithdrawInfo)
			if err != nil {
				_ = session.Rollback() // 回滚事务
				response.Err(c, err.Error())
				return
			}
		}

		_ = session.Commit()
		//覆盖用户钱包的数据
		if accountInfo.Id > 0 {
			_ = accountInfo.ResetCacheData(redis)
		}
		response.Ok(c)
	},
	Failure: func(c *gin.Context) { //失败
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		sql := "select a.*,b.vip, b.label from user_withdraws a left join users b on a.user_id=b.id where a.id=" + idStr
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		data, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		viewData := pongo2.Context{"r": data[0]}
		response.Render(c, "user_withdraws/failure.html", viewData)
	},
	FailureDo: func(c *gin.Context) { //失败
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		//防止多人同时更改
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		rKey := postedData["bill_no"].(string) + "_" + idStr + "_" + postedData["username"].(string)
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		if num > 1 {
			response.Err(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)

		userWithdrawInfo := &models.UserWithdraw{}
		if exists, err := models.UserWithdraws.FindById(platform, id, userWithdrawInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户提现订单失败")
			return
		}
		if userWithdrawInfo.Status != 1 {
			response.Err(c, "该订单已经被处理")
			return
		}
		accountInfo := &models.Account{}
		if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().
			And(builder.Eq{"user_id": userWithdrawInfo.UserId})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户账户信息失败")
			return
		}
		userInfo := &models.User{}
		if exists, err := models.Users.FindById(platform, int(userWithdrawInfo.UserId), userInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户信息失败")
			return
		}
		//事务操作
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		administrator := GetLoginAdmin(c)
		uMap := map[string]interface{}{
			"status":             3,
			"finance_process_at": time.Now().Format("2006-01-02 15:04:05"),
			"finance_admin":      administrator.Name,
			"cause_failure":      postedData["cause_failure"],
			"failure_reason":     postedData["failure_reason"],
			"remark":             postedData["remark"],
		}
		if _, err := session.Table("user_withdraws").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "更新失败")
			return
		}
		creatTime, _ := strconv.Atoi(postedData["created"].(string))
		consuming := time.Now().Unix() - int64(creatTime/1000000)
		iMap := map[string]interface{}{
			"bill_no":   postedData["bill_no"],
			"type":      1,
			"operating": "财务出款结束",
			"result":    "失败",
			"operator":  administrator.Name,
			"consuming": consuming,
			"remark":    postedData["remark"],
			"created":   tools.NowMicro(),
		}
		if _, err := session.Table("finance_logs").Insert(iMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "更新失败")
			return
		}
		transAction := &models.Transaction{}
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"description":   "会员提款-失败",
			"administrator": administrator.Name,
			"admin_user_id": administrator.Id,
			"serial_number": userWithdrawInfo.BillNo,
		}
		transType := consts.TransTypeRechargeWithdrawPlus
		if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, userWithdrawInfo.Money, extraMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, err.Error())
			return
		}
		//站内信通知
		if postedData["message"].(string) == "1" {
			mMap := map[string]interface{}{
				"type":        0,
				"send_type":   2,
				"send_target": postedData["username"],
				"title":       "会员提款",
				"first_admin": administrator.Name,
				"state":       2,
				"created":     tools.NowMicro(),
			}
			causeFailure := ""
			if postedData["cause_failure"].(string) == "1" {
				causeFailure = "打码量不足"
			} else if postedData["cause_failure"].(string) == "2" {
				causeFailure = "违规操作"
			} else if postedData["cause_failure"].(string) == "3" {
				causeFailure = "收款卡异常"
			} else if postedData["cause_failure"].(string) == "4" {
				causeFailure = "其他"
			}
			mMap["contents"] = causeFailure
			failureReason := strings.TrimSpace(postedData["failure_reason"].(string))
			if len(failureReason) > 0 { //不为空
				mMap["contents"] = causeFailure + "|" + failureReason
			}
			if _, err := session.Table("messages").Insert(mMap); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "更新失败")
				return
			}
		}
		_ = session.Commit()
		response.Ok(c)
	},
	GetStatus: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		userWithdrawInfo := &models.UserWithdraw{}
		platform := request.GetPlatform(c)
		if exists, err := models.UserWithdraws.FindById(platform, id, userWithdrawInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户提现订单失败")
			return
		}
		data := map[string]interface{}{
			"status": userWithdrawInfo.Status,
		}
		response.Result(c, data)
	},
	Notify: func(c *gin.Context) {
		response.Ok(c)
	},
	SaveConfig: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		postedData := request.GetPostedData(c)
		fieldArr := map[string]string{
			"auto_rate_withdraw":    "提款自动获取汇率",
			"withdraw_usdt_min":     "单次最小USDT提款金额",
			"withdraw_usdt_max":     "单次最大USDT提款金额",
			"withdraw_usdt_max_day": "单日最大USDT提款金额",
			"min_withdraw":          "单次最小提款金额",
			"max_withdraw":          "单次最大提款金额",
			"max_withdraw_day":      "单日最大提款金额",
			"withdraw_rate_float":   "提款浮动汇率",
			"withdraw_auto_rate":    "提款自动获取汇率",
			"withdraw_fixed_rate":   "提款固定汇率",
		}

		for k, v := range postedData {
			for fk, fv := range fieldArr {
				if fk == k {
					models.Parameters.SetValue(platform, k, v, 0, fv)
				}
			}
		}

		response.Ok(c)
	},
}
