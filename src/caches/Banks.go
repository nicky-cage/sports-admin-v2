package caches

import (
	"sort"
	models "sports-models"

	"xorm.io/builder"
)

const keyBanks = "banks"

// Banks 银行信息
var Banks = struct {
	Load func(string)
	All  func(string) []models.Bank
	Get  func(string, int) *models.Bank
}{
	Load: func(platform string) {
		banks := map[uint32]models.Bank{}
		var rs []models.Bank

		cond := builder.NewCond().And(builder.Eq{"state": 2})
		err := models.Banks.FindAllNoCount(platform, &rs, cond, "id ASC")
		if err != nil {
			return
		}

		for _, r := range rs {
			banks[r.Id] = r
		}

		_ = setCache(platform, keyBanks, banks)
	},
	All: func(platform string) []models.Bank {
		banks := map[uint32]models.Bank{}
		_ = getCache(platform, keyBanks, &banks)
		var ids []int
		for k := range banks {
			ids = append(ids, int(k))
		}
		sort.Ints(ids)

		var result []models.Bank
		for _, k := range ids {
			result = append(result, banks[uint32(k)])
		}

		return result
	},
	Get: func(platform string, id int) *models.Bank {
		banks := map[uint32]models.Bank{}
		_ = getCache(platform, keyBanks, &banks)
		realID := uint32(id)
		for k, v := range banks {
			if k == realID {
				return &v
			}
		}

		return nil
	},
}
