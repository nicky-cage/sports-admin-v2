package index

import (
	"sports-admin/caches"
	"sports-common/captchas"
	"sports-common/config"
	"sports-common/request"
	"sports-common/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// Index 默认首页
var Index = func(c *gin.Context) { //默认首页
	platform := request.GetPlatform(c)
	if platform == "" {
		response.ErrorHTML(c, "无效网址: "+c.Request.Host)
		return
	}
	img := captchas.NewCaptchaWith(platform, 28, 127, 5)
	imgInfo := img.GenerateCaptcha()
	if imgInfo["code"].(int) == 0 {
		response.ErrorHTML(c, "生成验证码失败")
		return
	}

	site := caches.PlatformSites.GetCurrent(c)
	data := response.ViewData{
		"captcha":     imgInfo["data"].(string),
		"captchaID":   imgInfo["captchaId"].(string),
		"site":        site,
		"username":    "",
		"password":    "",
		"google_code": "",
	}

	// 仅用于测试环境
	if whiteList := config.Get("login_white_list.ip_list"); whiteList != "" {
		if rList := strings.Split(whiteList, ","); len(rList) > 0 {
			clientIP := c.ClientIP()
			for _, v := range rList {
				if v == clientIP {
					data["username"] = "robin"
					data["password"] = "qwe123"
					data["google_code"] = "123123"
					break
				}
			}
		}
	}

	response.Render(c, "index.html", data)
}
