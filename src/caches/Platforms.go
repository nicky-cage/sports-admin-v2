package caches

import (
	"sports-common/consts"
	models "sports-models"
)

// 平台key
const keyPlatforms = "platforms"

// Platforms 平台信息
var Platforms = struct {
	Load func()
	Get  func(int) *models.Platform
	All  func() map[uint32]models.Platform
}{
	Load: func() {
		rows := map[uint32]models.Platform{}
		var rs []models.Platform
		err := models.Platforms.FindAllNoCount(consts.PlatformIntegrated, &rs, nil, "id ASC")
		if err != nil {
			return
		}
		for _, r := range rs {
			rows[r.Id] = r
		}

		_ = setCache(consts.PlatformIntegrated, keyPlatforms, rows)
	},
	Get: func(id int) *models.Platform {
		rows := map[uint32]models.Platform{}
		_ = getCache(consts.PlatformIntegrated, keyPlatforms, &rows)

		realID := uint32(id)
		for k, v := range rows {
			if k == realID {
				return &v
			}
		}

		return nil
	},
	All: func() map[uint32]models.Platform {
		rows := map[uint32]models.Platform{}
		_ = getCache(consts.PlatformIntegrated, keyPlatforms, &rows)
		return rows
	},
}
