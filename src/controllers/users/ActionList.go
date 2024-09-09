package users

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionList = &base_controller.ActionList{
	Model:    models.UserAccountInfos,
	ViewFile: "users/list.html",
	Rows: func() interface{} {
		return &[]models.UserAccountInfo{}
	},
	OrderBy: func(c *gin.Context) string {
		orderSort := "accounts.available"
		orderBy := "DESC"
		if ordSort := c.DefaultQuery("order_sort", ""); ordSort != "" {
			if ordSort == "created" { // 注册时间
				orderSort = "users.created"
			} else if ordSort == "vip" { // vip
				orderSort = "users.vip"
			} else if ordSort == "username" { // 钱包余额
				orderSort = "users.username"
			}
		}
		if ordBy := c.DefaultQuery("order_by", ""); ordBy != "" && (ordBy == "desc" || ordBy == "asc") {
			orderBy = ordBy
		}

		return orderSort + " " + orderBy
	},
	ExtendData: func(c *gin.Context) pongo2.Context { //钱包金额
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()

		// 有成功充值记录 - 一键全部回收
		sqlForDep := "SELECT DISTINCT(user_id) AS user_id FROM user_deposits WHERE money > 0 AND status = 2"
		rows, err := dbSession.QueryString(sqlForDep)
		var uArr []string
		if err == nil && len(rows) > 0 {
			for _, r := range rows {
				uArr = append(uArr, r["user_id"])
			}
		}

		onlineCount := len(base_controller.GetAllOnlineUser(platform))
		vipLevels := caches.UserLevels.All(platform)
		gameList, _ := models.GetGameListAllHandler(platform)
		return pongo2.Context{
			"vipLevels":   vipLevels,
			"onlineCount": onlineCount,
			"userIds":     strings.Join(uArr, ","),
			"gameList":    gameList,
		}
	},
	QueryCond: map[string]interface{}{
		"users.realname":      "%",
		"users.phone":         "%",
		"users.vip":           "=",
		"users.label":         "%",
		"users.status":        "=",
		"users.last_login_ip": "%",
		"users.register_ip":   "%",
	},
	GetQueryCond: func(c *gin.Context) builder.Cond {
		platform := request.GetPlatform(c)
		_ = base_controller.GetAllOnlineUser(platform) // 确保会执行获取在线用户数据方法
		cond := builder.NewCond()

		// 以,号隔开以精确匹配 以;号隔开以模糊匹配
		processConds := func(field string, val string) {
			qArr := strings.Split(val, ",")
			mArr := []string{}
			for _, v := range qArr {
				rv := strings.TrimSpace(v)
				if rv == "" {
					continue
				}
				if strings.Index(rv, ";") > 0 { // 模糊匹配部分
					pArr := strings.Split(rv, ";")
					for _, pv := range pArr {
						rpv := strings.TrimSpace(pv)
						if rpv == "" {
							continue
						}
						cond = cond.Or(builder.Like{field, rpv})
					}
				} else { // 精确匹配
					mArr = append(mArr, rv)
					// cond = cond.Or(builder.Eq{field: rv})
				}
			}
			if len(mArr) > 0 {
				cond = cond.And(builder.In(field, mArr))
			}
		}

		if userName := strings.TrimSpace(c.DefaultQuery("users.username", "")); userName == "" { // 如果有搜索用户名称 //单纯查注册日期人数
			if created := c.DefaultQuery("created", ""); created != "" {
				areas := strings.Split(created, " - ")
				start := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
				end := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
				cond = cond.And(builder.Gte{"users.created": start}).And(builder.Lte{"users.created": end})
			}
			if lastLogin := c.DefaultQuery("last_login_at", ""); lastLogin != "" {
				areas := strings.Split(lastLogin, " - ")
				start := tools.GetTimeStampByString(areas[0] + " 00:00:00")
				end := tools.GetTimeStampByString(areas[1] + " 23:59:59")
				cond = cond.And(builder.Gte{"users.last_login_at": start}).And(builder.Lte{"users.last_login_at": end})
			}
		} else { // 如果有用户名称 以,号隔开以精确匹配 以;号隔开以模糊匹配
			processConds("users.username", userName)
		}

		if topName := strings.TrimSpace(c.DefaultQuery("users.top_name", "")); topName != "" { // 如果有上级用户名称 ,号精确 ;号模糊
			processConds("users.top_name", topName)
		}

		// 关于用户在线的判断
		if onlineStatus := c.DefaultQuery("online", ""); onlineStatus == "1" || onlineStatus == "0" {
			onlineUserNames := base_controller.GetAllOnlineUser(platform)
			var userNames []interface{}
			for _, userName := range onlineUserNames {
				userNames = append(userNames, userName)
			}
			if onlineStatus == "1" {
				cond = cond.And(builder.In("users.username", userNames...))
			} else {
				cond = cond.And(builder.NotIn("users.username", userNames...))
			}
		}
		return cond
	},
	ProcessRow: func(c *gin.Context, rows interface{}) {
		rs := rows.(*[]models.UserAccountInfo)
		platform := request.GetPlatform(c)
		var uArr []string
		for k, r := range *rs {
			uArr = append(uArr, fmt.Sprintf("%d", r.Id))
			if ud, exists := base_controller.OnlineUserDatas.Data[platform]; exists {
				if _, has := ud.UserTokens[r.User.Username]; has {
					(*rs)[k].User.Online = true
				}
			}
		}
		if len(uArr) == 0 {
			return
		}
		mClient := common.Mysql(platform)
		defer mClient.Close()
		sql := "SELECT user_id, SUM(money) AS total FROM user_games WHERE user_id IN (" + strings.Join(uArr, ",") + ") GROUP BY user_id"
		tRows := []struct {
			UserId uint64  `json:"user_id"`
			Total  float64 `json:"total"`
		}{}
		if err := mClient.SQL(sql).Find(&tRows); err != nil {
			return
		}

		for rk, rv := range *rs {
			for _, tv := range tRows {
				if tv.UserId != rv.Id {
					continue
				}
				(*rs)[rk].AvailableTotal = tv.Total
			}
		}
	},
	AfterAction: func(c *gin.Context, data *response.ViewData) {
		rows := (*data)["rows"].(*[]models.UserAccountInfo)
		var rIps []string  // 注册ip
		var lIps []string  // 最后登录ip
		var nList []string // 姓名数组
		for _, r := range *rows {
			rIps = append(rIps, r.RegisterIp)
			lIps = append(lIps, r.LastLoginIp)
			nList = append(nList, r.RealName)
		}

		sql := strings.Join([]string{
			"(SELECT 'n' AS type, realname AS ip, COUNT(*) AS total FROM users WHERE realname IN ('" + strings.Join(nList, "','") + "') GROUP BY realname)",
			"(SELECT 'r' AS type, register_ip AS ip, COUNT(*) AS total FROM users WHERE register_ip IN ('" + strings.Join(rIps, "','") + "') GROUP BY register_ip)",
			"(SELECT 'l' AS type, last_login_ip AS ip, COUNT(*) AS total FROM users WHERE last_login_ip IN ('" + strings.Join(lIps, "','") + "') GROUP BY last_login_ip)",
		}, " UNION ")
		platform := request.GetPlatform(c)
		mConn := common.Mysql(platform)
		defer mConn.Close()
		lArr := map[string]int{}
		rArr := map[string]int{}
		nArr := map[string]int{}
		if rs, err := mConn.QueryString(sql); err == nil {
			for _, r := range rs {
				total, _ := strconv.Atoi(r["total"])
				if r["type"] == "l" {
					lArr[r["ip"]] = total // 登录IP
				} else if r["type"] == "r" {
					rArr[r["ip"]] = total // 注册IP
				} else if r["type"] == "n" {
					nArr[r["ip"]] = total
				}
			}
		}

		(*data)["registerArr"] = rArr
		(*data)["lastLoginArr"] = lArr
		(*data)["realNameArr"] = nArr
	},
}
