package index

import (
	"sports-common/response"
	"sports-common/tools"

	"github.com/gin-gonic/gin"
)

var Exchange = func(c *gin.Context) {
	exchangeRate := tools.GetExchangeRate()
	response.Result(c, exchangeRate)
}
