package index

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var Main = func(c *gin.Context) {
	admin := base_controller.GetLoginAdmin(c)
	platform := request.GetPlatform(c)
	role := caches.AdminRoles.Get(platform, int(admin.RoleId))
	site := caches.PlatformSites.GetCurrent(c)
	sta := models.FinanceMessages.Statistics(platform)
	response.Render(c, "index/main.html", response.ViewData{
		"role":  role,
		"admin": admin,
		"site":  site,
		"sta":   sta,
	})
}
