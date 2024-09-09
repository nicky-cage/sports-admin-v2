package agent_withdraw_hrs

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *AgentWithdrawHrs) List(c *gin.Context) {
	cond := builder.NewCond()
	timeStart, timeEnd := request.GetMicroTimesByQuery(c, "created") // 开始时间
	cond = cond.And(builder.Gte{"user_withdraws.created": timeStart}).And(builder.Lte{"user_withdraws.created": timeEnd})
	updatedStart, updatedEnd := request.GetMicroTimesByQuery(c, "updated") // 结束时间
	cond = cond.And(builder.Gte{"user_withdraws.updated": updatedStart}).And(builder.Lte{"user_withdraws.updated": updatedEnd})
	cond = cond.And(builder.Neq{"user_withdraws.status": 1}).And(builder.Eq{"user_withdraws.type": 2}) // 2 表示是代理提款
	request.QueryCond(c, &cond, map[string]map[string]string{
		">=": {
			"money_min": "user_withdraws.money",
		},
		"<=": {
			"money_max": "user_withdraws.money",
		},
		"%": {
			"username":       "user_withdraws.username",
			"bill_no":        "user_withdraws.bill_no",
			"risk_admin":     "user_withdraws.risk_admin",
			"finance_name":   "user_withdraws.finance_admin",
			"payment_method": "user_withdraws.payment_method",
			"business_type":  "user_withdraws.business_type",
		},
		"=": {
			"status": "user_withdraws.status",
		},
	})
	limit, offset := request.GetOffsets(c)
	agentWithdrawHrs := make([]AgentWithdrawHrsStruct, 0)
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	dbSession.Table("user_withdraws").Join("LEFT OUTER", "users", "user_withdraws.user_id = users.id").Where(cond).OrderBy("user_withdraws.id DESC")

	if request.IsExportExcel(c) { // 如果是导出报表
		err := dbSession.Find(&agentWithdrawHrs)
		if err != nil {
			response.Err(c, "导出代理提款数据有误")
			return
		}
		response.Result(c, agentWithdrawHrs)
		return
	}

	total, err := dbSession.Limit(limit, offset).FindAndCount(&agentWithdrawHrs)
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "获取列表错误")
		return
	}
	ss := new(AgentWithdrawHrSumStruct)
	sumTotal, _ := dbSession.Table("user_withdraws").Where(cond).Sum(ss, "money")
	pageSumMoney := 0.00
	for _, v := range agentWithdrawHrs {
		pageSumMoney += v.Money
	}
	viewData := pongo2.Context{
		"rows":           agentWithdrawHrs,
		"total":          total,
		"page_sum_money": pageSumMoney,
		"sum_money":      sumTotal,
	}
	viewFile := request.GetViewFile(c, "agent_withdraws/%shistory_records.html")
	response.Render(c, viewFile, viewData)
}
