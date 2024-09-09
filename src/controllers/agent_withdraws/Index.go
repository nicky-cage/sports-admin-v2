package agent_withdraws

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *AgentWithdraws) Index(c *gin.Context) {
	cond := builder.NewCond()
	var startAt int64
	var endAt int64
	if value, exists := c.GetQuery("created"); !exists {
		currentTime := time.Now().Unix()
		startAt = currentTime - currentTime%86400
		endAt = startAt + 86400
		cond = cond.And(builder.Gte{"user_withdraws.created": tools.SecondToMicro(startAt)}).And(builder.Lte{"user_withdraws.created": tools.SecondToMicro(endAt)})
	} else {
		areas := strings.Split(value, " - ")
		startAt = tools.GetTimeStampByString(areas[0])
		endAt = tools.GetTimeStampByString(areas[1])
		cond = cond.And(builder.Gte{"user_withdraws.created": tools.SecondToMicro(startAt)}).And(builder.Lte{"user_withdraws.created": tools.SecondToMicro(endAt)})
	}
	cond = cond.And(builder.Eq{"user_withdraws.status": 1}).And(builder.Eq{"user_withdraws.type": 2}).And(builder.Eq{"user_withdraws.process_step": 3})
	if min, ok := c.GetQuery("money_min"); ok {
		cond = cond.And(builder.Gte{"user_withdraws.money": min})
	}
	if max, ok := c.GetQuery("money_max"); ok {
		cond = cond.And(builder.Lte{"user_withdraws.money": max})
	}
	username := c.DefaultQuery("username", "")
	billNo := c.DefaultQuery("bill_no", "")
	agentAdmin := c.DefaultQuery("agent_admin", "")
	status := c.DefaultQuery("status", "")
	if len(username) > 0 {
		cond = cond.And(builder.Eq{"user_withdraws.username": username})
	}
	if len(billNo) > 0 {
		cond = cond.And(builder.Eq{"user_withdraws.bill_no": billNo})
	}
	if len(agentAdmin) > 0 {
		cond = cond.And(builder.Eq{"user_withdraws.agent_admin": agentAdmin})
	}
	if len(status) > 0 {
		cond = cond.And(builder.Eq{"user_withdraws.status": status})
	}
	limit, offset := request.GetOffsets(c)
	agentWithdraws := make([]AgentWithdrawsStruct, 0)
	platform := request.GetPlatform(c)
	engine := common.Mysql(platform)
	defer engine.Close()

	//出款方式
	depositCards := make([]models.DepositCard, 0)
	if err := engine.Table("deposit_cards").Find(&depositCards); err != nil {
		log.Logger.Error(err.Error())
	}
	depositCardsNew := make([]string, 0)
	for _, v := range depositCards {
		depositCardsNew = append(depositCardsNew, v.Byname)
	}
	depositCardsNew = append(depositCardsNew, "shipu_daifu")

	total, err := engine.Table("user_withdraws").Join("LEFT OUTER", "users", "user_withdraws.user_id = users.id").Where(cond).OrderBy("user_withdraws.id DESC").Limit(limit, offset).FindAndCount(&agentWithdraws)
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "获取列表错误")
		return
	}
	ss := new(AgentWithdrawSumStruct)
	sumTotal, _ := engine.Table("user_withdraws").Where(cond).Sum(ss, "money")
	pageSumMoney := 0.00
	for _, v := range agentWithdraws {
		pageSumMoney += v.Money
	}
	viewData := pongo2.Context{
		"rows":           agentWithdraws,
		"payment_method": depositCardsNew,
		"total":          total,
		"page_sum_money": pageSumMoney,
		"sum_money":      sumTotal,
	}
	viewFile := "agent_withdraws/agent_withdraws.html"
	if request.IsAjax(c) {
		viewFile = "agent_withdraws/_agent_withdraws.html"
	}
	response.Render(c, viewFile, viewData)
}
