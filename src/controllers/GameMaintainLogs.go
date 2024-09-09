package controllers

import (
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// GameMaintainLogs 维护设置日志
var GameMaintainLogs = struct {
	*ActionList
}{
	ActionList: &ActionList{
		Model:    models.GameVenueLogs,
		ViewFile: "game_maintain_logs/list.html",
		OrderBy: func(*gin.Context) string {
			return "created DESC"
		},
		Rows: func() interface{} {
			return &[]models.GameVenueLog{}
		},
		QueryCond: map[string]interface{}{
			"admin": "%",
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			cond := builder.NewCond()
			if idStr, exists := c.GetQuery("id"); exists {
				if id, err := strconv.Atoi(idStr); err == nil {
					cond = cond.And(builder.Eq{"venue_id": id})
				}
			}
			return cond
		},
		ExtendData: func(c *gin.Context) ViewData {
			return ViewData{
				"id": c.DefaultQuery("id", "0"),
			}
		},
	},
}
