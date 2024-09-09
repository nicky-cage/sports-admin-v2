package ads

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.AdLanuches,
	ViewFile: "ads/apps_edit.html",
	ExtendData: func(*gin.Context) pongo2.Context {
		//uploadFile := config.Get("internal.img_host_backend", "")
		return pongo2.Context{
			"method": "create",
		}
	},
}
