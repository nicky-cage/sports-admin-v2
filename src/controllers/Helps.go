package controllers

import (
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

type HelpsTest struct {
	Playname  string `json:"playname"`
	BetMoney  int    `json:"bet_money"`
	NetMoney  int    `json:"net_money"`
	UserID    int    `json:"user_id"`
	GameType  int    `json:"game_type"`
	CreatedAt uint32 `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

// Helps 帮助分类
var Helps = struct {
	*ActionList
	*ActionUpdate
	*ActionCreate
	*ActionSave
	*ActionDelete
	*ActionState
	Check  func(c *gin.Context)
	Detail func(c *gin.Context)
}{
	ActionCreate: &ActionCreate{
		Model:    models.HelpCategories,
		ViewFile: "helps/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"help_categories": caches.HelpCategories.All(platform),
				"venue_types":     consts.VenueTypes,
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Helps,
	},
	ActionList: &ActionList{
		Model:    models.Helps,
		ViewFile: "helps/list.html",
		OrderBy: func(*gin.Context) string {
			return "sort DESC"
		},
		Rows: func() interface{} {
			return &[]models.Help{}
		},
		QueryCond: map[string]interface{}{
			"name": "%",
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"venue_types":     consts.VenueTypes,
				"help_categories": caches.HelpCategories.All(platform),
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Helps,
		ViewFile: "helps/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			admin := GetLoginAdmin(c)
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"help_categories": caches.HelpCategories.All(platform),
				"venue_types":     consts.VenueTypes,
				"admin":           admin.Name,
			}
		},
		Row: func() interface{} {
			return &models.Help{}
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.Helps,
	},
	ActionState: &ActionState{
		Model: models.Helps,
	},
	Check: func(c *gin.Context) {
		var rows []models.HelpSimple
		platform := request.GetPlatform(c)
		if categoryIDStr, exists := c.GetQuery("category"); exists {
			if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
				_ = models.HelpSimples.FindAllNoCount(platform, &rows, builder.NewCond().And(builder.Eq{"category_id": categoryID}))
				response.Render(c, "helps/_list.html", ViewData{"rows": rows, "venue_types": consts.VenueTypes})
				return
			}
		}
		_ = models.HelpSimples.FindAllNoCount(platform, &rows)
		response.Render(c, "helps/_list.html", ViewData{"rows": rows, "venue_types": consts.VenueTypes})
	},
	Detail: func(c *gin.Context) {
		id := c.Query("id")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from helps where id='%s'"
		sqll := fmt.Sprintf(sql, id)
		res, err := dbSession.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
			return
		}
		ViewData := pongo2.Context{
			"r": res[0],
		}

		response.Render(c, "helps/detail.html", ViewData)
	},
}
