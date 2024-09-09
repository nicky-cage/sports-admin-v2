package controllers

import (
	"github.com/gin-gonic/gin"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"xorm.io/builder"
)

var SourceMaterials = struct {
	*ActionDelete
	*ActionList
	*ActionUpdate
	*ActionSave
	*ActionCreate
	*ActionState
}{
	ActionState: &ActionState{
		Model: models.SourceMaterials,
	},
	ActionSave: &ActionSave{
		Model: models.SourceMaterials,
	},
	ActionList: &ActionList{
		Model:    models.SourceMaterials,
		ViewFile: "source_materials/list.html",
		Rows: func() interface{} {
			return &[]models.SourceMaterial{}
		},
		QueryCond: map[string]interface{}{
			"state": "=",
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			cond := builder.NewCond()
			var startAt int64
			var endAt int64
			if value, exists := c.GetQuery("created"); exists {
				areas := strings.Split(value, " - ")
				startAt = tools.GetTimeStampByString(areas[0])
				endAt = tools.GetTimeStampByString(areas[1])
				cond = cond.And(builder.Gte{"time_start": startAt}).And(builder.Lte{"time_end": endAt})

			}
			return cond
		},
		OrderBy: func(c *gin.Context) string {
			return "sort ASC"
		},
	},
	ActionUpdate: &ActionUpdate{
		Model: models.SourceMaterials,
		Row: func() interface{} {
			return &models.SourceMaterial{}
		},
		ViewFile: "source_materials/updated.html",
	},
	ActionDelete: &ActionDelete{
		Model: models.SourceMaterials,
	},
	ActionCreate: &ActionCreate{
		ViewFile: "source_materials/created.html",
	},
}
