package controllers

import (
	"github.com/gin-gonic/gin"
	models "sports-models"
)

var ConfigPopulars = struct {
	List func(c *gin.Context)
	*ActionSave
}{
	List: func(c *gin.Context) {
		//内容再  config/updated
	},
	ActionSave: &ActionSave{
		Model: models.Configs,
	},
}
