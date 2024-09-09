package user_notes

import (
	"sports-admin/controllers/base_controller"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model:    models.UserNotes,
	ViewFile: "user_notes/edit.html",
	Row: func() interface{} {
		return &models.UserNote{}
	},
	ExtendData: func(c *gin.Context) response.ViewData {
		return getUserInfo(c)
	},
}
