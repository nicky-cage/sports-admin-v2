package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sports-admin/caches"
	"sports-admin/game_detail"
	common "sports-common"
	"sports-common/config"
	"sports-common/consts"
	"sports-common/es"
	"sports-common/log"
	"sports-common/pgsql"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/imroc/req"
	"github.com/olivere/elastic/v7"
)

// SumV 统计
type SumV struct {
	Value float64 `json:"value"`
}

// AllWagerRecord 记录
type AllWagerRecord struct {
	WagerRecord models.WagerRecord // 注单记录
	GameName    string             // 游戏名称
}

// Replenishment 补单
type Replenishment struct {
	ReplenishmentContent consts.ReplenishmentContent // 补单内容
	encoded              []byte                      // 码
	err                  error                       // 错误
}

// BuDanResp 补单
type BuDanResp struct {
	Errcode int    `json:"errcode"` // 错误代码
	Message string `json:"message"` // 信息
}

// CountByVenue 按场馆
type CountByVenue struct {
	GameCode    string  `json:"game_code"`
	TotalCount  int     `json:"total_count"`
	TotalUser   int     `json:"total_user"`
	Hot         float64 `json:"hot"`
	TotalBet    float64 `json:"total_bet"`
	TotalValid  float64 `json:"total_valid"`
	TotalRebate float64 `json:"total_rebate"`
	TotalWin    float64 `json:"total_win"`
	LastCreated int     `json:"last_created"` // 最后同步时间
	LastUpdated int     `json:"last_updated"` // 最后更新时间
	//LastPlay    int     `json:"last_play"`
}

// CountByUser 按用户
type CountByUser struct {
	UserName    string  `json:"user_name"`
	TotalCount  int     `json:"total_count"`
	TotalGame   int     `json:"total_game"`
	TotalBet    float64 `json:"total_bet"`
	TotalValid  float64 `json:"total_valid"`
	TotalRebate float64 `json:"total_rebate"`
	TotalWin    float64 `json:"total_win"`
}

// BetTotalInfo -- 定义数据类型, 并且返回统计信息
type BetTotalInfo struct {
	BetMoney    float64 `json:"bet_money"`    // 投注
	NetMoney    float64 `json:"net_money"`    // 输赢
	ValidMoney  float64 `json:"valid_money"`  // 有效投注
	RebateMoney float64 `json:"rebate_money"` // 返水
}

//
func (r *Replenishment) ensureEncoded() {
	if r.encoded == nil && r.err == nil {
		r.encoded, r.err = json.Marshal(r)
	}
}

// Encode 编码
func (r *Replenishment) Encode() ([]byte, error) {
	r.ensureEncoded()
	return r.encoded, r.err
}

// Length 长度
func (r *Replenishment) Length() int {
	r.ensureEncoded()
	return len(r.encoded)
}

// UserBets 用户投注
var UserBets = struct {
	List          func(*gin.Context)
	ListV1        func(*gin.Context)
	SetUp         func(*gin.Context)
	SaveDo        func(*gin.Context)
	Detail        func(*gin.Context)
	Sync          func(*gin.Context)
	CountByVenues func(*gin.Context)
	CountByUsers  func(*gin.Context)
	Verify        func(*gin.Context)
}{
	CountByVenues: func(c *gin.Context) {
		timeStart, timeEnd := request.GetTimesByQuery(c, "created") // 左开右闭
		platform := request.GetPlatform(c)
		userSQL := ""
		if val, exists := c.GetQuery("username"); exists {
			userSQL = "AND username = '" + val + "'"
		}
		sql := fmt.Sprintf("SELECT game_code, COUNT(id) AS total_count, "+
			"COUNT(DISTINCT(username)) AS total_user, 0 AS hot, "+
			"SUM(bet_money) AS total_bet, SUM(valid_money) AS total_valid, "+
			"SUM(net_money) AS total_win, SUM(rebate_money) AS total_rebate, "+
			"MAX(created_at) AS last_created, MAX(updated_at) AS last_updated "+
			//"MAX(created_at) AS last_play FROM wager_records GROUP BY game_code ORDER BY total_user DESC"
			"FROM wager_records WHERE created_at >= %d AND created_at <= %d %s GROUP BY game_code", timeStart, timeEnd, userSQL)
		pConn := pgsql.GetConnForReading(platform)
		defer pConn.Close()
		rows := []CountByVenue{}
		_, err := pConn.Query(&rows, sql)
		if err != nil {
			fmt.Println("查询统计信息出错:", err)
		}

		mConn := common.Mysql(platform)
		defer mConn.Close()
		if rs, err := mConn.QueryString("SELECT COUNT(*) AS total FROM users"); err == nil && len(rows) > 0 {
			userTotal, _ := strconv.Atoi(rs[0]["total"])
			if userTotal > 0 {
				for k, r := range rows {
					if r.TotalUser > 0 {
						hot := float64(r.TotalUser) / float64(userTotal) * 100.0
						rows[k].Hot = hot
					}
				}
			}
		}

		response.Render(c, "user_bets/_count_by_venues.html", ViewData{
			"currentTime": tools.Now(),
			"rows":        rows,
		})
	},
	CountByUsers: func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize := 15
		platform := request.GetPlatform(c)
		timeStart, timeEnd := request.GetTimesByQuery(c, "created")
		userSQL := ""
		if val, exists := c.GetQuery("username"); exists {
			userSQL = "AND username = '" + val + "'"
		}
		sql := fmt.Sprintf("SELECT username AS user_name, COUNT(id) AS total_count, "+
			"COUNT(DISTINCT(game_code)) AS total_game,  "+
			"SUM(bet_money) AS total_bet, SUM(valid_money) AS total_valid, "+
			"SUM(net_money) AS total_win, SUM(rebate_money) AS total_rebate "+
			"FROM wager_records WHERE created_at > %d AND created_at <= %d %s GROUP BY user_name ORDER BY total_win DESC LIMIT %d OFFSET %d",
			timeStart, timeEnd, userSQL, pageSize, (page-1)*pageSize)
		pConn := pgsql.GetConnForReading(platform)
		defer pConn.Close()
		rows := []CountByUser{}
		_, err := pConn.Query(&rows, sql)
		if err != nil {
			fmt.Println("查询统计信息出错:", err)
		}
		tRow := struct {
			Total int `json:"total"`
		}{}
		sql = fmt.Sprintf("SELECT COUNT(DISTINCT(username)) AS total FROM wager_records WHERE created_at >= %d AND created_at <= %d "+userSQL, timeStart, timeEnd)
		_, err = pConn.QueryOne(&tRow, sql)
		if err != nil {
			fmt.Println("获取用户总数出错:", err)
		}

		response.Render(c, "user_bets/_count_by_users.html", ViewData{
			"rows":   rows,
			"rTotal": tRow,
		})
	},
	Sync: func(c *gin.Context) { // 查询数据统计
		platform := request.GetPlatform(c)
		if value, exists := c.GetQuery("count"); exists && value == "1" {
			value, exists := c.GetQuery("query_created")
			if !exists {
				response.Err(c, "必须提供开始日期/结束日期")
				return
			}
			areas := strings.Split(value, " - ")
			timeStartStr := areas[0] + " 00:00:00"
			timeEndStr := areas[1] + " 23:59:59"
			startAt := tools.GetTimeStampByString(timeStartStr)
			endAt := tools.GetTimeStampByString(timeEndStr)
			boolQuery := elastic.NewBoolQuery()
			boolQuery.Filter(elastic.NewRangeQuery("created_at").Gte(startAt).Lte(endAt)) // 开始时间 - 结束时间
			indexName := platform + "_wagers"
			eClient := es.GetClient(platform)
			defer es.ReturnClient(platform, eClient)
			res, err := eClient.Search(indexName).Query(boolQuery).TrackTotalHits(true).Size(0).Do(context.Background())
			if err != nil {
				fmt.Println("*** 读取ES数据发生错误:", err)
				return
			}
			totalEs := res.TotalHits() // 总记录数量
			pConn := pgsql.GetConnForReading(platform)
			defer pConn.Close()
			totalPg, err := pConn.Model((*models.WagerRecord)(nil)).Where("created_at >= ? AND created_at <= ?", startAt, endAt).Count()
			if err != nil {
				fmt.Println("获取pg数量出错:", err)
			}

			response.Result(c, struct {
				DateArea string `json:"date_area"`
				TotalEs  int64  `json:"total_es"`
				TotalPg  int64  `json:"total_pg"`
			}{
				DateArea: timeStartStr + " - " + timeEndStr,
				TotalEs:  totalEs,
				TotalPg:  int64(totalPg),
			})
			return
		}

		// -- 以下是同步数据
		value, exists := c.GetQuery("created")
		if !exists {
			response.Err(c, "")
			return
		}
		areas := strings.Split(value, " - ")
		timeStart := areas[0]
		timeEnd := areas[1]
		total := models.WagerRecords.SyncToPg(platform, timeStart, timeEnd) // 总计通步数据数量
		response.Result(c, total)
	},
	List: func(c *gin.Context) {
		getTimes := func(field string) (int64, int64) {
			value, exists := c.GetQuery(field)
			if !exists {
				currentTime := time.Now().Format("2006-01-02")
				startAt := tools.GetTimeStampByString(currentTime + " 00:00:00")
				endAt := startAt + 86399
				return startAt, endAt
			}
			areas := strings.Split(value, " - ")
			startAt := tools.GetTimeStampByString(strings.TrimSpace(areas[0]))
			endAt := tools.GetTimeStampByString(strings.TrimSpace(areas[1]))
			return startAt, endAt
		}
		timeStart, timeEnd := getTimes("created")  // 投注时间
		pgQueryBuilder := pgsql.NewQueryBuilder(). // 查询管理器
								Gte("created_at", timeStart).
								Lte("created_at", timeEnd)
		if val, exists := c.GetQuery("ignore_updated"); exists && val == "0" { // 是否忽略结算时间
			timeFrom, timeTo := getTimes("updated") // 结算时间 - 只有附合条件时才加上
			pgQueryBuilder.Gte("updated_at", timeFrom).
				Lte("updated_at", timeTo)
		}

		wSQL, wParams := pgQueryBuilder.
			Queries(c, map[string]string{
				"game_code": "game_code",
				"username":  "username",
				"playname":  "playname",
				//"bill_no":   "bill_no",
			}).
			QueriesByInt(c, map[string]string{
				"game_type": "game_type",
				"status":    "status",
			}).
			Build()
		// 追加关于金额的查询
		if moneyMinStr := c.Query("money_min"); moneyMinStr != "" {
			if moneyMin, err := strconv.Atoi(moneyMinStr); err == nil {
				wSQL += " AND valid_money >= " + strconv.Itoa(moneyMin)
			}
		}
		if moneyMaxStr := c.Query("money_max"); moneyMaxStr != "" {
			if moneyMax, err := strconv.Atoi(moneyMaxStr); err == nil {
				wSQL += " AND valid_money <= " + strconv.Itoa(moneyMax)
			}
		}
		if val, exists := c.GetQuery("bill_no"); exists && val != "" { // 关于订单编号的处理
			wSQL += fmt.Sprintf(" AND bill_no LIKE '#%v#'", val)
		}
		platform := request.GetPlatform(c)
		pConn := pgsql.GetConnForReading(platform)
		if pConn == nil {
			fmt.Println("获取PG连接失败")
			response.Err(c, "获取PG连接失败")
			return
		}
		defer pConn.Close()
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize := 15

		rows := []models.WagerRecord{}
		tableName := "wager_records"
		if currMonth := time.Unix(timeStart, 0).Format("200601"); currMonth == time.Unix(timeEnd, 0).Format("200601") {
			tableName += "_" + currMonth
		}
		sql := fmt.Sprintf("SELECT * FROM %s WHERE "+wSQL+" ORDER BY created_at DESC LIMIT %d OFFSET %d", tableName, pageSize, (page-1)*pageSize)
		sql = strings.ReplaceAll(sql, "#", "%")
		_, err := pConn.Query(&rows, sql, wParams...)
		total := 0
		if err != nil {
			fmt.Println("获取记录信息出错:", err)
		}
		totalInfo := struct {
			Total int `json:"total"`
		}{}
		if _, err = pConn.QueryOne(&totalInfo, fmt.Sprintf("SELECT COUNT(*) AS total FROM %s WHERE "+wSQL, tableName), wParams...); err == nil {
			total = totalInfo.Total
		}

		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		gameVenus := make([]models.GameVenue, 0)
		if err := dbSession.Table("game_venues").Where("pid > 0 AND venue_type > 0").Find(&gameVenus); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "无法查找到场馆信息: "+err.Error())
			return
		}
		temp := make(map[string]string)
		for _, v := range gameVenus {
			temp[v.Code+"_"+strconv.Itoa(int(v.VenueType))] = v.Name
		}

		newRows := make([]AllWagerRecord, len(rows))
		for i, v := range rows { // 修正 IM 类型错误历史遗留问题
			if v.GameCode == "IM" {
				rows[i].GameType = consts.SportsType
				rows[i].GameCodeType = rows[i].GameCode + "-" + strconv.Itoa(rows[i].GameType)
				v.GameType = consts.SportsType
				v.GameCodeType = v.GameCode + "-" + strconv.Itoa(v.GameType)
			}
			newRows[i].WagerRecord = v
			newRows[i].GameName = temp[v.GameCode+"_"+strconv.Itoa(v.GameType)]
		}

		queryTotal := BetTotalInfo{}   // 本次查询总计
		tSQL := fmt.Sprintf("SELECT "+ // 用于总的查询
			"SUM(bet_money) AS bet_money, "+
			"SUM(valid_money) AS valid_money, "+
			"SUM(net_money) AS net_money, "+
			"SUM(rebate_money) AS rebate_money FROM %s WHERE "+
			wSQL, tableName)
		_, err = pConn.QueryOne(&queryTotal, tSQL, wParams...)
		if err != nil {
			fmt.Println("获取统计信息出错:", err)
		}
		pageTotal := BetTotalInfo{}
		for _, r := range rows {
			pageTotal.NetMoney += r.NetMoney
			pageTotal.BetMoney += r.BetMoney
			pageTotal.ValidMoney += r.ValidMoney
			pageTotal.RebateMoney += r.RebateMoney
		}

		viewData := ViewData{
			"rows":        newRows,
			"rows_total":  len(newRows),
			"total":       total,
			"page_total":  pageTotal,
			"query_total": queryTotal,
			"game_venue":  caches.GameVenues.All(platform),
		}
		SetLoginAdmin(c)
		viewFile := "user_bets/index.html"
		if request.IsAjax(c) {
			viewFile = "user_bets/_list_v2.html"
		}
		response.Render(c, viewFile, viewData)
	},
	ListV1: func(c *gin.Context) { //默认首页
		timeStart, timeEnd := func() (int64, int64) {
			value, exists := c.GetQuery("created")
			if !exists {
				currentTime := time.Now().Format("2006-01-02")
				startAt := tools.GetTimeStampByString(currentTime + " 00:00:00")
				endAt := startAt + 86399
				return startAt, endAt
			}
			areas := strings.Split(value, " - ")
			startAt := tools.GetTimeStampByString(areas[0])
			endAt := tools.GetTimeStampByString(areas[1])
			return startAt, endAt
		}()

		venue := c.DefaultQuery("venue", "")
		venueType := c.DefaultQuery("venue_type", "")
		username := c.DefaultQuery("username", "")
		billNo := c.DefaultQuery("bill_no", "")
		pageStr := c.DefaultQuery("page", "1")
		playName := c.DefaultQuery("playname", "")
		gameCode := c.DefaultQuery("game_code", "")
		status := c.DefaultQuery("status", "")
		page, _ := strconv.Atoi(pageStr)
		pageSize := 15

		platform := request.GetPlatform(c)
		esIndexName := platform + "_wagers"
		esClient, err := es.GetClientByPlatform(platform)
		defer es.ReturnClient(platform, esClient)
		// defer esClient.Stop()
		if err != nil {
			response.Err(c, "es连接无响应: "+err.Error())
			return
		}

		// 设置查询条件
		boolQuery := elastic.NewBoolQuery()
		if len(venue) > 0 {
			boolQuery.Must(elastic.NewMatchQuery("game_code", venue))
		}
		if len(venueType) > 0 {
			boolQuery.Must(elastic.NewMatchQuery("game_type", venueType))
		}
		if len(username) > 0 {
			boolQuery.Must(elastic.NewMatchQuery("username", username))
		}
		if len(billNo) > 0 {
			boolQuery.Must(elastic.NewMatchQuery("bill_no", billNo))
		}
		if len(playName) > 0 {
			boolQuery.Must(elastic.NewMatchQuery("playname", playName))
		}
		if len(gameCode) > 0 {
			boolQuery.Must(elastic.NewMatchQuery("game_code", gameCode))
		}
		if len(status) > 0 {
			boolQuery.Must(elastic.NewMatchQuery("status", status))
		}
		// 分页
		boolQuery.Filter(elastic.NewRangeQuery("created_at").Gte(timeStart).Lte(timeEnd))
		currentContext := context.Background() // 当前协程上下文
		res, err := esClient.Search(esIndexName).Query(boolQuery).From((page-1)*pageSize).Size(pageSize).Sort("created_at", false).Do(currentContext)
		if err != nil {
			response.Err(c, "es获取列表数据无响应~: "+err.Error())
			return
		}

		total := res.TotalHits()
		subTotal := BetTotalInfo{}
		var typ models.WagerRecord
		var rows []models.WagerRecord
		if res.Hits.TotalHits.Value > 0 {
			for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
				t := item.(models.WagerRecord)
				t.RebateMoney = tools.ToFixed(t.RebateMoney, 2)
				rows = append(rows, t)
				subTotal.BetMoney += tools.ToFixed(t.BetMoney, 2)     // 投注
				subTotal.NetMoney += tools.ToFixed(t.NetMoney, 2)     // 输赢
				subTotal.ValidMoney += tools.ToFixed(t.ValidMoney, 2) // 有效投注
				subTotal.RebateMoney += t.RebateMoney                 //返水
			}
		}

		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		gameVenus := make([]models.GameVenue, 0)
		if err := dbSession.Table("game_venues").Where("pid > 0 AND venue_type > 0").Find(&gameVenus); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "无法查找到场馆信息: "+err.Error())
			return
		}
		temp := make(map[string]string)
		for _, v := range gameVenus {
			temp[v.Code+"_"+strconv.Itoa(int(v.VenueType))] = v.Name
		}

		newRows := make([]AllWagerRecord, len(rows))
		for i, v := range rows {
			// 修正 IM 类型错误历史遗留问题
			if v.GameCode == "IM" {
				rows[i].GameType = consts.SportsType
				rows[i].GameCodeType = rows[i].GameCode + "-" + strconv.Itoa(rows[i].GameType)
				v.GameType = consts.SportsType
				v.GameCodeType = v.GameCode + "-" + strconv.Itoa(v.GameType)
			}
			newRows[i].WagerRecord = v
			newRows[i].GameName = temp[v.GameCode+"_"+strconv.Itoa(v.GameType)]
		}

		queryTotal := BetTotalInfo{} // 本次查询总计
		es.AggregationStat(esClient, esIndexName, boolQuery, func(res *elastic.SearchResult) {
			for _, item := range res.Each(reflect.TypeOf(queryTotal)) {
				info := item.(BetTotalInfo)
				queryTotal.BetMoney += tools.ToFixed(info.BetMoney, 2)
				queryTotal.NetMoney += tools.ToFixed(info.NetMoney, 2)
				queryTotal.ValidMoney += tools.ToFixed(info.ValidMoney, 2)
				queryTotal.RebateMoney += tools.ToFixed(info.RebateMoney, 2)
			}
		}, "bet_money", "net_money", "valid_money", "rebate_money")

		SetLoginAdmin(c)
		viewData := ViewData{
			"rows":        newRows,
			"rows_total":  len(newRows),
			"total":       total,
			"page_total":  subTotal,
			"query_total": queryTotal,
			"game_venue":  caches.GameVenues.All(platform),
			// "all_total":   allTotal,
		}
		viewFile := "user_bets/index.html"
		if request.IsAjax(c) {
			viewFile = "user_bets/_list.html"
		}
		response.Render(c, viewFile, viewData)
	},
	SetUp: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		gameVenus := make([]models.GameVenue, 0)
		// if err := engine.Table("game_venues").Where("pid = 0 and ename != 'CENTERWALLET'").Find(&gameVenus); err != nil {
		// 	log.Logger.Error(err.Error())
		// 	response.Err(c, "系统繁忙")
		// 	return
		// }
		if err := engine.Table("game_venues").Where("pid != 0 and ename != 'CENTERWALLET' and is_online=1").Find(&gameVenus); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		viewData := pongo2.Context{
			"rows": gameVenus,
		}
		viewFile := "user_bets/set_up.html"
		response.Render(c, viewFile, viewData)
	},
	SaveDo: func(c *gin.Context) { //补单
		postedData := request.GetPostedData(c)
		gameVenusId := postedData["game_venus"].(string)
		replenishmentTime := postedData["replenishment_time"].(string) //就是pid=0的code的值
		timeSlice := strings.Split(replenishmentTime, " - ")
		if len(timeSlice) < 2 {
			response.Err(c, "时间格式不对")
			return
		}
		startTimeArr := strings.Split(timeSlice[0], " ")
		endTimeArr := strings.Split(timeSlice[1], " ")
		if len(startTimeArr) < 2 || len(endTimeArr) < 2 {
			response.Err(c, "时间格式不对")
			return
		}
		if startTimeArr[0] != endTimeArr[0] {
			response.Err(c, "时间范围，必须是1天以内")
			return
		}
		platform := request.GetPlatform(c)
		redisClient := common.Redis(platform)
		defer common.RedisRestore(platform, redisClient)
		rKey := consts.ReplenishmenKey + gameVenusId
		val, err := redisClient.Get(rKey).Result()
		if err == redis.Nil { //key不存在
			_ = redisClient.Incr(rKey) //防止并发点击
			_ = redisClient.Expire(rKey, 5*time.Minute)
		} else if err != nil {
			log.Err(err.Error())
			response.Err(c, "redis系统繁忙")
			return
		} else if val != "" {
			response.Err(c, "该场馆正在补单中，请5分钟之后再尝试")
			return
		}

		val, _ = redisClient.Get(rKey).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt != 1 {
			response.Err(c, "不要一直点击，已经有再处理的补单程序，5分钟之后再尝试")
			return
		}
		req.SetTimeout(50 * time.Second)
		req.Debug = true
		headerB := req.Header{
			"Accept": "application/json",
		}
		paramB := req.Param{
			"game_venus":         gameVenusId,
			"replenishment_time": replenishmentTime,
		}
		baseGameUrl := config.Get("internal.internal_game_service", "")
		GameUrl := baseGameUrl + "/game/v1/internal/budan?platform=" + platform
		rB, err := req.Post(GameUrl, headerB, req.BodyJSON(paramB))
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统内部繁忙")
			return
		}
		BuDanResp := BuDanResp{}
		BuDanJsonErr := rB.ToJSON(&BuDanResp)
		if BuDanJsonErr != nil {
			log.Logger.Error(BuDanJsonErr.Error())
			response.Err(c, "补单失败")
			return
		}
		if BuDanResp.Errcode != 0 { //接口返回失败的处理
			response.Err(c, BuDanResp.Message)
			return
		}
		response.Message(c, BuDanResp.Message)
	},
	Detail: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		billNo := c.Query("id")
		gameCode := c.Query("game_code")

		d := game_detail.NewGameDetail(gameCode)
		temp, tempList := d.Data(billNo, platform)
		nums := len(tempList)

		response.Render(c, "user_bets/detail.html", pongo2.Context{"row": temp, "rows": tempList, "total": nums})
	},
	Verify: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		value, exists := c.GetQuery("created")
		if !exists || value == "" {
			response.Err(c, "")
			return
		}
		areas := strings.Split(value, " - ")
		timeStartStr := areas[0] + " 00:00:00"
		timeEndStr := areas[1] + " 23:59:59"
		startAt := tools.GetTimeStampByString(timeStartStr) // -- 开始时间
		endAt := tools.GetTimeStampByString(timeEndStr)     // 结束时间
		pConn := pgsql.GetConnForReading(platform)
		defer pConn.Close()

		tableName := func() string { // 获取表名
			monthStart := time.Unix(startAt, 0).Format("200601")
			monthEnd := time.Unix(endAt, 0).Format("200601")
			if monthStart == monthEnd {
				return "wager_records_" + monthStart
			}
			return "wager_records"
		}()
		sql := "SELECT COUNT(DISTINCT bill_no) AS total_real, COUNT(*) AS total FROM " + tableName
		res := struct {
			TotalReal int `json:"total_real"`
			Total     int `json:"total"`
		}{}
		_, err := pConn.QueryOne(&res, sql)
		if err != nil {
			response.Err(c, err.Error())
			return
		}
		if res.TotalReal == res.Total {
			response.Err(c, "数据正确, 无需校对")
			return
		}

		// 1. 先列出有问题的数据
		rows := []struct {
			BillNo string `json:"bill_no"`
		}{}
		sql = fmt.Sprintf("SELECT bill_no FROM %s GROUP BY bill_no HAVING COUNT(id) > 1", tableName)
		_, err = pConn.Query(&rows, sql)
		if err != nil {
			response.Err(c, err.Error())
			return
		}
		if len(rows) == 0 {
			response.Err(c, "数据读取错误, 结果不可能为0")
			return
		}

		// 2. 再挨个删除有问题的数据
		pGroup := pgsql.GetConnGroup(platform)
		defer pGroup.Close()
		sCount := 0
		for _, r := range rows {
			for _, conn := range pGroup.Conns {
				sql = fmt.Sprintf("DELETE FROM %s WHERE id IN "+
					"(SELECT id FROM %s WHERE bill_no = '%s' ORDER BY id DESC LIMIT 1)", tableName, tableName, r.BillNo) // 将重复的数据从数据库当中删除
				_, err = conn.Exec(sql)
				if err != nil {
					fmt.Println("处理数据出错:", err)
					response.Err(c, "处理数据出错:"+err.Error())
					return
				}
				sCount += 1
			}
		}

		message := fmt.Sprintf("成功删除 %d 条数据", sCount)
		response.Message(c, message)
	},
}
