package admin_authorizes

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AdminAuthorizes) Edit(c *gin.Context) { //编辑
	response.Render(c, "admin_authorizes/edit.html", pongo2.Context{})
}
