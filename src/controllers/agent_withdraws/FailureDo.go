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

func (ths *AgentWithdraws) FailureDo(c *gin.Context) { //失败
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
		response.Err(c, "无法获取提现订单信息")
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
	//事务操作
	session := common.Mysql(platform)
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()
		response.Err(c, "事务启动失败")
		return
	}
	administrator := base_controller.GetLoginAdmin(c)
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
		"type":      2,
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
		"description":   "代理提款-失败",
		"administrator": administrator.Name,
		"admin_user_id": administrator.Id,
		"serial_number": userWithdrawInfo.BillNo,
	}
	transType := consts.TransTypeRechargeAgentWithdrawPlus
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
			"title":       "代理提款",
			"first_admin": administrator.Name,
			"state":       2,
			"is_agent":    2,
			"created":     tools.NowMicro(),
		}
		causeFailure := ""
		if postedData["cause_failure"].(string) == "1" {
			causeFailure = "收款卡异常"
		} else if postedData["cause_failure"].(string) == "2" {
			causeFailure = "申请超时"
		} else if postedData["cause_failure"].(string) == "3" {
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
	//覆盖用户钱包的数据
	if accountInfo.Id > 0 {
		_ = accountInfo.ResetCacheData(redis)
	}
	response.Ok(c)
}
