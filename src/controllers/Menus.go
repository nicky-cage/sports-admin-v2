package controllers

import (
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// Menus 后台角色菜单
var Menus = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
}{
	ActionList: &ActionList{
		Model:    models.Menus,
		ViewFile: "menus/list.html",
		OrderBy: func(c *gin.Context) string {
			return "sort DESC"
		},
		Rows: func() interface{} {
			return &[]models.Menu{}
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			return pongo2.Context{
				"menuLevelIds": []uint8{1, 2, 3, 4, 5, 6},
			}
		},
		QueryCond: map[string]interface{}{
			"name":  "%",
			"level": "=",
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Menus,
		ViewFile: "menus/edit.html",
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Menus,
		ViewFile: "menus/edit.html",
		Row: func() interface{} {
			return &models.Menu{}
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			return pongo2.Context{
				"menuLevelIds": []uint8{1, 2, 3, 4, 5, 6},
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Menus,
	},
	ActionDelete: &ActionDelete{
		Model: models.Menus,
	},
}
