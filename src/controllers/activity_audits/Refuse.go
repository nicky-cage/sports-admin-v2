package activity_audits

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
	"github.com/imroc/req"
	"xorm.io/builder"
)

func (ths *ActivityAudits) Refuse(c *gin.Context) {
	postedData := request.GetPostedData(c)
	idStr, exists := postedData["id"].(string)
	if !exists || exists && idStr == "0" {
		response.Err(c, "id为空")
		return
	}
	id, _ := strconv.Atoi(idStr)
	dividend := &models.UserDividend{}
	platform := request.GetPlatform(c)
	b, err := models.UserDividends.FindById(platform, id, dividend)
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "查询错误")
		return
	}
	if !b {
		response.Err(c, "找不这条红利申请记录")
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
		log.Logger.Error(err.Error())
		response.Err(c, "redis系统繁忙")
		return
	}
	if num > 1 {
		response.Err(c, "请稍等片刻再试")
		return
	}
	defer redis.Del(rKey)
	//更改状态
	engine := common.Mysql(platform)
	defer engine.Close()
	uMap := map[string]interface{}{
		"reviewer_remark": postedData["reviewer_remark"],
		"state":           3,
		"reviewer":        base_controller.GetLoginAdmin(c).Name,
		"updated":         tools.NowMicro(),
	}
	if dividend.MoneyType == 1 { //中心钱包
		accountInfo := &models.Account{}
		_, _ = models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": dividend.UserId}))
		uMap["before_money"] = accountInfo.Available
		uMap["money"] = dividend.Money
		uMap["after_money"] = accountInfo.Available
	} else { //场馆钱包
		//查询余额
		req.SetTimeout(50 * time.Second)
		req.Debug = true
		header := req.Header{
			"Accept": "application/json",
		}
		param := req.Param{
			"game_code": dividend.Venue,
		}
		baseBalanceUrl := consts.InternalGameServUrl
		BalanceUrl := baseBalanceUrl + "/game/v1/internal/balance?user_id=" + strconv.Itoa(int(dividend.UserId))
		r, err := req.Post(BalanceUrl, header, param)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, err.Error())
			return
		}
		balanceResp := base_controller.BalanceResp{}
		_ = r.ToJSON(&balanceResp)
		if balanceResp.Errcode != 0 { //查询余额失败
			response.Err(c, "拒绝失败")
			return
		} else {
			uMap["before_money"] = balanceResp.Data.Balance
			uMap["money"] = dividend.Money
			beforeMoney := balanceResp.Data.Balance
			uMap["after_money"] = beforeMoney + dividend.Money
		}
	}
	if _, err := engine.Table("user_dividends").Where("id=?", id).Update(uMap); err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "操作失败")
		return
	}
	response.Ok(c)
}
