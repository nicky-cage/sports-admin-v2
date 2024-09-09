package user_detail_accounts

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *UserDetailAccounts) Index(c *gin.Context) {
	response.Render(c, "", pongo2.Context{})
}
