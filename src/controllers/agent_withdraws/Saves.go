package agent_withdraws

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *AgentWithdraws) Saves(c *gin.Context) {
	postedData := request.GetPostedData(c)
	idStr, exists := postedData["id"].(string)
	if !exists || exists && idStr == "0" {
		response.Err(c, "id为空")
		return
	}
	id, _ := strconv.Atoi(idStr)
	//防止多人同时更改
	platform := request.GetPlatform(c)
	rClient := common.Redis(platform)
	defer common.RedisRestore(platform, rClient)
	rKey := postedData["bill_no"].(string) + "_" + idStr + "_" + postedData["username"].(string)
	num, err := rClient.Incr(rKey).Result()
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "系统繁忙")
		return
	}
	if num > 1 {
		response.Err(c, "请稍等片刻再试")
		return
	}
	defer rClient.Del(rKey)

	UserWithdrawInfo := &models.UserWithdraw{}
	exists, err = models.UserWithdraws.FindById(platform, id, UserWithdrawInfo)
	if !exists || err != nil {
		response.Err(c, "此订单不存在")
		return
	}

	if UserWithdrawInfo.Status != 1 {
		response.Err(c, "该订单已经被处理")
		return
	}
	admin := base_controller.GetLoginAdmin(c)
	now := time.Now().Unix()
	timeStr := time.Unix(now, 0).Format("2006-01-02 15:04:05")
	myClient := common.Mysql(platform)
	defer myClient.Close()

	switch postedData["type"] {
	case "1": // 事务操作
		if err := myClient.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = myClient.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		uMap := map[string]interface{}{
			"process_step":    3,
			"risk_process_at": timeStr,
			"risk_admin":      admin.Name,
		}
		if _, err := myClient.Table("user_withdraws").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = myClient.Rollback()
			response.Err(c, "更新失败")
			return
		}
		_ = myClient.Commit()
		response.Result(c, "通过")
	case "2":
		sql := "update user_withdraws set `status`='%d', `risk_admin`='%s',`risk_process_at`='%s' ,`risk_remark`='%s'where id='%d'"
		sqll := fmt.Sprintf(sql, 3, admin.Name, timeStr, postedData["remark"].(string), id)
		_, err := myClient.QueryString(sqll)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		csql := "select * from user_withdraws where id=" + postedData["id"].(string)
		wRes, err := myClient.QueryString(csql)
		if err != nil {
			log.Err(err.Error())
		}
		//解冻金额。
		userId, _ := strconv.Atoi(wRes[0]["user_id"])

		accountInfo := &models.Account{}
		if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": userId})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户账户信息失败")
			return
		}
		userInfo := &models.User{}
		if exists, err := models.Users.FindById(platform, userId, userInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户信息失败")
			return
		}

		transType := consts.TransTypeRechargeAgentWithdrawPlus

		// 事务操作
		if err := myClient.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = myClient.Rollback()
			response.Err(c, "事务启动失败")
			return
		}

		transAction := &models.Transaction{}
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"description":   "代理提现拒绝,金额解冻",
			"administrator": admin.Name,
			"admin_user_id": admin.Id,
			"serial_number": tools.GetBillNo("j", 5),
		}
		money, _ := strconv.ParseFloat(wRes[0]["money"], 64)
		if _, err := transAction.AddTransaction(platform, myClient, rClient, userInfo, accountInfo, transType, money, extraMap); err != nil {
			log.Logger.Error(err.Error())
			_ = myClient.Rollback()
			response.Err(c, err.Error())
			return
		}

		_ = myClient.Commit()
		//覆盖用户钱包的数据
		if accountInfo.Id > 0 {
			_ = accountInfo.ResetCacheData(rClient)
		}
		created, _ := strconv.Atoi(wRes[0]["created"])
		cus := time.Now().Unix() - int64(created)
		//往财务日志写记录
		fSql := "insert into finance_logs(bill_no,type,operating,result,operator,remark,created,consuming)"
		_, err = myClient.Exec(fSql, wRes[0]["bill_no"], 2, "代理提款审核", "出款失败", admin.Name, postedData["risk_mark"], time.Now().Unix(), cus)
		if err != nil {
			log.Err(err.Error())
		}
		response.Result(c, "已拒绝")
	}
}
