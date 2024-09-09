package index

import (
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

var Password = func(c *gin.Context) {
	response.Render(c, "index/password.html")
}
