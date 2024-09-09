package agent_commissions

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
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *AgentCommissions) Grant(c *gin.Context) {
	paltfrom := request.GetPlatform(c)
	postData := request.GetPostedData(c)
	id, _ := strconv.Atoi(postData["id"].(string))
	money, _ := strconv.ParseFloat(postData["money"].(string), 64)
	//year, months, _ := time.Now().Date()
	//thisMonth := time.Date(year, months, 1, 0, 0, 0, 0, time.Local)
	//end := thisMonth.AddDate(0, 0, -1).Unix()
	//starts := thisMonth.AddDate(0, -1, 0).Format("2006-01")
	starts := time.Now().Format("2006-01")
	Imap := make(map[string]interface{})
	administrator := base_controller.GetLoginAdmin(c)
	Imap["risk_admin"] = administrator.Name
	Imap["status"] = 2
	Imap["updated"] = time.Now().UnixMicro()
	Imap["risk_process_at"] = time.Now().Unix()
	db := common.Mysql(paltfrom)
	defer db.Close()

	//查看用户钱包
	accountInfo := &models.Account{}
	if exists, err := models.Accounts.Find(paltfrom, accountInfo, builder.NewCond().And(builder.Eq{"user_id": id})); !exists || err != nil {
		if err != nil {
			log.Logger.Error(err.Error())
		}
		response.Err(c, "查找用户账户信息失败")
		return
	}
	//查询用户信息
	userInfo := &models.User{}
	if exists, err := models.Users.FindById(paltfrom, id, userInfo); !exists || err != nil {
		if err != nil {
			log.Logger.Error(err.Error())
		}
		response.Err(c, "查找用户信息失败")
		return
	}
	//查看是否已发放
	var checkAgentCommission models.AgentCommissionLog
	checkBool, _ := db.Where("month=? and user_id=? and status=2", starts, id).Get(&checkAgentCommission)
	if checkBool {
		response.Err(c, userInfo.Username+"本月佣金已发放")
		return
	}
	billNo := tools.GetBillNo("g", 5)
	redis := common.Redis(paltfrom)
	defer common.RedisRestore(paltfrom, redis)
	k := "agents_commissions_grant" + postData["id"].(string)
	num, err := redis.Incr(k).Result()

	if err != nil {
		log.Err(err.Error())
		response.Err(c, "请稍等片刻再试")
		return
	}
	if num > 1 {

		response.Err(c, "请稍等片刻再试")
		return
	}
	defer redis.Del(k)
	//
	//事务操作
	session := common.Mysql(paltfrom)
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Err(err.Error())
		_ = session.Rollback()
		response.Err(c, "事务启动失败")
		return
	}
	//更新代理账号
	extraMap := map[string]interface{}{
		"proxy_ip":      "",
		"ip":            c.ClientIP(),
		"description":   "发放佣金",
		"administrator": administrator.Name,
		"admin_user_id": administrator.Id,
		"serial_number": billNo,
	}
	transAction := &models.Transaction{}
	transType := consts.TransTypeRechargeAgentGrantCommissionPlus
	if _, err := transAction.AddTransaction(paltfrom, session, redis, userInfo, accountInfo, transType, money, extraMap); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()
		response.Err(c, err.Error())
		return
	}
	_ = session.Commit()
	if accountInfo.Id > 0 {
		_ = accountInfo.ResetCacheData(redis)
	}

	//佣金发放记录
	if _, err := db.Table("agent_commission_logs").Where("user_id=? and month=?", id, starts).Update(Imap); err != nil {
		log.Err(err.Error())
		response.Err(c, "更新失败")
		return
	}

	response.Result(c, "已发放")
}
