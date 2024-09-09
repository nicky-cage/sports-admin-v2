package controllers

import (
	"sports-admin/caches"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// Games 游戏列表
var Games = struct {
	*ActionList
	*ActionUpdate
	*ActionSave
	*ActionState
}{
	ActionList: &ActionList{
		Model:    models.Games,
		ViewFile: "games/list.html",
		OrderBy: func(*gin.Context) string {
			return "sort DESC"
		},
		Rows: func() interface{} {
			return &[]models.Game{}
		},
		QueryCond: map[string]interface{}{
			"en_name":        "=",
			"venue_type":     "=",
			"platform_types": "%",
			"is_online":      "=",
			"game_code":      "=",
			"venue_id":       "=",
			"cn_name":        "=",
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Games,
		ViewFile: "games/edit.html",
		Row: func() interface{} {
			return &models.Game{}
		},
		ExtendData: func(*gin.Context) pongo2.Context {
			return pongo2.Context{
				"venue_types": caches.GameVenues.All,
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Games,
		SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
			id := (*m)["id"].(string)
			sql := "select web_code,game_code from games where id=" + id
			platform := request.GetPlatform(c)
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			res, err := dbSession.QueryString(sql)
			if err != nil {
				log.Err(err.Error())
			}
			gameCode := res[0]["web_code"]
			redis := common.Redis(platform)
			defer common.RedisRestore(platform, redis)
			cacheKey := consts.EgameCacheKey + gameCode
			gameKey := consts.EgameCacheKey + res[0]["game_code"]
			redis.Del(cacheKey)
			redis.Del(gameKey)
		},
	},
	ActionState: &ActionState{
		Model: models.Games,
		Field: "is_online",
		StateAfter: func(c *gin.Context) {
			id := c.Query("id")
			sql := "select web_code,game_code from games where id=" + id
			platform := request.GetPlatform(c)
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			res, err := dbSession.QueryString(sql)
			if err != nil {
				log.Err(err.Error())
			}
			gameCode := res[0]["web_code"]

			redis := common.Redis(platform)
			defer common.RedisRestore(platform, redis)
			cacheKey := consts.EgameCacheKey + gameCode
			agKey := consts.EgameCacheKey + res[0]["game_code"]
			redis.Del(cacheKey)
			redis.Del(agKey)
		},
	},
}
