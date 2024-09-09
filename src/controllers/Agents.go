package controllers

import (
	"sports-admin/controllers/agent_commission_logs"
	"sports-admin/controllers/agent_commission_plans"
	"sports-admin/controllers/agent_commissions"
	"sports-admin/controllers/agent_domains"
	"sports-admin/controllers/agent_logs"
	"sports-admin/controllers/agent_users"
	"sports-admin/controllers/agent_withdraw_hrs"
	"sports-admin/controllers/agent_withdraw_logs"
	"sports-admin/controllers/agent_withdraw_records"
	"sports-admin/controllers/agent_withdraws"
	"sports-admin/controllers/agents"
)

type GameRecord struct {
	ID          string                 `json:"id"`
	Username    string                 `json:"username"`
	Playname    string                 `json:"playname"`
	UserID      int                    `json:"user_id"`
	GameCode    string                 `json:"game_code"`
	GameType    int                    `json:"game_type"`
	BillNo      string                 `json:"bill_no"`
	BetMoney    int                    `json:"bet_money"`
	ReturnMoney int                    `json:"return_money"`
	NetMoney    int                    `json:"net_money"`
	ValidMoney  int                    `json:"valid_money"`
	Status      int                    `json:"status"`
	CreatedAt   uint32                 `json:"created_at"`
	UpdatedAt   int                    `json:"updated_at"`
	ExtendStr   string                 `json:"extend_str"`
	TopName     string                 `json:"top_name"`
	Extend      map[string]interface{} `json:"extend"`
}

// AgentLogs 代理记录
var AgentLogs = agent_logs.AgentLogs{
	ActionEsList: agent_logs.ActionEsList,
}

// Agents 代理
var Agents = agents.Agents{
	ActionList:   agents.ActionList,
	ActionUpdate: agents.ActionUpdate,
	ActionSave:   agents.ActionSave,
	ActionCreate: agents.ActionCreate,
}

// AgentCommissions 代理佣金
var AgentCommissions = agent_commissions.AgentCommissions{
	ActionList:   agent_commissions.ActionList,
	ActionCreate: agent_commissions.ActionCreate,
	ActionSave:   agent_commissions.ActionSave,
}

// AgentCommissionsLogs 代理佣金日志
var AgentCommissionsLogs = agent_commission_logs.AgentCommissionsLogs{
	ActionList: agent_commission_logs.ActionList,
}

// AgentCommissionsPlan 代理佣金计划
var AgentCommissionsPlan = agent_commission_plans.AgentCommissionPlans{
	ActionCreate: agent_commission_plans.ActionCreate,
}

// AgentUsers 代理用户
var AgentUsers = agent_users.AgentUsers{
	ActionList:   agent_users.ActionList,
	ActionDetail: agent_users.ActionDetail,
	ActionUpdate: agent_users.ActionUpdate,
}

// AgentWithdraws 代理提款
var AgentWithdraws = agent_withdraws.AgentWithdraws{
	ActionList: agent_withdraws.ActionList,
}

// AgentWithdrawsRecords 代理提款记录
var AgentWithdrawsRecords = agent_withdraw_records.AgentWithdrawsRecords{
	ActionList: agent_withdraw_records.ActionList,
}

// AgentWithdrawHrs 代理提款管理-历史记录
var AgentWithdrawHrs = agent_withdraw_hrs.AgentWithdrawHrs{
	ActionExport: agent_withdraw_hrs.ActionExport,
}

// AgentWithdrawLogs 代理提款管理-日志记录
var AgentWithdrawLogs = agent_withdraw_logs.AgentWithdrawLogs{
	ActionList: agent_withdraw_logs.ActionList,
}

// AgentDomains 代理域名绑定
var AgentDomains = agent_domains.AgentDomains{
	ActionList:   agent_domains.ActionList,
	ActionCreate: agent_domains.ActionCreate,
	ActionUpdate: agent_domains.ActionUpdate,
	ActionSave:   agent_domains.ActionSave,
	ActionDelete: agent_domains.ActionDelete,
	ActionState:  agent_domains.ActionState,
}

// func GetVenueId(gameCode string, gameType int, arr []models.GameVenue) int {
// 	for _, v := range arr {
// 		if gameCode == v.Code && gameType == int(v.VenueType) {
// 			return int(v.Id)
// 		}
// 	}
// 	return 0
// }
//
// func GetGameRate(vip int, typeId int, venueId int, arr []models.UserRebateSetting) float64 {
// 	//rateTemp[int(v.VenueId)] = v.Ratio
// 	//rates[int(v.TypeId)] = rateTemp
// 	//rate[int(v.VipId)] = rates
// 	for _, v := range arr {
// 		if vip == int(v.VipId) && typeId == int(v.TypeId) && venueId == int(v.VenueId) {
// 			return v.Ratio
// 		}
// 	}
// 	return 0.0
// }
