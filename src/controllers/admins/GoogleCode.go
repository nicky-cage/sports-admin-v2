package admins

import (
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ths *Admins) GoogleCode(c *gin.Context) { // 修改谷歌验证
	idStr := c.DefaultQuery("id", "")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Err(c, "管理员编号有误")
		return
	}
	toStateStr := c.DefaultQuery("to_google_code", "")
	toState, err := strconv.Atoi(toStateStr)
	if err != nil {
		response.Err(c, "谷歌状态有误")
		return
	}
	if toState != 1 && toState != 2 {
		response.Err(c, "提交状态代码有误")
		return
	}

	// 修改google验证状态
	data := map[string]interface{}{
		"id":          id,
		"google_code": toState,
	}
	platform := request.GetPlatform(c)
	if err := models.Admins.Update(platform, data); err != nil {
		response.Err(c, "修改google验证失败")
		return
	}

	response.Ok(c)
}
