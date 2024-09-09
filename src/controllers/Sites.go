package controllers

import (
	"sports-admin/caches"
	"sports-common/consts"
	"sports-common/request"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// Sites 底部信息
var Sites = struct {
	Index func(*gin.Context)
	*ActionUpdate
}{
	Index: func(c *gin.Context) { //默认首页
		platform := request.GetPlatform(c)
		viewData := pongo2.Context{
			"help_categories": caches.HelpCategories.All(platform),
			"bottomTypes":     consts.BottomTypes,
		}
		response.Render(c, "sites/index.html", viewData)
	},
	ActionUpdate: &ActionUpdate{
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			r := caches.Configs.Get(platform, 1)
			return pongo2.Context{
				"r": r,
			}
		},
		ViewFile: "configs/edit.html",
	},
}
