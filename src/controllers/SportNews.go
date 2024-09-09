package controllers

import (
	"github.com/gin-gonic/gin"
	models "sports-models"
	"xorm.io/builder"
)

var SportNews = struct {
	*ActionSave
	*ActionList
	*ActionUpdate
	*ActionDelete
	*ActionCreate
}{
	ActionDelete: &ActionDelete{
		Model: models.SportsNews,
	},
	ActionCreate: &ActionCreate{
		ViewFile: "sport_news/created.html",
	},
	ActionSave: &ActionSave{
		Model: models.SportsNews,
	},
	ActionUpdate: &ActionUpdate{
		Model: models.SportsNews,
		Row: func() interface{} {
			return &models.SportNews{}
		},
		ViewFile: "sport_news/created.html",
	},
	ActionList: &ActionList{
		Model:    models.SportsNews,
		ViewFile: "sport_news/list.html",
		Rows: func() interface{} {
			return &[]models.SportNews{}
		},
		QueryCond: map[string]interface{}{
			"title":     "%",
			"author":    "%",
			"item_type": "=",
			"is_hot":    "=",
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			//
			cond := builder.NewCond()
			return cond
		},
	},
}
