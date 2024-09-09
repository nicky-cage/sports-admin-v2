package controllers

import (
	"errors"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var MessageAgent = struct {
	*ActionList
	*ActionSave
	*ActionUpdate
	*ActionCreate
	*ActionDelete
}{
	ActionList: &ActionList{
		Model:    models.Messages,
		ViewFile: "message_agent/list.html",
		Rows: func() interface{} {
			return &[]models.Message{}
		},
		ProcessRow: func(c *gin.Context, rs interface{}) {
			list := rs.(*[]models.Message)
			for i, v := range *list {
				tempContent := v.Contents
				tempContent = strings.Replace(tempContent, "<p>", "", -1)
				tempContent = strings.Replace(tempContent, "</p>", "\n", -1)
				tempContent = strings.Replace(tempContent, "&nbsp;", " ", -1)
				(*list)[i].Contents = tempContent
			}
		},
		QueryCond: map[string]interface{}{
			"title": "=",
			"state": "=",
		},
		OrderBy: func(*gin.Context) string {
			return "created DESC"
		},
		GetQueryCond: func(c *gin.Context) builder.Cond { //默认只查当月
			cond := builder.NewCond()
			var start_at int64
			var end_at int64
			if value, exists := c.GetQuery("created"); !exists {
				currentDayTime := time.Now().Format("2006-01-02")
				start_at = tools.GetMicroTimeStampByString(currentDayTime + " 00:00:00")
				end_at = tools.GetMicroTimeStampByString(currentDayTime + " 23:59:59")
				cond = cond.And(builder.Gte{"created": start_at}).And(builder.Lte{"created": end_at})
			} else {
				areas := strings.Split(value, " - ")
				start_at = tools.GetMicroTimeStampByString(areas[0])
				end_at = tools.GetMicroTimeStampByString(areas[1])
				cond = cond.And(builder.Gte{"created": start_at}).And(builder.Lte{"created": end_at})
			}
			cond = cond.And(builder.Eq{"is_agent": 2})
			return cond
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Messages,
		ViewFile: "message_agent/created.html",
	},
	ActionUpdate: &ActionUpdate{
		Model: models.Messages,
		Row: func() interface{} {
			return &models.Message{}
		},
		ViewFile: "message_agent/updated.html",
	},
	ActionSave: &ActionSave{
		Model: models.Messages,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			if (*m)["send_type"].(string) == "2" {
				if len((*m)["send_target"].(string)) <= 0 {
					return errors.New("会员不能为空")
				}
				sendTarget := strings.Split((*m)["send_target"].(string), ",")
				platform := request.GetPlatform(c)
				for _, v := range sendTarget {
					str := strings.Replace(v, " ", "", -1)
					if str == "" {
						continue
					}
					user := &models.User{}
					b, _ := models.Users.Find(platform, user, builder.NewCond().And(builder.Eq{"username": str}))
					if !b {
						return errors.New("有不存在的会员")
					}
				}
			} else {
				(*m)["send_target"] = ""
			}
			return nil
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.Messages,
	},
}
