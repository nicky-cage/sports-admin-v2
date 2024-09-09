package agent_commissions

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionCreate = &base_controller.ActionCreate{
	ViewFile: "agents/commissions_adjustment.html",
	ExtendData: func(c *gin.Context) pongo2.Context {
		platfrom := request.GetPlatform(c)
		id := c.Query("id")
		db := common.Mysql(platfrom)
		defer db.Close()
		//year, months, _ := time.Now().Date()
		//thisMonth := time.Date(year, months, 1, 0, 0, 0, 0, time.Local)
		//lastMonth := thisMonth.AddDate(0, -1, 0).Format("2006-01")
		lastMonth := time.Now().Format("2006-01")
		sql := "select a.id,a.user_id,a.month,a.commission_adjust,a.money,b.username from agent_commission_logs a join users b on a.user_id=b.id where a.user_id =? and a.month=?"
		res, err := db.QueryString(sql, id, lastMonth)
		if err != nil {
			log.Err(err.Error())
			return nil
		}
		admin := base_controller.GetLoginAdmin(c)
		money := c.Query("money")
		data := pongo2.Context{"r": res[0], "admin": admin.Name, "money": money}
		return data
	},
}
