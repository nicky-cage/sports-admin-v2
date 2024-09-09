package controllers

import (
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var UserFeedBack = struct {
	*ActionList
	*ActionSave
	*ActionUpdate
	*ActionDelete
}{
	ActionList: &ActionList{
		Model:    models.UserFeedbacks,
		ViewFile: "user_feedback/list.html",
		OrderBy: func(c *gin.Context) string {
			return "created DESC"
		},
		Rows: func() interface{} {
			return &[]models.UserFeedback{}
		},
		ProcessRow: func(c *gin.Context, rs interface{}) {
			list := rs.(*[]models.UserFeedback)
			for i, v := range *list {
				tempContent := v.Reply
				tempContent = strings.Replace(tempContent, "<p>", "", -1)
				tempContent = strings.Replace(tempContent, "</p>", "\n", -1)
				tempContent = strings.Replace(tempContent, "&nbsp;", " ", -1)
				(*list)[i].Reply = tempContent

				if v.ImageUrl != "" {
					iArr := strings.Split(v.ImageUrl, ",")
					(*list)[i].ImageLinks = iArr
				}
			}
		},
		QueryCond: map[string]interface{}{
			"login_id":     "=",
			"message_type": "=",
			"username":     "%",
			"admin":        "%",
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			cond := builder.NewCond()
			var startAt int64
			var endAt int64
			if value, exists := c.GetQuery("created"); !exists {
				//currentTime := time.Now().Unix()
				//startAt = currentTime - currentTime%86400
				//endAt = startAt + 86400
			} else {
				areas := strings.Split(value, " - ")

				if areas[0] == areas[1] {
					//startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
					//endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
				} else {
					startAt = tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
					endAt = tools.GetMicroTimeStampByString(areas[1] + " 00:00:00")
					cond = cond.And(builder.Gte{"created": startAt}).And(builder.Lt{"created": endAt})
				}
			}

			cond = cond.And(builder.Neq{"message_type": 0})
			return cond
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.UserFeedbacks,
		ViewFile: "user_feedback/edit.html",
		Row: func() interface{} {
			return &models.UserFeedback{}
		}, ExtendData: func(c *gin.Context) pongo2.Context {
			admin := GetLoginAdmin(c)
			return pongo2.Context{"admin": admin.Name}
		},
	},
	ActionSave: &ActionSave{
		Model: models.UserFeedbacks,
	},
	ActionDelete: &ActionDelete{
		Model: models.UserFeedbacks,
	},
}
