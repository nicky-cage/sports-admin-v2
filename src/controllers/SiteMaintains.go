package controllers

import (
	"fmt"
	common "sports-common"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// SiteMaintains 维护设置
var SiteMaintains = struct {
	*ActionList
	*ActionUpdate
	*ActionSave
}{
	ActionList: &ActionList{
		Model:    models.SiteMaintains,
		ViewFile: "site_maintains/list.html",
		OrderBy: func(c *gin.Context) string {
			return "id DESC"
		},
		Rows: func() interface{} {
			return &[]models.SiteMaintain{}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.SiteMaintains,
		ViewFile: "site_maintains/edit.html",
		Row: func() interface{} {
			return &models.SiteMaintain{}
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			admin := GetLoginAdmin(c)
			return pongo2.Context{
				"admin": admin.Name,
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.SiteMaintains,
		SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
			maintainTime := (*data)["maintain_time"]
			times := strings.Split(maintainTime.(string), " - ")
			strStart := times[0]
			strEnd := times[1]
			timeStart := tools.GetTimeStampByString(strStart)
			timeEnd := tools.GetTimeStampByString(strEnd)

			(*data)["time_start"] = uint32(timeStart) //开始时间
			(*data)["time_end"] = uint32(timeEnd)     //结束时间

			admin := GetLoginAdmin(c)
			(*data)["admin_id"] = admin.Id
			(*data)["admin_name"] = admin.Name
			return nil
		},
		SaveAfter: func(c *gin.Context, data *map[string]interface{}) {
			platform := request.GetPlatform(c)
			_, _ = models.SiteMaintainLogs.Create(platform, *data) //创建一条记录
			// maintainState := (*data)["state"].(string)   // 1:正在维护/2:正常状态
			maintainPlatform := (*data)["platform-name"] // 平台: pc/app/h5
			// maintainStartTime := (*data)["time_start"]   // 开始维护时间
			// maintainEndTime := (*data)["time_end"]       // 结束维护时间

			admin := GetLoginAdmin(c)
			cacheKey := fmt.Sprintf("maintain:%v", maintainPlatform)
			id, _ := strconv.Atoi((*data)["id"].(string))
			platformId, _ := strconv.Atoi((*data)["platform_id"].(string))
			state, _ := strconv.Atoi((*data)["state"].(string))
			cacheValue := (&models.SiteMaintain{
				Id:         uint32(id),
				PlatformId: uint8(platformId),
				TimeStart:  (*data)["time_start"].(uint32),
				TimeEnd:    (*data)["time_end"].(uint32),
				State:      uint8(state),
				AdminId:    uint32(admin.Id),
				AdminName:  admin.Name,
				Remark:     (*data)["remark"].(string),
			}).ToJson()

			cache := common.Redis(platform)
			defer common.RedisRestore(platform, cache)
			timeZero := time.Duration(0)
			cache.Set(cacheKey, cacheValue, timeZero)
		},
	},
}
