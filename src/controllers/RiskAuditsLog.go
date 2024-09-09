package controllers

import (
	"fmt"
	"sports-admin/filters"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// RiskAuditsLog 审核历史记录
var RiskAuditsLog = struct {
	//*ActionList
	Update func(c *gin.Context)
	List   func(c *gin.Context)
}{
	List: func(c *gin.Context) {
		var part = ""
		var cond = ""
		//参数判断
		if username := c.Query("username"); username != "" {
			part += " AND l.username LIKE '%" + username + "%' "
		}
		if billNo := c.Query("bill_no"); billNo != "" {
			part += " AND l.bill_no ='" + billNo + "' "
		}
		if status := c.Query("status"); status != "" {
			part += " AND l.status = '" + status + "' "
		}
		if bankRealName := c.Query("bank_real_name"); bankRealName != "" {
			part += " AND l.bank_real_name LIKE '%" + bankRealName + "%' "
		}
		if money := c.Query("money"); money != "" {
			part += " AND l.money = '" + money + "' "
		}
		timeStart, timeEnd := func() (int64, int64) {
			value := c.DefaultQuery("created", "")
			if value == "" {
				startAt := tools.GetTodayBegin()
				return tools.SecondToMicro(startAt), tools.SecondToMicro(startAt + 86400)
			}
			areas := strings.Split(value, " - ")
			startAt := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
			endAt := tools.GetMicroTimeStampByString(areas[1]+" 00:00:00") + 86400
			return startAt, endAt
		}()
		part += fmt.Sprintf(" AND (l.created >= %d AND l.created < %d) ", timeStart, timeEnd)

		page := c.Query("page")
		p, _ := strconv.Atoi(page)
		if p > 1 {
			p = (p - 1) * 15
		} else {
			p = 0
		}
		page = strconv.Itoa(p)
		limit, offset := request.GetOffsets(c)

		//sql查询
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		datas := make([]map[string]string, 0)

		//显示领取了。未处理的订单。
		riskSql := "select a.*, b.vip, b.label from user_withdraws a join users b on a.user_id = b.id where  a.risk_admin!='' and a.status=1  "
		riskData, err := dbSession.QueryString(riskSql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		wsql := "select arrive_money, confirm_at from  user_deposits where user_id=%s and status=2  order by confirm_at desc limit 1" // 存款
		aSql := "select money,updated from user_account_sets where user_id=%s and status=2 and type= 1 order by updated DESC limit 1" // 上分
		for _, v := range riskData {
			if v["bank_name"] == "其他银行" {
				sqlCard := "select bank_name from user_cards where card_number= " + v["bank_card"]
				dataCard, err := dbSession.QueryString(sqlCard)
				if err == nil && len(dataCard) > 0 && dataCard[0]["bank_name"] != "" {
					v["bank_name"] = dataCard[0]["bank_name"]
				}
			}
			tempSql := fmt.Sprintf(wsql, v["user_id"])
			res, err := dbSession.QueryString(tempSql)
			if err != nil {
				log.Err(err.Error())
				response.Err(c, "系统错误")
				return
			}
			aSqll := fmt.Sprintf(aSql, v["user_id"])
			aRes, err := dbSession.QueryString(aSqll)
			if err != nil {
				log.Err(err.Error())
				return
			}
			var atime int
			if len(aRes) > 0 {
				atime, _ = strconv.Atoi(aRes[0]["updated"])
			}
			var wtime int
			if len(res) > 0 {
				wtime, _ = strconv.Atoi(res[0]["confirm_at"])
			}

			if atime > wtime {
				v["last_money"] = aRes[0]["money"]
			} else {
				if wtime > 0 {
					v["last_money"] = res[0]["arrive_money"]
				} else {
					v["last_money"] = "0"
				}
			}

			id, _ := strconv.Atoi(v["id"])
			userId, _ := strconv.Atoi(v["user_id"])
			info := models.RefreshFlow(platform, id, userId)
			v["flow_total"] = fmt.Sprintf("%.2f", info["flow_total"].(float64))
			v["flow_current"] = fmt.Sprintf("%.2f", info["flow_current"].(float64))
			v["sys"] = "-"
			v["user_label"] = filters.GetUserLabels(platform, v["label"])
			v["status"] = "2"
			datas = append(datas, v)
		}

		cond = fmt.Sprintf(" ORDER BY l.risk_process_at DESC LIMIT %d, %d", limit, offset)
		sql := "SELECT l.*, w.withdraw_cost FROM risk_audit_logs AS l, user_withdraws AS w WHERE l.bill_no = w.bill_no AND 1 = 1 "
		data, err := dbSession.QueryString(sql + part + cond)
		if err != nil {
			log.Err(err.Error())
			return
		}

		for _, v := range data {
			if v["bank_name"] == "其他银行" {
				sqlCard := "SELECT bank_name FROM user_cards WHERE card_number = '" + v["bank_card"] + "'"
				dataCard, err := dbSession.QueryString(sqlCard)
				if err == nil && len(dataCard) >= 1 {
					if dataCard[0]["bank_name"] != "" {
						v["bank_name"] = dataCard[0]["bank_name"]
					}
				}
			}
			//获取标签
			sqlLabel := "SELECT label FROM users WHERE username = '" + v["username"] + "'"
			dataLabel, _ := dbSession.QueryString(sqlLabel)
			if err == nil && len(dataLabel) > 0 {
				v["UserLabel"] = filters.GetUserLabels(platform, dataLabel[0]["label"])
			}
			v["status"] = "2"
			datas = append(datas, v)
		}

		sqlCount := func() int { // 总结统计
			rows, err := dbSession.QueryString("SELECT COUNT(*) AS total FROM risk_audit_logs AS l WHERE 1 = 1 " + part)
			if err != nil {
				panic(err)
			}
			if len(rows) == 0 {
				return 0
			}
			total, _ := strconv.Atoi(rows[0]["total"])
			return total
		}()
		SetLoginAdmin(c)
		response.Render(c, "risk_audits/_history.html", ViewData{
			"rows":  datas,
			"total": sqlCount,
		})
	},
	//ActionList: &ActionList{
	//	ViewFile: "risk_audits/history.html",
	//	Model:    models.RiskAuditLogs,
	//	QueryCond: map[string]interface{}{
	//		"username":       "%",
	//		"bill_no":        "=",
	//		"status":         "=",
	//		"bank_real_name": "%",
	//		"money":          "=",
	//	},
	//	Rows: func() interface{} {
	//		return &[]models.RiskAuditLog{}
	//	},
	//	GetQueryCond: func(c *gin.Context) builder.Cond {
	//		cond := builder.NewCond()
	//		var startAt int64
	//		var endAt int64
	//		if value, exists := c.GetQuery("created"); !exists {
	//			currentTime := time.Now().Unix()
	//			startAt = currentTime - currentTime%86400
	//			endAt = startAt + 86400
	//		} else {
	//			areas := strings.Split(value, " - ")
	//			startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
	//			endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
	//			cond = cond.And(builder.Gte{"created": startAt}).And(builder.Lt{"created": endAt})
	//		}
	//		return cond
	//	},
	//	ExtendData: func(c *gin.Context) ViewData {
	//		return ViewData{
	//			"vipLevels": caches.UserLevels.All(),
	//		}
	//	},
	//	ProcessRow: func(list interface{}) {
	//		rows := list.(*[]models.RiskAuditLog)
	//		for k, v := range *rows {
	//			if rs, err := sess.QueryString("SELECT label FROM users WHERE username = ?", v.Username); err == nil && len(rs) > 0 {
	//				(*rows)[k].UserLabel = filters.GetUserLabels(rs[0]["label"])
	//			}
	//		}
	//	},
	//	OrderBy: func(*gin.Context) string {
	//		return "risk_process_at DESC"
	//	},
	//},
	//ActionUpdate: &ActionUpdate{
	//	Model:    models.RiskAuditLogs,
	//	ViewFile: "risk_audits/histories_detail.html",
	//	Row: func() interface{} {
	//		return &models.UserWithdraw{}
	//	},
	//},
	Update: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		id := c.Query("id")
		sql := "SELECT a.*, " +
			"b.finance_admin, b.money, b.withdraw_cost, b.transaction_fee, b.actual_money, b.card_number,b.finance_process_at,b.status as w_status,b.remark as f_remark, b.business_type " +
			"FROM risk_audit_logs a JOIN user_withdraws b ON a.bill_no = b.bill_no " +
			"WHERE a.bill_no = '" + id + "'"
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
		}
		response.Render(c, "risk_audits/histories_detail.html", pongo2.Context{"r": res[0]})
	},
}
