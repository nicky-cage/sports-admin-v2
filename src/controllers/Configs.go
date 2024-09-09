package controllers

import (
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/consts"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// Configs 配置信息
var Configs = struct {
	*ActionUpdate
	*ActionSave
}{
	ActionUpdate: &ActionUpdate{
		Model:    models.Configs,
		ViewFile: "sites/index.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			// res := []models.Config{}
			// _, err := models.Configs.FindAll(platform, &res)
			// if err != nil {
			// 	fmt.Println("获取配置信息失败")
			// 	log.Err(err.Error())
			// }

			allowMulDep := "1"
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			res, err := dbSession.QueryString("select * from deposit_limits_set LIMIT 1")
			if err == nil || len(res) > 0 {
				allowMulDep = res[0]["allow"]
			}

			var rows []models.HelpSimple
			err = models.HelpSimples.FindAllNoCount(platform, &rows)
			if err != nil {
				fmt.Println("修改帮助文章出错:", err)
			}
			return pongo2.Context{
				"help_categories": caches.HelpCategories.All(platform),
				"bottomTypes":     consts.BottomTypes,
				"customer_list":   consts.CustomerList,
				"channel_list":    consts.ChannelList,
				"rows":            rows,
				"venue_types":     consts.VenueTypes,
				"allow":           allowMulDep,
				//"t":               res[0],
			}
		},
		Row: func() interface{} {
			return &models.Config{}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Configs,
	},
}
