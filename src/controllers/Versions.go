package controllers

import (
	"sports-common/log"
	"sports-common/request"
	models "sports-models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

// downloadConfig
var downloadConfig = map[string]map[string]string{
	"shipu": {
		"branch_key":    "key_live_anZakG4YMZ6qAAMiwLrpLjjoAtde0A8v",
		"branch_secret": "secret_live_7Y2Rx2jsqaqz639is6rG6nizZfkqUVOH",
	},
	"venice": {
		"branch_key":    "key_live_kp1jZLjG82gT6EFWZHL1NkfmAtooInzJ",
		"branch_secret": "secret_live_87W01cO2VPKMHFagFxwaOzl2mpPKyung",
	},
	"xingkong": {
		"branch_key":    "key_live_ci1vfP9H1tAenHxkV1IZhfldACoUNU2S",
		"branch_secret": "secret_live_LPZnjGfqf11jxAB65Lu6YWXz846iK2Xd",
	},
}

// Versions 版本处理
var Versions = struct {
	*ActionList
	*ActionSave
	*ActionUpdate
	*ActionCreate
	*ActionDelete
}{
	ActionList: &ActionList{
		Model: models.SiteVersions,
		Rows: func() interface{} {
			return &[]models.SiteVersion{}
		},
		ViewFile: "site_versions/list.html",
	},
	ActionSave: &ActionSave{
		Model: models.SiteVersions,
		SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
			platform := request.GetPlatform(c)
			postData := request.GetPostedData(c) //获取IOS或者安卓Url
			url := "https://api2.branch.io//v1/app/" + downloadConfig[platform]["branch_key"]
			req.SetTimeout(30 * time.Second)
			req.Debug = false
			headerB := req.Header{
				"Accept": "application/json",
			}
			paramB := req.Param{
				"branch_secret": downloadConfig[platform]["branch_secret"],
				"dev_email":     "brandoneedyou@gmail.com",
			}
			if postData["app_type"] == "1" {
				paramB["android_url"] = postData["updated_url"]
			} else {
				paramB["ios_url"] = postData["updated_url"]
			}
			_, err := req.Put(url, headerB, req.BodyJSON(paramB))
			if err != nil {
				log.Err(err.Error())
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.SiteVersions,
		ViewFile: "site_versions/updated.html",
		Row: func() interface{} {
			return &models.SiteVersion{}
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.SiteVersions,
		ViewFile: "site_versions/add.html",
	},
	ActionDelete: &ActionDelete{
		Model: models.SiteVersions,
	},
}
