package index

import (
	"sports-common/captchas"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

var Captcha = func(c *gin.Context) {
	platform := request.GetPlatform(c)
	img := captchas.NewCaptchaWith(platform, 28, 127, 5)
	response.Result(c, img.GenerateCaptcha())
}
