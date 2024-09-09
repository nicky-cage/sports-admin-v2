package admin_authorizes

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AdminAuthorizes) Add(c *gin.Context) { //新增
	response.Render(c, "admin_authorizes/add.html", pongo2.Context{})
}
