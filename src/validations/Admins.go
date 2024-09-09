package validations

import (
	"sports-common/validation"
	"strconv"
)

// Admins 后台管理员
var Admins = func(platform string, data map[string]interface{}) error {
	isCreate := true
	if idStr, exists := data["id"]; exists {
		if id, err := strconv.Atoi(idStr.(string)); err == nil && id > 0 {
			isCreate = false
		}
	}

	validator := validation.New(data).
		Field("nickname").Length(3, 20, "输入昵称格式有误").
		Field("mail").Mail().
		Field("role_id").Uint("角色选择有误")

	// 只有添加的情况下, 才有用户名称、密码选项
	if isCreate {
		validator.Field("name").UserName("用户名称格式有误").
			Field("password").Password().
			Field("re_password").Password().Equal("password")
	}

	return validator.Validate()
}
