package agent_logs

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *AgentLogs) Deposits(c *gin.Context) {
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
	status := c.Query("status")
	if status != "" {
		part = part + " and status='" + status + "'"
	}
	topName := c.Query("top_name")
	if topName != "" {
		part = part + " and top_name='" + topName + "'"
	}
	Type := c.Query("channel_type")
	if Type != "" {
		part = part + " and channel_type='" + Type + "'"
	}
	created := c.Query("created")
	if created != "" {
		areas := strings.Split(created, " - ")

		part = part + "and confirm_at >UNIX_TIMESTAMP('" + areas[0] + " 00:00:00') and confirm_at <UNIX_TIMESTAMP('" + areas[1] + " 23:59:59')"

	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()
	order := " order  by confirm_at desc "
	sql := "select * from user_deposits  where status>1 "
	tsql := "select count(*) as total from user_deposits where status>1 "
	res, err := myClient.QueryString(sql + part + order + offset)
	if err != nil {
		log.Err(err.Error())
		return
	}
	totalRes, _ := myClient.QueryString(tsql + part)
	response.Render(c, "agent_logs/_deposits.html", pongo2.Context{"rows": res, "total": totalRes[0]["total"]})
}
