package caches

import (
	"sort"
	models "sports-models"
)

const keyUserLevels = "user_levels"

// UserLevels 用户等级
var UserLevels = struct {
	Load func(string)
	All  func(string) []models.UserLevel
	Get  func(string, int) *models.UserLevel
}{
	Load: func(platform string) {
		userLevels := map[uint8]models.UserLevel{}

		var rs []models.UserLevel
		err := models.UserLevels.FindAllNoCount(platform, &rs, nil, "id ASC")
		if err != nil {
			return
		}

		for _, r := range rs {
			userLevels[r.Digit] = r
		}

		_ = setCache(platform, keyUserLevels, userLevels)
	},
	All: func(platform string) []models.UserLevel {
		userLevels := map[uint8]models.UserLevel{}
		_ = getCache(platform, keyUserLevels, &userLevels)

		var ids []int
		for k := range userLevels {
			ids = append(ids, int(k))
		}
		sort.Ints(ids)

		var result []models.UserLevel
		for _, k := range ids {
			result = append(result, userLevels[uint8(k)])
		}

		return result
	},
	Get: func(platform string, id int) *models.UserLevel {
		userLevels := map[uint8]models.UserLevel{}
		_ = getCache(platform, keyUserLevels, &userLevels)

		realID := uint8(id)
		for k, v := range userLevels {
			if k == realID {
				return &v
			}
		}

		return nil
	},
}
