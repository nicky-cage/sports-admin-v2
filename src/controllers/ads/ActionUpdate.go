package ads

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model:    models.AdLanuches,
	ViewFile: "ads/apps_edit.html",
	Row: func() interface{} {
		return &models.AdLanuche{}
	},
	ExtendData: func(*gin.Context) pongo2.Context {
		//uploadFile := config.Get("internal.img_host_backend", "")
		return pongo2.Context{
			//"upload_file": uploadFile,
			"method": "update",
		}
	},
}
