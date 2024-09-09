package admins

import (
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ActionState 修改状态
var ActionState = &base_controller.ActionState{
	Model: models.Admins,
	StateAfter: func(c *gin.Context) { //修改状态后处理
		idStr := c.DefaultQuery("id", "0") //检测ID
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}
		platform := request.GetPlatform(c)
		models.LoginAdmins.Kick(platform, id) // 将用户T下线
	},
}
