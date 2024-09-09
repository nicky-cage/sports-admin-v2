package controllers

import (
	"errors"
	"sports-common/config"
	"sports-common/consts"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// Notices 系统公告
var Notices = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
	*ActionState
}{
	ActionList: &ActionList{
		Model:    models.Notices,
		ViewFile: "notices/list.html",
		Rows: func() interface{} {
			return &[]models.Notice{}
		},
		ProcessRow: func(c *gin.Context, rs interface{}) {
			list := rs.(*[]models.Notice)
			for i, v := range *list {
				temp := v.PlatformTypes
				temp = strings.Replace(temp, "1", "全站", -1)
				temp = strings.Replace(temp, "2", "体育", -1)
				temp = strings.Replace(temp, "3", "Web", -1)
				temp = strings.Replace(temp, "4", "H5", -1)
				tempContent := v.Contents
				tempContent = strings.Replace(tempContent, "<p>", "", -1)
				tempContent = strings.Replace(tempContent, "</p>", "\n", -1)
				tempContent = strings.Replace(tempContent, "&nbsp;", " ", -1)
				(*list)[i].PlatformTypes = temp
				(*list)[i].Contents = tempContent
			}
		},
		OrderBy: func(*gin.Context) string {
			return "id DESC"
		},
		QueryCond: map[string]interface{}{
			"title": "=",
			"state": "=",
		},
		GetQueryCond: func(c *gin.Context) builder.Cond { //默认只查当月
			cond := builder.NewCond()
			var startAt int64
			var endAt int64
			if value, exists := c.GetQuery("created"); !exists {
				currentDayTime := time.Now().Format("2006-01-02")
				startAt = tools.GetMicroTimeStampByString(currentDayTime + " 00:00:00")
				endAt = tools.GetMicroTimeStampByString(currentDayTime + " 23:59:59")
				cond = cond.And(builder.Gte{"created": startAt}).And(builder.Lte{"created": endAt})
			} else {
				areas := strings.Split(value, " - ")
				startAt = tools.GetMicroTimeStampByString(areas[0])
				endAt = tools.GetMicroTimeStampByString(areas[1])
				cond = cond.And(builder.Gte{"created": startAt}).And(builder.Lte{"created": endAt})
			}
			return cond
		}, ExtendData: func(c *gin.Context) pongo2.Context {
			return pongo2.Context{"messageType": consts.FeedbackTypes}
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Notices,
		ViewFile: "notices/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			return pongo2.Context{
				"static_url": config.Get("config.static_url"),
				"method":     "create",
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model: models.Notices,
		Row: func() interface{} {
			return &models.Notice{}
		},
		ViewFile: "notices/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			return pongo2.Context{
				"static_url": config.Get("config.static_url"),
				"method":     "update",
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Notices,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			startAt := (*m)["start_at"].(string)
			(*m)["start_at"] = tools.GetTimeStampByString(startAt)
			endAt := (*m)["end_at"].(string)
			(*m)["end_at"] = tools.GetTimeStampByString(endAt)
			if (*m)["start_at"].(int64) > (*m)["end_at"].(int64) {
				return errors.New("开始时间大于结束时间")
			}

			tempPlatformTypes := ""
			for i := 0; i <= 3; i++ {
				v, ok := (*m)["platform_types["+strconv.Itoa(i)+"]"]
				if ok {
					tempPlatformTypes += v.(string) + ","
				}
			}
			tempPlatformTypes = strings.TrimRight(tempPlatformTypes, ",")
			(*m)["platform_types"] = tempPlatformTypes

			tempVipIds := ""
			for i := 0; i <= 10; i++ {
				v, ok := (*m)["vip_ids["+strconv.Itoa(i)+"]"]
				if ok {
					tempVipIds += v.(string) + ","
				}
			}
			tempVipIds = strings.TrimRight(tempVipIds, ",")
			(*m)["vip_ids"] = tempVipIds
			return nil
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.Notices,
	},
	ActionState: &ActionState{
		Model: models.Notices,
	},
}
