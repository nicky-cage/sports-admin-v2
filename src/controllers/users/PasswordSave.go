package users

import (
	"sports-admin/validations"
	common "sports-common"
	"sports-common/consts"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PasswordSave 保存用户修改密码
func (ths *Users) PasswordSave(c *gin.Context) { //默认首页
	postedData := request.GetPostedData(c)
	if err := validations.UserPasswordSave(postedData); err != nil { //如果提交的检测失败
		response.Err(c, err.Error())
		return
	}
	userId, _ := strconv.Atoi(postedData["id"].(string))
	platform := request.GetPlatform(c)
	var user models.User
	if exists, err := models.Users.FindById(platform, userId, &user); !exists {
		response.Err(c, "相关用户信息并不存在")
		return
	} else if err != nil {
		response.Err(c, err.Error())
		return
	} else {
		password := postedData["password"].(string)
		if err = models.Users.Update(platform, map[string]interface{}{
			"id":       userId,
			"password": tools.MD5(tools.MD5(password)),
		}); err != nil {
			response.Err(c, err.Error())
			return
		}
		platform := request.GetPlatform(c)
		key := consts.CtxKeyLoginUser + postedData["id"].(string)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		redis.Del(key)
	}
	response.Ok(c)
}
