package activities

import (
	common "sports-common"
	"sports-common/config"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (*Activities) Edit(c *gin.Context) {
	idStr := c.DefaultQuery("id", "")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		response.Err(c, "id不正确")
		return
	}
	platform := request.GetPlatform(c)
	gameVenus := &[]models.GameVenue{}
	_ = models.GameVenues.FindAllNoCount(platform, gameVenus, builder.NewCond().And(builder.Neq{"venue_type": 0}))
	vips := &[]models.UserLevel{}
	_ = models.UserLevels.FindAllNoCount(platform, vips)
	activity := &models.Activity{}
	_, _ = models.Activities.FindById(platform, id, activity)

	myClient := common.Mysql(request.GetPlatform(c))
	defer myClient.Close()

	newGameCodeList := make([]string, 0)
	if len(activity.GameCodeList) > 0 {
		gameCodeList := strings.Split(activity.GameCodeList, ",")
		for _, v := range gameCodeList {
			gameVenusInfo := &models.GameVenue{}
			temp := strings.Split(v, "-")
			code := temp[0]
			venueType, _ := strconv.Atoi(temp[1])
			_, _ = myClient.Table("game_venues").Where("code=? and venue_type=?", code, venueType).Get(gameVenusInfo)
			newGameCodeList = append(newGameCodeList, strconv.Itoa(int(gameVenusInfo.Id)))
		}
	}
	gameCodeListStr := strings.Join(newGameCodeList, ",")
	viewData := pongo2.Context{
		"gameVenus":       gameVenus,
		"vips":            vips,
		"activity":        activity,
		"gameCodeListStr": gameCodeListStr,
		"STATIC_URL":      config.Get("internal.img_host_backend", ""),
		"activity_code":   "CODE_" + strings.ToUpper(tools.RandString(12)),
	}
	viewFile := "activities/edit.html"
	response.Render(c, viewFile, viewData)
}
