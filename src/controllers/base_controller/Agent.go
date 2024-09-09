package base_controller

import (
	"fmt"
	"math"
	common "sports-common"
	"sports-common/log"
	"sports-common/tools"
	models "sports-models"
	"strconv"
)

func GetMoney(platform string, userId string, start string, end string, venueList map[string]string, depositCost float64, money float64, rate float64, topId string) float64 {
	//已经获取了当前用户的， 现在要获取所有下级的, 普通模式的不抽，
	db := common.Mysql(platform)
	defer db.Close()
	lowSql := "select id,agent_commission from users where  top_id=?  and is_agent=1"
	res, _ := db.QueryString(lowSql, userId)

	//if topId != "1000" && topId != "1001" && topId != "1002" {
	//	//算出上级的抽成。
	//	upSql := "select agent_commission from users where id=" + topId
	//
	//	upRes, _ := db.QueryString(upSql)
	//	if len(upRes) > 0 {
	//		var pRes []models.AgentCommissionPlan
	//		db.Table("agent_commission_plans").Where("agent_commission=? and user_id=?", upRes[0]["agent_commission"], topId).OrderBy("created desc ,level_id asc").Find(&pRes)
	//		if len(pRes) > 0 && pRes[0].Type == 2 {
	//			money = money - money*(pRes[0].Rate-rate)
	//		}
	//	}
	//
	//}

	sumSql := "select sum(arrive_money) as deposit_money,(select sum(money)  from user_account_sets where status=2 and top_id = ? and created >= UNIX_TIMESTAMP(?) and  created< UNIX_TIMESTAMP(?) and status=2) as account_money," +
		" (select sum(money) as money from user_rebate_records where top_id=? and  created>= UNIX_TIMESTAMP(?) and  created< UNIX_TIMESTAMP(?)) as rebate_money," +
		" (select sum(adjust_money) from user_resets where top_id=? and  updated>= UNIX_TIMESTAMP(?) and  updated<UNIX_TIMESTAMP(?)  and adjust_method=1 and status=2 ) as adjust_add," +
		"(select sum(adjust_money)  from user_resets where top_id=? and  updated>= UNIX_TIMESTAMP(?) and  updated<UNIX_TIMESTAMP(?)  and adjust_method=2 and status=2) as adjust_dec," +
		" (select sum(money) from user_dividends where top_id=? and  created>= UNIX_TIMESTAMP(?) and  created< UNIX_TIMESTAMP(?)  and state =2) as dividend_money ," +
		"(select sum(award)  from activity_applies where top_id=? and  created>= UNIX_TIMESTAMP(?) and  created< UNIX_TIMESTAMP(?) and status=2) as activity_award," +
		"(select sum(money)  from user_activities where top_id=? and  created>= UNIX_TIMESTAMP(?) and  created< UNIX_TIMESTAMP(?) and status=2) as activity_money" +
		" from user_deposits where top_id=? and confirm_at >= UNIX_TIMESTAMP(?) and  confirm_at< UNIX_TIMESTAMP(?) and status=2 "

	for _, v := range res {

		res, err := db.QueryString(sumSql, v["id"], start, end, v["id"], start, end, v["id"], start, end, v["id"], start, end, v["id"], start, end, v["id"], start, end, v["id"], start, end, v["id"], start, end)
		if err != nil {
			log.Err(err.Error())
		}
		//存款金额

		//代客充值
		cusMoney, _ := strconv.ParseFloat(res[0]["account_money"], 64)
		dMoney, _ := strconv.ParseFloat(res[0]["deposit_money"], 64)
		rr1, _ := strconv.ParseFloat(res[0]["adjust_add"], 64)
		rr2, _ := strconv.ParseFloat(res[0]["adjust_dec"], 64)
		activityMoney, _ := strconv.ParseFloat(res[0]["activity_award"], 64)
		artificialMoney, _ := strconv.ParseFloat(res[0]["activity_money"], 64)
		reFloat, _ := strconv.ParseFloat(res[0]["rebate_money"], 64)
		diFloat, _ := strconv.ParseFloat(res[0]["dividend_money"], 64)

		var winFloat float64
		reportSql := "select sum(direct_net_money) as ,game_code from user_daily_reports where user_id=" + v["id"] + " and  DATE_FORMAT(day,'%s')>= DATE_FORMAT('" + start + "','%s') and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + end + "','%s') and game_code!='0' group by game_code"
		reportSqls := fmt.Sprintf(reportSql, "%y-%m-%d", "%y-%m-%d", "%y-%m-%d", "%y-%m-%d")
		reportRes, _ := db.QueryString(reportSqls)

		//场馆费  查找所有场馆，再去算所有场馆 的下注量
		var gameCost float64
		for _, v := range reportRes {
			gFloat, _ := strconv.ParseFloat(v["money"], 64)
			if gFloat < 0 { //当总输赢为负的时候算场馆费
				temp := math.Abs(gFloat)
				rateStr := venueList[v["game_code"]]
				rate, _ := strconv.ParseFloat(rateStr, 64)
				gameCost = gameCost + temp*rate
			}
			winFloat = winFloat + gFloat
		}

		//存款手续费
		DepositsCost := (cusMoney + dMoney) * depositCost

		//净输赢    总输赢-红 利-存款手续费，场馆非，-反水。
		profit := -1*winFloat - reFloat - diFloat - gameCost - DepositsCost - artificialMoney - activityMoney + rr2 - rr1
		if profit < 0 {
			continue
		}

		var pRes []models.AgentCommissionPlan
		db.Table("agent_commission_plans").Where("agent_commission=? and user_id=?", v["agent_commission"], v["id"]).OrderBy("created desc ,level_id asc").Find(&pRes)

		if len(pRes) == 0 {
			continue
		}
		var rates float64
		if pRes[0].Type != 2 { //当不是占成模式。
			continue
		}
		rates = pRes[0].Rate
		currMoney := tools.ToFixed(profit*rates, 2)

		//累加。
		money += currMoney * (rate - rates)

	}

	return money
}
