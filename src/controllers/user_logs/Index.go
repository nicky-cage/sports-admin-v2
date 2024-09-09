package user_logs

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *UserLogs) Index(c *gin.Context) { //默认首页
	response.Render(c, "user_logs/index.html", pongo2.Context{})
}
