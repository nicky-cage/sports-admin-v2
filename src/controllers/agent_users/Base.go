package agent_users

import "sports-admin/controllers/base_controller"

// SumAgent 统计
type SumAgent struct {
	Value float64 `json:"value"`
}

// AgentUsers 代理用户
type AgentUsers struct {
	*base_controller.ActionList
	*base_controller.ActionDetail
	*base_controller.ActionUpdate
}
