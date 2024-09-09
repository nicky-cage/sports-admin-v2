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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"xorm.io/builder"
)

func (ths *ActivityAudits) BatchRefuse(c *gin.Context) {
	postedData := request.GetPostedData(c)
	idStr, exists := postedData["id"].(string)
	if !exists || exists && idStr == "0" {
		response.Err(c, "id为空")
		return
	}
	idSilce := strings.Split(idStr, ",")
	platform := request.GetPlatform(c)
	redis := common.Redis(platform)
	defer common.RedisRestore(platform, redis)
	engine := common.Mysql(platform)
	defer engine.Close()
	administrator := base_controller.GetLoginAdmin(c)
	for _, v := range idSilce {
		dividend := &models.UserDividend{}
		vid, _ := strconv.Atoi(v)
		b, err := models.UserDividends.FindById(platform, vid, dividend)
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		//记录不存在或已经被处理过
		if !b || dividend.State != 1 {
			continue
		}
		rKey := dividend.BillNo + "_" + idStr + "_" + dividend.Username
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		//这条记录正在被处理的过滤
		if num > 1 {
			continue
		}
		uMap := map[string]interface{}{
			"reviewer_remark": postedData["reviewer_remark"],
			"state":           3,
			"reviewer":        administrator.Name,
			"updated":         tools.NowMicro(),
		}
		if dividend.MoneyType == 1 { //中心钱包
			accountInfo := &models.Account{}
			if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": dividend.UserId})); !exists || err != nil {
				if err != nil {
					log.Logger.Error(err.Error())
				}
				continue
			}
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
				continue
			}
			balanceResp := base_controller.BalanceResp{}
			_ = r.ToJSON(&balanceResp)
			if balanceResp.Errcode != 0 { //查询余额失败
				continue
			} else {
				uMap["before_money"] = balanceResp.Data.Balance
				uMap["money"] = dividend.Money
				beforeMoney := balanceResp.Data.Balance
				uMap["after_money"] = beforeMoney + dividend.Money
			}
		}
		if _, err := engine.Table("user_dividends").Where("id=?", v).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			continue
		}
		redis.Del(rKey)
	}
	response.Ok(c)
}
