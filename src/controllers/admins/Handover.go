package admins

import (
	"sports-admin/caches"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (ths *Admins) Handover(c *gin.Context) { //一键交接
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.ErrorHTML(c, err.Error())
		return
	}
	admin := models.Admin{}
	platform := request.GetPlatform(c)
	_, _ = models.Admins.FindById(platform, id, &admin)
	if !admin.IsDisabled() {
		response.ErrorHTML(c, "账号启用状态下不可进行一键交接")
		return
	}
	admins := []models.Admin{}
	cond := builder.NewCond().And(builder.Eq{"state": 2})
	_ = models.Admins.FindAllNoCount(platform, &admins, cond)
	response.Render(c, "admins/handover.html", pongo2.Context{
		"admin":  admin,
		"admins": admins,
		"roles":  caches.AdminRoles.All(platform),
	})
}
