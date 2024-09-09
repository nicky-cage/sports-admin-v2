package controllers

import (
	"sports-admin/controllers/admin_authorizes"
	"sports-admin/controllers/admin_roles"
	"sports-admin/controllers/admins"
	models "sports-models"
)

// AdminRoles 角色管理
var AdminRoles = admin_roles.AdminRoles{
	ActionList:   admin_roles.ActionList,
	ActionCreate: admin_roles.ActionCreate,
	ActionUpdate: admin_roles.ActionUpdate,
	ActionSave:   admin_roles.ActionSave,
	ActionDetail: admin_roles.ActionDetail,
	ActionDelete: admin_roles.ActionDelete,
}

// Admins 系统账号管理
var Admins = admins.Admins{
	ActionList:   admins.ActionList,
	ActionCreate: admins.ActionCreate,
	ActionUpdate: admins.ActionUpdate,
	ActionSave:   admins.ActionSave,
	ActionState:  admins.ActionState,
	ActionDelete: &ActionDelete{Model: models.Admins},
}

// AdminAuthorizes 访问授权
var AdminAuthorizes = admin_authorizes.AdminAuthorizes{}
