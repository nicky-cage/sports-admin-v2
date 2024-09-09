package user_tags

import (
	"errors"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.UserTags,
	SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
		platform := request.GetPlatform(c)
		caches.UserTagCategories.Load(platform)
	},
	CreateBefore: func(c *gin.Context, m *map[string]interface{}) error {
		platform := request.GetPlatform(c)
		if name, exists := (*m)["name"]; exists {
			cond := builder.NewCond().And(builder.Eq{"name": name})
			var row models.UserTag
			if ext, err := models.UserTags.Find(platform, &row, cond); ext || err != nil {
				return errors.New("标签名称不能重复")
			}
		} else {
			return errors.New("缺少标签名称")
		}
		return nil
	},
	UpdateBefore: func(c *gin.Context, m *map[string]interface{}) error {
		if _, exists := (*m)["name"]; !exists {
			return errors.New("缺少标签名称")
		}
		return nil
	},
}
