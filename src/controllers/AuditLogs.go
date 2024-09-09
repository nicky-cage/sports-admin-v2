package controllers

import (
	"github.com/gin-gonic/gin"
	models "sports-models"
	"xorm.io/builder"
)

var AuditLogs = struct {
	*ActionList
}{
	ActionList: &ActionList{
		Model: models.AuditLogs,
		OrderBy: func(c *gin.Context) string {
			return "created desc"
		},
		Rows: func() interface{} {
			return &[]models.AuditLog{}
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			cond := builder.NewCond()
			value := c.Query("type")
			cond = cond.And(builder.Eq{"type": value})
			return cond
		},
		ViewFile: "users/audit_logs.html",
	},
}
