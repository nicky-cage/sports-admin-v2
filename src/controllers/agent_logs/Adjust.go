package agent_logs

import (
	"fmt"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentLogs) Adjust(c *gin.Context) {
	var part string
	var offset string
	page := c.Query("page")
	total, _ := strconv.Atoi(page)
	if page == "1" || page == "" {
		offset = " limit 15"
		total = 1
	} else {
		start := (total-1)*15 + 1
		st := strconv.Itoa(start)
		offset = " limit " + st + "limit " + st + ", 15"
	}

	username := c.Query("username")
	if username != "" {
		part = part + " and  a.username='" + username + "'"
	}
	topName := c.Query("top_name")
	if topName != "" {
		part = part + "and   b.top_name='" + topName + "'"
	}
	vip := c.Query("vip")
	if vip != "" {
		part = part + " and  b.vip='" + vip + "'"
	}
	label := c.Query("label")
	if label != "" {
		arr := strings.Split(label, "|")
		part = part + " and  b.label like '%" + arr[0] + "%" + arr[1] + "%'"
	}

	created := c.Query("created")
	if created != "" {
		areas := strings.Split(created, " - ")
		startInt := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
		endInt := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")

		part = part + "and a.updated>" + fmt.Sprintf("%d", startInt) + " and a.updated<" + fmt.Sprintf("%d", endInt)

	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()
	order := " order by a.updated desc"
	sql := "select  a.*,b.top_name,b.vip,b.label from user_resets a join users b on a.user_id=b.id and a.adjust_type='3' where 1=1  "
	tsql := "select  count(*) as total ,sum(adjust_money) as money from user_resets a join users b on a.user_id=b.id and a.adjust_type='3'  "

	res, err := myClient.QueryString(sql + part + order + offset)
	if err != nil {
		log.Err(err.Error())
		return
	}

	for _, v := range res {
		temp, _ := strconv.Atoi(v["vip"])
		v["vip"] = strconv.Itoa(temp - 1)
	}
	totalRes, _ := myClient.QueryString(tsql + part)

	response.Render(c, "agent_logs/_adjust.html", pongo2.Context{"rows": res, "total": totalRes[0]["total"], "adjust_money": totalRes[0]["money"]})
}
