package users

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionSate = &base_controller.ActionState{
	Model: models.Users,
	Field: "status",
	StateAfter: func(c *gin.Context) {
		ids := c.Query("id")
		key := consts.CtxKeyLoginUser + ids
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		redis.Del(key)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select username from users where id=" + ids
		res, _ := dbSession.QueryString(sql)
		tokenKey := consts.CacheKeyUserNameToken + res[0]["username"]
		token, _ := redis.Get(tokenKey).Result()
		ss := consts.CacheKeyTokenPrefix + token
		redis.Del(ss)       // 删除用户信息
		redis.Del(tokenKey) // 删除用户缓存信息
	},
}
