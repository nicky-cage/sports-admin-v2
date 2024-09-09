package index

import (
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var GoogleBind = func(c *gin.Context) { // 绑定google验证码
	data := request.GetPostedData(c)
	secretOri, exists := data["secret"]
	if !exists {
		response.Err(c, "绑定google验证码失败")
		return
	}
	admin := base_controller.GetLoginAdmin(c) // 得到登录信息
	secret := secretOri.(string)
	saveData := map[string]interface{}{
		"id":            admin.Id,         // 后台编号
		"google_secret": secret,           // 验证密钥
		"udpated":       tools.NowMicro(), // 最后修改时间
	}
	platform := request.GetPlatform(c)
	if err := models.Admins.Update(platform, saveData); err != nil {
		response.Err(c, err.Error())
		return
	}
	response.Ok(c)
}
