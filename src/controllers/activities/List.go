package activities

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

func (*Activities) List(c *gin.Context) { //默认首页
	cond := builder.NewCond()
	if value, exists := c.GetQuery("show_time"); exists {
		areas := strings.Split(value, " - ")
		startAt := tools.GetTimeStampByString(areas[0])
		endAt := tools.GetTimeStampByString(areas[1])
		//如活动展示时间为7月10号-8月3号
		//搜索展示时间为8月2号到8月31号
		//则该活动依然被搜索到，即，在这个搜索展示时间内，活动依然处于展示状态的，都可以搜索到
		cond = cond.And(builder.Gte{"show_time_start": startAt}).And(builder.Lte{"show_time_end": endAt})
	}
	cond = cond.And(builder.Neq{"content_form": 3})
	if title := c.DefaultQuery("title", ""); title != "" {
		cond = cond.And(builder.Eq{"title": title})
	}
	if state := c.DefaultQuery("state", ""); state != "" {
		cond = cond.And(builder.Eq{"state": state})
	}
	if activityType := c.DefaultQuery("activity_type", ""); activityType != "" {
		cond = cond.And(builder.Eq{"activity_type": activityType})
	}

	limit, offset := request.GetOffsets(c)
	aRows := make([]models.Activity, 0)
	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()

	total, err := myClient.Table("activities").Where(cond).OrderBy("id DESC").Limit(limit, offset).FindAndCount(&aRows)
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "获取列表错误")
		return
	}
	rules := &models.InviteFriendsRule{}
	b, err := models.InviteFriendsRules.Find(platform, rules)
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(c, "系统错误")
		return
	}
	if !b {
		response.Err(c, "邀请好友规则不存在")
		return
	}
	viewData := pongo2.Context{
		"rows":  aRows,
		"total": total,
		"rules": rules,
	}
	viewFile := "activities/activities.html"
	if request.IsAjax(c) {
		viewFile = "activities/_activities.html"
	}
	base_controller.SetLoginAdmin(c)
	response.Render(c, viewFile, viewData)
}
