package index

import (
	"errors"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// 检测用户登录基本信息
var checkLogin = func(c *gin.Context) (*models.Admin, error) {
	postedData := request.GetPostedData(c)
	if _, exists := postedData["username"]; !exists {
		return nil, errors.New("缺少用户名称")
	}
	if _, exists := postedData["password"]; !exists {
		return nil, errors.New("缺少用户密码")
	}

	var admin models.Admin
	platform := request.GetPlatform(c)
	userName := postedData["username"].(string)
	exists, err := models.Admins.Find(platform, &admin, builder.NewCond().And(builder.Eq{"name": userName}))
	if !exists || err != nil {
		return nil, errors.New("用户信息并不存在")
	}
	if admin.IsDisabled() {
		return nil, errors.New("用户已经被禁")
	}

	password := postedData["password"].(string)
	realPassword := tools.GetPassword(password, admin.Salt)
	if realPassword != admin.Password {
		return nil, errors.New("用户密码输入有误")
	}
	if admin.Mail == "" {
		return nil, errors.New("没有设置邮箱, 无法发送验证码")
	}

	return &admin, nil
}
