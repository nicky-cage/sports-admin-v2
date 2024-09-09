package agent_commissions

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
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissions) Lists(c *gin.Context) {
	var month string
	var part string
	var cond string

	var start string
	var ends string
	var startInt int64
	var endInt int64
	month = c.Query("month")
	if len(month) > 0 {
		loc, _ := time.LoadLocation("Local") //获取时区
		tmp, _ := time.ParseInLocation("2006-01-02", month+"-01", loc)
		endMonth := tmp.AddDate(0, 1, 0)
		start = month + "-01"
		ends = endMonth.Format("2006-01-02")
		startInt = tools.GetTimeStampByString(start + " 00:00:00")
		endInt = tools.GetTimeStampByString(ends + " 00:00:00")
	} else {
		year, months, _ := time.Now().Date()
		thisMonth := time.Date(year, months, 1, 0, 0, 0, 0, time.Local)
		start = thisMonth.Format("2006-01-02")
		month = thisMonth.Format("2006-01")
		endMonth := thisMonth.AddDate(0, 1, 0)
		ends = endMonth.Format("2006-01-02")
		//	lastMonth = thisMonth.AddDate(0, -1, 0).Format("2006-01")
		startInt = tools.GetTimeStampByString(start + " 00:00:00")
		endInt = tools.GetTimeStampByString(ends + " 00:00:00")
	}

	startMicroInt := tools.SecondToMicro(startInt)
	endMicroInt := tools.SecondToMicro(endInt)

	platform := request.GetPlatform(c)
	username := c.Query("username")
	if username != "" {
		cond = cond + " and username='" + username + "'"
	}

	id := c.Query("id")
	if id != "" {
		cond = cond + " and user_id='" + id + "'"
	}
	agentType := c.Query("agent_type")
	if agentType != "" {
		cond = cond + " and type='" + agentType + "'"
	}
	topName := c.Query("top_name")
	if topName != "" {
		cond = cond + " and top_name='" + topName + "'"
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
	//if cond == "" {
	//	cond = " and (top_name='' or top_name='sys_test_agent') "
	//}
	// 说明已经将数据还没 获取上月的下月累计负盈利,,  修改一下。先改为本月。 ,negative_profit_next
	sql := "select user_id,commission_adjust from agent_commission_logs where month= '" + month + "'" + cond + " order by user_id asc " + part

	res, err := db.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		return
	}

	totalSql := "select user_id,commission_adjust from agent_commission_logs where month= '" + month + "'" + cond
	totalRes, _ := db.QueryString(totalSql)
	total := len(totalRes)

	//查找场馆code
	vsql := "select code,platform_rate from game_venues where pid=0 and is_online=1 and  code!='CENTERWALLET'"
	vRes, verr := db.QueryString(vsql)
	if verr != nil {
		log.Err(verr.Error())
	}
	gameRateList := make(map[string]string, len(vRes))
	for _, v := range vRes {
		gameRateList[v["code"]] = v["platform_rate"]
	}

	sumSql := " select sum(arrive_money) as deposit_money," +
		" (select sum(money)  from user_rebate_records where top_id=? and  created>= ? and  created< ?) as rebate_money," +
		" (select sum(adjust_money) from user_resets where top_id=? and  updated>= ? and  updated<?  and adjust_method=1 and status=2 ) as adjust_add," +
		"(select sum(adjust_money)  from user_resets where top_id=? and  updated>= ? and  updated<?  and adjust_method=2 and status=2) as adjust_dec," +
		" (select sum(money) from user_dividends where top_id=? and  created>= ? and  created< ?  and state =2) as dividend_money ," +
		"(select sum(award)  from activity_applies where top_id=? and  created>= ? and  created< ? and status=2) as activity_award," +
		"(select sum(money)  from user_activities where top_id=? and  created>= ? and  created< ? and status=2) as activity_money," +
		"(select sum(money) from user_withdraws where top_id=? and updated>= ? and  updated<? and status=2 ) as withdraw_money " +
		" from user_deposits where top_id=? and confirm_at >= ? and  confirm_at< ? and status=2 "

	for _, v := range res {

		usql := "select id,username,agent_type,top_name,agent_commission,top_id,(select count(id) from users where top_id=" + v["user_id"] + ") as lownum from users where id=" + v["user_id"]
		//获取所有代理的下属会员的id ……

		//获取会员信息
		uRes, uerr := db.QueryString(usql)
		if uerr != nil {
			log.Err(uerr.Error())
			return
		}
		//信息插入代理佣金的res中
		v["id"] = uRes[0]["id"]
		v["username"] = uRes[0]["username"]
		v["agent_type"] = uRes[0]["agent_type"]
		v["top_name"] = uRes[0]["top_name"]
		v["lower_num"] = uRes[0]["lownum"]
		v["month"] = month
		v["top_id"] = uRes[0]["top_id"]

		res, err := db.QueryString(sumSql, v["user_id"], startMicroInt, endMicroInt, v["user_id"], startMicroInt, endMicroInt, v["user_id"], startMicroInt, endMicroInt, v["user_id"], startMicroInt, endMicroInt, v["user_id"], startMicroInt, endMicroInt, v["user_id"], startMicroInt, endMicroInt, v["user_id"], startMicroInt, endMicroInt, v["user_id"], startInt, endInt)
		if err != nil {
			log.Err(err.Error())
		}

		//cusMoney, _ := strconv.ParseFloat(res[0]["account_money"], 64)
		dMoney, _ := strconv.ParseFloat(res[0]["deposit_money"], 64)
		v["deposits"] = strconv.FormatFloat(tools.ToFixed(dMoney, 2), 'f', -1, 64)

		//提款金额

		if res[0]["withdraw_money"] == "" {
			v["withdraws"] = "0"
		} else {
			v["withdraws"] = res[0]["withdraw_money"]
		}

		//activeArr := make([]map[string]string, 1)
		//activeArr[0] = make(map[string]string)
		//活跃会员

		//投注活跃 还要有时间的限制。
		ursql := "select user_id from user_daily_reports where top_id=" + v["user_id"] + " and  day>='" + start + "' and  day<'" + ends + "' group by user_id"
		urRes, _ := db.QueryString(ursql)

		activeNum := len(urRes)
		v["active_num"] = strconv.Itoa(activeNum)

		v["rebate"] = res[0]["rebate_money"]

		//输赢调整  获取调整增加的，
		//adsql := "select sum(adjust_money) as money from user_resets where top_id=" + v["user_id"] + " and  updated>= UNIX_TIMESTAMP('" + start + "') and  updated<UNIX_TIMESTAMP('" + ends + "')  and adjust_method=1 and status=2 "
		//r1, derr1 := db.QueryString(adsql)
		//if derr1 != nil {
		//	log.Err(derr1.Error())
		//	return
		//}
		//adsql1 := "select sum(adjust_money) as money from user_resets where top_id=" + v["user_id"] + " and  updated>= UNIX_TIMESTAMP('" + start + "') and  updated<UNIX_TIMESTAMP('" + ends + "')  and adjust_method=2 and status=2  "
		//r2, derr2 := db.QueryString(adsql1)
		//if derr2 != nil {
		//	log.Err(derr2.Error())
		//	return
		//}

		//增加 的减去 减少的

		rr1, _ := strconv.ParseFloat(res[0]["adjust_add"], 64)
		rr2, _ := strconv.ParseFloat(res[0]["adjust_dec"], 64)
		v["reset"] = strconv.FormatFloat(tools.ToFixed(rr2-rr1, 2), 'f', -1, 64)

		//红利

		diMoney, _ := strconv.ParseFloat(res[0]["dividend_money"], 64)
		activityMoney, _ := strconv.ParseFloat(res[0]["activity_award"], 64)
		artificialMoney, _ := strconv.ParseFloat(res[0]["activity_money"], 64)

		v["dividends"] = strconv.FormatFloat(tools.ToFixed(activityMoney+diMoney+artificialMoney, 2), 'f', -1, 64)

		//总输赢
		var winFloat float64
		reportSql := "select sum(net_money) as money ,game_code from user_daily_reports where top_id=" + v["user_id"] + " and  day>='" + start + "' and  day<'" + ends + "' and game_code!='0' group by game_code"

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
		//累计 负盈利

		//nPs, _ := strconv.ParseFloat(v["negative_profit_statistics"], 64)
		//if nPs > 0 {
		//	//ji负累计大于0 不
		//	nPs = 0
		//}

		reFloat, _ := strconv.ParseFloat(res[0]["rebate_money"], 64)

		//winFloat, _ := strconv.ParseFloat(res[0]["report_money"], 64)

		//存款手续费
		//depositCostSql := "select value from parameters where name='deposit_cost'"
		//depositCostRes, _ := db.QueryString(depositCostSql)
		//depositCost, _ := strconv.ParseFloat(depositCostRes[0]["value"], 64)
		//DepositsCost := (cusMoney + dMoney) * depositCost

		v["deposits_cost"] = "0" // strconv.FormatFloat(tools.ToFixed(DepositsCost, 2), 'f', 2, 64)
		//净输赢    总输赢-红 利-存款手续费，场馆非，-反水。
		profit := -1*winFloat - reFloat - diMoney - gameCost - artificialMoney - activityMoney

		v["final"] = strconv.FormatFloat(tools.ToFixed(-winFloat, 2), 'f', -1, 64)
		v["only_profit"] = strconv.FormatFloat(tools.ToFixed(profit, 2), 'f', -1, 64)

		//上月结余 已有。
		//冲正后净谁赢 暂无
		//佣金比例
		var pRes []models.AgentCommissionPlan
		perr := db.Table("agent_commission_plans").Where("agent_commission=? and user_id=?", uRes[0]["agent_commission"], v["user_id"]).OrderBy("created desc ,level_id asc").Find(&pRes)
		if perr != nil {
			log.Err(perr.Error())
			return
		}
		if len(pRes) == 0 {
			v["rate"] = "0%"
			v["money"] = "0"
			continue
		}
		var rate float64

		if pRes[0].Type == 2 {
			rate = pRes[0].Rate
			rateStr := strconv.FormatFloat(pRes[0].Rate*100, 'f', 2, 64)
			v["rate"] = rateStr + "%"

		} else {

			pNum := len(pRes)
			//判断属于哪一个层级。 净输赢 *佣金比例
			for k, v := range pRes {
				if activeNum < int(v.ActiveNum) {
					for key, val := range pRes {
						//如果活跃人数小于标准，说明他处于上一级
						if int(profit) < int(val.NegativeProfit) {
							//判断盈利是在那个阶段，
							if k == 0 || key == 0 || profit < 0 {
								//同时满足 才能有佣金比例
								rate = 0
								goto Loop
							}

							if k < key {
								rate = pRes[k-1].Rate
								goto Loop
							} else {
								rate = pRes[key-1].Rate
								goto Loop
							}
						} else {
							if pNum == key+1 {
								if k == 0 {
									rate = 0
									goto Loop
								}
								//当某个盈利太大
								rate = pRes[k-1].Rate
								goto Loop
							}
						}
					}
				} else {
					if pNum == k+1 {
						//当k最大时，只是比较K key
						for key, val := range pRes {
							if int(profit) < int(val.NegativeProfit) {
								if key == 0 {
									rate = 0
									goto Loop
								}
								rate = pRes[key-1].Rate
								goto Loop
							} else {
								//一般情况。最大对最大，
								if pNum == key+1 {
									rate = pRes[k].Rate
									goto Loop
								}
							}
						}
					}
				}

			}
		Loop:

			rateStr := strconv.FormatFloat(rate*100, 'f', 2, 64)
			v["rate"] = rateStr + "%"

		}

		//佣金调整 已有
		//本 月结余 已有
		//佣金
		//上月结余

		//negativeProfitStatistics, _ := strconv.ParseFloat(v["negative_profit_statistics"], 64)
		////var negativeProfit float64
		//if negativeProfitStatistics > 0 {
		//	negativeProfit = 0
		//} else {
		//	negativeProfit = negativeProfitStatistics
		//}
		//	adjust, _ := strconv.Atoi(v["commission_adjust"])
		//当有净输赢为正  才是佣金， 不然是欠的钱，不应该乘佣金比例
		//var profitNext float64

		money := tools.ToFixed(profit*rate, 2)

		//普通模式， 没有抽成，
		//if pRes[0].Type == 2 {
		//	money = GetMoney(v["user_id"], start, ends, gameRateList, depositCost, money, rate, v["top_id"])
		//}

		if money <= 0 {
			v["money"] = "0"
		} else {
			//tempMoney := tools.ToFixed(profit*rate, 2)
			tempMoney := base_controller.GetMoney(platform, v["user_id"], start, ends, gameRateList, 0, money, rate, v["top_id"])
			v["money"] = strconv.FormatFloat(tempMoney, 'f', 2, 64)
		}

		//	}
		//	//要不要加上佣金 调整。
		//	totalNum, err = db.Table("agent_commission_logs").Where("month=?", starts).Limit(limit, offset).FindAndCount(&list)
	}
	base_controller.SetLoginAdmin(c)
	if request.IsAjax(c) {
		response.Render(c, "agents/_commission_offer.html", pongo2.Context{"res": res, "total": total}) //, "total": total
		return
	}
	if topName != "" {
		response.Render(c, "agents/commission_detail.html", pongo2.Context{"res": res, "total": total})
	} else {
		response.Render(c, "agents/commissions.html", pongo2.Context{"res": res, "total": total})
	}
}
