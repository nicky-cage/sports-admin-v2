package users

import (
	"sports-admin/caches"
	"sports-common/request"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// AddTags 增加标签
func (ths *Users) AddTags(c *gin.Context) {
	platform := request.GetPlatform(c)
	response.Render(c, "users/add_tags.html", pongo2.Context{
		"tagCategories": caches.UserTagCategories.All(platform),
		"userIds":       c.Query("ids"),
	})
}
