package ad_carousels

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.AdCarousels,
	ViewFile: "ads/carousels_edit.html",
	ExtendData: func(*gin.Context) pongo2.Context {
		return pongo2.Context{
			"method": "create",
		}
	},
}
