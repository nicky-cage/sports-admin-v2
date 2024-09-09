package user_codes

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *UserCodes) Index(c *gin.Context) {
	response.Render(c, "", pongo2.Context{})
}
