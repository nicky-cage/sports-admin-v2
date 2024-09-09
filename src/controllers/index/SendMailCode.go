package index

import (
	"crypto/tls"
	"fmt"
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// SendMailCode 发送邮件
// Web_login 	https://bvault-me.awsapps.com/mail
// Email     	service@bvault.me
// Account     	service
// Password 	1qaz@WSX3edc
// IMAP         imap.mail.us-east-1.awsapps.com:993(SSL)
// SMTP(tls)   	smtp.mail.us-east-1.awsapps.com:465(SSL)
var SendMailCode = func(c *gin.Context) {
	platform := request.GetPlatform(c)
	admin, err := checkLogin(c)
	if err != nil {
		response.Err(c, err.Error())
		return
	}

	mailCode := tools.RandString(6)
	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	key := "admin_mail_code:" + admin.Name
	val, _ := rd.Get(key).Result()
	if val != "" {
		response.Err(c, "邮件验证密码已经发送,有效期5分钟")
		return
	}

	mailUser := config.Get("mail.username")
	mailPass := config.Get("mail.password")
	mailHost := config.Get("mail.host")
	mailPort := config.GetInt("mail.port")
	sendMail := func(mail, title, content string) error { //定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
		var conn = struct {
			User     string
			Password string
			Host     string
			Port     int
		}{
			User:     mailUser,
			Password: mailPass,
			Host:     mailHost,
			Port:     mailPort,
		}
		m := gomail.NewMessage()
		m.SetHeader("From", m.FormatAddress(conn.User, "自动增加授权IP")) //这种方式可以添加别名，即“XX官方”
		//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
		//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
		//m.SetHeader("From", mailConn["user"])
		m.SetHeader("To", mail)         //发送给多个用户
		m.SetHeader("Subject", title)   //设置邮件主题
		m.SetBody("text/html", content) //设置邮件正文
		d := gomail.NewDialer(conn.Host, conn.Port, conn.User, conn.Password)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		err := d.DialAndSend(m)
		return err
	}

	title := admin.Name + ", 请完成你的增加IP/邮箱验证"
	content := fmt.Sprintf("平台名称: 综合平台<br />"+
		"用户名称: %s<br />"+
		"验证邮件: %s<br />"+
		"上次登录: %s<br />"+
		"------------------------------------------------------------------<br />"+
		"你的邮箱验证密码是: %s, 请在5分钟内提交验证!<br /><br />(验证密码仅限当次有效, 请不要回复此邮件)", admin.Name, admin.Mail, admin.LastIp, mailCode)
	if err := sendMail(admin.Mail, title, content); err != nil {
		log.Logger.Error("发送邮件错误:", err)
		response.Err(c, "邮件发送发生错误:"+err.Error())
		return
	}

	// 如果邮件发送成功, 则写入redis信息
	_, _ = rd.Set(key, mailCode, 300*time.Second).Result() // 5m 内有效
	response.Ok(c)
}
