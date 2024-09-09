package validations

import (
	"sports-common/validation"
)

// UserPasswordSave 保存密码校验
var UserPasswordSave = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("id").Uint("必须提供用户编号").
		Field("password").Password().
		Field("re_password").Password().Equal("password", "2次密码必须一致").
		Validate()
}

// UserCreateSave 创建之前判断
var UserCreateSave = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("password").Password().Length(6, 20, "密码位数必须大于6小于20").
		Field("re_password").Password().Equal("password", "确认密码必须与密码一致").
		Validate()

}

// Users 保存用户信息
var Users = func(data map[string]interface{}) error {
	validator := validation.New(data).
		Field("name").UserName("请输入正确的用户名称").
		Field("phone").Mobile("请输入正确格工的手机号码").
		Field("vip").Uint("VIP值必须是有效的数字").
		Field("email").Mail().
		Field("birthday").Date("").
		Field("realname").Length(2, 10, "请输入真实姓名").
		Field("nickname").Length(5, 20, "请输入会员昵称").
		Field("gender").InValues([]string{"1", "2", "0"}, "性别输入有误")

	// 如果有密码提交, 则需要校验密码
	if idStr, exists := data["id"]; exists && idStr == "0" { // 表示是添加用户信息
		validator.Field("password").Password().
			Field("re_password").Password().Equal("password")
	}

	return validator.Validate()
}
