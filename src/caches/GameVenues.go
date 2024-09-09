package caches

import (
	models "sports-models"

	"xorm.io/builder"
)

const keyGameVenues = "game_venues"

// GameVenues 游戏场馆信息
var GameVenues = struct {
	Load func(string)
	All  func(string) map[uint32]models.GameVenue
}{
	Load: func(platform string) {
		gameVenues := map[uint32]models.GameVenue{}
		var rs []models.GameVenue
		cond := builder.NewCond().And(builder.Eq{"pid": 0}).And(builder.Eq{"is_online": 1})
		err := models.GameVenues.FindAllNoCount(platform, &rs, cond, "id DESC")
		if err != nil {
			return
		}

		for _, r := range rs {
			gameVenues[r.Id] = r
		}

		setCache(platform, keyGameVenues, gameVenues)
	},
	All: func(platform string) map[uint32]models.GameVenue {
		gameVenues := map[uint32]models.GameVenue{}
		_ = getCache(platform, keyGameVenues, &gameVenues)
		return gameVenues
	},
}
