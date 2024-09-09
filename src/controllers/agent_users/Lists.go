package agent_users

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentUsers) Lists(c *gin.Context) {
	var part string
	var cond string
	offset, limit := request.GetOffsets(c)
	if offset == 0 {
		cond = "limit 15"
	} else {
		temp := "limit %d,%d"
		cond = fmt.Sprintf(temp, limit, offset)
	}

	username := c.Query("username")
	if username != "" {
		part = part + " and username='" + username + "'"
	}

	status := c.Query("status")
	if status != "" {
		part = part + " and status='" + status + "'"
	}

	Id := c.Query("id")
	if Id != "" {
		part = part + " and top_id='" + Id + "'"
	}

	topName := c.Query("top_name")
	if topName != "" {
		part = part + " and top_name='" + topName + "'"
	}
	var strsAt string
	var strsEnd string
	var startInt int64
	var endInt int64
	created := c.Query("created")
	if created != "" {
		areas := strings.Split(created, " - ")
		strsAt = areas[0]
		strsEnd = areas[1]
		startInt = tools.GetTimeStampByString(strsAt + " 00:00:00")
		endInt = tools.GetTimeStampByString(strsEnd + " 23:59:59")
		//wform = "and DATE_FORMAT(finance_process_at,'%y-%m-%d %h:%i:%s')> DATE_FORMAT('" + strsAt + "','%y-%m-%d') and DATE_FORMAT(finance_process_at,'%y-%m-%d %h:%i:%s')< DATE_FORMAT('" + strsEnd + " 23:59:59','%y-%m-%d %h:%i:%s')"
		//dform = "and confirm_at>UNIX_TIMESTAMP('" + areas[0] + "') and confirm_at<UNIX_TIMESTAMP('" + areas[1] + " 23:59:59')"
	} else {
		todayDateStr := time.Now().Format("2006-01")
		strsAt = todayDateStr + "-1"
		strsEnd = time.Now().Format("2006-01-02")
		//wform = "and DATE_FORMAT(finance_process_at,'%y-%m-%d %h:%i:%s')> DATE_FORMAT('" + strsAt + "','%y-%m-%d')"
		//dfm = "and DATE_FORMAT(confirm_at,'%y-%m-%d %h:%i:%s')> DATE_FORMAT('" + strsAt + "','%y-%m-%d')"
		startInt = tools.GetTimeStampByString(strsAt + " 00:00:00")
		endInt = time.Now().Unix()
	}

	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	lsql := "select id,username,realname,phone,top_name,status,created from users where  top_name !='' and is_agent!=1 "

	res, err := dbSession.QueryString(lsql + part + cond)
	if err != nil {
		log.Err(err.Error())
		return
	}
	countSql := "select count(*) as total from users where top_name!='' and is_agent!=1 "
	countRes, counterr := dbSession.QueryString(countSql + part)
	if counterr != nil {
		log.Err(counterr.Error())
	}
	withSql := "select sum(money) as withdraws from user_withdraws where user_id=? and updated>= ? and  updated<? and status=2"
	depoSql := "select sum(arrive_money) as deposits from user_deposits where user_id=? and confirm_at >= ? and  confirm_at< ? and status=2"
	//withSql := "select sum(money) as withdraws from user_withdraws where  status=2 and  user_id='%s' "
	//depoSql := "select sum(arrive_money) as deposits from user_deposits where  status=2 and  user_id='%s' "
	acSql := "select balance from accounts where user_id='%s'"

	for _, v := range res {
		//获取总提款
		wRes, err := dbSession.QueryString(withSql, v["id"], tools.SecondToMicro(startInt), tools.SecondToMicro(endInt))
		if err != nil {
			log.Err(err.Error())
			return
		}

		if wRes[0]["withdraws"] == "" {
			v["withdraws"] = "0"
		} else {
			v["withdraws"] = wRes[0]["withdraws"]
		}

		//获取总存款
		dRes, derr := dbSession.QueryString(depoSql, v["id"], startInt, endInt)
		if derr != nil {
			log.Err(derr.Error())
			return
		}
		//代客充值
		//cusSql := "select sum(money) as money from user_account_sets where status=2 and user_id=" + v["id"]
		//cusRes, err := dbSession.QueryString(cusSql)
		//cusMoney, _ := strconv.ParseFloat(cusRes[0]["money"], 64)
		dMoney, _ := strconv.ParseFloat(dRes[0]["deposits"], 64)
		v["deposits"] = strconv.FormatFloat(dMoney, 'f', -1, 64)

		//获取中心钱包
		Asqll := fmt.Sprintf(acSql, v["id"])
		aRes, aerr := dbSession.QueryString(Asqll)
		if aerr != nil {
			log.Err(aerr.Error())
			return
		}
		if len(aRes) > 0 {
			v["balance"] = aRes[0]["balance"]
		} else {
			v["balance"] = "0"
		}

		//获取总输赢
		fsql := "select sum(net_money) as money from user_daily_reports where user_id=? and day>='" + strsAt + "' and day<='" + strsEnd + "' and game_code !='0' "
		fRes, err := dbSession.QueryString(fsql, v["id"])
		if err != nil {
			log.Err(err.Error())
		}
		if fRes[0]["money"] == "" {
			v["final"] = "0"
		} else {
			v["final"] = fRes[0]["money"]
		}
	}

	base_controller.SetLoginAdmin(c)
	if request.IsAjax(c) {
		response.Render(c, "agents/_users.html", pongo2.Context{"rows": res, "total": countRes[0]["total"]})
		return
	}
	response.Render(c, "agents/users.html", pongo2.Context{"rows": res, "total": countRes[0]["total"]})
}
