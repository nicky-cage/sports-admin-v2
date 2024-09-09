package dao

import (
	"errors"
	"sports-common/config"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// Index 常规方法 -
var Index = struct {
	Login          func(*gin.Context, string, *map[string]interface{}) (*models.Admin, error)
	UpdatePassword func(string, uint32, map[string]interface{}) error
}{
	Login: func(c *gin.Context, platform string, data *map[string]interface{}) (*models.Admin, error) {
		username := (*data)["username"].(string)
		cond := builder.NewCond().And(builder.Eq{"name": username})
		admin := models.Admin{}
		exists, err := models.Admins.Find(platform, &admin, cond)
		if !exists { // 如果记录不存在
			return nil, errors.New("用户名称或密码错误")
		}
		if err != nil { // 如果登录错误
			return nil, err
		}

		if admin.IsDisabled() { // 如果此用户已被禁用
			return nil, errors.New("此用户已被禁用")
		}
		password := (*data)["password"].(string)
		saltPassword := tools.GetPassword(password, admin.Salt)
		if saltPassword != admin.Password {
			return nil, errors.New("用户密码或者密码错误")
		}

		loginIP := strings.Split((*data)["last_ip"].(string), ":")[0]
		if !strings.Contains(","+admin.AllowIps+",", ","+loginIP+",") { // 判断是否在ip可控范围之内
			return nil, errors.New("此IP不在授权范围之内: " + loginIP)
		}

		// 检测google验证 - 只在生产环境开启
		if config.EnvIsProduct() && admin.GoogleEnable() { // 如果已启用
			if code, exists := (*data)["google_code"]; !exists {
				return nil, errors.New("必须输入谷歌验证密码")
			} else if chk, err := tools.NewGoogleAuth().VerifyCode(admin.GoogleSecret, code.(string)); err != nil {
				return nil, err
			} else if !chk {
				return nil, errors.New("谷歌验证密码错误")
			}
		}

		// 再检查是否在数据库范围之内 - middleware - CheckPermissionIp 是第一重保险, 这里是第二重保险
		condIP := builder.NewCond().And(builder.Eq{"ip": loginIP})
		permissionIP := models.PermissionIp{}
		if exists, err := models.PermissionIps.Find(platform, &permissionIP, condIP); err != nil || !exists { // 检测是否存在此IP授权: error or no this record
			return nil, errors.New("此IP不在可访问授权范围之内: " + loginIP)
		} else if !permissionIP.Enabled() { // 如果没有启用此ip
			return nil, errors.New("此IP未授权访问后台管理系统: " + loginIP)
		}

		saltNew := tools.Secret()
		passwordNew := tools.GetPassword(password, saltNew)
		dataNew := map[string]interface{}{
			"id":          admin.Id,             // 用户编号
			"salt":        saltNew,              // 新的密钥
			"password":    passwordNew,          // 新的密码
			"login_count": admin.LoginCount + 1, // 登录次数累加
			"last_login":  tools.Now(),          // 最后登录时间
			"last_ip":     loginIP,              // 最后登录ip
			"google_code": 2,                    // 将谷歌验证设为必须 - 每次登录,都自动变更为下次必须google验证码
		}
		if err := models.Admins.Update(platform, dataNew); err != nil {
			return nil, err
		}

		// 写入登录日志
		logData := map[string]interface{}{
			"admin_id":   admin.Id,
			"admin_name": admin.Name,
			"ip":         loginIP,
			"user_agent": c.Request.UserAgent(),
		}
		_, _ = models.AdminLoginLogs.Create(platform, logData)

		return &admin, nil
	},
	UpdatePassword: func(platform string, adminId uint32, data map[string]interface{}) error {
		password := data["password"].(string)        //旧的密码
		passwordNew := data["password_new"].(string) //新的密码
		passwordRep := data["password_rep"].(string) //得复密码
		if passwordNew != passwordRep {
			return errors.New("两次输入的密码不一致")
		}
		if passwordNew == password {
			return errors.New("旧的密码不能与新密码一致")
		}

		// 获取用户信息
		admin := models.Admin{}
		cond := builder.NewCond().And(builder.Eq{"id": adminId}) //拿到后台用户信息
		exists, err := models.Admins.Find(platform, &admin, cond)
		if err != nil {
			return err
		}
		if !exists { //记录不存在
			return errors.New("后台用户信息不存在")
		}
		passwordReal := tools.GetPassword(password, admin.Salt)
		if passwordReal != admin.Password {
			return errors.New("旧的密码输入有误")
		}

		secret := tools.Secret()
		passwordRealNew := tools.GetPassword(passwordNew, secret)
		updatedData := map[string]interface{}{
			"id":       adminId,
			"password": passwordRealNew,
			"salt":     secret,
		}
		if err = models.Admins.Update(platform, updatedData); err != nil {
			return err
		}

		return nil
	},
}
