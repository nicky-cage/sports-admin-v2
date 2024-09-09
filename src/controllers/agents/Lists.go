package agents

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *Agents) Lists(c *gin.Context) {
	var part string
	var offset string

	limit, size := request.GetOffsets(c)
	if size != 0 {
		temp := "limit %d,%d"
		offset = fmt.Sprintf(temp, size, limit)
	} else {
		offset = "limit 15"
	}

	username := c.Query("username")
	if username != "" {
		part = part + " and username='" + username + "'"
	}
	topName := c.Query("top_name")
	if topName != "" {
		part = part + " and top_name='" + topName + "'"
	}
	agentCommission := c.Query("agent_commission")
	if agentCommission != "" {
		part = part + " and agent_commission='" + agentCommission + "'"
	}

	status := c.Query("id")
	if status != "" {
		part = part + " and id='" + status + "'"
	}
	topId := c.Query("agent_type")
	if topId != "" {
		part = part + " and agent_type='" + topId + "'"
	}

	var start string
	var end string
	var startInt int64
	var endInt int64
	todayDateStr := time.Now().Format("2006-01")
	start = todayDateStr + "-01"
	tempInt, _ := time.Parse("2006-01-02", start)
	startInt = tempInt.Unix()
	endInt = time.Now().Unix()
	end = time.Now().Format("2006-01-02 15:04:05")
	created := c.Query("created")
	if created != "" {
		areas := strings.Split(created, " - ")
		startInt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
		endInt = tools.GetTimeStampByString(areas[1] + " 23:59:59")

	}

	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select id ,username,realname,top_name,agent_type,transform_agent,status,agent_commission from users where is_agent != 0"
	res, err := dbSession.QueryString(sql + part + " ORDER BY id DESC " + offset)
	if err != nil {
		log.Err(err.Error())
		return
	}
	countSql := "select count(*) as total from users where is_agent!=0 "
	countRes, counterr := dbSession.QueryString(countSql + part)
	if counterr != nil {
		log.Err(counterr.Error())
	}
	vSql := "select id  from users where top_id='%s'"
	sumw := new(models.UserWithdraw)
	sumd := new(models.UserDeposit)
	//sumcus := new(models.UserAccountSet)
	sumr := new(models.UserDailyReport)

	for _, v := range res {
		//var temp []string
		vSqll := fmt.Sprintf(vSql, v["id"])
		lRes, err := dbSession.QueryString(vSqll)
		if err != nil {
			log.Err(err.Error())
			return
		}

		// 统计取款
		s, err := dbSession.Table("user_withdraws").Where("top_id=?", v["id"]).And(builder.Gte{"finance_process_at": start}).And(builder.Lt{"finance_process_at": end}).And(builder.Eq{"status": 2}).Sum(sumw, "money")
		if err != nil {
			log.Err(err.Error())
			return
		}
		v["withdraws"] = strconv.FormatFloat(s, 'f', -1, 64)

		//统 计存款

		d, err := dbSession.Table("user_deposits").Where("top_id=?", v["id"]).And(builder.Gte{"confirm_at": startInt}).And(builder.Lt{"confirm_at": endInt}).And(builder.Eq{"status": 2}).Sum(sumd, "arrive_money")
		if err != nil {
			log.Err(err.Error())
			return
		}
		v["deposits"] = strconv.FormatFloat(d, 'f', -1, 64)

		//下线会员
		v["low_num"] = strconv.Itoa(len(lRes))
		//活跃人数
		var activeNum []models.UserDailyReport
		dbSession.Table("user_daily_reports").Where("top_id=?", v["id"]).And(builder.Gte{"day": start}.And(builder.Lt{"day": end})).GroupBy("user_id ").Find(&activeNum)
		v["active_num"] = strconv.Itoa(len(activeNum))
		//总输赢

		win, err := dbSession.Table("user_daily_reports").Where("top_id=? and game_code>'0'", v["id"]).And(builder.Gte{"day": start}).And(builder.Lt{"day": end}).Sum(sumr, "net_money")
		if err != nil {
			log.Err(err.Error())
		}
		v["final"] = strconv.FormatFloat(win*-1, 'f', -1, 64)

	}

	commissionSql := "select agent_commission from  agent_commission_plans group by agent_commission"
	comRes, err := dbSession.QueryString(commissionSql)
	if err != nil {
		log.Err(err.Error())
	}
	base_controller.SetLoginAdmin(c)
	viewData := pongo2.Context{"rows": res, "total": countRes[0]["total"], "agent_commission": comRes}
	if request.IsAjax(c) {
		response.Render(c, "agents/_index.html", viewData)
	} else {
		response.Render(c, "agents/index.html", viewData)
	}
}
