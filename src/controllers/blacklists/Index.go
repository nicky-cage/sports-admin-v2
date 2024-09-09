package blacklists

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"sports-common/response"
)

func (ths *BlackLists) Index(c *gin.Context) { //默认首页
	response.Render(c, "blacklists/index.html", pongo2.Context{})
}
