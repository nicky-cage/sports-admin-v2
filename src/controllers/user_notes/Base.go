package user_notes

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 得到用户相关信息
var getUserInfo = func(c *gin.Context) response.ViewData {
	userIDStr := c.DefaultQuery("uid", "0")
	var userName string = ""
	userID, err := strconv.Atoi(userIDStr)
	platform := request.GetPlatform(c)
	if err == nil {
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		rows, err := dbSession.Table("users").Where("id = ?", userID).QueryString()
		if err == nil && len(rows) > 0 {
			userName = rows[0]["username"]
		}
	}

	return response.ViewData{
		"user_id":   userID,
		"user_name": userName,
	}
}

type UserNotes struct {
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
}
