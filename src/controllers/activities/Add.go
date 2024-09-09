package activities

import (
	"sports-common/config"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (*Activities) Add(c *gin.Context) {
	platform := request.GetPlatform(c)
	gameVenus := &[]models.GameVenue{}
	_ = models.GameVenues.FindAllNoCount(platform, gameVenus, builder.NewCond().And(builder.Neq{"venue_type": 0}))
	vips := &[]models.UserLevel{}
	_ = models.UserLevels.FindAllNoCount(platform, vips)
	viewData := pongo2.Context{
		"gameVenus":     gameVenus,
		"vips":          vips,
		"activity":      &models.Activity{},
		"STATIC_URL":    config.Get("internal.img_host_backend", ""),
		"activity_code": "CODE_" + strings.ToUpper(tools.RandString(12)),
	}
	viewFile := "activities/edit.html"
	response.Render(c, viewFile, viewData)
}
