package controllers

import (
	"context"
	"fmt"
	"reflect"
	common "sports-common"
	"sports-common/es"
	"sports-common/log"
	"sports-common/pgsql"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"sports-common/utils"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"xorm.io/builder"
)

type ReportUsersStruct struct {
	Date               string  //统计日期
	UserId             int     // 用户编号
	Username           string  //会员账号
	TopName            string  //上级代理
	DepositAmount      float64 //存款金额
	FirstDepositAmount float64 //首存金额
	WithdrawalAmount   float64 //提款金额
	ValidBet           float64 //有效投注
	TotalNetMoney      float64 //总输赢
	DividendMoney      float64 //红利
	RebateMoney        float64 //返水
	WinLoseAdjustment  float64 //输赢调整
	BetNums            int     //投注笔数
	TimeStart          int     // 本条记录开始时间
	TimeEnd            int     // 本条记录结束时间
}

type ReportUsersSumStruct struct {
	DepositMoney      float64
	FirstDepositMoney float64
	WithdrawalMoney   float64
	ValidMoney        float64
	NetMoney          float64
	Dividend          float64
	Rebate            float64
	WinloseAdjustment float64
	GameCode          string
}

// 报表获取的临时性相关字段/结构
type ReportUserRow struct {
	CountUserId int     `json:"count_user_id"`
	CountDate   string  `json:"count_date"`
	RowType     string  `json:"row_type"`
	Total       int     `json:"total"`
	Val1        float64 `json:"val1"`
	Val2        float64 `json:"val2"`
	Val3        float64 `json:"val3"`
	Val4        float64 `json:"val4"`
}

func GetReportUsersPage(page, pageSize int, data []models.UserDailyReport) []models.UserDailyReport {
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
var ReportUsers = struct {
	List   func(*gin.Context)
	Detail func(*gin.Context)
	*ActionState
}{
	List: func(c *gin.Context) { //默认首页
		platform := request.GetPlatform(c) // 平台识别号
		// ------------------------------------------- 处理查询条件 - 开始 -------------------------------------------------//
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			page = 1
		}
		pageSize := 15
		cond := builder.NewCond() // 查询条件
		startAt, endAt, startDayAt, endDayAt := func() (int64, int64, string, string) {
			value, exists := c.GetQuery("created")
			if !exists {
				currentDayTime := time.Now().Format("2006-01-02")
				outStartAt := tools.GetMicroTimeStampByString(currentDayTime + " 00:00:00")
				outEndAt := tools.GetMicroTimeStampByString(currentDayTime + " 23:59:59")
				outStartDayAt := currentDayTime
				outEndDayAt := currentDayTime
				cond = cond.And(builder.Eq{"day": outStartDayAt})
				return outStartAt, outEndAt, outStartDayAt, outEndDayAt
			}
			areas := strings.Split(value, " - ")
			outStartAt := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
			outEndAt := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
			outStartDayAt := areas[0]
			outEndDayAt := areas[1]
			cond = cond.And(builder.Gte{"day": outStartDayAt}).And(builder.Lte{"day": outEndDayAt})
			return outStartAt, outEndAt, outStartDayAt, outEndDayAt
		}()

		// -- 条件处理
		findingUser := false
		topName := c.DefaultQuery("top_name", "")
		if username := c.DefaultQuery("username", ""); username != "" {
			cond = cond.And(builder.Eq{"username": username})
			findingUser = true
		} else {
			cond = cond.And(builder.Gt{"bet_money": 0})
		}
		orderBy := func() string {
			order := c.DefaultQuery("order_by", "")
			if order != "" {
				desc := c.Query("desc")
				return order + " " + desc
			}
			return "user_id DESC"
		}()
		if len(topName) > 0 {
			cond = cond.And(builder.Eq{"top_name": topName})
		}
		// ------------------------------------------- 处理查询条件 - 结束 -------------------------------------------------//

		// 数据库连接 - mysql
		conn := common.Mysql(platform)
		defer conn.Close()
		// 数据库连接 - pgsql
		pgConn := pgsql.GetConnForReading(platform)
		defer pgConn.Close()

		newPageList := make([]models.UserDailyReport, 0)                             // 用于处理数据
		cond = cond.And(builder.Eq{"game_code": "0"}).And(builder.Eq{"is_agent": 0}) // 大于
		if err := conn.Table("user_daily_reports").Where(cond).Cols("username,user_id,top_name,day").OrderBy(orderBy).Limit(pageSize, (page-1)*pageSize).Find(&newPageList); err != nil {
			fmt.Println("获取数据出错:", err)
		}
		totalNum := new(models.UserDailyReport)
		total, _ := conn.Table("user_daily_reports").Where(cond).Count(totalNum) // 计算投注人次/人的数量
		reportUserList := make([]ReportUsersStruct, 0)                           // 用于返回给 client
		userIdStr := ""
		idArr := func() []string { // 计算总的 id
			arr := []string{}
			for k, v := range newPageList {
				idStr := strconv.Itoa(int(v.UserId))
				arr = append(arr, idStr) // 如果是只查找单个用户, 则直接返回
				if findingUser && k == 0 {
					userIdStr = idStr
					break
				}
			}
			return arr
		}()
		totalRecords := 0 // 投注总单数
		if len(idArr) > 0 {
			ids := strings.Join(idArr, ",")

			// --------- 从 pg 里面取数据, 每天的抽注单数
			uRows := map[int]map[string]int{} // 用户id -> 哪天 -> 多少单
			sql := fmt.Sprintf("SELECT TO_TIMESTAMP(created_at)::DATE AS count_date, COUNT(*) AS count_total, user_id AS count_user_id "+
				"FROM wager_records WHERE created_at >= %d AND created_at <= %d AND user_id IN (%s) GROUP BY user_id, TO_TIMESTAMP(created_at)::DATE", startAt, endAt, ids)
			totalRows := []struct {
				CountUserId int    `json:"count_user_id"`
				CountDate   string `json:"count_date"`
				CountTotal  int    `json:"count_total"`
			}{}
			_, err = pgConn.Query(&totalRows, sql)
			if err != nil {
				fmt.Println("获取注单数据有误:", err)
			}
			for _, r := range totalRows {
				if val, exists := uRows[r.CountUserId]; exists {
					if _, ext := val[r.CountDate]; !ext {
						uRows[r.CountUserId][r.CountDate] = r.CountTotal
					}
				} else {
					uRows[r.CountUserId] = map[string]int{r.CountDate: r.CountTotal}
				}
			}

			// --------- 对于其他项目进行汇总
			tRows := []ReportUserRow{}
			getField := func(field string) string {
				if findingUser {
					return "FROM_UNIXTIME(" + field + ", '%Y-%m-%d') AS count_date, " + userIdStr + " AS count_user_id"
				}
				return "FROM_UNIXTIME(" + field + ", '%Y-%m-%d') AS count_date, user_id AS count_user_id"
			}
			getGroupBy := func(field string) string {
				return "user_id, FROM_UNIXTIME(" + field + ", '%Y-%m-%d')"
			}
			sqlArr := []string{
				// 存款金额
				fmt.Sprintf("(SELECT %s, 'user_deposits' AS row_type, SUM(arrive_money) AS val1, 0 AS val2, 0 AS val3, 0 AS val4 "+
					"FROM user_deposits WHERE status = 2 AND created >= %d AND created <= %d AND user_id IN (%s) GROUP BY %s)", getField("created"), startAt, endAt, ids, getGroupBy("created")),
				// 取款金额
				fmt.Sprintf("(SELECT %s, 'user_withdraws' AS row_type,  SUM(money) AS val1, 0 AS val2, 0 AS val3, 0 AS val4 "+
					"FROM user_withdraws WHERE status = 2 AND created >= %d AND created <= %d AND user_id IN (%s) GROUP BY %s)", getField("created"), startAt, endAt, ids, getGroupBy("created")),
				// 活动 - user_activities
				fmt.Sprintf("(SElECT %s, 'user_activities' AS row_type,  SUM(money) AS val1, 0 AS val2, 0 AS val3, 0 AS val4 "+
					"FROM user_activities WHERE status = 2 AND updated >= %d AND updated <= %d AND user_id IN (%s) GROUP BY %s)", getField("updated"), startAt, endAt, ids, getGroupBy("updated")),
				// 活动 - activity_applies
				fmt.Sprintf("(SElECT %s, 'activity_applies' AS row_type,  SUM(award) AS val1, 0 AS val2, 0 AS val3, 0 AS val4 "+
					"FROM activity_applies WHERE status = 2 AND updated >= %d AND updated <= %d AND user_id IN (%s) GROUP BY %s)", getField("updated"), startAt, endAt, ids, getGroupBy("updated")),
				// 红利
				fmt.Sprintf("(SElECT %s, 'user_dividends' AS row_type,  SUM(money) AS val1, 0 AS val2, 0 AS val3, 0 AS val4 "+
					"FROM user_dividends WHERE state = 2 AND updated >= %d AND updated <= %d AND user_id IN (%s) GROUP BY %s)", getField("updated"), startAt, endAt, ids, getGroupBy("updated")),
				// 调整 - 1/上分 - 2/下分
				fmt.Sprintf("(SElECT %s, 'user_resets' AS row_type,  SUM(IF(adjust_method = 1, adjust_money, 0)) AS val1, SUM(IF(adjust_method = 2, adjust_money, 0)) AS val2, 0 AS val3, 0 AS val4 "+
					"FROM user_resets WHERE status = 2 AND updated >= %d AND updated <= %d AND adjust_method = 1 AND user_id IN (%s) GROUP BY %s)", getField("updated"), startAt, endAt, ids, getGroupBy("updated")),
				fmt.Sprintf("(SELECT day AS count_date, user_id AS count_user_id, 'user_daily_reports' AS row_type, SUM(valid_money) AS val1, SUM(net_money) AS val2, SUM(dividend) AS val3, SUM(rebate) AS val4 "+
					"FROM user_daily_reports "+
					"WHERE day >= '%s' AND day <= '%s' AND game_code = '0' AND game_type = '0' AND user_id IN (%s) GROUP BY user_id, day)", startDayAt[0:10], endDayAt[0:10], ids),
			}
			// 数据统计 - 用户进行统计
			sql = "(" + strings.Join(sqlArr, " UNION ALL ") + ")"
			err = conn.SQL(sql).Find(&tRows)
			if err != nil {
				fmt.Println("获取用户统计信息有误:", err)
			}
			for _, v := range newPageList {
				r := ReportUsersStruct{}
				iUserId := int(v.UserId)
				r.UserId = iUserId      // 用户编号
				r.Date = v.Day          // r.Date = time.Unix(startAt, 0).Format("2006-01-02") + "~" + time.Unix(endAt, 0).Format("2006-01-02")
				r.Username = v.Username // 用户名称
				r.TopName = v.TopName   // 上级名称
				r.BetNums = func() int {
					if val, exists := uRows[r.UserId]; exists { // 计算投注总数 - 单个用户
						if v, ext := val[r.Date]; ext {
							return v
						}
					}
					return 0
				}()
				// if findingUser {
				// 	r.ValidBet = v.ValidMoney
				// 	r.RebateMoney = v.Rebate
				// 	r.TotalNetMoney = v.NetMoney
				// }
				r.TimeStart = int(tools.GetTimeStampByString(v.Day + " 00:00:00")) // 开始时间 - 当天
				r.TimeEnd = r.TimeStart + 86400                                    // 结束时间 - 当天
				for _, t := range tRows {
					//if (!findingUser && t.CountUserId == iUserId) || (findingUser && t.CountDate == r.Date) {
					if t.CountDate == r.Date && t.CountUserId == r.UserId { // 用户日期/用户编号相等
						if t.RowType == "user_deposits" { // 存款
							r.DepositAmount = t.Val1
						} else if t.RowType == "user_withdraws" { // 提款
							r.WithdrawalAmount = t.Val1
						} else if t.RowType == "user_daily_reports" { // 报表统计
							r.ValidBet += t.Val1      // 有效投注 - only 非单个用户查询
							r.TotalNetMoney += t.Val2 //  - only 非单个用户查询
							r.RebateMoney += t.Val4   //  - only 非单个用户查询
						} else if t.RowType == "user_resets" { // 用户调整
							r.WinLoseAdjustment = t.Val1 - t.Val2
						} else if t.RowType == "user_dividends" || t.RowType == "user_activities" || t.RowType == "activity_applies" { // 红利 + 活动(自动) + 活动(手动)
							r.DividendMoney += t.Val1 //
						}
					}
				}
				reportUserList = append(reportUserList, r)
			}
		}

		// 页面小计
		pageDepositAmount := 0.00
		pageFirstDepositAmount := 0.00
		pageWithdrawalAmount := 0.00
		pageValidBet := 0.00
		pageTotalNetMoney := 0.00
		pageDividendMoney := 0.00
		pageRebateMoney := 0.00
		pageWinLoseAdjustment := 0.00
		for _, v := range reportUserList {
			pageDepositAmount += v.DepositAmount
			pageFirstDepositAmount += v.FirstDepositAmount
			pageWithdrawalAmount += v.WithdrawalAmount
			pageValidBet += v.ValidBet
			pageTotalNetMoney += v.TotalNetMoney
			pageDividendMoney += v.DividendMoney
			pageRebateMoney += v.RebateMoney
			pageWinLoseAdjustment += v.WinLoseAdjustment
		}

		viewData := pongo2.Context{
			"rows":                   reportUserList, //[offSet:(offSet + pageSize)],
			"total":                  total,
			"pageBetNums":            totalRecords,
			"pageDepositAmount":      pageDepositAmount,
			"pageFirstDepositAmount": pageFirstDepositAmount,
			"pageWithdrawalAmount":   pageWithdrawalAmount,
			"pageValidBet":           pageValidBet,
			"pageTotalNetMoney":      pageTotalNetMoney,
			"pageDividendMoney":      pageDividendMoney,
			"pageRebateMoney":        pageRebateMoney,
			"pageWinLoseAdjustment":  pageWinLoseAdjustment,
		}
		viewFile := utils.IfString(request.IsAjax(c), "report_users/_list.html", "report_users/list.html")
		response.Render(c, viewFile, viewData)
	},
	ActionState: &ActionState{
		Model: models.GameVenues,
		Field: "maintain",
	},
	Detail: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		username := c.Query("id")
		paramsTime := c.Query("time")
		timeArr := strings.Split(paramsTime, "~")

		db := common.Mysql(platform)
		defer db.Close()
		sql := "select code,name,venue_type from game_venues where pid != 0 and is_online=1 and code!='CENTERWALLET'"
		gRes, _ := db.QueryString(sql)
		// -- 定义数据类型, 并且返回统计信息   注单条数， 投注额。 有效投注额。 输赢，game_code
		type totalInfo struct {
			BetNums    int     `json:"bet_nums"`
			BetMoney   float64 `json:"bet_money"`
			NetMoney   float64 `json:"net_money"`
			ValidMoney float64 `json:"valid_money"`
			GameName   string  `json:"game_name"`
			GameType   uint8   `json:"game_type"`
		}

		esClient, err := es.GetClientByPlatform(platform)
		if err != nil {
			log.Err("es初始化连接出错: %v\n", err)
			response.Err(c, "es初始化连接出错")
			return
		}
		startAt := tools.GetMicroTimeStampByString(timeArr[0] + " 00:00:00")
		endAt := tools.GetMicroTimeStampByString(timeArr[1] + " 23:59:59")
		esIndexName := platform + "_wagers"
		defer esClient.Stop()
		list := make([]totalInfo, 0)
		for _, v := range gRes {
			venueType, _ := strconv.Atoi(v["venue_type"])
			totalStat := totalInfo{}
			boolQuery := elastic.NewBoolQuery()
			boolQuery.Must(elastic.NewMatchQuery("username", username))
			boolQuery.Must(elastic.NewMatchQuery("game_code", v["code"]))
			boolQuery.Must(elastic.NewMatchQuery("game_type", venueType))
			boolQuery.Filter(elastic.NewRangeQuery("created_at").Gte(startAt).Lte(endAt))
			es.AggregationStat(esClient, esIndexName, boolQuery, func(res *elastic.SearchResult) {
				totalStat.GameName = v["name"]
				totalStat.GameType = uint8(venueType)
				for k, item := range res.Each(reflect.TypeOf(totalStat)) {
					info := item.(totalInfo)
					totalStat.BetMoney += info.BetMoney
					totalStat.NetMoney += info.NetMoney
					totalStat.ValidMoney += info.ValidMoney
					totalStat.BetNums = k + 1
				}
			}, "bet_money", "net_money", "valid_money")
			if totalStat.BetMoney > 0 {
				list = append(list, totalStat)
			}
		}
		viewData := ViewData{
			"rows": list,
		}
		response.Render(c, "report_users/detail.html", viewData)
	},
}

//
func UserReportSearch(username string, start string, end string, page string, platform string) pongo2.Context {
	startEnd := tools.GetTimeStampByString(start)
	startDate := tools.GetTimeStampByString(end)
	pageSize := 15
	if startDate > time.Now().Unix() {
		startDate = time.Now().Unix()
	}
	pageList := make([]string, 0)
	for i := startDate; i >= startEnd; i = i - 24*60*60 {
		pageList = append(pageList, time.Unix(i, 0).Format("2006-01-02"))
	}
	pages, _ := strconv.Atoi(page)
	//newPageList := make([]string, 0)
	//newPageList = GetDataPage(pages, 15, pageList)
	engine := common.Mysql(platform)
	defer engine.Close()
	reportUserList := make([]ReportUsersStruct, 0)
	total := len(pageList)
	esIndexName := platform + "_wagers"
	esClient, _ := es.GetClientByPlatform(platform)

	defer esClient.Stop()
	for _, v := range pageList {

		startAt := tools.GetTimeStampByString(v + " 00:00:00")
		endAt := tools.GetTimeStampByString(v + " 23:59:59")
		startMicroAt := tools.SecondToMicro(startAt)
		endMicroAt := tools.SecondToMicro(endAt)

		boolQuery := elastic.NewBoolQuery()
		temp := ReportUsersStruct{}
		boolQuery.Must(elastic.NewMatchQuery("username", username))
		boolQuery.Filter(elastic.NewRangeQuery("created_at").Gte(startMicroAt).Lte(endMicroAt))
		res, _ := esClient.Search(esIndexName).Query(boolQuery).Do(context.Background())

		temp.BetNums = int(res.Hits.TotalHits.Value)

		temp.Date = time.Unix(startAt, 0).Format("2006-01-02") + "~" + time.Unix(endAt, 0).Format("2006-01-02")
		temp.Username = username

		userInfo := models.User{}
		if _, err := engine.Table("users").Where("username=?", temp.Username).Get(&userInfo); err != nil {
			log.Logger.Error(err.Error())
		}
		temp.TopName = userInfo.TopName
		reportUsersSum := new(ReportUsersSumStruct)
		totals, err := engine.Table("user_daily_reports").Where(builder.NewCond().
			And(builder.Eq{"username": username}).
			And(builder.Gte{"created": startMicroAt}).
			And(builder.Lte{"created": endMicroAt}).
			And(builder.Eq{"game_code": "0"}).
			And(builder.Eq{"game_type": 0})).
			Sums(reportUsersSum, "valid_money", "net_money", "dividend", "rebate")
		if err != nil {
			log.Logger.Error(err.Error())
		}
		//存款金额
		var depositTemp models.UserDeposit
		depositRes, _ := engine.Table("user_deposits").Where("status=2 and confirm_at>=? and confirm_at<=? and user_id=?", startAt, endAt, userInfo.Id).Sum(&depositTemp, "money")
		temp.DepositAmount = depositRes
		//首存金额
		var firstDeposit models.UserDeposit
		_, _ = engine.Table("user_deposits").Where("status=2 and confirm_at>=? and confirm_at<=? and is_first_deposit=2 and user_id=?", startAt, endAt, userInfo.Id).Get(&firstDeposit)
		temp.FirstDepositAmount = firstDeposit.Money
		//提款金额
		var withdraw models.UserWithdraw
		withdrawRes, _ := engine.Table("user_withdraws").Where("status=2 and updated>=? and updated<=? and user_id=? ", startMicroAt, endMicroAt, userInfo.Id).Sum(&withdraw, "money")
		temp.WithdrawalAmount = withdrawRes
		//手动活动红利
		var huhanActivity models.UserActivity
		huhanActivityRes, _ := engine.Table("user_activities").Where("status=2 and updated>=? and updated<=? and user_id=? ", startMicroAt, endMicroAt, userInfo.Id).Sum(&huhanActivity, "money")
		//自動活動紅裏
		var devidend models.UserDividend
		devidendRes, _ := engine.Table("user_dividends").Where("state=2 and created>=? and created<=? and user_id=? ", startMicroAt, endMicroAt, userInfo.Id).Sum(&devidend, "money")
		var sysActivity models.ActivityApply
		sysActivityRes, _ := engine.Table("activity_applies").Where("status=2 and updated>=? and updated<=? and user_id=? ", startMicroAt, endMicroAt, userInfo.Id).Sum(&sysActivity, "award")

		temp.ValidBet = totals[0]                                            //有效投注
		temp.TotalNetMoney = totals[1]                                       //总输赢
		temp.DividendMoney = devidendRes + huhanActivityRes + sysActivityRes //红利
		temp.RebateMoney = totals[3]                                         //返水
		var reset models.UserReset                                           //输赢调整
		resetResAdd, _ := engine.Table("user_resets").Where("status=2 and updated>=? and updated<=? and adjust_method=1 and user_id=?", startMicroAt, endMicroAt, userInfo.Id).Sum(&reset, "adjust_money")
		resetResLost, _ := engine.Table("user_resets").Where("status=2 and updated>=? and updated<=? and adjust_method=2 and user_id=?", startMicroAt, endMicroAt, userInfo.Id).Sum(&reset, "adjust_money")

		temp.WinLoseAdjustment = resetResAdd - resetResLost
		//组装数据
		reportUserList = append(reportUserList, temp)
	}

	pageDepositAmount := 0.00
	pageFirstDepositAmount := 0.00
	pageWithdrawalAmount := 0.00
	pageValidBet := 0.00
	pageTotalNetMoney := 0.00
	pageDividendMoney := 0.00
	pageRebateMoney := 0.00
	pageBetNums := 0
	pageWinLoseAdjustment := 0.00
	for _, v := range reportUserList {
		pageDepositAmount += v.DepositAmount
		pageFirstDepositAmount += v.FirstDepositAmount
		pageWithdrawalAmount += v.WithdrawalAmount
		pageValidBet += v.ValidBet
		pageTotalNetMoney += v.TotalNetMoney
		pageDividendMoney += v.DividendMoney
		pageRebateMoney += v.RebateMoney
		pageBetNums += v.BetNums
		pageWinLoseAdjustment += v.WinLoseAdjustment
	}
	offSet := pageSize * (pages - 1)
	if offSet+pageSize > total {
		pageSize = total - offSet
	}
	if total == 0 {
		offSet = 0
		pageSize = 0
	}
	viewData := pongo2.Context{
		"rows":                   reportUserList[offSet:(offSet + pageSize)],
		"total":                  total,
		"pageDepositAmount":      pageDepositAmount,
		"pageFirstDepositAmount": pageFirstDepositAmount,
		"pageWithdrawalAmount":   pageWithdrawalAmount,
		"pageValidBet":           pageValidBet,
		"pageTotalNetMoney":      pageTotalNetMoney,
		"pageDividendMoney":      pageDividendMoney,
		"pageRebateMoney":        pageRebateMoney,
		"pageBetNums":            pageBetNums,
		"pageWinLoseAdjustment":  pageWinLoseAdjustment,
	}

	return viewData

}
