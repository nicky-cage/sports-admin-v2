package user_detail_commissions

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *UserDetailCommissions) Records(c *gin.Context) {
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	var list []models.AgentCommissionLog
	cond := builder.NewCond()
	month := c.DefaultQuery("month", "")
	if len(month) > 0 {
		cond = cond.And(builder.Eq{"month": month})
	}
	cond = cond.And(builder.Eq{"user_id": id})
	cond = cond.And(builder.Eq{"status": 2})
	total, err := dbSession.Table("agent_commission_logs").Where(cond).FindAndCount(&list)
	if err != nil {
		log.Err(err.Error())
		return
	}

	if month == "" {
		response.Render(c, "users/commission_records.html", pongo2.Context{"rows": list, "total": total, "id": id})
		return
	}
	response.Render(c, "users/_commission_records.html", pongo2.Context{"rows": list, "total": total, "id": id})
}
