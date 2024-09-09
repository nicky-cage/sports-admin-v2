package users

import (
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Password 修改用户密码
func (ths *Users) Password(c *gin.Context) { //默认首页
	userId := 0
	if idStr, exists := c.GetQuery("id"); exists {
		if id, err := strconv.Atoi(idStr); err == nil {
			userId = id
		}
	}
	response.Render(c, "users/password.html", response.ViewData{
		"r": map[string]int{"Id": userId},
	})
}
