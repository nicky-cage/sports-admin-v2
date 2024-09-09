package controllers

import (
	"sports-admin/caches"
	"sports-common/consts"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// PlatformSites 平台站点
var PlatformSites = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionState
	Config     func(*gin.Context)
	ConfigSave func(*gin.Context)
}{
	ActionList: &ActionList{
		Model:    models.Sites,
		ViewFile: "platform_sites/list.html",
		Rows: func() interface{} {
			return &[]models.Site{}
		},
		QueryCond: map[string]interface{}{
			"platform_id": "=",
			"name":        "%",
			"remark":      "%",
			"urls":        "%",
			"code":        "%",
		},
		ExtendData: func(c *gin.Context) ViewData {
			return ViewData{
				"platforms": caches.Platforms.All(),
			}
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Sites,
		ViewFile: "platform_sites/edit.html",
		ExtendData: func(c *gin.Context) ViewData {
			return ViewData{
				"platforms": caches.Platforms.All(),
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Sites,
		ViewFile: "platform_sites/edit.html",
		Row: func() interface{} {
			return &models.Site{}
		},
		ExtendData: func(c *gin.Context) ViewData {
			return ViewData{
				"platforms": caches.Platforms.All(),
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Sites,
	},
	ActionState: &ActionState{
		Model: models.Sites,
		Field: "status",
	},
	Config: func(c *gin.Context) {
		idStr := c.DefaultQuery("id", "0")
		ID, err := strconv.Atoi(idStr)
		if err != nil || ID <= 0 {
			response.Err(c, "站点信息编号有误")
			return
		}

		platform := consts.PlatformIntegrated
		r := models.Site{}
		exists, err := models.Sites.Find(platform, &r, builder.NewCond().And(builder.Eq{"id": ID}))
		if err != nil || !exists {
			response.Err(c, "查询站点信息失败")
			return
		}

		conf := models.Sites.GetConfigsBySiteID(platform, ID)
		viewData := ViewData{
			"platformID": r.PlatformId,
			"siteID":     ID,
			"conf":       conf,
		}
		response.Render(c, "platform_sites/config.html", viewData)
	},
	ConfigSave: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		platformIDStr := postedData["platform_id"].(string)
		platformID, err := strconv.Atoi(platformIDStr)
		if err != nil {
			response.Err(c, "平台编号有误")
			return
		}

		siteIDStr := postedData["site_id"].(string)
		siteID, err := strconv.Atoi(siteIDStr)
		if err != nil {
			response.Err(c, "站点编号信息有误")
			return
		}

		platform := consts.PlatformIntegrated
		conf := &models.PlatformSiteConf{}
		conf.PlatformID = platformID
		conf.SiteName = postedData["site_name"].(string)
		conf.Platform = postedData["platform"].(string)
		conf.ConnStrings = postedData["conn_strings"].(string)
		conf.KafkaStrings = postedData["kafka_strings"].(string)
		conf.ElasticStrings = postedData["elastic_strings"].(string)
		conf.RedisStrings = postedData["redis_strings"].(string)
		conf.PayStrings = postedData["pay_strings"].(string)
		conf.StaticURL = postedData["static_url"].(string)
		conf.UploadURL = postedData["upload_url"].(string)

		if err = models.Sites.SetConfigsBySiteID(platform, siteID, conf); err != nil {
			response.Err(c, "保存站点配置失败")
			return
		}

		response.Ok(c)
	},
}
