package controllers

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// UserAccountChanges 用户账户调整记录
var UserAccountChanges = struct {
	Index func(*gin.Context)
}{
	Index: func(c *gin.Context) { //默认首页
		response.Render(c, "user_account_changes/index.html", pongo2.Context{})
	},
}
