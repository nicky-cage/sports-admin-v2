package user_relations

import "sports-admin/controllers"

// 定义控制器
type userRelations struct {
	*controllers.ActionList
}

// UserRelations 控制器名称
var UserRelations = userRelations{}
