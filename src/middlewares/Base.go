package middlewares

import (
	"sports-common/request"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// Render 通用默认输出渲染
func Render(c *gin.Context, viewFile string, message string, args ...pongo2.Context) {
	viewData := pongo2.Context{}
	if len(args) >= 1 {
		viewData = args[0]
	}
	if request.IsAjax(c) {
		response.Err(c, message+" (ip: "+c.ClientIP()+")")
	} else {
		viewData["message"] = message          // 增加错误信息显示
		response.Render(c, viewFile, viewData) // 渲染页面内容
	}
	c.Abort()
}
