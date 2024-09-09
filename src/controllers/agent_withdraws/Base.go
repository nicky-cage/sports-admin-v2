package agent_withdraws

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

// AgentWithdraws 代理提款
type AgentWithdraws struct {
	*base_controller.ActionList
}

type AgentWithdrawsStruct struct {
	models.UserWithdraw `xorm:"extends"`
	Vip                 int32
}

type AgentWithdrawSumStruct struct {
	Money float64
}
