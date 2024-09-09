package admins

import (
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ths *Admins) Kick(c *gin.Context) { // 踢用户下线
	idStr := c.DefaultQuery("id", "")
	if idStr == "" {
		response.Err(c, "错误的用户编号")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	admin := models.Admin{}
	platform := request.GetPlatform(c)
	if exists, err := models.Admins.FindById(platform, id, &admin); err != nil {
		response.Err(c, err.Error())
		return
	} else if !exists {
		response.Err(c, "用户信息不存在")
		return
	}

	models.LoginAdmins.Kick(platform, int(admin.Id)) // 将用户T下线
	response.Ok(c)
}
