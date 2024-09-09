package controllers

import (
	"math"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

type ReportOperationsStruct struct {
	Date                string  //统计日期
	RegisteredNumber    int64   //注册人数
	FirstDepositNumber  int64   //首存人数
	ConversionRates     float64 //转化率
	FirstDepositAmount  float64 //首存金额
	FirstDepositPer     float64 //人均首存
	DepositsNumber      int64   //存款人数
	WithdrawalsNumber   int64   //提款人数
	DepositAmount       float64 //存款金额
	WithdrawalAmount    float64 //提款金额
	DepositOffer        float64 //存提差
	WithdrawalRate      float64 //提存率
	BetNumber           int64   //投注人数
	ValidBet            float64 //有效投注
	TotalBet            float64 //总投注额
	CompanyWinsLoses    float64 //公司输赢
	SurplusRatio        float64 //盈余比例
	DividendMoney       float64 //红利
	RebateMoney         float64 //返水
	DepositDiscount     float64 //存款优惠
	AgentCommission     float64 //代理佣金
	CompanyIncome       float64 //公司收入
	OnlineDepositNum    int
	OnlineDepositMoney  float64
	OffLineDepositNum   int
	OfflineDepositMoney float64
	HumanDepositNum     int
	HumanDepositMoney   float64
	HumanWithdrawNum    int
	HumanWithdrawMoney  float64
	VirtualDeposit      float64
	VirtualWithdraw     float64
}

type OperationsReportsStruct struct {
	BetMoney   float64
	ValidMoney float64
	NetMoney   float64
}

type OperationsDividendMoneyStruct struct {
	Money float64
}

type OperationsRebateMoneyStruct struct {
	Money float64
}

type OperationsDepositsMoneyStruct struct {
	Money    float64
	Discount float64
}
type OperationsUserAccountSetsMoneyStruct struct {
	Money float64
}
type OperationsWithdrawsMoneyStruct struct {
	Money float64
}
type OperationsFirstDepositAmountMoneyStruct struct {
	Money float64
}

type OperationsAgentCommissionLogsStruct struct {
	Money float64
}

func GetDataPage(page, pageSize int, data []string) []string {
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

// 经营报表
//投注人数文案调整为”活跃人数“ ”返水“文案调整为”活动奖金“
var ReportOperations = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) { //默认首页
		//	pageList := make([]string, 0)
		var startDay string
		var endDay string
		if value, exists := c.GetQuery("created"); !exists {
			startDay = time.Now().Format("2006-01-02")
			endDay = startDay
			//pageList = append(pageList, currentDayTime)
		} else {
			areas := strings.Split(value, " - ")
			endDay = areas[1][0:10]
			startDay = areas[0][0:10]
			//for i := startDate; i >= startEnd; i = i - 24*60*60 {
			//	pageList = append(pageList, time.Unix(i, 0).Format("2006-01-02"))
			//}
		}
		//pageStr := c.DefaultQuery("page", "1")
		//page, _ := strconv.Atoi(pageStr)
		//pageSize := 15
		//	total := len(pageList)
		//newPageList := make([]string, 0)
		//newPageList = GetDataPage(page, pageSize, pageList)

		//pageRegisteredNumber := int64(0)
		//pageFirstDepositNumber := int64(0)
		////pageFirstDepositAmount := 0.00
		//pageDepositsNumber := int64(0)
		//pageWithdrawalsNumber := int64(0)
		//pageDepositAmount := 0.00
		//pageWithdrawalAmount := 0.00
		//pageDepositOffer := 0.00
		//pageBetNumber := int64(0)
		//pageValidBet := 0.0
		//pageTotalBet := 0.00
		//pageCompanyWinsLoses := 0.00
		//pageDividendMoney := 0.00
		//pageRebateMoney := 0.00
		//pageDepositDiscount := 0.00
		//pageAgentCommission := 0.00
		//pageCompanyIncome := 0.00
		//
		//pageOnlineDepositNum := 0
		//pageOnlineDepositMoney := 0.00
		//pageOffLineDepositNum := 0
		//pageOfflineDepositMoney := 0.00
		//pageHumanDepositNum := 0
		//pageHumanDepositMoney := 0.00
		//pageHumanWithdrawNum := 0
		//pageHumanWithdrawMoney := 0.00
		//pageVirtualWithdrawMoney := 0.00
		//pageVirtualDepositMoney := 0.00
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		reportOperationsList := make([]ReportOperationsStruct, 0)
		//	for _, v := range newPageList {

		temp := ReportOperationsStruct{}
		temp.Date = startDay + " ~ " + endDay
		startAt := tools.GetMicroTimeStampByString(startDay + " 00:00:00")
		endAt := tools.GetMicroTimeStampByString(endDay + " 23:59:59")

		startSecondAt := tools.MicroToSecond(startAt)
		endSecondAt := tools.MicroToSecond(endAt)

		testSql := "select sum(money) as onlineDeposit," +
			"(select sum(money) from user_deposits where type=2 and status=2 and created>=? and created<?)as offLineDeposit," +
			"(select sum(money) from user_account_sets where status=2 and updated>=? and updated<? and type=1 and reason='人工代充') as humanAdd," +
			"(select sum(adjust_money)  from user_resets where status=2 and adjust_method=2 and updated>=? and updated<?) as humanReduce, " +
			"(select sum(amount) from transactions where description='存款-订单状态反转' and created>=? and created<?) as turnMoney," +
			"(SELECT count(user_id)  FROM `user_daily_reports` WHERE day>=? and day<=? and is_agent=0 and bet_money>0 ) as bet_num ," +
			"(select sum(money) from user_rebate_records where day>=? and day<=?) as rebate_money," +
			"(select sum(money) from user_dividends where state=2 and type!=5 and created>=? and created<?) as dividend_money ," +
			"(select sum(award) from activity_applies where status=2 and created>=? and created<?) as ActivityMoney," +
			"(select sum(money) from user_account_sets where status=2 and created>=? and created<? and reason!='人工代充' and type=1) as ArtificialMoney ," +
			"(SELECT count(*) FROM `users` WHERE created>=? AND created<=?) as num," +
			"(SELECT count(user_id) FROM `user_deposits` WHERE status=2  AND is_first_deposit=2 and confirm_at>=? and confirm_at<=?) as reFirstDeposit," +
			"(SELECT count(user_id) FROM `user_deposits` WHERE status=2 AND created>=? AND created<=?) as sumDepositsNumber," +
			"(SELECT count(user_id) FROM `user_withdraws` WHERE status=2 AND created>=? AND created<=? ) as sumWithdrawalNumber," +
			//"(select sum(money) from user_account_sets where status=2 and reason='活动红利' and updated>=? and updated<?) as TotalUserAccountSetsSum," +
			"(select sum(money) from user_deposits where type=4 and status=2 and created>=? and created<?)as VirtualDeposit," +
			"(select sum(money) from user_withdraws where status=2 and created>=? and created<? and wallet_id>0) as VirtualWithdraw," +
			"(select sum(money) from user_withdraws where status=2 and created>=? and created<?) as TotalUserWithdrawsSum" +
			" from user_deposits where type=1 and status=2 and confirm_at>=? and confirm_at<?"

		//在线存款人数,在线存款金额
		Res, err := engine.QueryString(testSql, startAt, endAt, startAt, endAt, startAt, endAt, startAt, endAt, startDay, endDay, startDay, endDay,
			startAt, endAt, startAt, endAt, startAt, endAt, startAt,
			endAt, startSecondAt, endSecondAt, startAt, endAt, startAt, endAt, startAt, endAt, startAt,
			endAt, startAt, endAt, startSecondAt, endSecondAt)
		if err != nil {
			log.Err(err.Error())
		}

		//注册人数

		TotalRegisteredNumber, _ := strconv.Atoi(Res[0]["num"])
		temp.RegisteredNumber = int64(TotalRegisteredNumber)
		//	pageRegisteredNumber += temp.RegisteredNumber
		//首存人数--is_first_deposit确定唯一

		sumFirstDepositNumber, _ := strconv.Atoi(Res[0]["reFirstDeposit"])
		temp.FirstDepositNumber = int64(sumFirstDepositNumber)

		//	pageFirstDepositNumber += temp.FirstDepositNumber
		//转化率

		//首存金额--is_first_deposit确定唯一

		//存款人数

		sumDepositsNumber, _ := strconv.Atoi(Res[0]["sumDepositsNumber"])
		temp.DepositsNumber = int64(sumDepositsNumber)
		//	pageDepositsNumber += temp.DepositsNumber
		//提款人数

		sumWithdrawalNumber, _ := strconv.Atoi(Res[0]["sumWithdrawalNumber"])
		temp.WithdrawalsNumber = int64(sumWithdrawalNumber)
		//pageWithdrawalsNumber += temp.WithdrawalsNumber

		//存款金额
		userDepositsSum := new(OperationsDepositsMoneyStruct)
		TotalUserDepositsSum, _ := engine.Table("user_deposits").Where(builder.NewCond().
			And(builder.Eq{"status": 2}).
			And(builder.Gte{"created": startAt}).
			And(builder.Lte{"created": endAt})).
			Sums(userDepositsSum, "money", "discount")

		//新增手动上分算存款   手动上分。 存款单在线。离线。 上分。

		//	TotalUserAccountSetsSum, _ := strconv.ParseFloat(Res[0]["TotalUserAccountSetsSum"], 64)

		onlineDepositMoney, _ := strconv.ParseFloat(Res[0]["onlineDeposit"], 64)
		temp.OnlineDepositMoney = onlineDepositMoney

		offLineDepositMoney, _ := strconv.ParseFloat(Res[0]["offLineDeposit"], 64)
		temp.OfflineDepositMoney = offLineDepositMoney

		onlineDepositNumSql := "select user_id from user_deposits where type=1 and status=2 and confirm_at>=? and confirm_at<? group by user_id"
		onlineDepositNum, _ := engine.QueryString(onlineDepositNumSql, startAt, endAt)
		temp.OnlineDepositNum = len(onlineDepositNum)
		//离线存款人数，离线存款金额

		offLineDepositNumSql := "select user_id from user_deposits where type!=1 and status=2 and confirm_at>=? and confirm_at<? group by user_id"
		offLineDepositNum, _ := engine.QueryString(offLineDepositNumSql, startAt, endAt)
		temp.OffLineDepositNum = len(offLineDepositNum)

		//人工存款人数，人工存款金额 后台人工上分和人工添加存款单算人工存款， 人數要去重
		humanAddMoney, _ := strconv.ParseFloat(Res[0]["humanAdd"], 64)
		temp.HumanDepositMoney = humanAddMoney //+ adjustAddMoney //offLineDepositMoney + onlineDepositMoney

		//存款金额=在线+离线+人工存款 -  代客充值的其他原因
		TotalArtificialMoneySum, _ := strconv.ParseFloat(Res[0]["ArtificialMoney"], 64)
		temp.DepositAmount = TotalUserDepositsSum[0] - TotalArtificialMoneySum
		//	pageDepositAmount += temp.DepositAmount

		humanAddNumSql := "select user_id from user_account_sets where status=2 and updated>=? and updated<? and type=1  and reason='人工代充' group by user_id"
		humanAddNumRes, _ := engine.QueryString(humanAddNumSql, startAt, endAt)

		humanAddNum := "select user_id from user_resets where status=2 and adjust_method=1 and updated>=? and updated<?  group by user_id"
		humanAddNums, _ := engine.QueryString(humanAddNum, startAt, endAt)

		tempNum := make(map[string]string, 0)
		//for _, v := range onlineDepositNum {
		//	tempNum[v["user_id"]] = v["user_id"]
		//}
		//for _, v := range offLineDepositNum {
		//	tempNum[v["user_id"]] = v["user_id"]
		//}
		for _, v := range humanAddNumRes {
			tempNum[v["user_id"]] = v["user_id"]
		}
		for _, v := range humanAddNums {
			tempNum[v["user_id"]] = v["user_id"]
		}
		temp.HumanDepositNum = len(tempNum)

		//人工提款人数，人工提款金额 人工下分算人工提款
		//存款反轉 也算下分

		temps, _ := strconv.ParseFloat(Res[0]["humanReduce"], 64)
		turn, _ := strconv.ParseFloat(Res[0]["turnMoney"], 64)
		humanReduceMoney := temps + math.Abs(turn)
		//人工提款金额
		temp.HumanWithdrawMoney = humanReduceMoney
		//提款金额
		TotalUserWithdrawsSum, _ := strconv.ParseFloat(Res[0]["TotalUserWithdrawsSum"], 64)
		temp.WithdrawalAmount = TotalUserWithdrawsSum + humanReduceMoney
		//	pageWithdrawalAmount += temp.WithdrawalAmount

		humanReduceNumSql := "select user_id from user_resets where status=2 and adjust_method=2 and updated>=? and updated<? group by user_id"
		humanReduceNum, _ := engine.QueryString(humanReduceNumSql, startAt, endAt)

		turnNumSql := "select user_id from transactions where  description='存款-订单状态反转' and created>=? and created<? group by user_id"
		turnNum, _ := engine.QueryString(turnNumSql, startAt, endAt)

		tempNums := make(map[string]string, 0)
		for _, v := range humanReduceNum {
			tempNums[v["user_id"]] = v["user_id"]
		}
		for _, v := range turnNum {
			tempNums[v["user_id"]] = v["user_id"]
		}

		//虚拟存款
		temp.VirtualDeposit, _ = strconv.ParseFloat(Res[0]["VirtualDeposit"], 64)
		temp.VirtualWithdraw, _ = strconv.ParseFloat(Res[0]["VirtualWithdraw"], 64)
		//	pageVirtualDepositMoney = pageVirtualDepositMoney + temp.VirtualDeposit
		//	pageVirtualWithdrawMoney = pageVirtualWithdrawMoney + temp.VirtualWithdraw
		temp.HumanWithdrawNum = len(tempNums)
		//虚拟提款
		//
		//pageOnlineDepositNum += len(onlineDepositNum)
		//pageOnlineDepositMoney += onlineDepositMoney
		//pageOffLineDepositNum += len(offLineDepositNum)
		//pageOfflineDepositMoney += offLineDepositMoney
		//pageHumanDepositNum += len(tempNum)
		//pageHumanDepositMoney += temp.HumanDepositMoney
		//pageHumanWithdrawNum += temp.HumanWithdrawNum
		//pageHumanWithdrawMoney += humanReduceMoney

		// 存提差
		temp.DepositOffer = temp.DepositAmount - temp.WithdrawalAmount
		//pageDepositOffer += temp.DepositOffer

		//有效投注人数--一天

		sumBetNumber, _ := strconv.Atoi(Res[0]["bet_num"])
		temp.BetNumber = int64(sumBetNumber)
		//pageBetNumber += temp.BetNumber

		//有效投注额--一天

		UserDailyReportsSum := new(OperationsReportsStruct)
		TotalUserDailyReportsSum, _ := engine.Table("user_daily_reports").Where(builder.NewCond().
			And(builder.Gte{"day": startDay}.And(builder.Lte{"day": endDay}))).And(builder.Neq{"game_code": "0"}).And(builder.Eq{"is_agent": 0}).
			Sums(UserDailyReportsSum, "valid_money", "net_money", "bet_money")

		temp.ValidBet = TotalUserDailyReportsSum[0]
		//pageValidBet += temp.ValidBet

		//公司输赢--玩家赢钱 我们就是输钱的
		temp.CompanyWinsLoses = -TotalUserDailyReportsSum[1]
		//	pageCompanyWinsLoses += temp.CompanyWinsLoses
		//总投注额。
		temp.TotalBet = TotalUserDailyReportsSum[2]
		//	pageTotalBet += TotalUserDailyReportsSum[2]
		//盈余比例
		if temp.ValidBet == 0 {
			temp.SurplusRatio = 0
		} else {
			temp.SurplusRatio = temp.CompanyWinsLoses / temp.ValidBet
		}
		//  红利

		TotalDividendMoneySum, _ := strconv.ParseFloat(Res[0]["dividend_money"], 64)
		//活动
		TotalActivityMoneySum, _ := strconv.ParseFloat(Res[0]["ActivityMoney"], 64)
		//人工活动

		temp.DividendMoney = TotalDividendMoneySum + TotalActivityMoneySum + TotalArtificialMoneySum

		//pageDividendMoney += temp.DividendMoney
		// 返水
		temp.RebateMoney, _ = strconv.ParseFloat(Res[0]["rebate_money"], 64)
		//	pageRebateMoney += temp.RebateMoney

		//存款优惠
		temp.DepositDiscount = TotalUserDepositsSum[1]
		//	pageDepositDiscount += temp.DepositDiscount

		//代理佣金
		//var agent []models.User
		//engine.Table("users").Where("is_agent=1").Cols("id,agent_commission").Find(&agent)
		//sumSql := "select sum(arrive_money) as deposit_money,(select sum(money)  from user_account_sets where status=2 and top_id = ? and created >= ? and  created<= ? and status=2) as account_money," +
		//	" (select sum(money) as money from user_rebate_records where top_id=? and  created>= ? and  created< ?) as rebate_money," +
		//	" (select sum(adjust_money) from user_resets where top_id=? and  updated>= ? and  updated<=?  and adjust_method=1 and status=2 ) as adjust_add," +
		//	"(select sum(adjust_money)  from user_resets where top_id=? and  updated>= ? and  updated<=?  and adjust_method=2 and status=2) as adjust_dec," +
		//	" (select sum(money) from user_dividends where top_id=? and  created>= ? and  created<= ?  and state =2) as dividend_money ," +
		//	"(select sum(award)  from activity_applies where top_id=? and  created>= ? and  created<= ? and status=2) as activity_award," +
		//	"(select sum(money)  from user_activities where top_id=? and  created>= ? and  created<=? and status=2) as activity_money" +
		//	" from user_deposits where top_id=? and confirm_at >= ? and  confirm_at<= ? and status=2 "
		//
		////场馆费率
		//vsql := "select code,platform_rate from game_venues where pid!=0 and is_online=1 and  code!='CENTERWALLET'"
		//vRes, verr := engine.QueryString(vsql)
		//if verr != nil {
		//	log.Err(verr.Error())
		//}
		////存款费率
		//depositCostSql := "select value from parameters where name='deposit_cost'"
		//depositCostRes, _ := engine.QueryString(depositCostSql)
		//depositCost, _ := strconv.ParseFloat(depositCostRes[0]["money"], 64)
		//
		//gameRateList := make(map[string]string, len(vRes))
		//for _, v := range vRes {
		//	gameRateList[v["code"]] = v["platform_rate"]
		//}
		//
		//for _, val := range agent {
		//	res, err := engine.QueryString(sumSql, val.Id, startAt, endAt, val.Id, startAt, endAt, val.Id, startAt, endAt, val.Id, startAt, endAt, val.Id, startAt, endAt, val.Id, startAt, endAt, val.Id, startAt, endAt, val.Id, startAt, endAt)
		//	if err != nil {
		//		log.Err(err.Error())
		//	}
		//
		//	cusMoney, _ := strconv.ParseFloat(res[0]["account_money"], 64)
		//	dMoney, _ := strconv.ParseFloat(res[0]["deposit_money"], 64)
		//	//投注活跃，没有
		//	rr1, _ := strconv.ParseFloat(res[0]["adjust_add"], 64)
		//	rr2, _ := strconv.ParseFloat(res[0]["adjust_dec"], 64)
		//
		//	reFloat, _ := strconv.ParseFloat(res[0]["rebate_money"], 64)
		//	diFloat, _ := strconv.ParseFloat(res[0]["dividend_money"], 64)
		//
		//	activityMoney, _ := strconv.ParseFloat(res[0]["activity_award"], 64)
		//	artificialMoney, _ := strconv.ParseFloat(res[0]["activity_money"], 64)
		//
		//	var winFloat float64
		//	reportSql := "select sum(direct_net_money) as money ,game_code from user_daily_reports where user_id=? and  day='" + v + "' and game_code!='0' group by game_code"
		//	reportRes, err := engine.QueryString(reportSql, val.Id)
		//	if err != nil {
		//		log.Err(err.Error())
		//	}
		//	//场馆费  查找所有场馆，再去算所有场馆 的下注量
		//	var gameCost float64
		//	for _, v := range reportRes {
		//		gFloat, _ := strconv.ParseFloat(v["money"], 64)
		//		if gFloat < 0 { //当总输赢为负的时候算场馆费
		//			temp := math.Abs(gFloat)
		//			rateStr := gameRateList[v["game_code"]]
		//			rate, _ := strconv.ParseFloat(rateStr, 64)
		//			gameCost = gameCost + temp*rate
		//		}
		//		winFloat = winFloat + gFloat
		//	}
		//
		//	DepositsCost := (cusMoney + dMoney) * depositCost
		//
		//	cost := reFloat + diFloat + gameCost + DepositsCost + artificialMoney + activityMoney - rr2 + rr1
		//	profit := -1*winFloat - cost
		//
		//	var pRes []models.AgentCommissionPlan
		//	perr := engine.Table("agent_commission_plans").Where("agent_commission=? and user_id=?", val.AgentCommission, val.Id).OrderBy("created desc ,level_id asc").Find(&pRes)
		//	if perr != nil {
		//		log.Err(perr.Error())
		//		return
		//	}
		//	var rate float64
		//	if len(pRes) > 0 && pRes[0].Type == 2 { //代理占成佣金比例。
		//		rate = pRes[0].Rate
		//	}
		//	//佣金没有负数
		//
		//	if val.Id == 1001 || val.Id == 1000 || val.Id == 1002 { //直属会员
		//		temp.CompanyIncome += tools.ToFixed(profit, 0)
		//
		//	} else {
		//		officeRate := 1 - rate
		//		tempMoney := profit * officeRate
		//		if tempMoney < 0 {
		//			temp.CompanyIncome += tools.ToFixed(profit, 0)
		//		} else {
		//			temp.CompanyIncome += tools.ToFixed(tempMoney, 0)
		//		}
		//		//if profit > 0 { //每日 每个代理
		//		temp.AgentCommission += tools.ToFixed(profit*rate, 0)
		//		//}
		//
		//	}
		//}
		//
		//pageAgentCommission += temp.AgentCommission
		temp.CompanyIncome = temp.CompanyWinsLoses - temp.DividendMoney - temp.RebateMoney - temp.AgentCommission - temp.DepositDiscount - temp.HumanWithdrawMoney //- TotalUserDepositsSums
		//pageCompanyIncome = pageCompanyIncome + temp.CompanyIncome
		reportOperationsList = append(reportOperationsList, temp)
		//	}

		//小计转化率

		//小计人均首存

		//小计提存率

		// 盈余比例

		viewData := pongo2.Context{
			"rows": reportOperationsList,
			//"total": total,
			//"pageRegisteredNumber":   pageRegisteredNumber,
			//"pageFirstDepositNumber": pageFirstDepositNumber,
			////"FirstDepositPer":        FirstDepositPer,
			////"pageFirstDepositAmount": pageFirstDepositAmount,
			//
			//"pageDepositsNumber":    pageDepositsNumber,
			//"pageWithdrawalsNumber": pageWithdrawalsNumber,
			//"pageDepositAmount":     pageDepositAmount,
			//"pageWithdrawalAmount":  pageWithdrawalAmount,
			//"pageDepositOffer":      pageDepositOffer,
			////"pageWithdrawalRate":      pageWithdrawalRate,
			//"pageBetNumber":            pageBetNumber,
			//"pageValidBet":             pageValidBet,
			//"pageTotalBet":             pageTotalBet,
			//"pageCompanyWinsLoses":     pageCompanyWinsLoses,
			//"pageDividendMoney":        pageDividendMoney,
			//"pageRebateMoney":          pageRebateMoney,
			//"pageDepositDiscount":      pageDepositDiscount,
			//"pageAgentCommission":      pageAgentCommission,
			//"pageCompanyIncome":        pageCompanyIncome,
			//"pageOnlineDepositNum":     pageOnlineDepositNum,
			//"pageOnlineDepositMoney":   pageOnlineDepositMoney,
			//"pageOffLineDepositNum":    pageOffLineDepositNum,
			//"pageOfflineDepositMoney":  pageOfflineDepositMoney,
			//"pageHumanDepositNum":      pageHumanDepositNum,
			//"pageHumanDepositMoney":    pageHumanDepositMoney,
			//"pageHumanWithdrawNum":     pageHumanWithdrawNum,
			//"pageHumanWithdrawMoney":   pageHumanWithdrawMoney,
			//"pageVirtualDepositMoney":  pageVirtualDepositMoney,
			//"pageVirtualWithdrawMoney": pageVirtualWithdrawMoney,
		}
		viewFile := "report_operations/list.html"
		if request.IsAjax(c) {
			viewFile = "report_operations/_list.html"
		}
		response.Render(c, viewFile, viewData)
	},
}
