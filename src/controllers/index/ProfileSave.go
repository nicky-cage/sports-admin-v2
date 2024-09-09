package index

import (
	"sports-admin/controllers/base_controller"
	"sports-admin/validations"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ProfileSave = func(c *gin.Context) {
	postedData := request.GetPostedData(c)
	if err := validations.UpdateProfile(postedData); err != nil {
		response.Err(c, err.Error())
		return
	}
	admin := base_controller.GetLoginAdmin(c)
	platform := request.GetPlatform(c)
	updateData := map[string]interface{}{
		"id":       admin.Id,
		"nickname": postedData["nickname"],
		"mail":     postedData["mail"],
	}
	if err := models.Admins.Update(platform, updateData); err != nil { //如果保存失败
		response.Err(c, err.Error())
		return
	}

	response.Ok(c)
}
