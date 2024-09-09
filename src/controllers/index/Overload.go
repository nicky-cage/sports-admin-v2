package index

import (
	"fmt"
	"sports-common/config"
	"sports-common/tools"
	models "sports-models"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"gopkg.in/olahol/melody.v1"
)

var Overload = func(sess *melody.Session, data interface{}) interface{} {
	finalData := func(message string) map[string]interface{} { // 如果出错返回的错误信息
		return map[string]interface{}{
			"error": message,
		}
	}

	// 检测cookie是否还存在, 如果不存在, 则表示用户已经离线
	adminCookieObj, err := sess.Request.Cookie("sports_login")
	if err != nil || adminCookieObj.Value == "" {
		return finalData("你的账号已经退出登录~")
	}
	adminCookie := models.AdminCookies.FromJSON(adminCookieObj.Value)

	host := sess.Request.Host
	platform := config.GetPlatformByURL(host)

	if adminCookie == nil { // 如果不存在cookie表示用户已经离线/超时
		return finalData("你的账号已经退出登录-")
	} else if admin := models.LoginAdmins.GetLogin(platform, adminCookie.ID); admin == nil {
		return finalData("你的账号已经退出登录|") // 表示redis里已经不骨用户数据
	} else if admin.Secret != adminCookie.Secret { // 只验证密钥是否存在, 不再验证是否相等
		return finalData("你的账号已经退出登录.")
	} // else { // 因为cookie可能被劫持, 所有再进一步进行安全判断
	//	clientIPList, exists := sess.Request.Header["X-Forwarded-For"]
	//	if !exists { // 再检测 X-Real-Ip
	//		clientIPList = sess.Request.Header["X-Real-Ip"]
	//	}
	//	if len(clientIPList) < 1 {
	//		fmt.Println("你的操作可能存在安全问题, 请重新登录")
	//	}

	//	IPList := strings.Split(clientIPList[0], ",")
	//	hasLogin := false
	//	for _, v := range IPList {
	//		if strings.TrimSpace(v) == admin.LoginIP {
	//			hasLogin = true
	//			break
	//		}
	//	}
	//	//fmt.Println("admin.LoginIP = ", admin.LoginIP)
	//	if !hasLogin {
	//		return finalData("你的账号已经退出登录!")
	//	}
	//}

	// 机器相关信息统计
	m, _ := mem.VirtualMemory()
	cpuCount, _ := cpu.Counts(true)
	lastTime := time.Unix(tools.Timestamp()-1, 0)
	duration := time.Since(lastTime)
	percents, _ := cpu.Percent(duration, true)
	avg := 0.0
	total := 0.0
	for _, v := range percents {
		total += v
	}
	avg = total / float64(len(percents))
	processList, _ := process.Pids()

	financeData := map[string]interface{}{
		"user_deposit":          0,
		"user_deposit_virtual":  0,
		"user_withdraw":         0,
		"user_withdraw_virtual": 0,
		"agent_apply":           0,
		"agent_withdraw":        0,
		"today_reg":             0,
	}
	result := map[string]interface{}{
		"cpu": map[string]interface{}{
			"total": cpuCount,                 //核心数量
			"avg":   fmt.Sprintf("%.2f", avg), //平均使用
		},
		"process": map[string]interface{}{
			"total": len(processList),
		},
		"memory": map[string]interface{}{
			"total":     fmt.Sprintf("%.2f", float64(m.Total)/1024.0/1024.0/1024.0),     //内存总计
			"available": fmt.Sprintf("%.2f", float64(m.Available)/1024.0/1024.0/1024.0), //可用内存
			"free":      fmt.Sprintf("%.2f", float64(m.Free)/1024.0/1024.0/1024.0),      //空闲内存
			"percent":   fmt.Sprintf("%.2f", m.UsedPercent),                             //使用比率
		},
		"finance": financeData,
	}
	// 财务个关信息 - 将财务信息放至websocket当中, 此处不再展示
	sta := models.FinanceMessages.Statistics(platform)
	if sta != nil {
		result["finance"] = map[string]interface{}{
			"user_deposit":          sta.UserDeposit,
			"user_deposit_virtual":  sta.UserDepositVirtual,
			"user_withdraw":         sta.UserWithdraw,
			"user_withdraw_virtual": sta.UserWithdrawVirtual,
			"agent_apply":           sta.AgentApply,
			"agent_withdraw":        sta.AgentWithdraw,
			"today_reg":             sta.TodayReg,
		}
	}
	return result
}
