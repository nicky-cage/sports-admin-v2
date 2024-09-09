package controllers

import (
	"errors"
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// Messages 站内信
var Messages = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
	*ActionState
	Top func(*gin.Context) //是否置顶
}{
	ActionList: &ActionList{
		Model:    models.Messages,
		ViewFile: "messages/list.html",
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
			return "id DESC"
		},
		GetQueryCond: func(c *gin.Context) builder.Cond { //默认只查当月
			cond := builder.NewCond()
			var startAt int64
			var endAt int64
			if value, exists := c.GetQuery("created"); !exists {
				currentDayTime := time.Now().Format("2006-01-02")
				startAt = tools.GetMicroTimeStampByString(currentDayTime + " 00:00:00")
				endAt = tools.GetMicroTimeStampByString(currentDayTime + " 23:59:59")
			} else {
				areas := strings.Split(value, " - ")
				startAt = tools.GetMicroTimeStampByString(areas[0])
				endAt = tools.GetMicroTimeStampByString(areas[1])
			}
			cond = cond.And(builder.Gte{"created": startAt}).And(builder.Lte{"created": endAt})
			cond = cond.And(builder.Neq{"is_agent": 2})
			return cond
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.Messages,
		ViewFile: "messages/create.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			return pongo2.Context{
				"static_url": config.Get("config.static_url"),
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model: models.Messages,
		Row: func() interface{} {
			return &models.Message{}
		},
		ViewFile: "messages/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			return pongo2.Context{
				"static_url": config.Get("config.static_url"),
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.Messages,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			if (*m)["method"] == "add" {
				if strings.TrimSpace((*m)["contents"].(string)) == "" {
					return errors.New("新增内容不能为空")
				}
			}
			platform := request.GetPlatform(c)
			if (*m)["send_type"].(string) == "2" {
				if len((*m)["send_target"].(string)) <= 0 {
					return errors.New("会员不能为空")
				}
				sendTarget := strings.Split((*m)["send_target"].(string), ",")
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
		DeleteBefore: func(c *gin.Context, data interface{}) error {
			idStr := c.DefaultQuery("id", "0") //检测ID
			id, _ := strconv.Atoi(idStr)
			platform := request.GetPlatform(c)
			engine := common.Mysql(platform)
			defer engine.Close()
			umap := map[string]interface{}{
				"is_del": 1, //停用 is_del=1不显示
			}
			if _, err := engine.Table("user_messages").Where("msg_id=?", id).Update(umap); err != nil {
				log.Logger.Error(err.Error())
			}
			return nil
		},
	},
	ActionState: &ActionState{
		Model: models.Messages,
		StateAfter: func(c *gin.Context) {
			idStr := c.DefaultQuery("id", "0") //检测ID
			id, _ := strconv.Atoi(idStr)
			toStateStr, _ := c.GetQuery("to_state")
			toState, _ := strconv.Atoi(toStateStr)
			platform := request.GetPlatform(c)
			engine := common.Mysql(platform)
			defer engine.Close()
			//启用 is_del=0 显示
			umap := map[string]interface{}{
				"is_del": 0,
			}
			if toState == 1 { //停用 is_del=1不显示
				umap["is_del"] = 1
			}
			if _, err := engine.Table("user_messages").Where("msg_id=?", id).Update(umap); err != nil {
				log.Logger.Error(err.Error())
			}
		},
	},
	Top: func(c *gin.Context) {
		isTop := &ActionState{
			Field: "is_top",
			Model: models.Messages,
		}
		idStr := c.DefaultQuery("id", "0") //检测ID
		id, _ := strconv.Atoi(idStr)
		toStateStr, _ := c.GetQuery("to_is_top")
		toState, _ := strconv.Atoi(toStateStr)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		umap := map[string]interface{}{
			"is_top": 2, //置顶
		}
		//is_top　1不置顶　2置顶
		if toState == 1 { //取消置顶
			umap["is_top"] = 1
		}
		if _, err := engine.Table("user_messages").Where("msg_id=?", id).Update(umap); err != nil {
			log.Logger.Error(err.Error())
		}
		isTop.State(c)
	},
}
