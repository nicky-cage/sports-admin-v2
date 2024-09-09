package user_detail_accounts

import (
	"fmt"
	common "sports-common"
	"sports-common/pgsql"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (ths *UserDetailAccounts) Wins(c *gin.Context) {
	id := c.Query("id")
	platform := request.GetPlatform(c)
	updated := c.DefaultQuery("updated_at", "")
	var startAt int64
	var endAt int64
	if len(updated) > 0 {
		areas := strings.Split(updated, " - ")
		startAt = tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
		endAt = tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
	} else {
		temp := time.Now().Format("2006-01-02")
		startAt = tools.GetMicroTimeStampByString(temp + " 00:00:00")
		endAt = startAt + tools.SecondToMicro(24*60*60)
	}
	sql := "select code,name,venue_type from game_venues where pid != 0 and is_online=1 and code!='CENTERWALLET'"
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	gRes, _ := dbSession.QueryString(sql)
	// -- 定义数据类型, 并且返回统计信息   注单条数， 投注额。 有效投注额。 输赢，game_code
	type totalInfo struct {
		BetNums    int     `json:"bet_nums"`
		BetMoney   float64 `json:"bet_money"`
		NetMoney   float64 `json:"net_money"`
		ValidMoney float64 `json:"valid_money"`
		GameName   string  `json:"game_name"`
		GameType   uint8   `json:"game_type"`
	}

	list := make([]totalInfo, 0)
	var totalBetMoney float64
	var totalNetMoney float64
	var totalValidMoney float64

	pConn := pgsql.GetConnForReading(platform)
	if pConn == nil {
		fmt.Println("获取PG连接失败")
		response.ErrCodeAndMsg(c, 1001, "网络错误")
		return
	}
	defer pConn.Close()
	for _, v := range gRes {
		totalStat := totalInfo{}
		cond := " user_id=" + id + " and  created_at>=%d and created_at<%d"
		cond = fmt.Sprintf(cond, startAt, endAt)
		cond += " and game_code='" + v["code"] + "' "
		cond += " and game_type=" + v["venue_type"]
		sumSql := "select sum(net_money) as net_money,sum(valid_money) as valid_money,sum(bet_money) as bet_money,count(user_id) as bet_nums from wager_records" +
			" where " + cond
		_, err := pConn.Query(&totalStat, sumSql)
		if err != nil {
			fmt.Println("获取记录信息出错:", err)
		}
		if totalStat.BetMoney > 0 {
			temp, _ := strconv.Atoi(v["venue_type"])
			totalStat.GameType = uint8(temp)
			totalStat.GameName = v["name"]
			list = append(list, totalStat)
			totalBetMoney = totalBetMoney + totalStat.BetMoney
			totalNetMoney = totalNetMoney + totalStat.NetMoney
			totalValidMoney = totalValidMoney + totalStat.ValidMoney
		}
	}

	viewData := response.ViewData{
		"rows":            list,
		"id":              id,
		"totalBetMoney":   totalBetMoney,
		"totalNetMoney":   totalNetMoney,
		"totalValidMoney": totalValidMoney,
	}
	if updated != "" {
		response.Render(c, "users/_detail_wins.html", viewData)
	} else {
		response.Render(c, "users/detail_wins.html", viewData)
	}
}
