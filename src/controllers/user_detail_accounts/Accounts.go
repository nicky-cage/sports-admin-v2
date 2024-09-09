package user_detail_accounts

import (
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *UserDetailAccounts) Accounts(c *gin.Context) {
	platform := request.GetPlatform(c)
	userIdStr := c.Query("id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		response.ErrorHTML(c, "用户编号有误")
		return
	}
	sql := fmt.Sprintf("(SELECT playname, game_code, money, created, updated FROM user_games WHERE user_id = %d and game_code!='IME') "+
		"UNION ALL "+
		"(SELECT username AS playname, 'CENTERWALLET' AS game_code, available AS money,  created, updated FROM accounts WHERE user_id = %d)", userId, userId)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	gameAccounts := []GameAccount{}
	_ = dbSession.SQL(sql).Find(&gameAccounts)
	data := pongo2.Context{
		"userId":       userId,                          // 用户编号
		"gameVenues":   caches.GameVenues.All(platform), // 所有游戏场馆
		"gameAccounts": gameAccounts,
	}
	response.Render(c, "users/detail_accounts.html", data)
}
