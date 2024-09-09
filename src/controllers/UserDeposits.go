package controllers

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/redis"
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

type UserDepositsStruct struct {
	models.UserDeposit `json:"user_deposit" xorm:"extends"`
	Label              string `json:"label" `
	Vip                int32  `json:"vip"`
}

type UserDepositSumStruct struct {
	Money    float64 `json:"money"`
	Discount float64 `json:"discount"`
}

// UserDeposits 存款管理-存款列表
var UserDeposits = struct {
	*ActionUpdate
	*ActionSave
	List        func(*gin.Context)
	ConfirmDo   func(*gin.Context) // 人工确认
	GetStatus   func(*gin.Context) // 获取状态
	AddSlip     func(*gin.Context) // 添加存款单页面
	AddSlipSave func(*gin.Context) // 保存存款单
	OrderInfo   func(*gin.Context) // 查询订单信息
	UserInfo    func(*gin.Context) // 用户信息
	*ActionExport
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		if value, exists := c.GetQuery("created"); exists {
			areas := strings.Split(value, " - ")
			startAt, endAt := tools.GetMicroTimeStampByString(areas[0]), tools.GetMicroTimeStampByString(areas[1])
			cond = cond.And(builder.Gte{"user_deposits.created": startAt}).And(builder.Lte{"user_deposits.created": endAt})
		}
		isOnline := func() bool { // 判断是否在线存款
			if isOnline, exists := c.Get("is_online"); exists {
				if isRealOnline, ok := isOnline.(bool); ok && isRealOnline {
					return true
				}
			}
			return false
		}()
		if isOnline {
			cond = cond.And(builder.Eq{"user_deposits.type": 1}) // 1: 在线
		} else {
			if val, exists := c.GetQuery("type"); exists && (val == "2" || val == "4") {
				depositType, _ := strconv.Atoi(val)
				cond = cond.And(builder.Eq{"user_deposits.type": depositType}) // 2: 离线
			} else {
				cond = cond.And(builder.In("user_deposits.type", 2)) // 2: 离线
			}
		}
		// cond = cond.And(builder.Neq{"user_deposits.type": 3}) // 不是代客充值
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
		// 默认条件
		queryNoStatus, _ := dbSession.Table("user_deposits").Where(cond).Sums(&queryRow, "money", "discount")
		// 加上 成功状态
		cond = cond.And(builder.Eq{"status": 2})
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
		totalType := "2, 4"
		if isOnline {
			totalType = "1"
		}
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

		// 获取可用的支付方式 - 在线
		payments := make([]models.Payment, 0)
		if err := dbSession.Table("payments").Where("is_online=2").Find(&payments); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取商户信息错误")
			return
		}

		paymentChannels := &[]models.Payment{}
		depositCards := &[]models.DepositCard{}
		_ = models.Payments.FindAllNoCount(platform, paymentChannels)
		_ = models.DepositCards.FindAllNoCount(platform, depositCards)
		viewData := pongo2.Context{
			"merchant": payments,
			"rows":     userDeposits,
			"total":    queryTotal,

			"page_total":              pageInfo.Total,                                              // page
			"page_money":              pageInfo.Money,                                              // page
			"page_success":            pageInfo.Success,                                            // page
			"page_success_money":      pageInfo.SuccessMoney,                                       // page
			"page_success_discount":   pageInfo.SuccessDiscount,                                    // page
			"page_success_rate":       float64(pageInfo.Success) / float64(pageInfo.Total) * 100.0, //
			"page_success_money_rate": pageInfo.SuccessMoney / pageInfo.Money * 100.0,              //
			"page_discount":           pageInfo.Discount,

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

			"channelTypes":    consts.PaymentTypes,
			"vipLevels":       caches.UserLevels.All(platform),
			"depositCards":    depositCards,
			"paymentChannels": paymentChannels,
			"depositVirtuals": models.DepositVirtuals.GetVirtuals(platform),
		}

		viewFile := func() string {
			isOnlineFile := func() string {
				if isOnline {
					return "_online"
				}
				return ""
			}()
			if request.IsAjax(c) {
				return fmt.Sprintf("user_deposit%ss/_user_deposits.html", isOnlineFile)
			}
			return fmt.Sprintf("user_deposit%ss/user_deposits.html", isOnlineFile)
		}()
		SetLoginAdmin(c)
		response.Render(c, viewFile, viewData)
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.UserDeposits,
		ViewFile: "user_deposits/user_deposits_edit.html",
		Row: func() interface{} {
			return &models.UserDeposit{}
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			depositCards := []models.DepositCard{}
			platform := request.GetPlatform(c)
			engine := common.Mysql(platform)
			defer engine.Close()
			if err := engine.Table("deposit_cards").Where("status=2").Find(&depositCards); err != nil {
				log.Logger.Error(err.Error())
			}
			payments := make([]models.Payment, 0)
			if err := engine.Table("payments").Where("is_online=2").Find(&payments); err != nil {
				log.Logger.Error(err.Error())
			}
			userLabel := func() string {
				rId, err := strconv.Atoi(c.DefaultQuery("id", "0"))
				if err != nil || rId <= 0 {
					return ""
				}
				sql := fmt.Sprintf("SELECT users.label FROM users, user_deposits WHERE users.id = user_deposits.user_id AND user_deposits.id = %d", rId)
				if rows, err := engine.QueryString(sql); err == nil && len(rows) > 0 {
					return rows[0]["label"]
				}
				return ""
			}()
			depositVirtuals := models.DepositVirtuals.GetVirtuals(platform)
			return pongo2.Context{
				"rows":             depositCards,
				"rs":               payments,
				"user_label":       userLabel,
				"deposit_virtuals": depositVirtuals,
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.UserDeposits,
		SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
			var startAt int64
			var endAt int64
			if value, exists := (*data)["created"].(string); !exists {
				currentTime := time.Now().Unix()
				startAt = currentTime - currentTime%86400
				endAt = startAt + 86400
			} else {
				areas := strings.Split(value, " - ")
				startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
				endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
			}
			(*data)["time_start"] = startAt
			(*data)["time_end"] = endAt
			delete(*data, "created")
			return nil
		},
	},
	AddSlip: func(c *gin.Context) {
		depositCards := make([]models.DepositCard, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if err := engine.Table("deposit_cards").Where("status=2").Find(&depositCards); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查村收款银行卡错误")
			return
		}
		viewData := pongo2.Context{
			"rows":         depositCards,
			"payTypes":     consts.PaymentTypes,
			"channelTypes": models.Payments.All(platform),
		}
		response.Render(c, "user_deposits/add_slips.html", viewData)
	},
	AddSlipSave: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		stype := postedData["type"].(string)
		username := postedData["username"].(string)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		userInfo := &models.User{}
		b, err := engine.Table("users").Where("username=?", username).Get(userInfo)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查询用户失败")
			return
		}
		if !b {
			response.Err(c, "用户不存在失败")
			return
		}
		administrator := GetLoginAdmin(c)
		imap := map[string]interface{}{
			"order_no":  tools.GetBillNo("D", 5),
			"user_id":   userInfo.Id,
			"top_code":  userInfo.TopCode,
			"top_name":  userInfo.TopName,
			"top_id":    userInfo.TopId,
			"type":      stype,
			"money":     postedData["money"],
			"username":  postedData["username"],
			"comment":   postedData["comment"],
			"status":    1,
			"applicant": administrator.Name,
			"created":   tools.NowMicro(),
			"updated":   tools.NowMicro(),
		}
		if stype == "1" { //在线存款
			payCode := postedData["pay_code"].(string)
			payCodeInfo := &models.Payment{}
			b, err := engine.Table("payments").Where("code=?", payCode).Get(payCodeInfo)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "查询支付编码错误")
				return
			}
			if !b {
				response.Err(c, "支付编码不存在")
				return
			}
			imap["channel_type"] = postedData["channel_type"]
			imap["pay_code"] = payCode
			imap["business_id"] = payCodeInfo.Id
			imap["business_name"] = payCode
		} else { //离线存款
			if accountByNameVal, exists := postedData["account_by_name"]; exists { // 如果有存款卡号相关信息
				accountByName := accountByNameVal.(string)
				if accountByName != "" {
					tempStr := strings.Split(accountByName, "-")
					depositCardInfo := &models.DepositCard{}
					if _, err := engine.Table("deposit_cards").Where("bank_card=?", tempStr[2]).Get(depositCardInfo); err != nil {
						log.Logger.Error(err.Error())
						response.Err(c, "查询收款银行卡错误")
						return
					}
					imap["card_number_id"] = depositCardInfo.Id
					imap["card_number"] = tempStr[2]
				}
				imap["account_by_name"] = accountByName
			}
			if depositName, exists := postedData["deposit_name"]; exists { // 如果有存款姓名
				imap["deposit_name"] = depositName
			}
		}
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		if _, err := session.Table("user_deposits").Insert(imap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "添加失败")
			return
		}
		iMap := map[string]interface{}{
			"bill_no":   postedData["order_no"],
			"type":      0,
			"operating": "后台手动添加",
			"result":    "成功",
			"operator":  administrator.Name,
			"created":   tools.NowMicro(),
		}
		if v, exists := postedData["remark"]; exists { // 如果有备注, 再添加备注
			iMap["remark"] = v
		}
		if _, err := session.Table("finance_logs").Insert(iMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "操作失败")
			return
		}
		_ = session.Commit()
		response.Ok(c)
	},
	ConfirmDo: func(c *gin.Context) { //手动确认
		postedData := request.GetPostedData(c)
		// 关于基础数据的验证 ----------------------------------------------------------------------------------
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, idErr := strconv.Atoi(idStr)
		if idErr != nil {
			response.Err(c, "存款id提交错误")
			return
		}
		submit := postedData["submit"].(string) // 这一行偶尔会提示, nil 错误

		arriveMoneyStr := postedData["arrive_money"].(string)
		arriveMoney, err := strconv.ParseFloat(arriveMoneyStr, 64)
		if err != nil || arriveMoney < 0 {
			response.Err(c, "到账金额有误: "+arriveMoneyStr)
			return
		}
		confirmMoneyStr := postedData["confirm_money"].(string)
		confirmMoney, err := strconv.ParseFloat(confirmMoneyStr, 64)
		if err != nil || confirmMoney < 0 {
			response.Err(c, "确认金误有误: "+confirmMoneyStr)
			return
		}

		platform := request.GetPlatform(c)
		//防止多人同时更改
		rKey, err := redis.Lock(platform, "confirm-deposit-"+idStr)
		if err != nil {
			log.Err(err.Error())
			fmt.Println("缓存服务器加锁失败: ", err)
			response.Err(c, "请不要同一时间内多次提交")
			return
		}
		defer redis.Unlock(platform, rKey)

		administrator := GetLoginAdmin(c)
		r := &models.UserDeposit{}
		if exists, err := models.UserDeposits.FindById(platform, id, r); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查询用户存款信息失败")
			return
		}
		if r.Status != 1 {
			response.Err(c, "该订单已经被操作过")
			return
		}

		// 用于保存到用户存款记录表
		depositData := map[string]interface{}{
			"arrive_money":     arriveMoney,
			"confirm_money":    confirmMoney,
			"remark":           postedData["remark"],
			"status":           3,
			"finance_admin":    administrator.Name,
			"updated":          time.Now().UnixMicro(),
			"is_first_deposit": 1,
		}
		newAccName := postedData["account_by_name"].(string)
		fRemark := postedData["remark"].(string)
		if newAccName != r.AccountByName {
			depositData["account_by_name"] = newAccName
			fRemark += ", 收款账号/钱包由 [" + r.AccountByName + "] 改为 [" + newAccName + "]"
		}

		// 用于保存到财务日志表
		financeData := map[string]interface{}{
			"bill_no":   postedData["order_no"],
			"type":      0,
			"operating": "存款结束",
			"result":    "失败",
			"operator":  administrator.Name,
			"consuming": time.Now().Unix() - int64(r.Created/1000000),
			"remark":    fRemark,
			"created":   tools.NowMicro(),
		}

		if submit == "2" { //失败按钮 -----------------------------
			session := common.Mysql(platform)
			defer session.Close()
			if err := session.Begin(); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "事务启动失败")
				return
			}
			depositData["status"] = 3
			if _, err := session.Table("user_deposits").Where("id=?", id).Update(depositData); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "操作失败")
				return
			}
			if _, err := session.Table("finance_logs").Insert(financeData); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "操作失败")
				return
			}

			if r.CardNumberId > 0 { // 如果是失败,则对相应的银行卡进行累减
				if err := models.DepositCards.ReduceUsedMoney(platform, int(r.CardNumberId), r.Money, int(r.Created), session); err != nil {
					log.Logger.Error(err)
					_ = session.Rollback()
					response.Err(c, "扣减银行卡相关信息失败")
				}
			}
			_ = session.Commit()
			response.Message(c, "已将订单状态设置为失败")
			return
		}

		// 以下处理成功按钮 ---------------------------------------------
		if err := saveConfirmDeposit(platform, r, depositData, financeData, c); err != nil {
			response.Err(c, err.Error())
			return
		}

		response.Message(c, "操作成功")
	},
	GetStatus: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		userDepositInfo := &models.UserDeposit{}
		platform := request.GetPlatform(c)
		if exists, err := models.UserDeposits.FindById(platform, id, userDepositInfo); !exists || err != nil {
			response.Err(c, "查找用户存款信息失败")
			return
		}
		data := map[string]interface{}{
			"status": userDepositInfo.Status,
		}
		response.Result(c, data)
	},
	OrderInfo: func(c *gin.Context) {
		orderNumber := c.DefaultQuery("order_number", "")
		if orderNumber == "" {
			response.Err(c, "缺少订单号码")
			return
		}
		result, err := models.PaymentThirds.OrderInfo(orderNumber)
		if err != nil {
			response.Err(c, err.Error())
			return
		}

		response.Result(c, result)
	},
	UserInfo: func(c *gin.Context) {
		Platforms := request.GetPlatform(c)
		db := common.Mysql(Platforms)
		defer db.Close()
		username := c.Query("username")
		sql := "select a.realname,b.balance from users a join accounts b on a.username=b.username where a.username='" + username + "'"
		res, _ := db.QueryString(sql)
		if len(res) == 0 {
			response.Err(c, "用户不存在")
		} else {
			data := map[string]interface{}{
				"realname": res[0]["realname"],
				"money":    res[0]["balance"],
			}
			response.Result(c, data)
		}
	},
	ActionExport: &ActionExport{
		Columns: []ExportHeader{
			{"序号", "id"},
			{"订单编号", "order_no"},
			{"会员编号", "user_id"},
			{"会员名称", "username"},
			{"会员等级", "vip"},
			{"订单金额", "money"},
			{"到账金额", "arrive_money"},
			{"存款优惠", "discount"},
			{"订单时间", "created"},
			{"收款银行/姓名/卡号", "account_by_name"},
			{"操作人", "finance_admin"},
			{"状态", "status"},
		},
		ProcessRawData: func(fields []string, rArr *[]map[string]interface{}, c *gin.Context) {
			base_controller.ExportRawDataReset(rArr)
		},
		ProcessRow: func(m *map[string]interface{}, c *gin.Context) {
			(*m)["id"] = (*m)["user_deposit.id"]
			(*m)["order_no"] = (*m)["user_deposit.order_no"]
			(*m)["username"] = (*m)["user_deposit.username"]
			(*m)["vip"] = base_controller.FieldToUserVip(c, (*m)["vip"])
			(*m)["money"] = fmt.Sprintf("%.2f", (*m)["user_deposit.money"].(float64))
			(*m)["arrive_money"] = fmt.Sprintf("%.2f", (*m)["user_deposit.arrive_money"].(float64))
			(*m)["discount"] = fmt.Sprintf("%.2f", (*m)["user_deposit.discount"].(float64))
			(*m)["created"] = func() string {
				if (*m)["user_deposit.created"] == nil {
					return ""
				}
				return base_controller.FieldToDateTime(fmt.Sprintf("%d", int((*m)["user_deposit.created"].(float64))))
			}()
			(*m)["finance_admin"] = (*m)["user_deposit.finance_admin"]
			(*m)["account_by_name"] = (*m)["user_deposit.account_by_name"]
			(*m)["status"] = func() string {
				switch int((*m)["user_deposit.status"].(float64)) {
				case 1:
					return "待确认"
				case 2:
					return "成功"
				case 3:
					return "失败"
				}
				return "未知"
			}()
		},
	},
}
