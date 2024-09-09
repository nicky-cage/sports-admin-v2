package controllers

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// UserWithdrawsCoinStruct 用户提现
type UserWithdrawsCoinStruct struct {
	models.UserWithdraw `xorm:"extends"`
	Label               string `jsong:"label"`
	Vip                 int32
}

// UserWithdrawCoinSumStruct 用户提现统计
type UserWithdrawCoinSumStruct struct {
	Money float64
}

// UserWithdrawCoins 提款管理
var UserWithdrawCoins = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); !exists {
			currentTime := time.Now().Unix()
			startAt = currentTime - currentTime%86400
			endAt = startAt + 86400
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0])
			endAt = tools.GetTimeStampByString(areas[1])
		}
		cond = cond.And(builder.Gte{"user_withdraws.created": tools.SecondToMicro(startAt)}).
			And(builder.Lte{"user_withdraws.created": tools.SecondToMicro(endAt)})
		cond = cond.And(builder.Eq{"process_step": 3}).
			And(builder.Eq{"user_withdraws.type": 1}).
			And(builder.Eq{"user_withdraws.status": 1})
		if min, ok := c.GetQuery("money_min"); ok {
			cond = cond.And(builder.Gte{"user_withdraws.money": min})
		}
		if max, ok := c.GetQuery("money_max"); ok {
			cond = cond.And(builder.Lte{"user_withdraws.money": max})
		}
		username := c.DefaultQuery("username", "")
		billNo := c.DefaultQuery("bill_no", "")
		riskAdmin := c.DefaultQuery("risk_admin", "")
		status := c.DefaultQuery("status", "")
		vipLevel := c.DefaultQuery("vip", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.username": username})
		}
		if len(billNo) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.bill_no": billNo})
		}
		if len(riskAdmin) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.risk_admin": riskAdmin})
		}
		if len(status) > 0 {
			cond = cond.And(builder.Eq{"user_withdraws.status": status})
		}
		if vipLevel != "" {
			cond = cond.And(builder.Eq{"users.vip": vipLevel})
		}
		cond = cond.And(builder.Gt{"user_withdraws.wallet_id": 0}) // 钱包id > 0
		limit, offset := request.GetOffsets(c)
		userWithdraws := make([]UserWithdrawsStruct, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()

		//出款方式
		depositCards := make([]models.DepositCard, 0)
		if err := engine.Table("deposit_cards").Find(&depositCards); err != nil {
			log.Logger.Error(err.Error())
		}
		depositCardsNew := make([]string, 0)
		for _, v := range depositCards {
			depositCardsNew = append(depositCardsNew, v.Byname)
		}
		depositCardsNew = append(depositCardsNew, "shipu_daifu")
		total, err := engine.Table("user_withdraws").Join("LEFT OUTER", "users", "user_withdraws.user_id = users.id").Where(cond).
			OrderBy("user_withdraws.id DESC").
			Limit(limit, offset).FindAndCount(&userWithdraws)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ss := new(UserWithdrawSumStruct)
		totalApply, _ := engine.Table("user_withdraws").Where(cond).Sum(ss, "money")
		totalCost, _ := engine.Table("user_withdraws").Where(cond).Sum(ss, "withdraw_cost")
		subtotalApply := 0.0
		subtotalCost := 0.0
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		for i := 0; i < len(userWithdraws); i++ {
			subtotalApply += userWithdraws[i].Money
			subtotalCost += userWithdraws[i].WithdrawCost
			if userWithdraws[i].BankName == "其他银行" {
				sqlCard := "SELECT bank_name FROM user_cards WHERE card_number= '" + userWithdraws[i].BankCard + "'"
				dataCard, err := dbSession.QueryString(sqlCard)
				if err == nil && len(dataCard) >= 1 && dataCard[0]["bank_name"] != "" {
					userWithdraws[i].BankName = dataCard[0]["bank_name"]
				}
			}
		}
		viewData := pongo2.Context{
			"rows":                  userWithdraws,
			"payment_method":        depositCardsNew,
			"total":                 total,
			"subtotal_apply":        subtotalApply,
			"subtotal_cost":         subtotalCost,
			"subtotal_actual":       subtotalApply - subtotalCost,
			"total_apply":           totalApply,
			"total_cost":            totalCost,
			"total_actual":          totalApply - totalCost,
			"vipLevels":             caches.UserLevels.All(platform),
			"rate":                  6.43, //  tools.GetExchangeRate(),
			"min_withdraw":          models.Parameters.GetValueByFloat(platform, "min_withdraw", 100.0, "单次最小提款金额"),
			"max_withdraw":          models.Parameters.GetValueByFloat(platform, "max_withdraw", 49999.0, "单次最大提款金额"),
			"max_withdraw_day":      models.Parameters.GetValueByFloat(platform, "max_withdraw_day", 999999.0, "单日最大提款金额"),
			"withdraw_usdt_min":     models.Parameters.GetValueByFloat(platform, "withdraw_usdt_min", 20.0, "单次最小USDT提款金额"),
			"withdraw_usdt_max":     models.Parameters.GetValueByFloat(platform, "withdraw_usdt_max", 5000.0, "单次最大USDT提款金额"),
			"withdraw_usdt_max_day": models.Parameters.GetValueByFloat(platform, "withdraw_usdt_max_day", 10000.0, "单日最大USDT提款金额"),
			"withdraw_auto_rate":    models.Parameters.GetValueByInt(platform, "withdraw_auto_rate", 2, "提款自动获取汇率"),
			"withdraw_rate_float":   models.Parameters.GetValueByFloat(platform, "withdraw_rate_float", 0.02, "提款浮动汇率"),
			"withdraw_fixed_rate":   models.Parameters.GetValueByFloat(platform, "withdraw_fixed_rate", 1.0, "提款固定汇率"),
		}

		base_controller.SetLoginAdmin(c)
		viewFile := "user_withdraw_coins/user_withdraws.html"
		if request.IsAjax(c) {
			viewFile = "user_withdraw_coins/_user_withdraws.html"
		}
		response.Render(c, viewFile, viewData)
	},
}
