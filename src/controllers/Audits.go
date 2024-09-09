package controllers

import (
	"fmt"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var Audits = struct {
	Deposit  func(c *gin.Context)
	Dividend func(c *gin.Context)
	*ActionCreate
	Save func(c *gin.Context)
}{
	Deposit: func(c *gin.Context) {
		id := c.Query("id")
		offset, limit := request.GetOffsets(c)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from  audits where user_id= " + id + " and type=1 order by created limit %d,%d"
		csql := "select count(*) as total ,sum(deposit_money) as money from audits where user_id= " + id + " and type=1 "
		sqll := fmt.Sprintf(sql, limit, offset)
		res, err := dbSession.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
		}
		cRes, _ := dbSession.QueryString(csql)
		//有效投注， 跟状态。 总有效投注，总状态
		response.Render(c, "users/audits.html", pongo2.Context{"id": id, "total": cRes[0]["total"], "deposit_money": cRes[0]["money"], "rows": res})
	},
	Dividend: func(c *gin.Context) {
		id := c.Query("id")
		offset, limit := request.GetOffsets(c)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from  audits where user_id= " + id + " and type=2 order by created limit %d,%d"
		csql := "select count(*) as total ,sum(dividend_money) as money from audits where user_id= " + id + " and type=2 "
		sqll := fmt.Sprintf(sql, limit, offset)
		res, err := dbSession.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
		}
		cRes, _ := dbSession.QueryString(csql)
		// 有效投注， 跟状态。 总有效投注，总状态,所需流水
		response.Render(c, "users/audit_dividend.html", pongo2.Context{"id": id, "total": cRes[0]["total"], "deposit_money": cRes[0]["money"], "rows": res})
	},
	Save: func(c *gin.Context) {
		postData := request.GetPostedData(c)
		if postData["status"] == nil {
			response.Err(c, "请选择调整状态")
			return
		}
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		csql := "select id from  audits where user_id=" + postData["id"].(string) + " and is_finish=1"
		res, err := dbSession.QueryString(csql)
		if err != nil {
			log.Err(err.Error())
		}
		var str string
		for _, v := range res {
			str = str + "<p>序号" + v["id"] + ",不通过->通过</p>"
		}

		sql := "update audits set status = 2 where user_id=" + postData["id"].(string)
		_, err = dbSession.Exec(sql)
		if err != nil {
			response.Err(c, "系统错误")
			return
		}

		admin := GetLoginAdmin(c)
		logSql := " insert into audit_logs(user_id,remark,operation_info,admin,created,type) values(?,?,?,?,?,?)"
		_, err = dbSession.Exec(logSql, postData["id"].(string), postData["remark"].(string), str, admin.Name, time.Now().Unix(), postData["type"].(string))

		if err != nil {
			response.Err(c, "系统错误")
			return
		}

		response.Ok(c)
	},
	ActionCreate: &ActionCreate{
		ViewFile: "users/audit_created.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			id := c.Query("id")
			return pongo2.Context{"id": id}
		},
	},
}
