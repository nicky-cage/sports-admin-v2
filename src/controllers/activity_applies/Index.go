package activity_applies

import (
	"fmt"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (ths *ActivityApplies) List(c *gin.Context) {
	orderBy := func() string { // 得到排序条件
		limit, offset := request.GetOffsets(c)
		return fmt.Sprintf(" ORDER BY a.created DESC LIMIT %d, %d", offset, limit)
	}()
	queryCond := func() string { // 得到查询条件
		cond := ""
		if userName := c.DefaultQuery("username", ""); userName != "" {
			cond += " AND a.user_name LIKE '%" + userName + "%'"
		}
		if topName := c.DefaultQuery("top_name", ""); topName != "" {
			cond += " AND u.top_name LIKE '%" + topName + "%'"
		}
		if status := c.DefaultQuery("status", ""); status != "" {
			if status != "2" {
				cond += " AND a.status!=2"
			} else {
				cond += " AND a.status=2"
			}
		}

		applied := c.DefaultQuery("applied", "")
		if applyArr := strings.Split(applied, " - "); len(applyArr) == 2 {
			startTime := tools.GetMicroTimeStampByString(applyArr[0] + " 00:00:00")
			cond += fmt.Sprintf(" AND a.updated >= %d", startTime)
			endTime := tools.GetMicroTimeStampByString(applyArr[1] + " 23:59:59")
			cond += fmt.Sprintf(" AND a.updated <= %d", endTime)
		}
		if activityTitle := c.DefaultQuery("activity_title", ""); activityTitle != "" {
			cond += " AND c.title LIKE '%" + activityTitle + "%'"
		}
		if activityType := c.DefaultQuery("activity_type", ""); activityType != "" {
			if typeInt, err := strconv.Atoi(activityType); err == nil {
				cond += fmt.Sprintf(" AND c.activity_type = '%d'", typeInt)
			}
		}
		return cond
	}()
	sql := "SELECT a.*, a.title as info ,ROUND((a.money + a.award) * a.flow_multiple, 2) AS flow_need, u.top_name, c.title, c.activity_type " +
		"FROM activity_applies AS a JOIN users AS u on a.user_id = u.id JOIN activities AS c on a.activity_id = c.id  " +
		" WHERE 1 = 1  " + queryCond + orderBy

	sqlTotal := "SELECT COUNT(a.id) AS total, SUM(a.money) AS money, SUM(a.award) AS award " +
		"FROM activity_applies AS a join users as u on a.user_id = u.id join activities as c on a.activity_id = c.id " +
		"WHERE a.user_id <> 0  and a.status=2 " + queryCond
	platform := request.GetPlatform(c)
	session := common.Mysql(platform)
	defer session.Close()
	rows, _ := session.QueryString(sql)
	getFinishFlow := func(createdStr, gameCode string, userId string) string {
		gameInfo := strings.Split(gameCode, "-")
		if len(gameInfo) < 2 {
			return "0.00"
		}
		gameType, _ := strconv.Atoi(gameInfo[1])
		createdInt, _ := strconv.ParseInt(createdStr, 10, 64)

		// GetFlowCurrent为Create At查询
		createdInt = tools.MicroToSecond(createdInt)

		userIdInt, _ := strconv.Atoi(userId)
		flowCurrent := models.GetFlowCurrent(platform, userIdInt, createdInt, gameType)
		return fmt.Sprintf("%.2f", flowCurrent)
		//sql := fmt.Sprintf("SELECT IFNULL(SUM(valid_money), 0) AS total "+
		//	"FROM user_daily_reports "+
		//	"WHERE game_code = '%s' AND game_type = '%s' AND  user_id = '%s' AND (created >= '%s' OR updated >= '%s')", gameInfo[0], gameInfo[1], userId, createdStr, createdStr)
		//fmt.Println("sql: ", sql)
		//if rowsFlow, err := session.QueryString(sql); err == nil && len(rowsFlow) >= 1 {
		//	return rowsFlow[0]["total"]
		//}
		//return "0.00"
	}
	for _, r := range rows {
		if r["status"] == "2" {
			temp, _ := strconv.Atoi(r["updated"])
			r["updated"] = time.UnixMicro(int64(temp)).Format("2006-01-02 15:04:05")
		} else {
			r["updated"] = " "
		}

		flowFinished, _ := strconv.ParseFloat(getFinishFlow(r["created"], r["game_code"], r["user_id"]), 64)
		if flowFinished <= 0.0 {
			r["flow_finished"] = "0.00"
			continue
		}
		flowNeed, _ := strconv.ParseFloat(r["flow_need"], 64)
		if flowNeed > flowFinished { // 已打流水小于所需流水
			r["flow_finished"] = fmt.Sprintf("<span style='color: red'>%.2f</span>", flowFinished)
		} else {
			r["flow_finished"] = fmt.Sprintf("<span style='color: green'>%.2f</span>", flowFinished)
		}
	}

	totalRows, totalMoney, totalAward := func() (string, string, string) {
		rowTotal, _ := session.QueryString(sqlTotal)
		if len(rowTotal) > 0 {
			return rowTotal[0]["total"], rowTotal[0]["money"], rowTotal[0]["award"]
		}
		return "0", "0.00", "0.00"
	}()

	response.Render(c, "user_changes/_activity_applies.html", response.ViewData{
		"rows":        rows,
		"total":       totalRows,
		"total_money": totalMoney,
		"award_total": totalAward,
	})
}
