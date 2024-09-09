package controllers

import (
	"fmt"
	"math"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

type totalCount struct {
	Deposits  float64
	Withdraws float64
	ActiveNum int
	Rebate    float64
	Dividends float64
	GameCost  float64
	Final     float64
	Profit    float64
	Money     float64
	Num       int
}

var CommissionTeam = struct {
	List   func(c *gin.Context)
	Detail func(c *gin.Context)
}{
	List: func(c *gin.Context) {
		var month string
		var part string
		var cond string

		var start string
		var ends string
		var startInt int64
		var endInt int64
		month = c.Query("created")
		if len(month) > 0 {
			temp := strings.Split(month, " - ")
			start = temp[0]
			ends = temp[1]
			startInt = tools.GetMicroTimeStampByString(temp[0] + " 00:00:00")
			endInt = tools.GetMicroTimeStampByString(temp[1] + " 23:59:59")
		} else {
			start = time.Now().Format("2006-01-02")
			ends = start
			startInt = tools.GetMicroTimeStampByString(start + " 00:00:00")
			endInt = tools.GetMicroTimeStampByString(start + " 23:59:59")
		}

		platform := request.GetPlatform(c)
		username := c.Query("username")
		if username != "" {
			cond = cond + " and username='" + username + "' "
		}

		size, currPage := request.GetOffsets(c)

		if currPage != 0 {
			temp := "limit %d ,%d"
			part = fmt.Sprintf(temp, currPage, size)
		} else {
			part = "limit 15"
		}

		db := common.Mysql(platform)
		defer db.Close()
		sql := "select id as user_id from users where is_agent=1 and id not in(1000,1001,1002) " + cond + part

		res, _ := db.QueryString(sql)
		totalSql := "select count(*) as num  from users where is_agent=1 " + cond
		totalRes, _ := db.QueryString(totalSql)
		total := totalRes[0]["num"]

		//查找场馆code
		vSql := "select code,platform_rate from game_venues where pid=0 and is_online=1 and  code!='CENTERWALLET'"
		vRes, _ := db.QueryString(vSql)

		gameRateList := make(map[string]string, len(vRes))
		for _, v := range vRes {
			gameRateList[v["code"]] = v["platform_rate"]
		}
		var totalAccount totalCount
		sumSql := " select sum(arrive_money) as deposit_money," +
			" (select sum(money)  from user_rebate_records where top_id=? and  created>= ? and  created< ?) as rebate_money," +
			" (select sum(adjust_money) from user_resets where top_id=? and  updated>= ? and  updated<?  and adjust_method=1 and status=2 ) as adjust_add," +
			"(select sum(adjust_money)  from user_resets where top_id=? and  updated>= ? and  updated<?  and adjust_method=2 and status=2) as adjust_dec," +
			" (select sum(money) from user_dividends where top_id=? and  created>= ? and  created< ?  and state =2) as dividend_money ," +
			"(select sum(award)  from activity_applies where top_id=? and  created>= ? and  created< ? and status=2) as activity_award," +
			"(select sum(money)  from user_activities where top_id=? and  created>= ? and  created< ? and status=2) as activity_money," +
			"(select sum(money) from user_withdraws where top_id=? and updated>= ? and  updated<? and status=2 ) as withdraw_money " +
			" from user_deposits where top_id=? and confirm_at >= ? and  confirm_at< ? and status=2 "

		for _, v := range res { //获取所有代理的下属会员的id ……

			usql := "select id,username,agent_type,top_name,agent_commission,top_id,(select count(id) from users where top_id=" + v["user_id"] + ") as lownum from users where id=" + v["user_id"]
			//获取会员信息
			uRes, uerr := db.QueryString(usql)
			if uerr != nil {
				log.Err(uerr.Error())
				return
			}
			v["start_time"] = start
			v["end_time"] = ends
			//信息插入代理佣金的res中
			v["id"] = uRes[0]["id"]
			v["username"] = uRes[0]["username"]
			v["agent_type"] = uRes[0]["agent_type"]
			v["top_name"] = uRes[0]["top_name"]
			v["lower_num"] = uRes[0]["lownum"]
			v["month"] = month
			v["top_id"] = uRes[0]["top_id"]
			tempNum, _ := strconv.Atoi(uRes[0]["lownum"])
			totalAccount.Num += tempNum
			res, err := db.QueryString(sumSql, v["user_id"], startInt, endInt, v["user_id"], startInt, endInt, v["user_id"], startInt, endInt, v["user_id"], startInt, endInt, v["user_id"], startInt, endInt, v["user_id"], startInt, endInt, v["user_id"], startInt, endInt, v["user_id"], tools.MicroToSecond(startInt), tools.MicroToSecond(endInt))
			if err != nil {
				log.Err(err.Error())
			}
			dMoney, _ := strconv.ParseFloat(res[0]["deposit_money"], 64)
			v["deposits"] = strconv.FormatFloat(tools.ToFixed(dMoney, 2), 'f', -1, 64)
			totalAccount.Deposits += dMoney
			//提款金额

			if res[0]["withdraw_money"] == "" {
				v["withdraws"] = "0"
			} else {
				v["withdraws"] = res[0]["withdraw_money"]
			}
			tempWithdraws, _ := strconv.ParseFloat(res[0]["withdraw_money"], 64)
			totalAccount.Withdraws += tempWithdraws
			//投注活跃 还要有时间的限制。
			ursql := "select count(distinct user_id) as num from user_daily_reports where top_id=" + v["user_id"] + " and  day>='" + start + "' and  day<='" + ends + "'"
			urRes, _ := db.QueryString(ursql)

			v["active_num"] = urRes[0]["num"]
			tempActiveNum, _ := strconv.Atoi(urRes[0]["num"])
			totalAccount.ActiveNum += tempActiveNum
			v["rebate"] = res[0]["rebate_money"]
			tempRebate, _ := strconv.ParseFloat(res[0]["rebate_money"], 64)
			totalAccount.Rebate += tempRebate
			rr1, _ := strconv.ParseFloat(res[0]["adjust_add"], 64)
			rr2, _ := strconv.ParseFloat(res[0]["adjust_dec"], 64)
			v["reset"] = strconv.FormatFloat(tools.ToFixed(rr2-rr1, 2), 'f', -1, 64)

			//红利

			diMoney, _ := strconv.ParseFloat(res[0]["dividend_money"], 64)
			activityMoney, _ := strconv.ParseFloat(res[0]["activity_award"], 64)
			artificialMoney, _ := strconv.ParseFloat(res[0]["activity_money"], 64)

			v["dividends"] = strconv.FormatFloat(tools.ToFixed(activityMoney+diMoney+artificialMoney, 2), 'f', -1, 64)

			totalAccount.Dividends = totalAccount.Dividends + (activityMoney + diMoney + artificialMoney)
			//总输赢
			var winFloat float64
			reportSql := "select sum(net_money) as money ,game_code from user_daily_reports where top_id=" + v["user_id"] + " and  day>='" + start + "' and  day<='" + ends + "' and game_code!='0' group by game_code"

			reportRes, err := db.QueryString(reportSql)
			if err != nil {
				log.Err(err.Error())
			}
			//场馆费  查找所有场馆，再去算所有场馆 的下注量
			var gameCost float64
			for _, v := range reportRes {
				gFloat, _ := strconv.ParseFloat(v["money"], 64)
				if gFloat < 0 { //当总输赢为负的时候算场馆费
					temp := math.Abs(gFloat)
					rateStr := gameRateList[v["game_code"]]
					rate, _ := strconv.ParseFloat(rateStr, 64)
					gameCost = gameCost + temp*rate
				}
				winFloat = winFloat + gFloat
			}
			v["game_cost"] = strconv.FormatFloat(tools.ToFixed(gameCost, 2), 'f', -1, 64)

			totalAccount.GameCost += gameCost
			reFloat, _ := strconv.ParseFloat(res[0]["rebate_money"], 64)

			profit := -1*winFloat - reFloat - diMoney - gameCost - artificialMoney - activityMoney
			if profit == -0 {
				profit = 0
			}
			v["final"] = strconv.FormatFloat(tools.ToFixed(-winFloat, 2), 'f', -1, 64)
			v["only_profit"] = strconv.FormatFloat(tools.ToFixed(profit, 2), 'f', -1, 64)
			totalAccount.Final += -winFloat
			totalAccount.Profit += profit
			var pRes []models.AgentCommissionPlan
			db.Table("agent_commission_plans").Where("agent_commission=? and user_id=?", uRes[0]["agent_commission"], v["user_id"]).OrderBy("created desc ,level_id asc").Find(&pRes)

			if len(pRes) == 0 {
				v["rate"] = "0%"
				v["money"] = "0"
				continue
			}
			var rate float64

			rate = pRes[0].Rate
			rateStr := strconv.FormatFloat(pRes[0].Rate*100, 'f', 2, 64)
			v["rate"] = rateStr + "%"

			money := tools.ToFixed(profit*rate, 2)

			if money <= 0 {
				money = 0
				v["money"] = "0"
			} else {
				tempMoney := base_controller.GetMoney(platform, v["user_id"], start, ends, gameRateList, 0, money, rate, v["top_id"])
				v["money"] = strconv.FormatFloat(tempMoney, 'f', 2, 64)
			}

			totalAccount.Money += money

		}

		SetLoginAdmin(c)
		if request.IsAjax(c) {
			response.Render(c, "team/_commission.html", pongo2.Context{"res": res, "total": total, "totalAccount": totalAccount}) //, "total": total
		} else {
			response.Render(c, "team/commission.html", pongo2.Context{"res": res, "total": total, "totalAccount": totalAccount})
		}

	},
	Detail: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		id := c.Query("id")
		start := c.Query("start_time")
		end := c.Query("end_time")
		data := make(map[string]string, 13)
		data["username"] = c.Query("username")
		data["top_name"] = c.Query("top_name")
		data["lower_num"] = c.Query("lower_num")
		data["active_num"] = c.Query("active_num")
		data["deposit"] = c.Query("deposit")
		data["withdraws"] = c.Query("withdraws")
		data["final"] = c.Query("final")
		data["game_cost"] = c.Query("game_cost")
		data["dividends"] = c.Query("dividends")
		data["net_profit"] = c.Query("net_profit")
		data["rebate"] = c.Query("rebate")
		data["money"] = c.Query("money")
		var pRes []models.AgentCommissionPlan
		dbSession.Table("agent_commission_plans").Where(" user_id=?", id).Find(&pRes)
		data["rate"] = strconv.FormatFloat(pRes[0].Rate, 'f', 2, 64)
		type Game struct {
			GameName    string  `json:"game_name"`
			BetUsers    int64   `json:"bet_users"`
			ActiveUsers int64   `json:"active_users"`
			BetMoney    float64 `json:"bet_money"`
			ValidMoney  float64 `json:"vaild_money"`
			NetMoney    float64 `json:"net_money"`
		}
		type AllCount struct {
			BetNum     int     `json:"bet_number" xorm:"bet_number"`
			ActivesNum int     `json:"actives_num" xorm:"actives_num"`
			BetMoney   float64 `json:"bet_money" xorm:"bet_money"`
			Win        float64 `json:"bet_win" xorm:"bet_win"`
			ValidMoney float64 `json:"valid_money" xorm:"valid_money"`
		}

		gameInfo := make([]Game, 0)
		var allCount AllCount
		defer dbSession.Close()
		////  暂时获取这3个
		gameVenues := make([]models.GameVenue, 0)
		_ = dbSession.Where("pid!=0 and is_online=1").Find(&gameVenues)

		for _, gv := range gameVenues {
			var dayReport []models.UserDailyReport
			dbSession.Where("day>=? and day<=?  and game_code=? and game_type=? and top_id=?", start, end, gv.Code, gv.VenueType, id).Find(&dayReport)

			var test Game
			for _, v := range dayReport {
				test.BetMoney = tools.ToFixed(test.BetMoney+v.BetMoney, 2)
				test.ValidMoney = tools.ToFixed(test.ValidMoney+v.ValidMoney, 2)
				test.NetMoney = tools.ToFixed(test.NetMoney+v.NetMoney, 2)
				allCount.BetMoney += v.BetMoney
				allCount.ValidMoney += v.ValidMoney
				allCount.Win = v.NetMoney
			}

			var userActive []models.UserDailyReport
			_ = dbSession.Where("day>=? and day<=? and top_id=?  and is_agent=0  and game_code=? and game_type=?", start, end, id, gv.Code, gv.VenueType).GroupBy("user_id").Find(&userActive)

			num := len(userActive)
			test.GameName = gv.Name
			test.ActiveUsers = int64(num)
			test.BetUsers = int64(num)
			allCount.BetNum += num
			allCount.ActivesNum += num
			gameInfo = append(gameInfo, test)
		}
		viewFile := "team/commission_detail.html"

		SetLoginAdmin(c)
		response.Render(c, viewFile, ViewData{
			"data":     gameInfo,
			"r":        data,
			"allCount": allCount,
		})
	},
}
