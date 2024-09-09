package ad_carousels

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model:    models.AdCarousels,
	ViewFile: "ads/carousels_edit.html",
	Row: func() interface{} {
		return &models.AdCarousel{}
	},
	ExtendData: func(*gin.Context) pongo2.Context {
		return pongo2.Context{
			"method": "update",
		}
	},
}
