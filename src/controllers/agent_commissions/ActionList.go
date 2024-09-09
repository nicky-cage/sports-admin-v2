package agent_commissions

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model: models.AgentCommissions,
	Rows: func() interface{} {
		return &[]models.AgentCommission{}
	},
	QueryCond: map[string]interface{}{
		"money":       "=",
		"Date":        "=",
		"username":    "%",
		"userid":      "=",
		"beforemoney": "=",
		"aftermoney":  "=",
	},
	ViewFile: "agents/commissions.html",
}
