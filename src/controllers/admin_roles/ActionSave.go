package admin_roles

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
	Model: models.AdminRoles,
	SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
		platform := request.GetPlatform(c)
		caches.AdminRoles.Load(platform, true) // 刷新缓存
	},
	CreateBefore: func(c *gin.Context, m *map[string]interface{}) error {
		name, exists := (*m)["name"]
		if !exists {
			return errors.New("缺少角色名称")
		}
		row := models.AdminRole{}
		cond := builder.NewCond().And(builder.Eq{"name": name})
		platform := request.GetPlatform(c)
		if exists, err := models.AdminRoles.Find(platform, &row, cond); err != nil || exists {
			return errors.New("添加失败: 相同的角色名称已经存在")
		}
		return nil
	},
	UpdateBefore: func(c *gin.Context, m *map[string]interface{}) error {
		name, exists := (*m)["name"]
		if !exists {
			return errors.New("缺少角色名称")
		}
		id, exists := (*m)["id"]
		if !exists {
			return errors.New("缺少编号字段")
		}
		row := models.AdminRole{}
		platform := request.GetPlatform(c)
		cond := builder.NewCond().And(builder.Eq{"name": name}).And(builder.Neq{"id": id})
		if exists, err := models.AdminRoles.Find(platform, &row, cond); err != nil || exists {
			return errors.New("修改失败: 要修改的角色名称已经存在")
		}
		return nil
	},
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		platform := request.GetPlatform(c)
		// -- 检查google验证码
		if err := base_controller.CheckGoogleCode(platform, c, m); err != nil {
			return err
		}
		return nil
	},
}
