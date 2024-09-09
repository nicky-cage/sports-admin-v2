package index

import (
	"sports-common/captchas"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

var AddIP = func(c *gin.Context) { // 增加ip
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
	data := response.ViewData{
		"captcha":   imgInfo["data"].(string),
		"captchaID": imgInfo["captchaId"].(string),
	}
	response.Render(c, "index/add_ip.html", data)
}
