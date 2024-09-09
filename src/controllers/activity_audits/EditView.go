package activity_audits

import (
	"sports-common/log"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *ActivityAudits) EditView(c *gin.Context) {
	idStr, exists := c.GetQuery("id")
	if !exists || idStr == "" {
		log.Err("无法获取id信息!\n")
		return
	}
	typeStr, _ := c.GetQuery("type")
	data := make(map[string]string)
	data["id"] = idStr
	viewData := pongo2.Context{"r": data}
	if typeStr == "agree" {
		response.Render(c, "dividend_managements/view_agree.html", viewData)
	} else if typeStr == "refuse" {
		response.Render(c, "dividend_managements/view_refuse.html", viewData)
	} else if typeStr == "batch_agree" {
		response.Render(c, "dividend_managements/view_batch_agree.html", viewData)
	} else {
		response.Render(c, "dividend_managements/view_batch_refuse.html", viewData)
	}
}
