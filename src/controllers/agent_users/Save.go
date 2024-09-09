package agent_users

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/pgsql"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"time"

	"github.com/gin-gonic/gin"
)

func (ths AgentUsers) Save(c *gin.Context) {
	//只有一個功能 就是轉代，記錄轉貸時間，
	dataPost := request.GetPostedData(c)
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	topName := dataPost["top_name"].(string)
	remark := dataPost["remark"].(string)
	before := dataPost["before_agent"].(string)
	username := dataPost["username"].(string)
	id := dataPost["id"].(string)
	transfer := time.Now().Unix()
	admin := base_controller.GetLoginAdmin(c)
	usql := "select id,is_agent from users where username='" + topName + "'"
	uRes, _ := dbSession.QueryString(usql)
	if len(uRes) == 0 {
		response.Err(c, "该代理不存在")
		return
	}
	if uRes[0]["is_agent"] != "1" {
		response.Err(c, "转代后所属必须是代理")
		return
	}
	pgGroup := pgsql.GetConnGroup(platform)
	defer pgGroup.Close()
	//获取转代时的钱包额度
	asql := "select available from accounts where user_id=" + id
	aRes, _ := dbSession.QueryString(asql)
	if username == topName {
		response.Err(c, "用户信息不完整")
		return
	}

	//加个判断，所转代理    不能为子级。
	lowerSql := "select username from users where top_name='" + username + "'"
	lRes, _ := dbSession.QueryString(lowerSql)
	for _, v := range lRes {
		if v["username"] == topName {
			response.Err(c, "不能将下线会员做为上级")
			return
		}
	}

	sql := "update users set trans_agent_money=?,trans_agent_admin=?,transform_agent=?,trans_before_agent=?,top_name=?,top_id=?,remark=? where id=?"
	_, err := dbSession.Exec(sql, aRes[0]["available"], admin.Name, transfer, before, topName, uRes[0]["id"], remark, id)
	if err != nil {
		log.Err(err.Error())
		response.Err(c, "系统错误")
	}
	//更改本月红利。 反水，存款，提款，月  涉及到top_id的、 输赢调整。

	d := time.Now().Format("2006-01-02")
	start := tools.GetTimeStampByString(d + " 00:00:00")

	//	mothTemps := time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, d.Location())
	//monthTime := mothTemps.Unix()
	//monthDay := mothTemps.Format("2006-01-02")
	//now := time.Now().Unix()
	//dividendSql := "update user_dividends set top_id=? where user_id=? and created>=? and created<?"
	//_, err1 := dbSession.Exec(dividendSql, uRes[0]["id"], id, monthTime, now)
	//if err1 != nil {
	//	log.Err(err1.Error())
	//}
	//rebateSql := "update user_rebate_records set top_id=? where user_id=? and created>=? and created<?"
	//_, err2 := dbSession.Exec(rebateSql, uRes[0]["id"], id, monthTime, now)
	//if err2 != nil {
	//	log.Err(err2.Error())
	//}
	//
	//depositsSql := "update user_deposits set top_id=? where user_id=? and updated>=? and updated<? and status=2"
	//
	//_, err3 := dbSession.Exec(depositsSql, uRes[0]["id"], id, monthTime, now)
	//if err3 != nil {
	//	log.Err(err3.Error())
	//}
	//withdrawSql := "update user_withdraws set top_id=? where user_id=? and updated>=? and updated<? and status=2"
	//_, err4 := dbSession.Exec(withdrawSql, uRes[0]["id"], id, monthTime, now)
	//if err4 != nil {
	//	log.Err(err4.Error())
	//}
	//resetSql := "update user_resets set top_id=? where user_id=? and updated>=? and updated<? and status=2"
	//_, err5 := dbSession.Exec(resetSql, uRes[0]["id"], id, monthTime, now)
	//if err5 != nil {
	//	log.Err(err5.Error())
	//}
	////月报表
	reportSql := "update user_daily_reports set top_id=?,top_name=? where user_id=? and day>= '" + d + "'"
	_, err6 := dbSession.Exec(reportSql, uRes[0]["id"], topName, id)
	if err6 != nil {
		log.Err(err6.Error())
	}
	//更新注单数据
	pgUpdateSql := "update wager_records set top_id=%s ,top_name='%s' where user_id=%s and created_at>=%d "
	pgUpdateSql = fmt.Sprintf(pgUpdateSql, uRes[0]["id"], topName, id, start)

	for _, v := range pgGroup.Conns {
		res, err := v.Exec(pgUpdateSql)
		if err != nil {
			log.Logger.Error("pgsql:", err)
		}
		if res.RowsAffected() == 0 {
			log.Logger.Error("pgsql:" + pgUpdateSql)
		}
	}
	//老代理更新 多条
	key := consts.CtxKeyLoginUser + id
	conn := common.Redis(platform)
	defer common.RedisRestore(platform, conn)
	conn.Del(key)
	response.Ok(c)
}
