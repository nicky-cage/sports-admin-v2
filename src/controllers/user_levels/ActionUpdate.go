package user_levels

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model:    models.UserLevels,
	ViewFile: "user_levels/edit.html",
	Row: func() interface{} {
		return &models.UserLevel{}
	},
	ExtendData: func(c *gin.Context) pongo2.Context {
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		var venue []models.GameVenue
		err := dbSession.Table("game_venues").Where("pid != ? and is_online=1", 0).Find(&venue)
		if err != nil {
			log.Err(err.Error())
		}
		return pongo2.Context{"res": venue}
	},
}
