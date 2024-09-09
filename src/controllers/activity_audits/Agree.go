package activity_audits

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"xorm.io/builder"
)

func (ths *ActivityAudits) Agree(c *gin.Context) {
	postedData := request.GetPostedData(c)
	idStr, exists := postedData["id"].(string)
	if !exists || exists && idStr == "0" {
		response.Err(c, "id为空")
		return
	}
	platform := request.GetPlatform(c)
	id, _ := strconv.Atoi(idStr)
	dividend := &models.UserDividend{}
	b, err := models.UserDividends.FindById(platform, id, dividend)
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "查询错误")
		return
	}
	if !b {
		response.Err(c, "找不这条记录")
		return
	}
	if dividend.State != 1 {
		response.Err(c, "这条红利申请已经被处理")
		return
	}
	//防止多人同时更改
	redis := common.Redis(platform)
	defer common.RedisRestore(platform, redis)
	rKey := dividend.BillNo + "_" + idStr + "_" + dividend.Username
	num, err := redis.Incr(rKey).Result()
	if err != nil {
		log.Err(err.Error())
		response.Err(c, "redis系统繁忙")
		return
	}
	if num > 1 {
		response.Err(c, "请稍等片刻再试")
		return
	}
	defer redis.Del(rKey)

	uMap := make(map[string]interface{})
	userInfo := &models.User{}
	if exists, err := models.Users.FindById(platform, int(dividend.UserId), userInfo); !exists || err != nil {
		if err != nil {
			log.Logger.Error(err.Error())
		}
		response.Err(c, "查找用户信息失败")
		return
	}
	accountInfo := &models.Account{}
	if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": dividend.UserId})); !exists || err != nil {
		if err != nil {
			log.Logger.Error(err.Error())
		}
		response.Err(c, "查找用户账户信息失败")
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
	successDo := false
	if dividend.MoneyType == 1 { //中心钱包
		transAction := &models.Transaction{}
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"description":   "红利审核-同意",
			"administrator": administrator.Name,
			"admin_user_id": administrator.Id,
			"serial_number": dividend.BillNo,
		}
		transType := consts.TransTypeAdjustmentDividendPlus
		if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, dividend.Money, extraMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, err.Error())
			return
		}
		uMap["before_money"] = accountInfo.Available
		uMap["money"] = dividend.Money
		uMap["after_money"] = accountInfo.Available + dividend.Money
		successDo = true
	} else { //场馆钱包
		//如果外接请求成功successDo = true
		req.SetTimeout(50 * time.Second)
		req.Debug = true
		//先查余额
		headerB := req.Header{
			"Accept": "application/json",
		}
		paramB := req.Param{
			"game_code": dividend.Venue,
		}
		baseBalanceUrl := consts.InternalGameServUrl
		BalanceUrl := baseBalanceUrl + "/game/v1/internal/balance?user_id=" + strconv.Itoa(int(dividend.UserId))
		rB, err := req.Post(BalanceUrl, headerB, paramB)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, err.Error())
			return
		}
		balanceResp := base_controller.BalanceResp{}
		BalanceJsonErr := rB.ToJSON(&balanceResp)
		if BalanceJsonErr != nil {
			log.Logger.Error(BalanceJsonErr.Error())
			response.Err(c, "同意失败")
			return
		}
		if balanceResp.Errcode != 0 { //查询余额失败
			response.Err(c, "同意失败")
			return
		}
		//再转账
		header := req.Header{
			"Accept": "application/json",
		}
		param := req.Param{
			"out_code": "CENTERWALLET",
			"in_code":  dividend.Venue,
			"all_in":   2,
			"money":    dividend.Money,
		}
		baseTransferUrl := consts.InternalGameServUrl
		TransferUrl := baseTransferUrl + "/game/v1/internal/transfer?user_id=" + strconv.Itoa(int(dividend.UserId))
		r, err := req.Post(TransferUrl, header, param)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, err.Error())
			return
		}
		transferResp := base_controller.TransferResp{}
		transferJsonErr := r.ToJSON(&transferResp)
		if transferJsonErr != nil {
			log.Logger.Error(transferJsonErr.Error())
			response.Err(c, "同意失败")
			return
		}
		if transferResp.Errcode != 0 {
			response.Err(c, "同意失败")
			return
		} else {
			uMap["before_money"] = balanceResp.Data.Balance
			uMap["money"] = dividend.Money
			beforeMoney := balanceResp.Data.Balance
			uMap["after_money"] = beforeMoney + dividend.Money
			successDo = true
		}
	}
	if successDo {
		uMap["reviewer_remark"] = postedData["reviewer_remark"]
		uMap["state"] = 2
		uMap["reviewer"] = base_controller.GetLoginAdmin(c).Name
		uMap["updated"] = time.Now().Unix()
		if _, err := session.Table("user_dividends").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "操作失败")
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
