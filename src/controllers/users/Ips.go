package users

import (
	"fmt"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"

	"github.com/gin-gonic/gin"
)

// Ips IP地址列表
func (ths *Users) Ips(c *gin.Context) {
	ip := c.DefaultQuery("ip", "")
	realName := c.DefaultQuery("realname", "")
	if ip == "" && realName == "" { // 无法获取ip
		response.Err(c, "IP地址和真实姓名不能同时为空")
		return
	}
	if ip != "" && !tools.IsIPv4(ip) && !tools.IsIPv6(ip) { // 如果不是v4/v6则退出
		response.Err(c, "IP地址格式有误")
		return
	}

	sql := "SELECT id, username, realname, status, created, last_login_at, register_ip, last_login_ip FROM users WHERE "
	if ip != "" {
		sql += fmt.Sprintf("register_ip = '%s' OR last_login_ip = '%s'", ip, ip)
	} else if realName != "" {
		sql += fmt.Sprintf("realname = '%s'", realName)
	}
	var rows []struct {
		Id          int    `json:"id"`
		Username    string `json:"username"`
		Realname    string `json:"realname"`
		Created     int    `json:"created"`
		Status      int    `json:"status"`
		LastLoginAt int    `json:"last_login_at"`
		RegisterIp  string `json:"register_ip"`
		LastLoginIp string `json:"last_login_ip"`
	}
	platform := request.GetPlatform(c)
	mConn := common.Mysql(platform)
	defer mConn.Close()

	if err := mConn.SQL(sql).Find(&rows); err != nil {
		response.Err(c, err.Error())
		return
	}

	rCount := 0
	lCount := 0
	for _, r := range rows {
		if ip != "" {
			if r.RegisterIp == ip {
				rCount += 1
			}
			if r.LastLoginIp == ip {
				lCount += 1
			}
		}
	}
	nCount := len(rows)

	viewData := response.ViewData{
		"ip":             ip,
		"realName":       realName,
		"rows":           rows,
		"rowsCount":      len(rows),
		"registerCount":  rCount,
		"lastLoginCount": lCount,
		"totalCount":     rCount + lCount,
		"realNameCount":  nCount,
	}
	response.Render(c, "users/used_ips.html", viewData)
}
