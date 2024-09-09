package index

import (
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var Test = func(c *gin.Context) {
	viewData := pongo2.Context{}
	provinces := &[]models.Province{}
	platform := request.GetPlatform(c)
	if err := models.Provinces.FindAllNoCount(platform, provinces); err == nil {
		for k, p := range *provinces {
			cond := builder.NewCond().And(builder.Eq{"province_code": p.Code})
			cities := &[]models.City{}
			if err := models.Cities.FindAllNoCount(platform, cities, cond); err == nil {
				for ck, c := range *cities {
					cond := builder.NewCond().And(builder.Eq{"city_code": c.Code})
					districts := &[]models.District{}
					_ = models.Districts.FindAllNoCount(platform, districts, cond)
					(*cities)[ck].Districts = *districts
				}
			}
			(*provinces)[k].Cities = *cities
		}
	}

	viewData["provinces"] = provinces
	response.Render(c, "index/test.html", viewData)
}
