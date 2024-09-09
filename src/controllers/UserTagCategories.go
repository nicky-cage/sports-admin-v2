package controllers

import (
	"sports-admin/caches"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"xorm.io/builder"
)

// UserTagCategories 会员标签分类
var UserTagCategories = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
}{
	ActionList: &ActionList{
		Model:    models.UserTagCategories,
		ViewFile: "user_tag_categories/list.html",
		QueryCond: map[string]interface{}{
			"name":   "%",
			"remark": "%",
		},
		Rows: func() interface{} {
			return &[]models.UserTagCategory{}
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.UserTagCategories,
		ViewFile: "user_tag_categories/edit.html",
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.UserTagCategories,
		ViewFile: "user_tag_categories/edit.html",
		Row: func() interface{} {
			return &models.UserTagCategory{}
		},
	},
	ActionSave: &ActionSave{
		Model: models.UserTagCategories,
		SaveAfter: func(c *gin.Context, _m *map[string]interface{}) {
			platform := request.GetPlatform(c)
			caches.UserTagCategories.Load(platform)
		},
		// 添加时检测是否存在相同名称
		CreateBefore: func(c *gin.Context, m *map[string]interface{}) error {
			platform := request.GetPlatform(c)
			if name, exists := (*m)["name"]; exists {
				cond := builder.NewCond().And(builder.Eq{"name": name})
				var row models.UserTagCategory
				if ext, err := models.UserTagCategories.Find(platform, &row, cond); ext || err != nil {
					return errors.New("标签分类名称不能重复")
				}
			} else {
				return errors.New("缺少标签分类名称")
			}
			return nil
		},
		// 修改时不能修改为已存在的标签分类名称
		UpdateBefore: func(c *gin.Context, m *map[string]interface{}) error {
			platform := request.GetPlatform(c)
			if name, exists := (*m)["name"]; exists {
				id := (*m)["id"]
				cond := builder.NewCond().And(builder.Eq{"name": name}).And(builder.Neq{"id": id})
				var row models.UserTagCategory
				if ext, err := models.UserTagCategories.Find(platform, &row, cond); ext || err != nil {
					return errors.New("相同标签分类名称已经存在")
				}
			} else {
				return errors.New("缺少标签分类名称")
			}
			return nil
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.UserTagCategories,
		DeleteAfter: func(c *gin.Context, _i interface{}) {
			platform := request.GetPlatform(c)
			caches.UserTagCategories.Load(platform)
		},
	},
}
