package index

import (
	"sports-admin/dao"
	"sports-admin/validations"
	"sports-common/captchas"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var Login = func(c *gin.Context) {
	// 以下仅用于初始化密码
	// salt := tools.Secret()
	// saltPassword := tools.GetPassword("qwe123", salt)
	// sql := "UPDATE admins SET password = '" + saltPassword + "', salt = '" + salt + "' WHERE id = 200"
	// 或者直接将 admin 用户密码重置为qwe123
	// update admins set password = 'd279f8bbfaa599d9d6d0daa2aca170ff', salt = '5rNuA23OhLDLdmrscFmdfkIhUejbGas5' where id = 204;
	// fmt.Println(sql)

	// 校验验证密码/
	postedData := request.GetPostedData(c)
	platform := request.GetPlatform(c)
	if verifyCode, exists := postedData["verify_code"]; !exists {
		response.Err(c, "必须输入图像校验密码")
		return
	} else if verifyID, exists := postedData["captchaID"]; !exists {
		response.Err(c, "必须输入图像验证密码")
		return
	} else {
		isNotLocal := func() bool { // 判断是否是本地环境, 如果是, 则直接跳过图形验证码
			return c.Request.Host != "admin.sports" && c.ClientIP() != "127.0.0.1"
		}
		if isNotLocal() && !captchas.Capt(platform).Verify(verifyID.(string), verifyCode.(string)) {
			response.Err(c, "图片验证校验失败")
			return
		}
	}

	// 关于登录判断
	if err := validations.UserLogin(postedData); err != nil {
		response.Err(c, err.Error())
		return
	}

	postedData["last_ip"] = c.ClientIP()
	admin, err := dao.Index.Login(c, platform, &postedData)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	loginAdmin := models.LoginAdmin{
		Id:     admin.Id,
		Name:   postedData["username"].(string),
		RoleId: uint32(admin.RoleId),
	}
	loginAdmin.SaveInCache(c)
	response.Ok(c)
}
