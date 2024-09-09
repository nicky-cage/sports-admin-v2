package controllers

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-admin/validations"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

type UserAccountSetStruct struct {
	models.User `xorm:"extends"`
	Available   float64
}

// UserAccountSets 上下分-首页
var UserAccountSets = struct {
	List      func(*gin.Context)
	TopMoney  func(*gin.Context) //上分界面
	DownMoney func(*gin.Context) //下分界面
	*ActionSave
	*ActionExport
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		username := c.DefaultQuery("username", "")
		limit, offset := request.GetOffsets(c)
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"users.username": username})
		}
		userAccountSets := make([]UserAccountSetStruct, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("users").Join("LEFT OUTER", "accounts", "users.id = accounts.user_id").Where(cond).
			OrderBy("users.id ASC").Limit(limit, offset).FindAndCount(&userAccountSets)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		viewData := pongo2.Context{
			"rows":  userAccountSets,
			"total": total,
		}
		viewFile := "user_account_sets/user_account_sets.html"
		if request.IsAjax(c) {
			viewFile = "user_account_sets/_user_account_sets.html"
		}
		SetLoginAdmin(c)
		response.Render(c, viewFile, viewData)
	},
	TopMoney: func(c *gin.Context) { //上分
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		sql := "select a.*,b.available from users a left join accounts b on a.id=b.user_id where a.id=" + idStr
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		data, err := dbSession.QueryString(sql)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		viewData := pongo2.Context{"r": data[0]}
		response.Render(c, "user_account_sets/top_money.html", viewData)
	},
	DownMoney: func(c *gin.Context) { //下分
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		sql := "select a.*,b.available from users a left join accounts b on a.id=b.user_id where a.id=" + idStr
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		data, err := dbSession.QueryString(sql)
		if err != nil {
			log.Logger.Error(err.Error())
			return
		}
		viewData := pongo2.Context{"r": data[0]}
		response.Render(c, "user_account_sets/down_money.html", viewData)
	},
	ActionSave: &ActionSave{
		Model: models.UserAccountSets,
		CreateBefore: func(c *gin.Context, data *map[string]interface{}) error {
			(*data)["bill_no"] = tools.GetBillNo("D", 5)
			(*data)["applicant"] = GetLoginAdmin(c).Name
			return nil
		},
		Validator: validations.UpAndDowns,
	},
	ActionExport: &ActionExport{
		Columns: []ExportHeader{
			{Name: "序号", Field: "id"},
			{Name: "订单编号", Field: "bill_no"},
			{Name: "会员账号", Field: "username"},
			{Name: "会员等级", Field: "user_vip"},
			{Name: "操作类型", Field: "type"},
			{Name: "原因", Field: "reason"},
			{Name: "操作金额", Field: "money"},
			{Name: "投注倍数", Field: "bet_times"},
			{Name: "审核备注", Field: "audit_remark"},
			{Name: "操作时间", Field: "updated"},
			{Name: "操作人", Field: "audit"},
			{Name: "状态", Field: "status"},
		},
		ProcessRow: func(m *map[string]interface{}, c *gin.Context) {
			theTime := int((*m)["updated"].(float64))
			(*m)["updated"] = base_controller.FieldToDateTime(fmt.Sprintf("%d", theTime))
			val, err := strconv.Atoi(fmt.Sprintf("%v", (*m)["user_vip"]))
			realVal := val - 1
			platform := request.GetPlatform(c)
			userLevel := caches.UserLevels.Get(platform, realVal)
			if err == nil {
				(*m)["user_vip"] = userLevel.Name
			}
			nval := (*m)["status"].(float64)
			if nval == 1 {
				(*m)["status"] = "未处理"
			} else if nval == 2 {
				(*m)["status"] = "成功"
			} else if nval == 3 {
				(*m)["status"] = "失败"
			} else {
				(*m)["status"] = "未知"
			}
		},
	},
}
