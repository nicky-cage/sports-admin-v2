package activity_managements

import (
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *ActivityManagements) Index(c *gin.Context) { //默认首页
	downExcelFile := config.Get("internal.img_host_backend", "") + "/uploads/Excel/dividend.xlsx"
	uploadExcelFile := config.Get("internal.img_host_backend", "")
	gameVenues := make([]models.GameVenue, 0)
	platform := request.GetPlatform(c)
	engine := common.Mysql(platform)
	defer engine.Close()
	if err := engine.Table("game_venues").Where("pid=? and maintain=? and code!=?", 0, 1, "CENTERWALLET").Find(&gameVenues); err != nil {
		log.Logger.Error(err.Error())
	}
	response.Render(c, "dividend_managements/index.html", pongo2.Context{
		"down_excel_url":    downExcelFile,
		"upload_excel_file": uploadExcelFile,
		"game_venus":        gameVenues,
	})
}
