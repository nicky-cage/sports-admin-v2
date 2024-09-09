package controllers

import (
	"sports-common/consts"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// SiteBottoms 底部信息
var SiteBottoms = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
	*ActionState
}{
	ActionList: &ActionList{
		Model: models.SiteBottoms,
		Rows: func() interface{} {
			return &[]models.SiteBottom{}
		},
		QueryCond: map[string]interface{}{
			"content_type": "=",
			"bottom_type":  "=",
		},
		ViewFile: "site_bottoms/list.html",
	},
	ActionCreate: &ActionCreate{
		Model:    &models.SiteBottom{},
		ViewFile: "site_bottoms/edit.html",
		ExtendData: func(*gin.Context) pongo2.Context {
			return pongo2.Context{
				"bottomTypes": consts.BottomTypes,
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model: models.SiteBottoms,
		Row: func() interface{} {
			return &models.SiteBottom{}
		},
		ExtendData: func(*gin.Context) pongo2.Context {
			return pongo2.Context{
				"bottomTypes": consts.BottomTypes,
			}
		},
		ViewFile: "site_bottoms/edit.html",
	},
	ActionSave: &ActionSave{
		Model: models.SiteBottoms,
	},
	ActionDelete: &ActionDelete{
		Model: models.SiteBottoms,
	},

	ActionState: &ActionState{
		Model: models.SiteBottoms,
	},
}
