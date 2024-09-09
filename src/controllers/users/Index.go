package users

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *Users) Index(c *gin.Context) { //默认首页
	response.Render(c, "users/index.html", pongo2.Context{})
}
