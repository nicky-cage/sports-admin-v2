package caches

import (
	models "sports-models"

	"xorm.io/builder"
)

// 所有用户标签分类
const keyUserTagCategories = "user_tag_categories"

// UserTagCategories 用户标签分类
var UserTagCategories = struct {
	Load func(string)
	All  func(string) map[uint32]models.UserTagCategory
	Get  func(string, int) *models.UserTagCategory
}{
	Load: func(platform string) {
		userTagCategories := map[uint32]models.UserTagCategory{}
		var rs []models.UserTagCategory
		err := models.UserTagCategories.FindAllNoCount(platform, &rs, nil, "id ASC")
		if err != nil {
			return
		}
		for _, r := range rs {
			var tags []models.UserTag
			cond := builder.NewCond().And(builder.Eq{"category_id": r.Id})
			_ = models.UserTags.FindAllNoCount(platform, &tags, cond)
			r.Tags = tags
			userTagCategories[r.Id] = r
		}
		_ = setCache(platform, keyUserTagCategories, userTagCategories)
	},
	All: func(platform string) map[uint32]models.UserTagCategory {
		userTagCategories := map[uint32]models.UserTagCategory{}
		_ = getCache(platform, keyUserTagCategories, &userTagCategories)
		return userTagCategories
	},
	Get: func(platform string, id int) *models.UserTagCategory {
		userTagCategories := map[uint32]models.UserTagCategory{}
		_ = getCache(platform, keyUserTagCategories, &userTagCategories)

		realID := uint32(id)
		for k, v := range userTagCategories {
			if k == realID {
				return &v
			}
		}
		return nil
	},
}
