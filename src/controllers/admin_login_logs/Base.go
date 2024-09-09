package admin_login_logs

import (
	"sports-admin/controllers"
)

// 定义控制器
type adminLoginLogs struct {
	*controllers.ActionList
}

// AdminLoginLogs 控制器名称
var AdminLoginLogs = adminLoginLogs{
	ActionList: actionList,
}
