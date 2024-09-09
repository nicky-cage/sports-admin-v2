package users

import (
	"sports-admin/controllers/base_controller"

	"github.com/gin-gonic/gin"
)

type UserIdRow struct {
	UserId int `json:"user_id"`
}

// Username 用户
type Username struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct {
		Exist string `json:"exist"`
	}
}

// Register 注册
type Register struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct{}
}

// UserInfoImport 用户导入数据结构
type UserInfoImport struct {
	UserName   string  // 用户名称
	State      int     // 状态
	RealName   string  // 真实姓名
	Gender     int     // 性别
	Birth      string  // 生日
	Account    float64 // 钱包余额
	Level      int     // 层级
	AgentTop   string  // 总代
	AgentTopId int     // 总代id
	Agent      string  // 代理
	AgentId    int     // 上级代理id
	Created    int64   // 注册时间
	LastLogin  int64   // 最后登录
	Phone      string  // 电话
	Mail       string  // 邮箱
	QQ         string  // qq
	WeChat     string  //wechat
	Remark     string  // 备注
}

type Users struct {
	Transfers func(*gin.Context)
	*base_controller.ActionCreate
	*base_controller.ActionUpdate
	*base_controller.ActionSave
	*base_controller.ActionList
	*base_controller.ActionState
	*base_controller.ActionExport
}
