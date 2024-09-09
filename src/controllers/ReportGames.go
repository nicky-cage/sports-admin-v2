package controllers

import (
	"math"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"xorm.io/builder"
)

type ListReportGamesStruct struct {
	Date             string
	BetNumber        int64   //投注人数
	BetFrequency     int64   //投注次数
	ValidBet         float64 //有效投注
	CompanyWinsLoses float64 //公司输赢
	SurplusRatio     float64 //盈余比例
	GameCost         float64 //场馆费
}

type PageReportGamesStruct struct {
	PageBetNumber        int64   //每页投注人数
	PageBetFrequency     int64   //每页投注次数
	PageValidBet         float64 //每页有效投注
	PageCompanyWinsLoses float64 //每页公司输赢
	PageSurplusRatio     float64 //每页盈余比例
	PageGameCost         float64 //每页场馆费
}

type ReportGamesDataStruct struct {
	Name string
	List []ListReportGamesStruct
	Page PageReportGamesStruct
}

func GetGameTypeList(engine *xorm.Session, pageList []string, gameCode string, gameType int) ReportGamesDataStruct {
	ReportGamesData := ReportGamesDataStruct{}
	PageReportGames := PageReportGamesStruct{}
	ReportGamesList := make([]ListReportGamesStruct, 0)
	//场馆费率
	var gameFee []models.GameVenue
	engine.Table("game_venues").Where("pid=0").Select("code,platform_rate").Find(&gameFee)
	gameFeeList := make(map[string]float64, len(gameFee))
	for _, v := range gameFee {
		gameFeeList[v.Code] = v.PlatformRate
	}

	for _, v := range pageList {
		ReportGames := ListReportGamesStruct{}
		ReportGames.Date = v
		//有效投注人数--一天
		sqlBetNumber := "SELECT user_id FROM `user_daily_reports` WHERE day=? and game_code=? and game_type=? group by user_id"
		reBetNumber, _ := engine.QueryString(sqlBetNumber, v, gameCode, gameType)
		sumBetNumber := len(reBetNumber)
		ReportGames.BetNumber = int64(sumBetNumber)
		PageReportGames.PageBetNumber += ReportGames.BetNumber
		//有效投注次数--一天
		UserDailyReportsSum := new(OperationsReportsStruct)
		TotalUserDailyReportsSum, _ := engine.Table("user_daily_reports").Where(builder.NewCond().
			And(builder.Eq{"game_type": gameType}).
			And(builder.Eq{"game_code": gameCode}).
			And(builder.Eq{"day": v})).
			Sums(UserDailyReportsSum, "number_bet", "bet_money", "net_money") //net_money 改爲bet_money
		ReportGames.BetFrequency = int64(TotalUserDailyReportsSum[0])
		PageReportGames.PageBetFrequency += ReportGames.BetFrequency
		//有效投注额--一天
		ReportGames.ValidBet = TotalUserDailyReportsSum[1]
		PageReportGames.PageValidBet += ReportGames.ValidBet
		//公司输赢--一天
		if float64(TotalUserDailyReportsSum[2]) != 0.0 {
			ReportGames.CompanyWinsLoses = -1 * TotalUserDailyReportsSum[2]
		}
		//场馆费
		if float64(TotalUserDailyReportsSum[2]) < 0 {
			ReportGames.GameCost = math.Abs(float64(TotalUserDailyReportsSum[2]) * gameFeeList[gameCode])
			PageReportGames.PageGameCost += ReportGames.GameCost
		}
		PageReportGames.PageCompanyWinsLoses += ReportGames.CompanyWinsLoses
		//盈余比例
		if TotalUserDailyReportsSum[1] == 0 {
			ReportGames.SurplusRatio = 0.00
		} else {
			ReportGames.SurplusRatio = TotalUserDailyReportsSum[2] / TotalUserDailyReportsSum[1]
		}
		ReportGamesList = append(ReportGamesList, ReportGames)
	}
	//盈余比例
	if PageReportGames.PageValidBet == 0 {
		PageReportGames.PageSurplusRatio = 0.00
	} else {
		PageReportGames.PageSurplusRatio = PageReportGames.PageCompanyWinsLoses / PageReportGames.PageValidBet
	}
	ReportGamesData.Name = gameCode
	ReportGamesData.List = ReportGamesList
	ReportGamesData.Page = PageReportGames
	return ReportGamesData
}

// 游戏报表
//投注人数文案调整为”活跃人数“ ”返水“文案调整为”活动奖金“
var ReportGames = struct {
	List func(*gin.Context)
}{
	List: func(c *gin.Context) { //默认首页-游戏报表
		pageList := make([]string, 0)
		isFirst := true
		if value, exists := c.GetQuery("created"); !exists {
			currentDayTime := time.Now().Format("2006-01-02")
			pageList = append(pageList, currentDayTime)
		} else {
			areas := strings.Split(value, " - ")
			startDate := tools.GetTimeStampByString(areas[1])
			startEnd := tools.GetTimeStampByString(areas[0])
			for i := startDate; i >= startEnd; i = i - 24*60*60 {
				pageList = append(pageList, time.Unix(i, 0).Format("2006-01-02"))
			}
			isFirst = false
		}
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		viewData := pongo2.Context{}
		sport := make([]ReportGamesDataStruct, 0)
		realPerson := make([]ReportGamesDataStruct, 0)
		lottery := make([]ReportGamesDataStruct, 0)
		chess := make([]ReportGamesDataStruct, 0)
		game := make([]ReportGamesDataStruct, 0)
		electronic := make([]ReportGamesDataStruct, 0)
		fishing := make([]ReportGamesDataStruct, 0)
		tempSport := []string{"IM_1", "BTI_1", "SABA_1", "XJ188_1", "BAOLI_1"}
		tempRealPerson := []string{"AG_3", "EBET_3", "BBIN_3"}                         //,, "BG_3"
		tempLottery := []string{"VR_6", "SG_6"}                                        //"SG_6", "LB_6",
		tempChess := []string{"KY_7", "LEG_7", "VG_7", "DT_7", "HL_7"}                 //"SG_7", , "GF_7"
		tempGame := []string{"LEIHUO_2", "IM_2", "AVIA_2"}                             //"IM_2", "FY_2",
		tempElectronic := []string{"AG_4", "MG_4", "BBIN_4", "PT_4", "JDB_4", "CQ9_4"} //"PG_4", "JDB_4",
		tempFishing := []string{"JDB_5"}                                               //FG_5 JDB_5

		if !isFirst {
			tempSport = make([]string, 0)
			for i := 0; i <= 3; i++ {
				temp := c.DefaultQuery("sport["+strconv.Itoa(i)+"]", "")
				if len(temp) > 0 {
					tempSport = append(tempSport, temp)
				}
			}
			tempRealPerson = make([]string, 0)
			for i := 0; i <= 3; i++ {
				temp := c.DefaultQuery("realPerson["+strconv.Itoa(i)+"]", "")
				if len(temp) > 0 {
					tempRealPerson = append(tempRealPerson, temp)
				}
			}
			tempLottery = make([]string, 0)
			for i := 0; i <= 2; i++ {
				temp := c.DefaultQuery("lottery["+strconv.Itoa(i)+"]", "")
				if len(temp) > 0 {
					tempLottery = append(tempLottery, temp)
				}
			}
			tempChess = make([]string, 0)
			for i := 0; i <= 4; i++ {
				temp := c.DefaultQuery("chess["+strconv.Itoa(i)+"]", "")
				if len(temp) > 0 {
					tempChess = append(tempChess, temp)
				}
			}
			tempGame = make([]string, 0)
			for i := 0; i <= 2; i++ {
				temp := c.DefaultQuery("game["+strconv.Itoa(i)+"]", "")
				if len(temp) > 0 {
					tempGame = append(tempGame, temp)
				}
			}
			tempElectronic = make([]string, 0)
			for i := 0; i <= 4; i++ {
				temp := c.DefaultQuery("electronic["+strconv.Itoa(i)+"]", "")
				if len(temp) > 0 {
					tempElectronic = append(tempElectronic, temp)
				}
			}
			tempFishing = make([]string, 0)
			for i := 0; i <= 1; i++ {
				temp := c.DefaultQuery("fishing["+strconv.Itoa(i)+"]", "")
				if len(temp) > 0 {
					tempFishing = append(tempFishing, temp)
				}
			}

		}

		for _, v := range tempSport {
			tempArr := strings.Split(v, "_")
			tempType, _ := strconv.Atoi(tempArr[1])
			tempList := GetGameTypeList(engine, pageList, tempArr[0], tempType)
			sport = append(sport, tempList)
		}

		for _, v := range tempRealPerson {
			tempArr := strings.Split(v, "_")
			tempType, _ := strconv.Atoi(tempArr[1])
			tempList := GetGameTypeList(engine, pageList, tempArr[0], tempType)
			realPerson = append(realPerson, tempList)
		}
		for _, v := range tempLottery {
			tempArr := strings.Split(v, "_")
			tempType, _ := strconv.Atoi(tempArr[1])
			tempList := GetGameTypeList(engine, pageList, tempArr[0], tempType)
			lottery = append(lottery, tempList)
		}
		for _, v := range tempChess {
			tempArr := strings.Split(v, "_")
			tempType, _ := strconv.Atoi(tempArr[1])
			tempList := GetGameTypeList(engine, pageList, tempArr[0], tempType)
			chess = append(chess, tempList)
		}
		for _, v := range tempGame {
			tempArr := strings.Split(v, "_")
			tempType, _ := strconv.Atoi(tempArr[1])
			tempList := GetGameTypeList(engine, pageList, tempArr[0], tempType)
			game = append(game, tempList)
		}
		for _, v := range tempElectronic {
			tempArr := strings.Split(v, "_")
			tempType, _ := strconv.Atoi(tempArr[1])
			tempList := GetGameTypeList(engine, pageList, tempArr[0], tempType)
			electronic = append(electronic, tempList)
		}
		for _, v := range tempFishing {
			tempArr := strings.Split(v, "_")
			tempType, _ := strconv.Atoi(tempArr[1])
			tempList := GetGameTypeList(engine, pageList, tempArr[0], tempType)
			fishing = append(fishing, tempList)
		}
		viewData["sport"] = sport
		viewData["realPerson"] = realPerson
		viewData["lottery"] = lottery
		viewData["chess"] = chess
		viewData["game"] = game
		viewData["electronic"] = electronic
		viewData["fishing"] = fishing
		viewFile := "report_games/list.html"
		if request.IsAjax(c) {
			viewFile = "report_games/_list.html"
		}
		response.Render(c, viewFile, viewData)
	},
}
