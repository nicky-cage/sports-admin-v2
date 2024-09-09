package index

import (
	"net/url"
	"sports-admin/controllers/base_controller"
	"sports-common/response"
	"sports-common/tools"

	"github.com/gin-gonic/gin"
)

var GoogleCode = func(c *gin.Context) { // 显示google二维码
	auth := tools.NewGoogleAuth()
	secret := auth.GetSecret()
	admin := base_controller.GetLoginAdmin(c)
	// code, _ := auth.GetCode(secret)
	secretString := "otpauth://totp/后台_" + admin.Name + "@" + c.Request.Host + "?secret=" + secret
	response.Render(c, "index/google_code.html", response.ViewData{
		"code":           url.QueryEscape(secretString),
		"secret":         secret,
		"current_second": 30 - tools.Now()%30, // 除以30的秒数
	})
}
