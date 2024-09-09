package index

import (
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/captchas"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

var AddIPSave = func(c *gin.Context) {
	platform := request.GetPlatform(c)
	postedData := request.GetPostedData(c)
	verifyID, existsVerifyID := postedData["captchaID"]
	verifyCode, existsVerifyCode := postedData["verify_code"]
	if !existsVerifyID || !existsVerifyCode {
		response.Err(c, "缺少图形验证码")
		return
	}
	if !captchas.Capt(platform).Verify(verifyID.(string), verifyCode.(string)) {
		response.Err(c, "图形验证码校验失败")
		return
	}

	mailCode, existsMailCode := postedData["mail_code"]
	if !existsMailCode {
		response.Err(c, "缺少邮箱验证码")
	}

	admin, err := checkLogin(c)
	if err != nil {
		response.Err(c, err.Error())
		return
	}

	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	key := "admin_mail_code:" + admin.Name
	val, _ := rd.Get(key).Result()
	if val == "" {
		response.Err(c, "邮箱验证密码超时")
		return
	}
	if val != mailCode {
		response.Err(c, "邮箱验证密码错误")
		return
	}
	defer rd.Del(key).Result()

	realIP := c.ClientIP()
	// 修改管理员ip
	sql := "UPDATE admins SET allow_ips = ?, updated = UNIX_TIMESTAMP() WHERE id = ? LIMIT 1"
	dbSession := common.Mysql(platform)
	defer dbSession.Close()

	_, _ = dbSession.Exec(sql, realIP, admin.Id)

	_, _ = dbSession.Exec("DELETE FROM permission_ips WHERE ip = ?", realIP)
	remark := fmt.Sprintf("%s: 邮件绑定", admin.Name)
	_, _ = dbSession.Exec("INSERT INTO permission_ips (ip, remark, state, created, updated) "+
		"VALUES (?, ?, 2, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())", realIP, remark)

	caches.PermissionIps.Load(platform) // 刷新缓存
	response.Ok(c)
}
