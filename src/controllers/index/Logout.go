package index

import (
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

var Logout = func(c *gin.Context) {
	loginAdmin := base_controller.GetLoginAdmin(c) // 拿到登录管理
	platform := request.GetPlatform(c)
	loginAdmin.Logout(platform) // 退出登录
	response.Ok(c)
}
