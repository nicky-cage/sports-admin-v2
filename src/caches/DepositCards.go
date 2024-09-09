package caches

import (
	"sort"
	models "sports-models"
)

const keyDepositCards = "deposit_cards"

// DepositCards 银行信息
var DepositCards = struct {
	Load func(string)
	All  func(string) []models.DepositCard
	Get  func(string, int) *models.DepositCard
}{
	Load: func(platform string) {
		depositCards := map[uint32]models.DepositCard{}
		var rs []models.DepositCard

		err := models.DepositCards.FindAllNoCount(platform, &rs, nil, "id ASC")
		if err != nil {
			return
		}

		for _, r := range rs {
			depositCards[r.Id] = r
		}

		_ = setCache(platform, keyDepositCards, depositCards)
	},
	All: func(platform string) []models.DepositCard {
		depositCards := map[uint32]models.DepositCard{}
		_ = getCache(platform, keyDepositCards, &depositCards)

		var ids []int
		for k := range depositCards {
			ids = append(ids, int(k))
		}
		sort.Ints(ids)

		var result []models.DepositCard
		for _, k := range ids {
			result = append(result, depositCards[uint32(k)])
		}

		return result
	},
	Get: func(platform string, id int) *models.DepositCard {
		depositCards := map[uint32]models.DepositCard{}
		_ = getCache(platform, keyDepositCards, &depositCards)

		realID := uint32(id)
		for k, v := range depositCards {
			if k == realID {
				return &v
			}
		}

		return nil
	},
}
