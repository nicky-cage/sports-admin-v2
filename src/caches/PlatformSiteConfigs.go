package caches

import (
	"sports-common/consts"
	models "sports-models"
)

const keyPlatformSiteConfigs = "platform_site_configs"

// PlatformSiteConfigs 站点配置信息
var PlatformSiteConfigs = struct {
	Load func()
	Get  func(int) *models.SiteConfig
	All  func() map[uint32]models.SiteConfig
}{
	Load: func() {
		rows := map[uint32]models.SiteConfig{}
		var rs []models.SiteConfig
		err := models.SiteConfigs.FindAllNoCount(consts.PlatformIntegrated, &rs, nil, "id ASC")
		if err != nil {
			return
		}
		for _, r := range rs {
			rows[r.Id] = r
		}

		_ = setCache(consts.PlatformIntegrated, keyPlatformSiteConfigs, rows)
	},
	Get: func(id int) *models.SiteConfig {
		rows := map[uint32]models.SiteConfig{}
		_ = getCache(consts.PlatformIntegrated, keyPlatformSiteConfigs, &rows)

		realID := uint32(id)
		for k, v := range rows {
			if k == realID {
				return &v
			}
		}

		return nil
	},
	All: func() map[uint32]models.SiteConfig {
		rows := map[uint32]models.SiteConfig{}
		_ = getCache(consts.PlatformIntegrated, keyPlatformSiteConfigs, &rows)
		return rows
	},
}
