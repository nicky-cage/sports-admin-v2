package user_tags

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionDelete = &base_controller.ActionDelete{
	Model: models.UserTags,
	DeleteAfter: func(c *gin.Context, _i interface{}) {
		platform := request.GetPlatform(c)
		caches.UserTagCategories.Load(platform)
	},
}
