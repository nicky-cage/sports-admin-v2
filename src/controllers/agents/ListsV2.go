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

	"github.com/gin-gonic/gin"
)

// ListsV2 代理列表
func (ths *Agents) ListsV2(c *gin.Context) {
	// 分页信息
	limit, size := request.GetOffsets(c)
	offset := fmt.Sprintf("LIMIT %d, %d", size, limit)

	// 处理查询条件
	cond := ""
	if username := c.Query("username"); username != "" {
		cond += " AND username = '" + username + "'"
	}
	if topName := c.Query("top_name"); topName != "" {
		cond += " AND top_name='" + topName + "'"
	}
	if agentCommission := c.Query("agent_commission"); agentCommission != "" {
		cond += " AND agent_commission = '" + agentCommission + "'"
	}
	if status := c.Query("id"); status != "" {
		cond += " AND  id='" + status + "'"
	}
	if topId := c.Query("agent_type"); topId != "" {
		cond += " AND agent_type='" + topId + "'"
	}

	todayTime := tools.CurrentTime()
	todayDateStr := todayTime.Format("2006-01")   // 今天时间
	start := todayDateStr + "-01"                 // 本月1号
	tempInt, _ := time.Parse("2006-01-02", start) // 本月1号 timestamp
	startInt := tempInt.Unix()                    // 同上
	endInt := todayTime.Unix()                    // 当前时间
	end := todayTime.Format("2006-01-02 15:04:05")
	if created := c.Query("created"); created != "" {
		areas := strings.Split(created, " - ")
		startInt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
		endInt = tools.GetTimeStampByString(areas[1] + " 23:59:59")
	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()

	// 先取出用户表相关信息
	listSQL := "SELECT id, username, realname, top_name, agent_type, transform_agent, status, agent_commission " +
		"FROM users WHERE is_agent != 0 "
	agents, err := myClient.QueryString(listSQL + cond + " ORDER BY id DESC " + offset)
	if err != nil {
		log.Err("获取全部代理信息有误:", err)
		return
	}
	// 总的记录数量统计
	totalRow := models.TotalInt{}
	if _, countErr := myClient.SQL("SELECT count(*) AS total FROM users WHERE is_agent != 0 " + cond).Get(&totalRow); countErr != nil {
		log.Err("获取统计信息有误:", countErr)
	}
	// 获取用户id信息
	ids := func() string {
		rArr := []string{}
		for _, r := range agents {
			rArr = append(rArr, r["id"])
		}
		return strings.Join(rArr, ",")
	}()

	countTypes := struct {
		Reg      string
		Withdraw string
		Deposit  string
		Active   string
		Final    string
		Domain   string
	}{
		Reg:      "reg",
		Withdraw: "withdraw",
		Deposit:  "deposit",
		Active:   "active",
		Final:    "final",
		Domain:   "domain",
	}
	rows := []struct {
		TopId   int     `json:"top_id"`
		Total   float64 `json:"total"`
		Type    string  `json:"type"`
		Content string  `json:"content"`
	}{}
	unionSQL := "(" + strings.Join([]string{
		// 下线会员总数
		fmt.Sprintf("SELECT '' AS content, COUNT(id) AS total, 'reg' AS type, top_id "+
			"FROM users WHERE top_id IN (%s) "+
			"GROUP BY top_id", ids),
		// 提款总数
		fmt.Sprintf("SELECT '' AS CONTENT, SUM(money) AS total, 'withdraw' AS type, top_id "+
			"FROM user_withdraws WHERE top_id IN (%s) AND finance_process_at >= '%s' AND finance_process_at < '%s' AND status = 2 "+
			"GROUP BY top_id", ids, start, end),
		// 存款总数
		fmt.Sprintf("SELECT '' AS content, SUM(arrive_money) AS total, 'deposit' AS type, top_id "+
			"FROM user_deposits WHERE top_id IN (%s) AND confirm_at >= %d AND confirm_at < %d AND status = 2 "+
			"GROUP BY top_id", ids, startInt, endInt),
		// 活跃人数
		fmt.Sprintf("SELECT '' AS content, COUNT(DISTINCT user_id) AS total, 'active' AS type, top_id "+
			"FROM user_daily_reports WHERE top_id IN (%s) AND day >= '%s' AND day < '%s' "+
			"GROUP BY top_id", ids, start, end),
		// 总输赢
		fmt.Sprintf("SELECT '' AS content, SUM(net_money) * -1 AS total, 'final' AS type, top_id "+
			"FROM user_daily_reports WHERE top_id IN (%s) AND day >= '%s' AND day < '%s' "+
			"GROUP BY top_Id", ids, start, end),
		// 域名
		fmt.Sprintf("SELECT GROUP_CONCAT(domain) AS content, 0 AS total, 'domain' AS type, user_id AS top_id "+
			"FROM agent_domains WHERE user_id IN (%s) AND state = 2 "+
			"GROUP BY user_id", ids),
	}, ") UNION ALL (") + ")"
	if err := myClient.SQL(unionSQL).Find(&rows); err != nil {
		log.Err("获取信息出错:", err)
		response.Err(c, err.Error())
		return
	}

	for _, agent := range agents {
		// 以下设置默认值
		agent["low_num"] = "0"      // 下线
		agent["withdraws"] = "0.00" // 提款
		agent["deposits"] = "0.00"  // 存款
		agent["final"] = "0.00"     // 总输赢
		agent["active_num"] = "0"   // 活跃用户
		agent["domains"] = ""       // 域名

		agentID, _ := strconv.Atoi(agent["id"])
		for _, r := range rows {
			if r.TopId != agentID {
				continue
			}
			if r.Type == countTypes.Reg { // 下线会员
				agent["low_num"] = fmt.Sprintf("%d", int(r.Total))
			} else if r.Type == countTypes.Withdraw { // 提款总额
				agent["withdraws"] = fmt.Sprintf("%.2f", r.Total)
			} else if r.Type == countTypes.Deposit { // 存款总额
				agent["deposits"] = fmt.Sprintf("%.2f", r.Total)
			} else if r.Type == countTypes.Final { // 最终输赢
				agent["final"] = fmt.Sprintf("%.2f", r.Total)
			} else if r.Type == countTypes.Active { // 活跃人数
				agent["active_num"] = fmt.Sprintf("%d", int(r.Total))
			} else if r.Type == countTypes.Domain { // 下属域名
				agent["domains"] = r.Content
			}
		}
	}

	base_controller.SetLoginAdmin(c)                                                           // 写入登录信息
	commSQL := "SELECT agent_commission FROM agent_commission_plans GROUP BY agent_commission" // 获取佣金模板
	comRes, err := myClient.QueryString(commSQL)
	if err != nil {
		log.Err("获取佣金模板出错", err.Error())
	}

	viewData := response.ViewData{
		"rows":             agents,
		"total":            totalRow.Total,
		"agent_commission": comRes,
	}
	if request.IsAjax(c) {
		response.Render(c, "agents/_index.html", viewData)
		return
	}

	response.Render(c, "agents/frames.html", viewData)
}

/*
(SELECT COUNT(id) AS total, 'reg' AS type, top_id FROM users WHERE top_id IN (168460,168458,168457,168456,168451,168449,168412,168407,168399,168398,168302,168254,168250,168248,168237) GROUP BY top_id)
UNION ALL
(SELECT SUM(money) AS total, 'withdraw' AS type, top_id FROM user_withdraws WHERE top_id IN (168460,168458,168457,168456,168451,168449,168412,168407,168399,168398,168302,168254,168250,168248,168237) AND finance_process_at >= '2022-05-01' AND finance_process_at < '2022-05-27 09:42:38' AND status = 2 GROUP BY top_id)
UNION ALL
(SELECT SUM(arrive_money) AS total, 'deposit' AS type, top_id FROM user_deposits WHERE top_id IN (168460,168458,168457,168456,168451,168449,168412,168407,168399,168398,168302,168254,168250,168248,168237) AND confirm_at >= 1651363200 AND confirm_at < 1653615758 AND status = 2 GROUP BY top_id)
UNION ALL
(SELECT COUNT(DISTINCT user_id) AS total, 'active' AS type, top_id FROM user_daily_reports WHERE top_id IN (168460,168458,168457,168456,168451,168449,168412,168407,168399,168398,168302,168254,168250,168248,168237) AND day >= '2022-05-01' AND day < '2022-05-27 09:42:38' GROUP BY top_id)
UNION ALL
(SELECT SUM(net_money) * -1 AS total, 'final' AS type, top_id FROM user_daily_reports WHERE top_id IN (168460,168458,168457,168456,168451,168449,168412,168407,168399,168398,168302,168254,168250,168248,168237) AND day >= '2022-05-01' AND day < '2022-05-27 09:42:38' GROUP BY top_Id)
*/
