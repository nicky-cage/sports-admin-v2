package controllers

import (
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// SiteMaintainLogs 维护设置日志
var SiteMaintainLogs = struct {
	*ActionList
}{
	ActionList: &ActionList{
		Model:    models.SiteMaintainLogs,
		ViewFile: "site_maintain_logs/list.html",
		//OrderBy: func(*gin.Context) string {
		//	return "id DESC"
		//},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			cond := builder.NewCond()
			if idStr, exists := c.GetQuery("id"); exists {
				if id, err := strconv.Atoi(idStr); err == nil {
					cond = cond.And(builder.Eq{"platform_id": id})
				}
			}
			return cond
		},
		Rows: func() interface{} {
			return &[]models.SiteMaintainLog{}
		},
	},
}
