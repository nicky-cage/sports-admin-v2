package admin_authorizes

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"sports-common/response"
)

func (ths *AdminAuthorizes) Index(c *gin.Context) { //默认首页
	response.Render(c, "admin_authorizes/index.html", pongo2.Context{})
}
