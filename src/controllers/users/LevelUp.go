package users

import (
	"sports-admin/caches"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LevelUp UsersLevelUp 修改用户vip等级 - 不给奖励
func (ths *Users) LevelUp(c *gin.Context) {
	platform := request.GetPlatform(c)
	currentVIP, _ := strconv.Atoi(c.DefaultQuery("vip", "0"))
	userLevels := caches.UserLevels.All(platform) // 用户等级
	viewData := response.ViewData{
		"levels":     userLevels,
		"userID":     c.DefaultQuery("id", "0"), // 用户编号
		"currentVIP": currentVIP,
	}
	response.Render(c, "users/level_up.html", viewData)
}
