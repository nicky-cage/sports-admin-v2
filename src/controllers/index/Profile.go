package index

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var Profile = func(c *gin.Context) {
	login := base_controller.GetLoginAdmin(c)
	admin := models.Admin{}
	platform := request.GetPlatform(c)
	cond := builder.NewCond().And(builder.Eq{"id": login.Id})
	_, _ = models.Admins.Find(platform, &admin, cond)
	response.Render(c, "index/profile.html", pongo2.Context{
		"admin":      login,
		"r":          admin,
		"adminRoles": caches.AdminRoles.All(platform),
	})
}
