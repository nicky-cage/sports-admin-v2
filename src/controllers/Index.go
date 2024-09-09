package controllers

import (
	"sports-admin/controllers/index"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

// Index 默认页
var Index = struct {
	Index        func(*gin.Context)                             //后台默认主页 - 登录页面
	Main         func(*gin.Context)                             //后台主要界面 - 登录进入之后
	Right        func(*gin.Context)                             //默认右侧页面 - 登录之后
	Profile      func(*gin.Context)                             //显示资料 - 管理员资料
	ProfileSave  func(*gin.Context)                             //保存资料 - 管理员保存资料
	Password     func(*gin.Context)                             //修改密码
	PasswordSave func(*gin.Context)                             //保存密码修改
	Test         func(*gin.Context)                             //测试
	Login        func(*gin.Context)                             //登录处理
	Logout       func(*gin.Context)                             //退出登录
	Overload     func(*melody.Session, interface{}) interface{} //后台负载信息 - WebSocket
	GoogleCode   func(*gin.Context)                             // 验证google验证码
	GoogleBind   func(*gin.Context)                             // 绑定google验证码
	QRCode       func(*gin.Context)                             // 二维码编码
	Captcha      func(*gin.Context)                             // 生成图形验证码
	AddIP        func(*gin.Context)                             // 增加ip
	AddIPSave    func(*gin.Context)                             // 增加ip
	SendMailCode func(*gin.Context)                             // 发送邮箱验证码
	Exchange     func(*gin.Context)                             // 汇率
}{
	Index:        index.Index,
	Main:         index.Main,
	Right:        index.Right,
	Profile:      index.Profile,
	ProfileSave:  index.ProfileSave,
	Password:     index.Password,
	PasswordSave: index.PasswordSave,
	Login:        index.Login,
	Logout:       index.Logout,
	Test:         index.Test,
	Overload:     index.Overload,
	GoogleCode:   index.GoogleCode,
	GoogleBind:   index.GoogleBind,
	QRCode:       index.QRCode,
	Captcha:      index.Captcha,
	AddIP:        index.AddIP,
	AddIPSave:    index.AddIPSave,
	SendMailCode: index.SendMailCode,
	Exchange:     index.Exchange,
}
