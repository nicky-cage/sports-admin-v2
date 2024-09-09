package users

import (
	common "sports-common"
	"sports-common/consts"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// DisableAll 禁用用户
func (ths *Users) DisableAll(c *gin.Context) {
	postedData := request.GetPostedData(c)
	if _, exists := postedData["ids"]; !exists {
		response.Err(c, "缺少用户编号")
		return
	}
	ids := postedData["ids"].(string)
	platform := request.GetPlatform(c)
	strIds := strings.Split(ids, ",") // 拆分用户编号
	if len(strIds) == 0 {
		response.Err(c, "批量禁用的用户信息有误")
		return
	}

	// redis
	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	for _, strId := range strIds {
		userId, err := strconv.Atoi(strId)
		if err != nil { // 如果转换为整型出错
			response.Err(c, "用户编号信息有误")
			return
		}
		var user models.User
		if exists, err := models.Users.FindById(platform, userId, &user); err != nil || !exists {
			response.Err(c, "缺少用户相关信息")
			return
		}
		data := map[string]interface{}{
			"id":     userId,
			"status": 1, // 1:禁用; 2:正常;
		}
		if err := models.Users.Update(platform, data); err != nil {
			response.Err(c, "修改用户状态失败")
			return
		}

		// 从缓存服务器删除用户token
		cacheKey := consts.CacheKeyUserNameToken + user.Username // 删除用户缓存信息
		rd.Del(cacheKey)                                         // 删除用户token
	}
	response.Ok(c)
}
