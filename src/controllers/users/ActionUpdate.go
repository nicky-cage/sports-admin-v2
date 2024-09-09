package users

import (
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	models "sports-models"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model:    models.Users,
	ViewFile: "users/edit.html",
	Row: func() interface{} {
		return &models.User{}
	},
	ExtendData: func(c *gin.Context) pongo2.Context {
		var balance string
		id := c.Query("id")
		sql := "select balance from accounts where user_id=" + id
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
		}
		if len(res) > 0 {
			balance = res[0]["balance"]
		}
		var pRes []models.Province
		perr := dbSession.Table("provinces").Find(&pRes)
		if perr != nil {
			log.Err(perr.Error())
		}

		var cRes []models.City
		cerr := dbSession.Table("cities").Find(&cRes)
		if cerr != nil {
			log.Err(cerr.Error())
		}

		var dRes []models.District
		derr := dbSession.Table("districts").Find(&dRes)
		if derr != nil {
			log.Err(derr.Error())
		}
		caches.UserTagCategories.Load(platform)

		var paymentGroups = []models.PaymentGroup{}
		_ = models.PaymentGroups.FindAllNoCount(platform, &paymentGroups)

		return pongo2.Context{
			"vipLevels":     caches.UserLevels.All(platform),
			"tagCategories": caches.UserTagCategories.All(platform),
			"balance":       balance,
			"provinces":     pRes,
			"city":          cRes,
			"districts":     dRes,
			"paymentGroups": paymentGroups,
		}
	},
}
