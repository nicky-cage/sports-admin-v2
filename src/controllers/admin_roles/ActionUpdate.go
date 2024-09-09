package admin_roles

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model:    models.AdminRoles,
	ViewFile: "admin_roles/edit.html",
	Row: func() interface{} {
		return &models.AdminRole{}
	},
	ExtendData: func(c *gin.Context) response.ViewData {
		platform := request.GetPlatform(c)
		menus := caches.Menus.List(platform)
		allMenus := caches.Menus.LayMenusByJson(platform)
		return response.ViewData{
			"menus":         menus,
			"menusJsonList": models.LayMenus.ToListJson(menus),
			"menusJsonData": allMenus,
		}
	},
}
