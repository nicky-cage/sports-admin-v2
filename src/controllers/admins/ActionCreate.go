package admins

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.Admins,
	ViewFile: "admins/edit.html",
	ExtendData: func(c *gin.Context) pongo2.Context {
		platform := request.GetPlatform(c)
		return pongo2.Context{
			"admin_roles": caches.AdminRoles.All(platform),
		}
	},
}
