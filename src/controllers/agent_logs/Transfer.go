package agent_logs

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *AgentLogs) Transfer(c *gin.Context) {
	var part string
	var offset string
	page := c.Query("page")
	total, _ := strconv.Atoi(page)
	if page == "1" || page == "" {
		offset = "limit 15"
		total = 1
	} else {
		start := (total-1)*15 + 1
		st := strconv.Itoa(start)

		offset = "limit " + st + ", 15"
	}

	username := c.Query("username")
	if username != "" {
		part = part + " and username='" + username + "'"
	}
	topName := c.Query("top_name")
	if topName != "" {
		part = part + " and top_name='" + topName + "'"
	}

	BetopName := c.Query("trans_before_agent")
	if BetopName != "" {
		part = part + " and trans_before_agent='" + BetopName + "'"
	}
	created := c.Query("created")
	if created != "" {
		areas := strings.Split(created, " - ")
		part = part + "and transform_agent>UNIX_TIMESTAMP('" + areas[0] + "') and transform_agent<UNIX_TIMESTAMP('" + areas[1] + " 23:59:59')"
	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()
	order := " order by transform_agent DESC "
	sql := "select id, transform_agent,username,top_name,trans_before_agent,trans_agent_money,trans_agent_admin,remark from users where  transform_agent!=0 and is_agent!=1 "
	tsql := "select count(*) as total from users where  transform_agent!=0 and is_agent!=1 "
	res, err := myClient.QueryString(sql + part + order + offset)
	if err != nil {
		log.Err(err.Error())
		return
	}
	totalRes, _ := myClient.QueryString(tsql + part)
	sumw := new(models.UserWithdraw)
	sumd := new(models.UserDeposit)
	sumr := new(models.UserDailyReport)
	for _, v := range res {
		trans, _ := strconv.Atoi(v["transform_agent"])
		tm := time.Unix(int64(trans), 0)
		timeStr := tm.Format("2006-01-02 15:04:05")
		//存款
		d, err := myClient.Where("user_id=?", v["id"]).And(builder.Gt{"confirm_at": timeStr}).Sum(sumd, "money")
		if err != nil {
			log.Err(err.Error())
			return
		}
		v["deposits"] = strconv.FormatFloat(d, 'f', -1, 64)
		//存款
		s, err := myClient.Where("user_id=?", v["id"]).And(builder.Gt{"finance_process_at": timeStr}).Sum(sumw, "money")
		if err != nil {
			log.Err(err.Error())
			return
		}
		v["withdraws"] = strconv.FormatFloat(s, 'f', -1, 64)

		//有效投注， 输赢，
		win, err := myClient.Where("user_id=?", v["id"]).And(builder.Gt{"day": timeStr}).Sums(sumr, "valid_money", "net_money")
		if err != nil {
			log.Err(err.Error())
		}

		v["valid_money"] = strconv.FormatFloat(win[0], 'f', -1, 64)
		v["net_money"] = strconv.FormatFloat(win[1], 'f', -1, 64)
	}

	response.Render(c, "agent_logs/_transfer_agent.html", pongo2.Context{"rows": res, "total": totalRes[0]["total"]})
}
