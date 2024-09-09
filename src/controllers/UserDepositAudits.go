package controllers

import (
	"encoding/json"
	"sports-admin/caches"
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

type UserDepositAuditsStruct struct {
	models.UserDeposit `xorm:"extends"`
	Label              string `json:"label"`
	Vip                int32
}

type UserDepositAuditSumStruct struct {
	Money    float64
	Discount float64
}

// UserDepositAudits 存款管理-审核列表
var UserDepositAudits = struct {
	List func(*gin.Context)
	// *****************************************************************
	// -*- [警告] -*-
	// *****************************************************************
	// 如果修改此处请务必修改相关联功能代码, 包括但不限于以下几个位置:
	// 1. models/src/PaymentsNotify.go => ProcessNotify
	// *****************************************************************
	Agree  func(*gin.Context) //同意
	Refuse func(*gin.Context) //拒绝
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
		cond = cond.And(builder.Gte{"user_deposits.created": startAt}).And(builder.Lte{"user_deposits.created": endAt})
		cond = cond.And(builder.Eq{"user_deposits.status": 4})
		cond = cond.And(builder.Neq{"user_deposits.type": 3})
		/*if min, ok := c.GetQuery("money_min"); ok {
			cond = cond.And(builder.Gte{"user_deposits.money": min})
		}
		if max, ok := c.GetQuery("money_max"); ok {
			cond = cond.And(builder.Lte{"user_deposits.money": max})
		}*/
		username := c.DefaultQuery("username", "")
		applicant := c.DefaultQuery("applicant", "")
		status := c.DefaultQuery("status", "")
		orderNo := c.DefaultQuery("order_no", "")
		sType := c.DefaultQuery("type", "")
		vipLevel := c.DefaultQuery("vip", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"user_deposits.username": username})
		}
		if len(applicant) > 0 {
			cond = cond.And(builder.Eq{"user_deposits.applicant": applicant})
		}
		if len(status) > 0 {
			cond = cond.And(builder.Eq{"user_deposits.status": status})
		}
		if len(orderNo) > 0 {
			cond = cond.And(builder.Eq{"user_deposits.order_no": orderNo})
		}
		if len(sType) > 0 {
			cond = cond.And(builder.Eq{"user_deposits.type": sType})
		}
		if vipLevel != "" {
			cond = cond.And(builder.Eq{"users.vip": vipLevel})
		}
		limit, offset := request.GetOffsets(c)
		userAuditsDeposits := make([]UserDepositAuditsStruct, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("user_deposits").
			Join("LEFT OUTER", "users", "user_deposits.user_id = users.id").
			Where(cond).
			OrderBy("user_deposits.id DESC").
			Limit(limit, offset).
			FindAndCount(&userAuditsDeposits)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ss := new(UserDepositSumStruct)
		sumTotal, _ := engine.Table("user_deposits").Where(cond).Sums(ss, "money", "discount")
		pageSumMoney := 0.00
		pageSumDiscount := 0.00
		for _, v := range userAuditsDeposits {
			pageSumMoney += v.Money
			pageSumDiscount += v.Discount
		}
		viewData := pongo2.Context{
			"rows":              userAuditsDeposits,
			"total":             total,
			"page_sum_money":    pageSumMoney,
			"page_sum_discount": pageSumDiscount,
			"sum_money":         sumTotal[0],
			"sum_discount":      sumTotal[1],
			"vipLevels":         caches.UserLevels.All(platform),
		}
		viewFile := "user_deposits/audits.html"
		if request.IsAjax(c) {
			viewFile = "user_deposits/_audits.html"
		}
		response.Render(c, viewFile, viewData)
	},
	Agree: func(c *gin.Context) { //同意
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
		rKey := "user_deposits_" + idStr + "_agree"
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

		userDepositInfo := &models.UserDeposit{}
		if exists, err := models.UserDeposits.FindById(platform, id, userDepositInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查询用户存款信息失败")
			return
		}
		if userDepositInfo.Status != 4 {
			if userDepositInfo.Status == 1 {
				response.Err(c, "这笔存款不能被审核")
				return
			} else {
				response.Err(c, "这笔存款已经被处理过了")
				return
			}
		}
		accountInfo := &models.Account{}
		if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": userDepositInfo.UserId})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户账户信息失败")
			return
		}
		userInfo := &models.User{}
		if exists, err := models.Users.FindById(platform, int(userDepositInfo.UserId), userInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户信息失败")
			return
		}
		//判断是否首存
		isFirstDeposit := 1
		engine := common.Mysql(platform)
		defer engine.Close()
		isFirstDepositInfo := &models.UserDeposit{}
		b, err := engine.Table("user_deposits").Where("status=2").Get(isFirstDepositInfo)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查找用户是否首存失败")
			return
		}
		if !b {
			isFirstDeposit = 2
		}
		//事务操作
		session := engine
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		administrator := GetLoginAdmin(c)
		uMap := map[string]interface{}{
			"top_money":        userDepositInfo.ArriveMoney + userDepositInfo.Discount,
			"status":           2,
			"confirm_at":       tools.Now(),
			"finance_admin":    administrator.Name,
			"is_first_deposit": isFirstDeposit,
		}
		if _, err := session.Table("user_deposits").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "确认失败")
			return
		}
		creatTime := userDepositInfo.Updated
		consuming := time.Now().Unix() - int64(creatTime)
		iMap := map[string]interface{}{
			"bill_no":   userDepositInfo.OrderNo,
			"type":      0,
			"operating": "存款结束",
			"result":    "成功",
			"operator":  administrator.Name,
			"consuming": consuming,
			"created":   time.Now().UnixMicro(),
		}
		if _, err := session.Table("finance_logs").Insert(iMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "确认失败")
			return
		}
		transAction := &models.Transaction{}
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"description":   "存款手动审核确认",
			"administrator": administrator.Name,
			"admin_user_id": administrator.Id,
			"serial_number": userDepositInfo.OrderNo,
		}
		transType := consts.TransTypeRechargeOnline
		if userDepositInfo.Type == 2 {
			transType = consts.TransTypeRechargeOffline
		}
		if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, userDepositInfo.ArriveMoney, extraMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, err.Error())
			return
		}

		//覆盖用户钱包的数据
		if accountInfo.Id > 0 {
			_ = accountInfo.ResetCacheData(redis)
		}
		//优惠
		savedMoney := userDepositInfo.ArriveMoney + userDepositInfo.Discount //上分金额=到账金额+存款优惠
		saveDiscountMoney := userDepositInfo.Discount
		if saveDiscountMoney > 0 {
			depositDiscounts := &models.UserDepositDiscount{}
			if _, err := session.Table("user_deposit_discounts").Where("payment_type=?", userDepositInfo.ChannelType).Get(depositDiscounts); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}
			offerContent := []byte(depositDiscounts.OfferContent)
			depositDiscountInfo := DepositDiscount{}
			_ = json.Unmarshal(offerContent, &depositDiscountInfo)
			flowMultiple := 1 //流水倍数
			//dayMaxDiscount := 0.00 //每日最高优惠
			for _, vv := range depositDiscountInfo {
				userVip := "VIP" + (strconv.Itoa(int(userInfo.Vip - 1)))
				if vv.Vip == userVip { //获取用户的vip
					tempFlowMultiple, _ := strconv.Atoi(vv.Multiple)
					//tempDayMaxDiscount, _ := strconv.ParseFloat(vv.DayMaxDiscount, 64)
					flowMultiple = tempFlowMultiple
					/*dayMaxDiscount = tempDayMaxDiscount
					if saveDiscountMoney > dayMaxDiscount { //优惠大于每日最高优惠
						saveDiscountMoney = dayMaxDiscount
					}*/
					break
				}
			}
			divideMap := map[string]interface{}{
				"bill_no":          userDepositInfo.OrderNo,
				"username":         userInfo.Username,
				"user_id":          userInfo.Id,
				"top_name":         userInfo.TopName,
				"top_id":           userInfo.TopId,
				"type":             5, //存款红利
				"venue":            "",
				"before_money":     accountInfo.Available,
				"money":            saveDiscountMoney,
				"after_money":      accountInfo.Available + saveDiscountMoney,
				"money_type":       1, //中心钱包
				"operation_type":   2, //单会员发放
				"flow_limit":       2, //需要流水
				"flow_multiple":    flowMultiple,
				"applicant":        GetLoginAdmin(c).Name,
				"applicant_remark": "存款优惠红利",
				"vip":              userInfo.Vip,
				"turnover_amount":  (savedMoney + saveDiscountMoney) * float64(flowMultiple), //所需要的流水金额
				"reviewer_remark":  "存款优惠自动送红利",
				"state":            2,
				"reviewer":         GetLoginAdmin(c).Name,
				"created":          tools.NowMicro(),
				"updated":          tools.NowMicro(),
			}
			if _, err := session.Table("user_dividends").Insert(divideMap); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}
			accountInfos := &models.Account{}
			if _, err := session.Table("accounts").Where("user_id=?", userDepositInfo.UserId).Get(accountInfos); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}
			transActions := &models.Transaction{}
			extraMaps := map[string]interface{}{
				"proxy_ip":      "",
				"ip":            c.ClientIP(),
				"description":   "存款红利",
				"administrator": administrator.Name,
				"admin_user_id": administrator.Id,
				"serial_number": tools.GetBillNo("hl", 5),
			}
			transTypes := consts.TransTypeAdjustmentDividendPlus
			if _, err := transActions.AddTransaction(platform, session, redis, userInfo, accountInfos, transTypes, saveDiscountMoney, extraMaps); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}
			if accountInfos.Id > 0 { //覆盖用户钱包的数据
				_ = accountInfos.ResetCacheData(redis)
			}
		}
		_ = session.Commit()
		response.Message(c, "操作成功")
	},
	Refuse: func(c *gin.Context) { //拒绝
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
		rKey := "user_deposits_" + idStr + "_refuse"
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

		userDepositInfo := &models.UserDeposit{}
		if exists, err := models.UserDeposits.FindById(platform, id, userDepositInfo); !exists || err != nil {
			response.Err(c, "查找用户存款信息失败")
			return
		}
		if userDepositInfo.Status != 4 {
			if userDepositInfo.Status == 1 {
				response.Err(c, "这笔存款不能被审核")
				return
			} else {
				response.Err(c, "这笔存款已经被处理过了")
				return
			}
		}
		administrator := GetLoginAdmin(c)
		uMap := map[string]interface{}{
			"status":        3,
			"finance_admin": administrator.Name,
			"updated":       tools.NowMicro(),
		}
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		if _, err := session.Table("user_deposits").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "操作失败")
			return
		}
		creatTime := userDepositInfo.Updated
		consuming := time.Now().Unix() - int64(creatTime)
		imap := map[string]interface{}{
			"bill_no":   userDepositInfo.OrderNo,
			"type":      0,
			"operating": "存款结束",
			"result":    "失败",
			"operator":  administrator.Name,
			"consuming": consuming,
			"created":   tools.NowMicro(),
		}
		if _, err := session.Table("finance_logs").Insert(imap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "确认失败")
			return
		}
		_ = session.Commit()
		response.Message(c, "操作成功")
	},
}
