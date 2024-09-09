package index

import (
	"sports-admin/controllers/base_controller"
	"sports-admin/dao"
	"sports-admin/validations"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

var PasswordSave = func(c *gin.Context) {
	postedData := request.GetPostedData(c)
	if err := validations.UpdatePassword(postedData); err != nil {
		response.Err(c, err.Error())
		return
	}
	admin := base_controller.GetLoginAdmin(c)
	platform := request.GetPlatform(c)
	if err := dao.Index.UpdatePassword(platform, admin.Id, postedData); err != nil {
		response.Err(c, err.Error()) //修改密码错误
		return
	}
	response.Ok(c) //修改密码成功
}
