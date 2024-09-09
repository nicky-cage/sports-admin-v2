package admin_roles

import (
	"sports-admin/caches"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

func (ths *AdminRoles) Configs(c *gin.Context) {
	platform := request.GetPlatform(c)
	menus := caches.Menus.List(platform)
	response.Render(c, "admin_roles/configs.html", response.ViewData{
		"menus": menus,
	})
}
