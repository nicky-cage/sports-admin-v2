package users

import (
	"sports-admin/controllers/base_controller"
	"sports-admin/dao"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ths *Users) Detail(c *gin.Context) { //默认首页
	idStr, exists := c.GetQuery("id")
	if !exists {
		response.ErrorHTML(c, "缺少必要参数")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ErrorHTML(c, "用户编号错误")
		return
	}

	detail, err := dao.Users.Detail(c, id)
	if err != nil {
		response.ErrorHTML(c, err.Error())
		return
	}
	if detail == nil {
		response.ErrorHTML(c, "无法获取用户详情信息")
		return
	}

	var lastUserNoteID uint32 = 0
	userNoteCount := len(detail.Notes)

	if userNoteCount > 0 {
		lastUserNoteID = detail.Notes[userNoteCount-1].Id
	}

	admin := base_controller.GetLoginAdmin(c) // 管理员信息
	response.Render(c, "users/detail.html", response.ViewData{
		"user":         detail.User,
		"agent":        detail.Agent,
		"cards":        detail.Cards,
		"money":        detail.Money,
		"lower":        detail.Lower,
		"notes":        detail.Notes,
		"noteID":       lastUserNoteID,
		"virtual_rows": detail.VirtualRows,
		"ADMIN":        admin,
	})
}
