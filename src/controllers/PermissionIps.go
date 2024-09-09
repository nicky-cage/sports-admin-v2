package controllers

import (
	"sports-admin/caches"
	common "sports-common"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// PermissionIps 授权IP
var PermissionIps = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
	*ActionState
}{
	ActionList: &ActionList{
		Model:    models.PermissionIps,
		ViewFile: "permission_ips/list.html",
		Rows: func() interface{} {
			return &[]models.PermissionIp{}
		},
		QueryCond: map[string]interface{}{
			"name":   "%",
			"ip":     "%",
			"remark": "%",
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.PermissionIps,
		ViewFile: "permission_ips/edit.html",
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.PermissionIps,
		ViewFile: "permission_ips/edit.html",
		Row: func() interface{} {
			return &models.PermissionIp{}
		},
	},
	ActionSave: &ActionSave{
		Model: models.PermissionIps,
		SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
			platform := request.GetPlatform(c)
			if ip, ok := (*data)["ip"].(string); ok {
				dbSession := common.Mysql(platform)
				defer dbSession.Close()
				dbSession.Exec("DELETE FROM permission_ips WHERE ip = ?", ip) // 先删掉原来重复的ip
			}
			// -- 检测google验证密码
			if err := CheckGoogleCode(platform, c, data); err != nil {
				return err
			}
			return nil
		},
		SaveAfter: func(c *gin.Context, data *map[string]interface{}) {
			platform := request.GetPlatform(c)
			caches.PermissionIps.Load(platform) // 刷新缓存
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.PermissionIps,
		DeleteAfter: func(c *gin.Context, _ interface{}) {
			platform := request.GetPlatform(c)
			caches.PermissionIps.Load(platform) // 刷新缓存
		},
	},
	ActionState: &ActionState{
		Model: models.PermissionIps,
		Field: "state",
		StateAfter: func(c *gin.Context) {
			platform := request.GetPlatform(c)
			caches.PermissionIps.Load(platform) // 刷新缓存
		},
	},
}
