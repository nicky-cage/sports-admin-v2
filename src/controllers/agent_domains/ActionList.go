package agent_domains

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.AgentDomains,
	ViewFile: "agent_domains/list.html",
	Rows: func() interface{} {
		return &[]models.AgentDomain{}
	},
	QueryCond: map[string]interface{}{
		"user_id":  "=", // 用户编号
		"username": "%", // 用户名称
		"domain":   "%", // 域名
		"state":    "=", // 状态
	},
}
