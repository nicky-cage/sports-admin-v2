package controllers

import (
	"encoding/json"
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-admin/filters"
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

// RiskAudits 风控审核
var RiskAudits = struct {
	View        func(*gin.Context)
	Detail      func(*gin.Context)
	SysAudits   func(*gin.Context)
	Refuse      func(*gin.Context)
	HandUp      func(*gin.Context)
	Receive     func(*gin.Context)
	ReceiveSave func(*gin.Context)
	Saves       func(*gin.Context)
	List        func(*gin.Context)
	SysDetail   func(c *gin.Context)
	CreateSave  func(c *gin.Context)
	State       func(c *gin.Context)
	*ActionCreate
	*ActionSave
	*ActionExport
}{
	View: func(c *gin.Context) {
		response.Render(c, "risk_audits/view.html", pongo2.Context{})
	},
	Refuse: func(c *gin.Context) {
		id := c.Query("id")
		result := c.Query("sys_result")
		vip := c.Query("vip")
		lastMoney := c.Query("last_money")
		response.Render(c, "risk_audits/refuse.html", pongo2.Context{"id": id, "result": result, "vip": vip, "last_money": lastMoney})
	},
	HandUp: func(c *gin.Context) {
		id := c.Query("id")
		response.Render(c, "risk_audits/handup.html", pongo2.Context{"id": id})
	},
	SysAudits: func(c *gin.Context) {
	},
	SysDetail: func(c *gin.Context) {
		response.Render(c, "risk_audits/sysaudits.html", pongo2.Context{})
	},
	Detail: func(c *gin.Context) {
		id := c.Query("id")
		result := c.Query("sys_result")
		vip := c.Query("vip")
		lastMoney := c.Query("last_money")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "SELECT w.username, w.withdraw_cost, w.bank_name, w.bank_realname, w.bank_card, w.bank_address, w.money, w.transaction_fee, w.actual_money,w.wallet_id,w.id, u.label " +
			"FROM user_withdraws AS w, users AS u " +
			"WHERE w.user_id = u.id AND w.id = " + id
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "系统错误: 无法查找相关订单")
		}
		for _, v := range res {
			if v["bank_name"] == "其他银行" {
				sqlCard := "select bank_name from user_cards where card_number = " + v["bank_card"]
				dataCard, err := dbSession.QueryString(sqlCard)
				if err == nil && len(dataCard) > 0 && dataCard[0]["bank_name"] != "" {
					v["bank_name"] = dataCard[0]["bank_name"]
				}
			}
		}

		res[0]["id"] = id
		var cost float64
		temp, _ := strconv.ParseFloat(res[0]["money"], 64)
		temp1, _ := strconv.ParseFloat(res[0]["withdraw_cost"], 64)
		var actualMoney float64 = 0.0
		actualMoney = temp - temp1
		if res[0]["wallet_id"] != "0" {
			var num float64
			rate := tools.GetExchangeRate()
			if models.Parameters.GetValueByInt(platform, "withdraw_auto_rate", 0) == 1 { // 表示不是自动获取汇率
				num = models.Parameters.GetValueByFloat(platform, "withdraw_fixed_rate", 0)
			} else { // 表示自动获取汇率
				temp := models.Parameters.GetValueByFloat(platform, "withdraw_rate_float", 0.02, "提款浮动汇率")
				num = (rate*100 + temp*100) / 100
			}
			cost = actualMoney / num
		}

		viewData := pongo2.Context{"r": res[0], "id": id, "result": result, "vip": vip, "last_money": lastMoney, "cost": tools.ToFixed(cost, 2), "actual_money": actualMoney}
		response.Render(c, "risk_audits/view.html", viewData)
	},
	ActionSave: &ActionSave{
		Model: models.UserWithdraws,
	},
	//人工审核
	Receive: func(c *gin.Context) {
		part := ""
		if username := c.Query("username"); username != "" {
			part += "AND a.username = '" + username + "' "
		}
		if billNo := c.Query("bill_no"); billNo != "" {
			part += "AND a.bill_no = '" + billNo + "' "
		}
		if realName := c.Query("bank_real_name"); realName != "" {
			part += "AND a.bank_realname = '" + realName + "' "
		}
		if money := c.Query("money"); money != "" {
			part += "AND a.money = '" + money + "' "
		}
		if vip := c.Query("vip"); vip != "" {
			part += "AND b.vip = '" + vip + "' "
		}
		if walletId := c.Query("wallet_id"); walletId == "2" {
			part += "AND a.wallet_id > 0 "
		}
		created := c.Query("created")
		if created != "" {
			areas := strings.Split(created, " - ")
			timeStart := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
			timeEnd := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
			part += fmt.Sprintf("AND (a.created >= %d AND a.created < %d) ", timeStart, timeEnd)
		}
		if list := c.Query("list"); list != "" {
			var arr map[string]interface{}
			_ = json.Unmarshal([]byte(list), &arr)
			var str string
			if len(arr) > 0 {
				for _, v := range arr {
					str = str + "'" + v.(string) + "',"
				}
				newStr := strings.Trim(str, ",")
				part += " AND a.bill_no NOT IN (" + newStr + ") "
			}
		}

		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()

		// cond := " GROUP BY a.user_id ORDER BY a.created DESC  LIMIT " + func() string {
		cond := " ORDER BY a.created DESC  LIMIT " + func() string {
			if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err != nil {
				fmt.Println("页面信息获取有误:", err)
				return "0, 15"
			} else {
				return fmt.Sprintf("%d, 15", (page-1)*15)
			}
		}()
		// SELECT @@GLOBAL.sql_mode
		// SELECT @@SESSION.sql_mode
		// dbSession.Exec("SET @@GLOBAL.sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'")
		// dbSession.Exec("SET @@SESSION.sql_mode='STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION'")
		processType := c.DefaultQuery("type", "1") // 处理状态
		wSQL := "WHERE a.status = 1 AND a.risk_admin = '' AND a.process_step = " + processType + " AND a.type = 1 " + part
		sql := "SELECT a.*, b.vip, b.label FROM user_withdraws a INNER JOIN users b ON a.user_id = b.id " + wSQL + cond
		data, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		// 总的统计信息
		tRow := models.TotalInt{}
		_, err = dbSession.SQL("SELECT COUNT(a.id) AS total FROM user_withdraws AS a INNER JOIN users AS b ON a.user_id = b.id " + wSQL).Get(&tRow)
		if err != nil {
			log.Err(err.Error())
		}
		wsql := "SELECT arrive_money, confirm_at FROM user_deposits WHERE user_id = %s AND status = 2 ORDER BY confirm_at DESC LIMIT 1"    // 存款
		aSql := "SELECT money, updated FROM user_account_sets WHERE user_id = %s AND status = 2 AND type= 1 order by updated DESC LIMIT 1" // 上分
		for _, v := range data {
			if v["bank_name"] == "其他银行" {
				sqlCard := "select bank_name from user_cards where card_number= " + v["bank_card"]
				dataCard, err := dbSession.QueryString(sqlCard)
				if err == nil && len(dataCard) > 0 && dataCard[0]["bank_name"] != "" {
					v["bank_name"] = dataCard[0]["bank_name"]
				}
			}
			tempSql := fmt.Sprintf(wsql, v["user_id"])
			res, err := dbSession.QueryString(tempSql)
			if err != nil {
				log.Err(err.Error())
				response.Err(c, "系统错误")
				return
			}
			aSqll := fmt.Sprintf(aSql, v["user_id"])
			aRes, err := dbSession.QueryString(aSqll)
			if err != nil {
				log.Err(err.Error())
				return
			}
			var atime int
			if len(aRes) > 0 {
				atime, _ = strconv.Atoi(aRes[0]["updated"])
				atime = int(tools.MicroToSecond(int64(atime)))
			}
			var wtime int
			if len(res) > 0 {
				wtime, _ = strconv.Atoi(res[0]["confirm_at"])
			}
			if atime > wtime {
				v["last_money"] = aRes[0]["money"]
			} else {
				if wtime > 0 {
					v["last_money"] = res[0]["arrive_money"]
				} else {
					v["last_money"] = "0"
				}
			}
			v["sys"] = "-"
			v["user_label"] = filters.GetUserLabels(platform, v["label"])
		}

		SetLoginAdmin(c)
		viewFile := "risk_audits/_receive_v2.html"
		if processType == "2" {
			viewFile = "risk_audits/_listreceive_v2.html"
		}
		response.Render(c, viewFile, ViewData{"rows": data, "total": tRow.Total})
	},
	ReceiveSave: func(c *gin.Context) { // 领取
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "缺少参数")
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
			response.Result(c, "系统繁忙")
			return
		}
		if num > 1 {
			response.Result(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)
		now := time.Now().Unix()
		timeStr := time.Unix(now, 0).Format("2006-01-02 15:04:05")
		admin := GetLoginAdmin(c)
		imap := map[string]interface{}{
			"risk_admin":       admin.Name,
			"risk_dispatch_at": timeStr,
		}
		session := common.Mysql(platform)
		defer session.Close()
		if _, err := session.Table("user_withdraws").Where("id=?", id).Update(imap); err != nil {
			log.Err(err.Error())
			_ = session.Rollback()
			response.Result(c, "领取失败")
			return
		}
		response.Result(c, "领取成功")
	},
	//处理
	Saves: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, ok := postedData["id"].(string)
		if !ok {
			response.Err(c, "提款单号有误!")
			return
		}
		withdrawID, err := strconv.Atoi(idStr)
		if err != nil {
			response.Err(c, "错误的提款订单号")
			return
		}
		now := time.Now().Unix()
		timeStr := time.Unix(now, 0).Format("2006-01-02 15:04:05")
		admin := GetLoginAdmin(c)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		redisLockKey, err := redis.Lock(platform, "withdraw-"+idStr)
		if err != nil {
			response.Err(c, "请不要重复提交订单")
			return
		}
		defer redis.Unlock(platform, redisLockKey)

		wRes, err := dbSession.QueryString("SELECT * FROM user_withdraws WHERE id = ? ", withdrawID)
		if err != nil || len(wRes) == 0 {
			log.Err("获取提现订单信息有误: %s", err.Error())
			response.Err(c, "无法获取提现订单信息")
			return
		}
		uw := wRes[0]
		if uw["status"] != "1" || uw["process_step"] == "3" {
			response.Err(c, "此提款订单已经被处理过了")
			return
		}

		var lastMoney float64
		if postedData["last_money"] != nil {
			lastMoney, _ = strconv.ParseFloat(postedData["last_money"].(string), 64)
		}

		// 事务处理以下各项
		if err := dbSession.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = dbSession.Rollback()
			response.Err(c, "事务启动失败")
			return
		}

		// 写审核日志
		logSql := "INSERT INTO risk_audit_logs " +
			"(bill_no,user_id, username,vip,money,created,bank_real_name,bank_name,bank_card,address, last_money,sys_result,risk_admin,risk_process_at,remark) " +
			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		switch postedData["type"].(string) {
		case "1": // 通过
			sql := fmt.Sprintf("update user_withdraws set `risk_admin`='%s',`process_step`='%d',`risk_process_at`='%s',`risk_remark`='风控审核通过' where id='%s'", admin.Name, 3, timeStr, idStr)
			_, err := dbSession.QueryString(sql)
			if err != nil {
				dbSession.Rollback()
				response.Err(c, "程序错误: 无法更新提款订单状态")
				return
			}
			_, err = dbSession.Exec(logSql, wRes[0]["bill_no"], wRes[0]["user_id"], wRes[0]["username"], postedData["vip"], wRes[0]["money"], wRes[0]["created"],
				wRes[0]["bank_realname"], wRes[0]["bank_name"], wRes[0]["bank_card"], wRes[0]["address"], lastMoney, postedData["sys_result"], admin.Name, time.Now().Unix(), postedData["remark"])
			if err != nil {
				dbSession.Rollback()
				response.Err(c, "程序错误: 无法写入审核日志")
				return
			}
			dbSession.Commit()
			response.Result(c, "操作成功: 通过成功")
		case "2": // 挂起
			sql := fmt.Sprintf("update user_withdraws set `risk_admin`='',`process_step`='%d',`risk_process_at`='%s',`failure_reason`='%s' where id='%s' ",
				2, timeStr, postedData["failure_reason"].(string), idStr)
			_, err := dbSession.QueryString(sql)
			if err != nil {
				dbSession.Rollback()
				response.Err(c, "程序错误: 无法挂起提款订单")
				return
			}
			dbSession.Commit()
			response.Result(c, "操作成功: 挂起成功")
		case "3": // 拒绝
			sql := fmt.Sprintf("update user_withdraws set `risk_admin`='%s',`status`='%d',`risk_process_at`='%s',`failure_reason`='%s',`risk_remark`='%s',remark='%s' where `id`='%s' ",
				admin.Name, 3, timeStr, postedData["failure_reason"].(string), postedData["remark"].(string), postedData["remark"].(string), idStr)
			_, err := dbSession.QueryString(sql)
			if err != nil {
				dbSession.Rollback()
				response.Err(c, "程序错误: 无法修改提款订单")
				return
			}

			if postedData["message"].(string) == "1" {
				mMap := map[string]interface{}{
					"type":        0,
					"send_type":   2,
					"send_target": wRes[0]["username"],
					"title":       "提款",
					"contents":    postedData["failure_reason"].(string),
					"first_admin": admin.Name,
					"state":       2,
					"created":     tools.NowMicro(),
				}
				if _, err := dbSession.Table("messages").Insert(mMap); err != nil {
					dbSession.Rollback()
					response.Err(c, "程序错误: 无法发送通知给用户")
					return
				}
			}

			_, err = dbSession.Exec(logSql, wRes[0]["bill_no"], wRes[0]["user_id"], wRes[0]["username"], postedData["vip"], wRes[0]["money"], wRes[0]["created"], wRes[0]["bank_realname"],
				wRes[0]["bank_name"], wRes[0]["bank_card"], wRes[0]["address"], lastMoney, postedData["sys_result"], admin.Name, time.Now().Unix(), postedData["remark"])
			if err != nil {
				dbSession.Rollback()
				response.Err(c, "程序错误: 无法写入日志")
				return
			}
			id, _ := strconv.Atoi(wRes[0]["user_id"])
			money, _ := strconv.ParseFloat(wRes[0]["money"], 64)

			//解冻金额。
			accountInfo := &models.Account{}
			if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": id})); !exists || err != nil {
				if err != nil {
					log.Logger.Error(err.Error())
				}
				dbSession.Rollback()
				response.Err(c, "查找用户账户信息失败")
				return
			}
			userInfo := &models.User{}
			if exists, err := models.Users.FindById(platform, id, userInfo); !exists || err != nil {
				if err != nil {
					log.Logger.Error(err.Error())
				}
				dbSession.Rollback()
				response.Err(c, "查找用户信息失败")
				return
			}

			transType := consts.TransTypeWithdrawReturn //事务操作
			transAction := &models.Transaction{}
			extraMap := map[string]interface{}{
				"proxy_ip":      "",
				"ip":            c.ClientIP(),
				"description":   "提现拒绝,金额解冻",
				"administrator": admin.Name,
				"admin_user_id": admin.Id,
				"serial_number": tools.GetBillNo("j", 5),
			}
			redisClient := common.Redis(platform)
			defer common.RedisRestore(platform, redisClient)
			if _, err := transAction.AddTransaction(platform, dbSession, redisClient, userInfo, accountInfo, transType, money, extraMap); err != nil {
				log.Logger.Error("拒绝出款错误:", err)
				_ = dbSession.Rollback()
				response.Err(c, "修改账变信息有误: "+err.Error())
				return
			}

			//覆盖用户钱包的数据
			if accountInfo.Id > 0 {
				_ = accountInfo.ResetCacheData(redisClient)
			}
			created, _ := strconv.Atoi(wRes[0]["created"])
			cus := time.Now().Unix() - int64(created)
			fSql := "INSERT INTO finance_logs(bill_no,type,operating,result,operator,remark,created,consuming)" +
				"VALUES (?, ?, ?, ?, ?, ?, ?, ?)" //往财务日志写记录
			_, err = dbSession.Exec(fSql, wRes[0]["bill_no"], 1, "风控审核", "出款失败", admin.Name, postedData["risk_mark"], time.Now().Unix(), cus)
			if err != nil {
				dbSession.Rollback()
				log.Logger.Error("写入财务日志出错:", err)
				response.Err(c, "写财务日志出错: "+err.Error())
				return
			}
			err = dbSession.Commit()
			if err != nil {
				dbSession.Rollback()
				response.Err(c, "拒绝出错错误: "+err.Error())
				return
			}
			response.Result(c, "操作成功: 拒绝出款")
		}
	},
	List: func(c *gin.Context) {
		var cond string
		limit, offset := request.GetOffsets(c)
		platform := request.GetPlatform(c)
		if offset == 0 {
			cond = " limit 15"
		} else {
			temp := "limit %d,%d"
			cond = fmt.Sprintf(temp, limit, offset)
		}
		admin := GetLoginAdmin(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "SElECT * FROM user_withdraws WHERE status = 1 AND risk_admin = '%s' AND process_step = 1 "
		sqll := fmt.Sprintf(sql, admin.Name)
		res, err := dbSession.QueryString(sqll + cond)
		if err != nil {
			log.Err(err.Error())
			return
		}

		vSql := "SELECT vip, label FROM users WHERE id = '%s' "
		lastSql := "SELECT arrive_money, confirm_at FROM user_deposits WHERE user_id = '%s' AND status = 2 ORDER BY confirm_at DESC LIMIT 1 " // 最后存款时间
		//检查上分的最后一条记录，对比时间。
		aSql := "SELECT money, updated FROM user_account_sets WHERE user_id = %s AND status = 2 AND type = 1 ORDER BY updated DESC LIMIT 1" // 上分
		for _, v := range res {
			vsqll := fmt.Sprintf(vSql, v["user_id"])
			vRes, err := dbSession.QueryString(vsqll)
			if err != nil {
				log.Err(err.Error())
				return
			}
			aSqll := fmt.Sprintf(aSql, v["user_id"])
			aRes, err := dbSession.QueryString(aSqll)
			if err != nil {
				log.Err(err.Error())
				return
			}

			vip, _ := strconv.Atoi(vRes[0]["vip"])
			v["vip"] = strconv.Itoa(vip - 1)
			v["user_label"] = filters.GetUserLabels(platform, vRes[0]["label"])
			lastsqll := fmt.Sprintf(lastSql, v["user_id"])
			lastRes, err := dbSession.QueryString(lastsqll)
			if err != nil {
				log.Err(err.Error())
				return
			}
			var atime int
			if len(aRes) > 0 {
				atime, _ = strconv.Atoi(aRes[0]["updated"])
			}
			var wtime int
			if len(lastRes) > 0 {
				wtime, _ = strconv.Atoi(lastRes[0]["confirm_at"])
			}

			if atime > wtime {
				v["last_money"] = aRes[0]["money"]
			} else {
				if wtime > 0 {
					v["last_money"] = lastRes[0]["arrive_money"]
				} else {
					v["last_money"] = "0"
				}
			}

			id, _ := strconv.Atoi(v["id"])
			userId, _ := strconv.Atoi(v["user_id"])
			info := models.RefreshFlow(platform, id, userId)
			v["flow_total"] = fmt.Sprintf("%.2f", info["flow_total"].(float64))
			v["flow_current"] = fmt.Sprintf("%.2f", info["flow_current"].(float64))
			if v["bank_name"] == "其他银行" {
				sqlCard := "SELECT bank_name FROM user_cards WHERE card_number = '" + v["bank_card"] + "'"
				dataCard, err := dbSession.QueryString(sqlCard)
				if err == nil && len(dataCard) > 0 && dataCard[0]["bank_name"] != "" {
					v["bank_name"] = dataCard[0]["bank_name"]
				}
			}
		}

		if request.IsExportExcel(c) { // 如果是导出
			response.Result(c, res)
			return
		}

		SetLoginAdmin(c)
		viewData := pongo2.Context{
			"rows":      res,
			"total":     15,
			"vipLevels": caches.UserLevels.All(platform),
		}
		viewFile := request.GetViewFile(c, "risk_audits/%sindex.html")
		response.Render(c, viewFile, viewData)
	},
	CreateSave: func(c *gin.Context) {
		postData := request.GetPostedData(c)
		id := postData["id"].(string)
		deposits := postData["deposits"].(string)
		withdrawTimes := postData["withdraw_times"].(string)
		maxMoney := postData["max_money"].(string)
		//一个个处理的把。
		proportion := postData["deposits_proportion"].(string)
		proportionInt, _ := strconv.ParseFloat(proportion, 64)
		var isDepositsProportion string
		if postData["is_deposits_proportion"] == nil {
			isDepositsProportion = "1"
		} else {
			isDepositsProportion = "2"
		}
		//
		var isDeposits string
		if postData["is_deposits"] == nil {
			isDeposits = "1"
		} else {
			isDeposits = "2"
		}
		var isFirstDeposits string
		if postData["is_first_deposits"] == nil {
			isFirstDeposits = "1"
		} else {
			isFirstDeposits = "2"
		}
		var FirstDeposits string
		if postData["first_deposits"] == nil {
			FirstDeposits = "1"
		} else {
			FirstDeposits = "2"
		}
		var isMaxMoney string
		if postData["is_max_money"] == nil {
			isMaxMoney = "1"
		} else {
			isMaxMoney = "2"
		}
		var isMuchDevices string
		if postData["is_much_devices"] == nil {
			isMuchDevices = "1"
		} else {
			isMuchDevices = "2"
		}
		var isRiskLabel string
		if postData["is_risk_label"] == nil {
			isRiskLabel = "1"
		} else {
			isRiskLabel = "2"
		}
		var isWithdrawTimes string
		if postData["is_withdraw_times"] == nil {
			isWithdrawTimes = "1"
		} else {
			isWithdrawTimes = "2"
		}
		var muchDevices string
		if postData["much_devices"] == nil {
			muchDevices = "1"
		} else {
			muchDevices = "2"
		}
		var riskLabel string
		if postData["risk_label"] == nil {
			riskLabel = "1"
		} else {
			riskLabel = "2"
		}
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "update risk_conditions set deposits_proportion=?,is_deposits_proportion=?,withdraw_times=?,is_withdraw_times=?,deposits=?,is_deposits=?,first_deposits=?,is_first_deposits=?,much_devices=?,is_much_devices=?,max_money=?,is_max_money=?,risk_label=?,is_risk_label=? where id=?"
		_, err := dbSession.Exec(sql, proportionInt/100, isDepositsProportion, withdrawTimes, isWithdrawTimes, deposits, isDeposits, FirstDeposits, isFirstDeposits, muchDevices, isMuchDevices, maxMoney, isMaxMoney, riskLabel, isRiskLabel, id)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "修改失败")
			return
		}
		caches.RiskConditions.Load(platform)
		response.Ok(c)
	},
	ActionCreate: &ActionCreate{
		ViewFile: "risk_audits/updated.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			risk := caches.RiskConditions.All(platform)
			return pongo2.Context{"r": risk}
		},
	},
	State: func(c *gin.Context) {},
	ActionExport: &ActionExport{
		Columns: []ExportHeader{
			{Name: "订单编号", Field: "bill_no"},
			{Name: "会员账号", Field: "username"},
			{Name: "会员等级", Field: "vip"},
			{Name: "提款金额", Field: "money"},
			{Name: "行政费用", Field: "withdraw_cost"},
			{Name: "实需出款", Field: "actual_money"},
			{Name: "申请时间", Field: "created"},
			{Name: "银行卡信息", Field: "bank_info"},
			{Name: "最后存款", Field: "last_money"},
			{Name: "总需流水", Field: "flow_total"},
			{Name: "完成流水", Field: "flow_current"},
			{Name: "仍需流水", Field: "flow_left"},
			{Name: "系统审核结果", Field: "risk_admin"},
		},
		ProcessRow: func(m *map[string]interface{}, c *gin.Context) {
			(*m)["bank_info"] = fmt.Sprintf("%s/%s/%s", (*m)["bank_name"], (*m)["bank_realname"], (*m)["bank_card"])
			(*m)["risk_admin"] = ""
			(*m)["actual_money"] = func() string { // 真实费用
				money, err := strconv.ParseFloat((*m)["money"].(string), 64)
				if err != nil {
					return "0.00"
				}
				withdrawCost, err := strconv.ParseFloat((*m)["withdraw_cost"].(string), 64)
				if err != nil {
					return "0.00"
				}
				return fmt.Sprintf("%.2f", money-withdrawCost)
			}()
			(*m)["flow_left"] = func() string { // 余下流水
				flowTotal, err := strconv.ParseFloat((*m)["flow_total"].(string), 64)
				if err != nil {
					return "0.00"
				}
				flowCurrent, err := strconv.ParseFloat((*m)["flow_current"].(string), 64)
				if err != nil {
					return "0.00"
				}
				return fmt.Sprintf("%.2f", flowTotal-flowCurrent)
			}()
			(*m)["created"] = base_controller.FieldToDateTime((*m)["created"].(string))
			(*m)["vip"] = base_controller.FieldToUserVip(c, (*m)["vip"], 0) // 有时不需要减去vip等级
		},
	},
}
