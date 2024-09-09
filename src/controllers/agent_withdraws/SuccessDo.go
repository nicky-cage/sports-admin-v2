package agent_withdraws

import (
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

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *AgentWithdraws) SuccessDo(c *gin.Context) { //成功
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

	userWithdrawInfo := &models.UserWithdraw{}
	if exists, err = models.UserWithdraws.FindById(platform, id, userWithdrawInfo); !exists || err != nil {
		if err != nil {
			log.Logger.Error(err.Error())
		}
		response.Err(c, "此订单不存在")
		return
	}
	if userWithdrawInfo.Status != 1 {
		response.Err(c, "该订单已经被处理")
		return
	}
	accountInfo := &models.Account{}
	if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": userWithdrawInfo.UserId})); !exists || err != nil {
		if err != nil {
			log.Logger.Error(err.Error())
		}
		response.Err(c, "查找用户账户信息失败")
		return
	}
	userInfo := &models.User{}
	if exists, err = models.Users.FindById(platform, int(userWithdrawInfo.UserId), userInfo); !exists || err != nil {
		if err != nil {
			log.Logger.Error(err.Error())
		}
		response.Err(c, "查询用户信息失败")
		return
	}
	administrator := base_controller.GetLoginAdmin(c)
	uMap := map[string]interface{}{
		"status":             2,
		"finance_process_at": time.Now().Format("2006-01-02 15:04:05"),
		"finance_admin":      administrator.Name,
		"remark":             postedData["remark"],
	}
	if postedData["pay_method"] == "bank_card" { //银行卡出款
		bankCardNumStr := postedData["bank_card_num"].(string)
		bankCardNum, _ := strconv.Atoi(bankCardNumStr)
		paymentMethodTemp := ""
		moneyTemp := 0.00
		transactionFee := 0.00
		for i := 1; i <= bankCardNum; i++ {
			paymentMethodTemp += postedData["de_bank_card["+strconv.Itoa(i-1)+"]"].(string) + ","
			deMoney, _ := strconv.Atoi(postedData["de_money["+strconv.Itoa(i-1)+"]"].(string))
			moneyTemp += float64(deMoney)
			deTransactionFee, _ := strconv.Atoi(postedData["de_transaction_fee["+strconv.Itoa(i-1)+"]"].(string))
			if deTransactionFee > deMoney {
				response.Err(c, "手续费必须小于出款金额")
				return
			}
			transactionFee += float64(deTransactionFee)
		}
		fMoney, _ := strconv.ParseFloat(postedData["money"].(string), 64)
		if moneyTemp != fMoney {
			response.Err(c, "银行卡出款金额和用户出款金额不一致")
			return
		}
		uMap["payment_method"] = strings.TrimRight(paymentMethodTemp, ",")
		uMap["card_number"] = strings.TrimRight(paymentMethodTemp, ",")
		uMap["transaction_fee"] = transactionFee
	} else { //代付出款-暂时只有一种
		uMap["business_type"] = 1 //shipu代付
		uMap["payment_method"] = postedData["daifu"].(string)
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
	consuming := time.Now().Unix() - int64(creatTime)
	imap := map[string]interface{}{
		"bill_no":   postedData["bill_no"],
		"type":      2,
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
		"description":   "代理提款-成功",
		"administrator": administrator.Name,
		"admin_user_id": administrator.Id,
		"serial_number": userWithdrawInfo.BillNo,
	}
	transType := consts.TransTypeRechargeAgentWithdrawLess
	var account models.Account
	if account, err = transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, userWithdrawInfo.Money, extraMap); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()
		response.Err(c, err.Error())
		return
	}
	_ = session.Commit()
	//覆盖用户钱包的数据
	if account.Id > 0 {
		_ = account.ResetCacheData(redis)
	}
	response.Ok(c)
}
