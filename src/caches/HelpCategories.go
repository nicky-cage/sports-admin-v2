package caches

import (
	models "sports-models"
)

const keyHelpCategories = "help_categories"

// HelpCategories 配置信息
var HelpCategories = struct {
	Load func(string)
	Get  func(string, int) *models.HelpCategory
	All  func(string) map[uint32]models.HelpCategory
}{
	Load: func(platform string) {
		helpCategories := map[uint32]models.HelpCategory{}
		var rs []models.HelpCategory
		err := models.HelpCategories.FindAllNoCount(platform, &rs, nil, "id ASC")
		if err != nil {
			return
		}
		for _, r := range rs {
			helpCategories[r.Id] = r
		}

		_ = setCache(platform, keyHelpCategories, helpCategories)
	},
	Get: func(platform string, id int) *models.HelpCategory {
		helpCategories := map[uint32]models.HelpCategory{}
		_ = getCache(platform, keyHelpCategories, &helpCategories)

		realID := uint32(id)
		for k, v := range helpCategories {
			if k == realID {
				return &v
			}
		}

		return nil
	},
	All: func(platform string) map[uint32]models.HelpCategory {
		helpCategories := map[uint32]models.HelpCategory{}
		_ = getCache(platform, keyHelpCategories, &helpCategories)
		return helpCategories
	},
}
