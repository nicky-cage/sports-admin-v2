package users

import (
	"fmt"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GameUserIds 所有游戏用户 - 玩过某个游戏的所有用户
func (ths *Users) GameUserIds(c *gin.Context) {

	gameID, gErr := strconv.Atoi(c.DefaultQuery("game_id", "0"))
	if gErr != nil {
		log.Logger.Error("无法获取游戏编号:", gErr)
		response.Err(c, "获取游戏编号失败")
		return
	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()

	// 先获取游戏信息
	gSQL := "SELECT * FROM game_venues "
	if gameID > 0 {
		gSQL += fmt.Sprintf("WHERE id = %d", gameID)
	}
	gRows := []models.GameVenue{}
	if err := myClient.SQL(gSQL).Find(&gRows); err != nil {
		log.Logger.Error("获取游戏信息出错:", err)
		response.Err(c, "获取游戏信息出错")
		return
	}
	gArr := []string{}
	for _, g := range gRows {
		gArr = append(gArr, "'"+g.Code+"'")
	}
	if len(gArr) == 0 {
		log.Logger.Error("获取游戏相们信息为 0")
		response.Err(c, "无法获取游戏信息")
		return
	}

	// 再获最近半月投过注
	uRows := []UserIdRow{}
	dSQL := fmt.Sprintf("SELECT DISTINCT(ud.user_id) AS user_id "+
		"FROM user_deposits AS ud INNER JOIN user_games AS ug ON ud.user_id = ug.user_id "+
		"WHERE ud.status = 2 AND ud.money > 0 AND ug.game_code IN (%s)", strings.Join(gArr, ","))
	if err := myClient.SQL(dSQL).Find(&uRows); err != nil {
		log.Logger.Error("查找存款用户信息失败:", err)
		response.Err(c, "查找存款用户信息失败")
		return
	}

	// iArr := []string{}
	// for _, r := range uRows {
	// 	iArr = append(iArr, fmt.Sprintf("%d", r.UserId))
	// }
	// uIds := strings.Join(iArr, ",")
	// timeEnd := tools.Now()
	// timeBegin := timeEnd - 86400*15
	// iRows := []UserIdRow{}
	// bSQL := fmt.Sprintf("SELECT DISTINCT(user_id) AS user_id "+
	// 	"FROM wager_records WHERE user_id IN (%s) AND created_at >= %d AND created_at < %d", uIds, timeBegin, timeEnd)
	// pgClient := pgsql.GetConnForReading(platform)
	// defer pgClient.Close()
	// if _, err := pgClient.Query(&iRows, bSQL); err != nil {
	// 	log.Logger.Error("获取投过注的用户失败:", err)
	// 	response.Err(c, "获取投注用户失败")
	// 	return
	// }

	// 拼接字符串
	dArr := []int{}
	for _, r := range uRows {
		dArr = append(dArr, r.UserId)
	}

	response.Result(c, dArr)
}
