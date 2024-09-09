package caches

import (
	"sports-common/consts"
	models "sports-models"
	"strings"

	"github.com/gin-gonic/gin"
)

const keyPlatformSites = "platform_sites"

// PlatformSites 平台站点
var PlatformSites = struct {
	Load       func()
	Get        func(int) *models.Site
	GetCurrent func(c *gin.Context) *models.Site
	All        func() map[uint32]models.Site
}{
	Load: func() {
		rows := map[uint32]models.Site{}
		var rs []models.Site
		err := models.Sites.FindAllNoCount(consts.PlatformIntegrated, &rs, nil, "id ASC")
		if err != nil {
			return
		}
		for _, r := range rs {
			rows[r.Id] = r
		}

		_ = setCache(consts.PlatformIntegrated, keyPlatformSites, rows)
	},
	Get: func(id int) *models.Site {
		rows := map[uint32]models.Site{}
		_ = getCache(consts.PlatformIntegrated, keyPlatformSites, &rows)

		realID := uint32(id)
		for k, v := range rows {
			if k == realID {
				return &v
			}
		}

		return nil
	},
	GetCurrent: func(c *gin.Context) *models.Site {
		hostName := c.Request.Host
		rows := map[uint32]models.Site{}
		err := getCache(consts.PlatformIntegrated, keyPlatformSites, &rows)
		if err != nil {
			return nil
		}

		row := models.Site{}
		for _, r := range rows {
			domainNames := strings.Split(r.AdminUrl, ",")
			for _, d := range domainNames {
				domain := strings.TrimSpace(d)
				if domain == hostName {
					row = r
				}
			}
		}

		return &row
	},
	All: func() map[uint32]models.Site {
		rows := map[uint32]models.Site{}
		_ = getCache(consts.PlatformIntegrated, keyPlatformSites, &rows)
		return rows
	},
}
