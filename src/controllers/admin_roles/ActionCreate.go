package admin_roles

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.AdminRoles,
	ViewFile: "admin_roles/edit.html",
	ExtendData: func(c *gin.Context) response.ViewData {
		platform := request.GetPlatform(c)
		menus := caches.Menus.List(platform)
		return response.ViewData{
			"menus": menus,
		}
	},
}
