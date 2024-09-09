package agent_withdraw_hrs

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

// AgentWithdrawHrsStruct 代理提现
type AgentWithdrawHrsStruct struct {
	models.UserWithdraw `xorm:"extends"`
	Vip                 int32
}

// AgentWithdrawHrSumStruct 代理记录汇总
type AgentWithdrawHrSumStruct struct {
	Money float64
}

type AgentWithdrawHrs struct {
	*base_controller.ActionExport
}
