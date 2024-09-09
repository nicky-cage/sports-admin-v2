package controllers

import (
	"sports-admin/controllers/user_codes"
	"sports-admin/controllers/user_detail_accounts"
	"sports-admin/controllers/user_detail_commissions"
	"sports-admin/controllers/user_level_changes"
	"sports-admin/controllers/user_levels"
	"sports-admin/controllers/user_login_logs"
	"sports-admin/controllers/user_logs"
	"sports-admin/controllers/user_notes"
	"sports-admin/controllers/user_tags"
	"sports-admin/controllers/users"
)

// Users 默认页
var Users = users.Users{
	ActionList:   users.ActionList,
	ActionCreate: users.ActionCreate,
	ActionUpdate: users.ActionUpdate,
	ActionSave:   users.ActionSave,
	ActionState:  users.ActionSate,
	ActionExport: users.ActionExport,
}

// UserDetailAccounts 用户详情 - 账户相关
var UserDetailAccounts = user_detail_accounts.UserDetailAccounts{}

// UserDetailCommissions 用户详情 - 佣金部分
var UserDetailCommissions = user_detail_commissions.UserDetailCommissions{}

// UserNotes 用户备注
var UserNotes = user_notes.UserNotes{
	ActionCreate: user_notes.ActionCreate,
	ActionUpdate: user_notes.ActionUpdate,
	ActionSave:   user_notes.ActionSave,
}

// UserTags 用户标签
var UserTags = user_tags.UserTags{
	ActionList:   user_tags.ActionList,
	ActionCreate: user_tags.ActionCreate,
	ActionUpdate: user_tags.ActionUpdate,
	ActionSave:   user_tags.ActionSave,
	ActionDelete: user_tags.ActionDelete,
}

// UserLevels 用户级别
var UserLevels = user_levels.UserLevels{
	ActionList:   user_levels.ActionList,
	ActionUpdate: user_levels.ActionUpdate,
	ActionSave:   user_levels.ActionSave,
}

// UserLoginLogs 会员等级
//elastic search列表展示公用函数，处理row的返回值以及List方法内的断言就可以
var UserLoginLogs = user_login_logs.UserLoginLogs{
	ActionEsList: user_login_logs.ActionEsList,
}

// UserLogs 默认页
var UserLogs = user_logs.UserLogs{}

// UserCodes 用户验证码
var UserCodes = user_codes.UserCodes{}

// UserLevelChanges 用户等级记录
var UserLevelChanges = user_level_changes.UserLevelChanges{
	ActionList: user_level_changes.ActionList,
}
