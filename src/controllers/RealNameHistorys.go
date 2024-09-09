package controllers

import (
	"github.com/gin-gonic/gin"
	models "sports-models"
	"xorm.io/builder"
)

var RealNameHistories = struct {
	*ActionList
}{
	ActionList: &ActionList{
		Model: models.UserBanksNameLogs,
		Rows: func() interface{} {
			return &[]models.UserBanksNameLog{}
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			id := c.Query("id")
			cond := builder.NewCond()
			cond = cond.And(builder.Eq{"user_id": id})
			return cond
		},
		ViewFile: "users/real_name_histories.html",
	},
}
