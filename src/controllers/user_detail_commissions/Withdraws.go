package user_detail_commissions

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *UserDetailCommissions) Withdraws(c *gin.Context) {
	cond := builder.NewCond()
	id := c.Query("id")
	status := c.Query("status")
	if status != "" {
		cond = cond.And(builder.Eq{"status": status})
	}
	created := c.Query("created")
	if created != "" {
		areas := strings.Split(created, " - ")
		start := tools.GetMicroTimeStampByString(areas[0])
		end := tools.GetMicroTimeStampByString(areas[1])
		cond = cond.And(builder.Gte{"created": start}).And(builder.Lte{"created": end})
	}
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	var list []models.UserWithdraw

	//.And(builder.Eq{"type": 2})
	limit, offset := request.GetOffsets(c)
	cond = cond.And(builder.Eq{"user_id": id})
	to, _ := dbSession.Table("user_withdraws").Where(cond).Limit(limit, offset).OrderBy("created DESC").FindAndCount(&list)
	//sql := "select * from user_withdraws where user_id= " + id
	ssql := "select sum(money) as money from user_withdraws where user_id= " + id + " and status=2"

	sRes, err := dbSession.QueryString(ssql)
	if err != nil {
		log.Err(err.Error())
		return
	}
	total := sRes[0]["money"]

	if created != "" {
		//temp := map[string]interface{}{}
		//temp["rows"] = list
		//temp["id"] = id
		//temp["total"] = id
		//response.Result(c, temp)
		response.Render(c, "users/_commission_withdraws.html", pongo2.Context{"rows": list, "total": to})
		return
	}
	response.Render(c, "users/commission_withdraws.html", pongo2.Context{"rows": list, "sum": total, "total": to, "id": id})
}
