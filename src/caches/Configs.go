package caches

import (
	models "sports-models"
)

const keyConfigs = "configs"

// Configs 配置信息
var Configs = struct {
	Load func(string)
	Get  func(string, int) *models.Config
}{
	Load: func(platform string) {
		configs := map[uint32]models.Config{}
		var rs []models.Config
		err := models.Configs.FindAllNoCount(platform, &rs, nil, "id ASC")
		if err != nil {
			return
		}

		for _, r := range rs {
			configs[r.Id] = r
		}

		_ = setCache(platform, keyConfigs, configs)
	},
	Get: func(platform string, id int) *models.Config {
		configs := map[uint32]models.Config{}
		_ = getCache(platform, keyConfigs, &configs)

		realID := uint32(id)
		for k, v := range configs {
			if k == realID {
				return &v
			}
		}

		return nil
	},
}
