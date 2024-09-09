package validations

import "sports-common/validation"

// AdminRoles 管理角色
var AdminRoles = func(data map[string]interface{}) error {
	return validation.New(data).
		Field("name").Length(2, 20, "必须输入角色名称").
		Field("menu_ids").Length(1, 500, "角色不能没有权限菜单").
		Validate()
}
