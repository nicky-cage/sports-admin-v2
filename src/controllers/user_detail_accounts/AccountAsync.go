package user_detail_accounts

import (
	"sports-admin/dao"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AccountAsync 账户同步
func (ths *UserDetailAccounts) AccountAsync(c *gin.Context) {
	platform := request.GetPlatform(c)
	userIdStr := c.DefaultQuery("id", "0")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		response.Err(c, "缺少用户信息")
		return
	}
	processCode := c.DefaultQuery("code", "")
	if processCode == "game" { // 同步所有游戏账户余额
		accounts := dao.UserGames.GetAccounts(platform, userId)
		response.Result(c, accounts)
		return
	}
	dbSession := common.Mysql(platform)
	defer dbSession.Close()

	sql := "select available from accounts where user_id = " + userIdStr
	res, err := dbSession.QueryString(sql)
	if err == nil && len(res) > 0 {
		response.Result(c, res[0]["available"])
		return
	}
	response.Result(c, 0)
}
