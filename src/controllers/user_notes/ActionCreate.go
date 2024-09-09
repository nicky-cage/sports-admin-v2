package user_notes

import (
	"sports-admin/controllers/base_controller"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.UserNotes,
	ViewFile: "user_notes/edit.html",
	ExtendData: func(c *gin.Context) response.ViewData {
		return getUserInfo(c)
	},
}
