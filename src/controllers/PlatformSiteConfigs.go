package controllers

import (
	"encoding/json"
	"sports-admin/caches"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// PlatformSiteConfigs 平台 - 站点 - 配置
var PlatformSiteConfigs = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionState
}{
	ActionList: &ActionList{
		Model:    models.SiteConfigs,
		ViewFile: "platform_site_configs/list.html",
		Rows: func() interface{} {
			return &[]models.SiteConfig{}
		},
		QueryCond: map[string]interface{}{
			"platform_id": "=",
			"name":        "%",
			"value":       "%",
			"remark":      "%",
		},
		ExtendData: func(*gin.Context) ViewData {
			return ViewData{
				"platforms": caches.Platforms.All(),
			}
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.SiteConfigs,
		ViewFile: "platform_site_configs/edit.html",
		ExtendData: func(*gin.Context) ViewData {
			platforms := models.Platforms.Related()
			bytes, _ := json.Marshal(platforms)
			return ViewData{
				"platforms":     platforms,
				"platformsJSON": string(bytes),
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.SiteConfigs,
		ViewFile: "platform_site_configs/edit.html",
		Row: func() interface{} {
			return &models.SiteConfig{}
		},
		ExtendData: func(*gin.Context) ViewData {
			platforms := models.Platforms.Related()
			bytes, _ := json.Marshal(platforms)
			return ViewData{
				"platforms":     platforms,
				"platformsJSON": string(bytes),
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.SiteConfigs,
	},
	ActionState: &ActionState{
		Model: models.SiteConfigs,
		Field: "status",
	},
}
