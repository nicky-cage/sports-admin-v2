package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"sports-admin/dao"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/pgsql"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"xorm.io/builder"
)

// AuditCleanLimit 限额
const AuditCleanLimit = 5.0 // 稽核小于固定金额则进行清零

// AuditTrigger 稽核检测
var AuditTrigger = struct {
	Audit func(*xorm.Session, string, int) error // 中心钱包余额
}{
	Audit: func(session *xorm.Session, platform string, userId int) error { // 统计账户总的余额
		// -- 先获取用户中心钱包的余额, 如果余额哆了, 则不再往下进行
		sql := fmt.Sprintf("SELECT balance FROM accounts WHERE user_id = %d", userId)
		res, err := session.QueryString(sql)
		if err != nil || len(res) == 0 {
			log.Logger.Error("[audit]获取用户余额出错:", err)
			return err
		}
		if val, err := strconv.ParseFloat(res[0]["balance"], 64); err != nil {
			log.Logger.Error("[audit]获取用户中心钱包余额转化出错:", err)
			return err
		} else if val >= AuditCleanLimit { // 如果余额大于则不再进行往下执行
			log.Logger.Info("[audit]用户还没有输光, 不需要获取三方账户余额")
			return nil
		}

		// -- 检查用户是否有过存款, 如果没有存款, 就不可能有投注
		var dInfo = struct {
			Total   int   `json:"total"`   // 总计记录数
			Created int64 `json:"created"` // 最后存款时间
		}{}
		if exists, err := session.SQL(fmt.Sprintf("SELECT COUNT(*) AS total, MAX(created) AS created "+
			"FROM user_deposits WHERE user_id = %d AND status = 2", userId)).Get(&dInfo); err != nil {
			log.Logger.Error("[audit]获取用户存款次数出错:", err)
			return err
		} else if !exists || dInfo.Total == 0 { // 表示之前没有存款, 直接返回
			log.Logger.Info("[audit]之前没有存过款,不需要往下再处理")
			return nil // 则不再往下处理
		}
		lastDepositTime := dInfo.Created // 最后存款时间

		// -- 检查是否还有未结算的注单
		pgConn := pgsql.GetConnForReading(platform)
		defer pgConn.Close()
		currentTime := tools.NowMicro()
		var countInfo = struct {
			Total int `json:"total"`
		}{}
		// -- 有存款的话, 先检查是否有过投注
		_, err = pgConn.Query(&countInfo, "SELECT COUNT(*) AS total "+
			"FROM wager_records WHERE user_id = ? AND created_at <= ? AND created_at > ?", userId, currentTime, lastDepositTime)
		if err != nil {
			log.Logger.Error("[Audit] 获取注单统计有误: ", err)
			return err
		} else if countInfo.Total == 0 { // 还没有投过注
			log.Logger.Info("[Audit] 用户还没有投过注")
			return nil
		}

		// 再检检查是否有未结算的注单
		var totalInfo = struct {
			Total int `json:"total"`
		}{}
		_, err = pgConn.Query(&totalInfo, "SELECT COUNT(*) AS total "+
			"FROM wager_records WHERE user_id = ? AND status = 0 AND created_at <= ? AND created_at > ?", userId, currentTime, lastDepositTime)
		if err != nil {
			log.Logger.Error("[audit]从pg获取数据有误:", err)
			return err
		} else if totalInfo.Total > 0 { // 如果有未结算订单 - 则不作处理
			log.Logger.Info("[audit]用户还有未结算注单, 不再往下统计")
			return nil
		}

		// -- 再获取三方账户额, 如果余额够了 则不再往下继续
		var accTotal = 0.0 // 用于从三方账户获取总的余额
		accounts := dao.UserGames.GetAccounts(platform, userId)
		for _, v := range accounts {
			val, err := strconv.ParseFloat(v.Account, 64)
			if err != nil {
				log.Logger.Error("格式化用户余额出错:", err)
				return nil
			} else if val > 0.0 {
				accTotal += val
			}
		}
		if accTotal > AuditCleanLimit { // 如果还没有输光/或小于固定额度, 则不进行处理
			log.Logger.Info("[audit]用户余额:", accTotal, ", 还没有输光, 不需要清空稽核")
			return nil
		}

		// -- 如果没有未结算注单, 并且所有账户余额之后小于 5.0, 则清空稽核
		if totalInfo.Total == 0 { // 得到上次注单清零之后的参与活动或存款
			auditTime := lastDepositTime
			sql := fmt.Sprintf("SELECT id FROM user_audits "+
				"WHERE user_id = %d AND audit_time = %d AND audit_type = 101 LIMIT 1", userId, auditTime)
			if rows, err := session.QueryString(sql); err == nil && len(rows) > 0 { // 已经有相关记录, 则不再写入
				log.Logger.Info("[audit]已经有相关记录, 不再写入")
				return nil
			} else if err != nil {
				log.Logger.Error("[audit]获取最后稽核时间有误:", err)
				return err
			}

			created := tools.NowMicro()    // - 1 // 表示这个时间点之前的全部清零, 不再计算稽核
			updated := created             // 最后更新时间
			depositTime := lastDepositTime // 存款时间
			sql = fmt.Sprintf("INSERT INTO user_audits "+
				"(user_id, audit_id, audit_type, admin_id, admin_name, created, updated, audit_time, deposit_time, remark) VALUES "+
				"(%d, 	   0,        101,        0,        'system',   %d,      %d,      %d,         %d,           '余额%.2f, 不足%.2f, 稽核清零')",
				userId, created, updated, auditTime, depositTime, accTotal, AuditCleanLimit)
			if res, err := session.Exec(sql); err != nil {
				log.Logger.Error("[audit]更新用户稽核信息失败:", err)
				return err
			} else if effected, err := res.RowsAffected(); err != nil || effected == 0 {
				log.Logger.Error("[audit]更新用户稽核信息失败:", err, ", 或者影响行数为零.")
				return err
			}
		}

		return nil
	},
}

// 确认并且保存订单信息 -> 成功
var saveConfirmDeposit = func(platform string, r *models.UserDeposit, depositData, financeData map[string]interface{}, c *gin.Context, args ...*xorm.Session) error {
	currentTime := tools.NowMicro() // 记录最开始时间
	transStarted := false
	var session *xorm.Session
	if len(args) >= 1 { // 表示事务已经启动
		session = args[0]
		transStarted = true
	} else { // 未启动事务, 则需要新建事务
		session = common.Mysql(platform)
		defer session.Close()
	}
	accountInfo := &models.Account{}
	if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": r.UserId})); !exists || err != nil {
		log.Logger.Error("无法获取用户账户信息: ", err)
		return errors.New("查找用户账户信息失败")
	}
	userInfo := &models.User{}
	if exists, err := models.Users.FindById(platform, int(r.UserId), userInfo); !exists || err != nil {
		log.Logger.Error("无法获取用户相关信息: ", err)
		return errors.New("查找用户信息失败")
	}

	// *** 对于用户以前数据进行稽核, 以判断是否需要清空稽核, 如果总余额小于一定金额, 则清空稽核
	if err := AuditTrigger.Audit(session, platform, int(r.UserId)); err != nil {
		return err
	}

	// 判断是否首存
	var userDeposits models.UserDeposit
	depositsBool, _ := session.Table("user_deposits").Where("user_id = ? and status = 2", r.UserId).Get(&userDeposits)
	if !depositsBool { // 判断是否首存
		depositData["is_first_deposit"] = 2
	}

	savedMoney := r.Money + r.Discount
	depositData["status"] = 2               // 表示成功
	depositData["confirm_at"] = currentTime // 确认时间
	depositData["top_money"] = savedMoney   // 真实上分金额
	if !transStarted {                      // 如果还没有启动
		if err := session.Begin(); err != nil { //事务操作
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			return errors.New("事务启动失败")
		}
	}
	if _, err := session.Table("user_deposits").Where("id = ?", r.Id).Update(depositData); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()
		return errors.New("操作失败")
	}

	financeData["result"] = "成功"
	if _, err := session.Table("finance_logs").Insert(financeData); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()
		return errors.New("操作失败")
	}
	transAction := &models.Transaction{}
	administrator := GetLoginAdmin(c)
	extraMap := map[string]interface{}{
		"proxy_ip":      "",
		"ip":            c.ClientIP(),
		"description":   "存款人工确认",
		"administrator": administrator.Name,
		"admin_user_id": administrator.Id,
		"serial_number": tools.GetBillNo("ck", 5),
	}
	transType := consts.TransTypeRechargeOnline
	if r.Type == 2 {
		transType = consts.TransTypeRechargeOffline
	}
	redisSession := common.Redis(platform)
	defer common.RedisRestore(platform, redisSession)
	if _, err := transAction.AddTransaction(platform, session, redisSession, userInfo, accountInfo, transType, r.Money, extraMap); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()
		return err
	}
	if accountInfo.Id > 0 { //覆盖用户钱包的数据
		_ = accountInfo.ResetCacheData(redisSession)
	}

	saveDiscountMoney := r.Discount
	if saveDiscountMoney > 0 { //有优惠
		depositDiscounts := &models.UserDepositDiscount{}
		if _, err := session.Table("user_deposit_discounts").Where("payment_type=?", r.ChannelType).Get(depositDiscounts); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			return err
		}
		offerContent := []byte(depositDiscounts.OfferContent)
		depositDiscountInfo := DepositDiscount{}
		_ = json.Unmarshal(offerContent, &depositDiscountInfo)
		flowMultiple := 1 //流水倍数
		for _, vv := range depositDiscountInfo {
			userVip := "VIP" + (strconv.Itoa(int(userInfo.Vip - 1)))
			if vv.Vip == userVip { //获取用户的vip
				tempFlowMultiple, _ := strconv.Atoi(vv.Multiple)
				flowMultiple = tempFlowMultiple
				break
			}
		}
		divideMap := map[string]interface{}{
			"bill_no":          r.OrderNo,
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
			"created":          currentTime,
			"updated":          currentTime,
		}
		if _, err := session.Table("user_dividends").Insert(divideMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			return err
		}
		accountInfos := &models.Account{}
		if _, err := session.Table("accounts").Where("user_id=?", r.UserId).Get(accountInfos); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			return err
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
		if _, err := transActions.AddTransaction(platform, session, redisSession, userInfo, accountInfos, transTypes, saveDiscountMoney, extraMaps); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			return err
		}
		if accountInfos.Id > 0 { // 覆盖用户钱包的数据
			_ = accountInfos.ResetCacheData(redisSession)
		}
	}

	if !transStarted { // 表明有外部事务调用
		if err := session.Commit(); err != nil {
			fmt.Println("调用存款相关事务处理出错:", err)
			session.Rollback()
		}
	}

	return nil
}
