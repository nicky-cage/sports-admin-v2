package validations

import (
	"sports-common/validation"
)

// UserLogin 用户登录
var UserLogin = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("username").UserName("输入的用户名称格式有误").
		Field("password").Password("输入的密码密码格式有误").
		Validate()
}

// UpdatePassword 检测修改后台用户密码
var UpdatePassword = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("password").Password("输入密码格式有误").
		Field("password_new").Password("输入新的密码格式有误").
		Field("password_rep").Password("输入重复密码格式有误").
		Validate()
}

// UpdateProfile 检测修改管理员资料
var UpdateProfile = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("nickname").Length(2, 20, "用户昵称格式输入有误").
		Field("mail").Mail().
		Validate()
}
