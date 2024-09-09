package middlewares

import (
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// CheckLogin 检测用户是否已经登录
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Host == "admin.ip.vhost" {
			c.Next()
			return
		}

		admin, err := models.LoginAdmins.GetLoginByRequest(c)
		if err != nil {
			Render(c, "unauthorized.html", err.Error())
			return
		}
		if admin == nil {
			Render(c, "unauthorized.html", "操作失败, 请先登录后台管理系统!")
			return
		}

		admin.Refresh(c) // 刷新缓存

		c.Next()
	}
}
