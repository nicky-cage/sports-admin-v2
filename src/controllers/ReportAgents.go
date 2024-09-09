package controllers

import (
	"fmt"
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
	"xorm.io/builder"
)

type ReportAgentsStruct struct {
	Date               string  //统计日期
	Id                 uint64  //代理编号
	Username           string  //代理账号
	AgentType          int32   //代理类型
	RegisteredNumber   int64   //注册人数
	BetNumber          int64   //投注人数
	FirstDepositNumber int64   //首存人数
	DepositsNumber     int64   //存款人数
	WithdrawalsNumber  int64   //提款人数
	DepositAmount      float64 //存款金额
	FirstDepositAmount float64 //首存金额
	WithdrawalAmount   float64 //提款金额
	ValidBet           float64 //有效投注
	TotalNetMoney      float64 //总输赢
	DividendMoney      float64 //红利
	RebateMoney        float64 //返水
	WinLoseAdjustment  float64 //输赢调整
}

type ReportAgentsIntSumStruct struct {
	DirectRegisterNumber   int64
	DirectDepositNumber    int64
	DirectFirstDeposit     int64
	DirectWithdrawalNumber int64
}

type ReportAgentsFloatSumStruct struct {
	DirectDepositMoney      float64
	DirectFirstDepositMoney float64
	DirectWithdrawalMoney   float64
	DirectValidMoney        float64
	DirectNetMoney          float64
	DirectDividend          float64
	DirectRebate            float64
	DirectWinloseAdjustment float64
}

func GetReportAgentsPage(page, pageSize int, data []models.UserDailyReport) []models.UserDailyReport {
	start := (page - 1) * pageSize
	stop := start + pageSize
	if start > len(data) {
		return nil
	}
	if stop > len(data) {
		stop = len(data)
	}
	return data[start:stop]
}

// 会员报表
//投注人数文案调整为”活跃人数“ ”返水“文案调整为”活动奖金“
var ReportAgents = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) { //默认首页
		var startAt int64
		var endAt int64
		var startDayAt string
		var endDayAt string
		cond := builder.NewCond()
		limit, size := request.GetOffsets(c)

		if value, exists := c.GetQuery("created"); !exists {
			currentDayTime := time.Now().Format("2006-01-02")
			startAt = tools.GetTimeStampByString(currentDayTime + " 00:00:00")
			endAt = startAt + 24*60*60
			startDayAt = currentDayTime
			endDayAt = time.Unix(endAt, 0).Format("2006-01-02")
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetTimeStampByString(areas[0])
			endAt = tools.GetTimeStampByString(areas[1])
			startDayAt = areas[0]
			endDayAt = areas[1]
		}

		startDayMicro := fmt.Sprintf("%d", tools.SecondToMicro(startAt))
		endDayMicro := fmt.Sprintf("%d", tools.SecondToMicro(endAt+86399))

		username := c.Query("username")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"username": username})
		}
		cond = cond.And(builder.Eq{"is_agent": 1})
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()

		var users []models.User
		engine.Table("users").Where(cond).Limit(limit, size).Cols("username,id").Find(&users)

		var usersCount []models.User
		engine.Table("users").Where(cond).Cols("id").Find(&usersCount)
		total := len(usersCount)

		mapSql := "select count(*) as registerNum, (SELECT count(DISTINCT user_id) FROM `user_daily_reports` WHERE  top_name=? AND day=? and bet_money>0 ) as sumBetNumber ," +
			"(select count(DISTINCT user_id)  from user_withdraws where top_id=? and  updated>=" + startDayMicro + " and updated<" + endDayMicro + " and status=2 ) as WithdrawalsNumber ," +
			"(select sum(arrive_money)  from user_deposits where top_id=? and confirm_at >= " + startDayMicro + " and  confirm_at< " + endDayMicro + " and status=2) as dMoney ," +
			"(select sum(money)  from user_withdraws where top_id=? and  updated>=" + startDayMicro + " and updated<" + endDayMicro + "  and status=2) as WithdrawMoney," +
			"(select sum(money)  from user_dividends where top_id=? and  created>= " + startDayMicro + " and  created< " + endDayMicro + "  and state =2) as DividendMoney," +
			"(select sum(award) from activity_applies where top_id=? and  created>= " + startDayMicro + " and  created< " + endDayMicro + " and status=2 ) as activityMoney," +
			"(select sum(money) from user_activities where top_id=? and  created>= " + startDayMicro + " and  created< " + endDayMicro + " and status=2 ) as artificialMoney," +
			"(select sum(money)  from user_rebate_records where top_id=? and  created>= " + startDayMicro + " and  created< " + endDayMicro + ") as RebateMoney," +
			"(select sum(adjust_money)  from user_resets where top_id=? and  updated>= " + startDayMicro + " and  updated<" + endDayMicro + "  and adjust_method=1 and status=2) as rr1 ," +
			"(select sum(adjust_money)  from user_resets where top_id=? and  updated>= " + startDayMicro + " and  updated<" + endDayMicro + "  and adjust_method=2 and status=2) as rr2 " +
			"from users where top_id=? and created>=? and created<?"
		reportAgentsList := make([]ReportAgentsStruct, 0)
		for _, v := range users {
			temp := ReportAgentsStruct{}
			temp.Date = time.Unix(startAt, 0).Format("2006-01-02") + " ~ " + time.Unix(endAt, 0).Format("2006-01-02")

			temp.Username = v.Username
			userId := strconv.Itoa(int(v.Id))

			mapRes, err := engine.QueryString(mapSql, v.Username, startDayAt, userId, userId, userId, userId, userId, userId, userId, userId, userId, userId, startDayMicro, endDayMicro)
			if err != nil {
				log.Err(err.Error())
			}

			registerNum, _ := strconv.Atoi(mapRes[0]["registerNum"])
			//注册人数
			temp.RegisteredNumber = int64(registerNum)

			//投注人数
			sumBetNumber, _ := strconv.Atoi(mapRes[0]["sumBetNumber"])
			temp.BetNumber = int64(sumBetNumber)

			//存款人数
			depositNumSql := "select user_id from user_deposits where top_id=" + userId + " and confirm_at >= UNIX_TIMESTAMP('" + startDayAt + "') and  confirm_at< UNIX_TIMESTAMP('" + endDayAt + " 23:59:59') and status=2 group by user_id"
			depositNumRes, depositNumErr := engine.QueryString(depositNumSql)
			if depositNumErr != nil {
				log.Err(depositNumErr.Error())
				return
			}
			//还有上分的。
			accountNumSql := "select user_id from user_account_sets where top_id =" + userId + " and created >=" + startDayMicro + " and  created< " + endDayMicro + " and status=2 group by user_id"
			accountNumRes, accountNumErr := engine.QueryString(accountNumSql)
			if accountNumErr != nil {
				log.Err(depositNumErr.Error())
				return
			}
			DepositsNumber := make(map[string]string, 0)
			for _, v := range depositNumRes {
				DepositsNumber[v["user_id"]] = v["user_id"]
			}
			for _, v := range accountNumRes {
				DepositsNumber[v["user_id"]] = v["user_id"]
			}
			DepositsNumbers := len(DepositsNumber)
			temp.DepositsNumber = int64(DepositsNumbers)

			//首存人数--is_first_deposit确定唯一
			firstNumSql := "select sum(money) as money ,count(*) as num from user_deposits where top_id=" + userId + " and confirm_at >= UNIX_TIMESTAMP('" + startDayAt + "') and  confirm_at< UNIX_TIMESTAMP('" + endDayAt + " 23:59:59') and status=2 and is_first_deposit=2"
			firstNumRes, depositNumErr := engine.QueryString(firstNumSql)
			if depositNumErr != nil {
				log.Err(depositNumErr.Error())
				return
			}

			if len(firstNumRes) > 0 {
				FirstDepositNumber, _ := strconv.Atoi(firstNumRes[0]["num"])
				temp.FirstDepositNumber = int64(FirstDepositNumber)
			}

			//提款人数

			WithdrawalsNumber, _ := strconv.Atoi(mapRes[0]["WithdrawalsNumber"])
			temp.WithdrawalsNumber = int64(WithdrawalsNumber)

			//存款金额
			dMoney, _ := strconv.ParseFloat(mapRes[0]["dMoney"], 64)
			temp.DepositAmount = dMoney

			//首存金额--is_first_deposit确定唯一
			FirstDepositMoney, _ := strconv.ParseFloat(firstNumRes[0]["money"], 64)
			temp.FirstDepositAmount = FirstDepositMoney
			//提款金额

			WithdrawMoney, _ := strconv.ParseFloat(mapRes[0]["WithdrawMoney"], 64)
			temp.WithdrawalAmount = WithdrawMoney
			//有效投注
			reportSql := "select sum(valid_money) as valid_money,sum(net_money) as net_money from user_daily_reports where top_id=? and day>=? and day<=? and game_code!='0'"
			reportRes, reportErr := engine.QueryString(reportSql, v.Id, startDayAt, endDayAt)

			if reportErr != nil {
				log.Err(reportErr.Error())
			}
			if len(reportRes) > 0 {
				ValidBet, _ := strconv.ParseFloat(reportRes[0]["valid_money"], 64)
				temp.ValidBet = ValidBet
				//总输赢
				TotalNetMoney, _ := strconv.ParseFloat(reportRes[0]["net_money"], 64)
				temp.TotalNetMoney = -1 * TotalNetMoney
			}

			//红利

			DividendMoney, _ := strconv.ParseFloat(mapRes[0]["DividendMoney"], 64)
			// 活动奖金
			//人工红利
			activityMoney, _ := strconv.ParseFloat(mapRes[0]["activityMoney"], 64)
			artificialMoney, _ := strconv.ParseFloat(mapRes[0]["artificialMoney"], 64)
			temp.DividendMoney = DividendMoney + activityMoney + artificialMoney
			//返水

			RebateMoney, _ := strconv.ParseFloat(mapRes[0]["RebateMoney"], 64)
			temp.RebateMoney = RebateMoney
			//输赢调整

			//增加 的减去 减少的
			rr1, _ := strconv.ParseFloat(mapRes[0]["rr1"], 64)
			rr2, _ := strconv.ParseFloat(mapRes[0]["rr2"], 64)
			temp.WinLoseAdjustment = rr1 - rr2
			//组装数据
			reportAgentsList = append(reportAgentsList, temp)
		}
		pageDepositAmount := 0.00
		pageFirstDepositAmount := 0.00
		pageWithdrawalAmount := 0.00
		pageValidBet := 0.00
		pageTotalNetMoney := 0.00
		pageDividendMoney := 0.00
		pageRebateMoney := 0.00
		pageWinLoseAdjustment := 0.00
		pageRegisteredNumber := int64(0)
		pageBetNumber := int64(0)
		pageFirstDepositNumber := int64(0)
		pageDepositsNumber := int64(0)
		pageWithdrawalsNumber := int64(0)
		for _, v := range reportAgentsList {
			pageDepositAmount += v.DepositAmount
			pageFirstDepositAmount += v.FirstDepositAmount
			pageWithdrawalAmount += v.WithdrawalAmount
			pageValidBet += v.ValidBet
			pageTotalNetMoney += v.TotalNetMoney
			pageDividendMoney += v.DividendMoney
			pageRebateMoney += v.RebateMoney
			pageWinLoseAdjustment += v.WinLoseAdjustment
			pageRegisteredNumber += v.RegisteredNumber
			pageBetNumber += v.BetNumber
			pageFirstDepositNumber += v.FirstDepositNumber

			pageDepositsNumber += v.DepositsNumber

			pageWithdrawalsNumber += v.WithdrawalsNumber
		}
		viewData := pongo2.Context{
			"rows":                   reportAgentsList,
			"total":                  total,
			"pageDepositAmount":      pageDepositAmount,
			"pageFirstDepositAmount": pageFirstDepositAmount,
			"pageWithdrawalAmount":   pageWithdrawalAmount,
			"pageValidBet":           pageValidBet,
			"pageTotalNetMoney":      pageTotalNetMoney,
			"pageDividendMoney":      pageDividendMoney,
			"pageRebateMoney":        pageRebateMoney,
			"pageWinLoseAdjustment":  pageWinLoseAdjustment,
			"pageRegisteredNumber":   pageRegisteredNumber,
			"pageBetNumber":          pageBetNumber,
			"pageFirstDepositNumber": pageFirstDepositNumber,
			"pageDepositsNumber":     pageDepositsNumber,
			"pageWithdrawalsNumber":  pageWithdrawalsNumber,
		}
		viewFile := "report_agents/list.html"
		if request.IsAjax(c) {
			viewFile = "report_agents/_list.html"
		}
		response.Render(c, viewFile, viewData)
	},
}
