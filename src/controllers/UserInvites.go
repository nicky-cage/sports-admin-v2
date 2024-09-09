package controllers

import (
	common "sports-common"
	"sports-common/config"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

type InviteListStruct struct {
	models.InviteList `xorm:"extends" json:"invite_list"`
	models.User       `xorm:"extends" json:"user"`
}

// UserInvites 邀请好友
var UserInvites = struct {
	List        func(*gin.Context)
	RuleSetting func(*gin.Context)
	SaveDo      func(*gin.Context)
	Enable      func(*gin.Context)
	Agree       func(*gin.Context)
	Refuse      func(*gin.Context)
	AgreeDo     func(*gin.Context)
	RefuseDo    func(*gin.Context)
}{
	List: func(c *gin.Context) { //默认首页
		cond := builder.NewCond()
		username := c.DefaultQuery("username", "")
		introduceUsername := c.DefaultQuery("introduce_username", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"invite_lists.username": username})
		}
		if len(introduceUsername) > 0 {
			cond = cond.And(builder.Eq{"invite_lists.invite_uname": introduceUsername})
		}
		if min, ok := c.GetQuery("money_min"); ok {
			cond = cond.And(builder.Gte{"invite_lists.first_deposit": min})
		}
		if max, ok := c.GetQuery("money_max"); ok {
			cond = cond.And(builder.Lte{"invite_lists.first_deposit": max})
		}
		status := c.DefaultQuery("status", "")
		if len(status) > 0 {
			cond = cond.And(builder.Eq{"invite_lists.status": status})
		}
		limit, offset := request.GetOffsets(c)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		inviteFriends := make([]InviteListStruct, 0)
		total, err := engine.Table("invite_lists").
			Join("LEFT OUTER", "users", "invite_lists.user_id = users.id").
			Where(cond).OrderBy("invite_lists.id DESC").
			Limit(limit, offset).
			FindAndCount(&inviteFriends)
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
			"rows":  inviteFriends,
			"r":     rules,
			"total": total,
		}
		viewFile := "activities/_user_invites.html"
		if request.IsAjax(c) {
			viewFile = "activities/_user_invites.html"
		}
		response.Render(c, viewFile, viewData)
	},
	RuleSetting: func(c *gin.Context) { //规则设置
		rules := &models.InviteFriendsRule{}
		platform := request.GetPlatform(c)
		b, err := models.InviteFriendsRules.Find(platform, rules)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统错误")
			return
		}
		if !b {
			response.Err(c, "记录不存在")
			return
		}
		viewData := pongo2.Context{
			"r":          rules,
			"STATIC_URL": config.Get("internal.img_host_backend", ""),
		}
		viewFile := "activities/rule_setting.html"
		response.Render(c, viewFile, viewData)
	},
	SaveDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		umap := map[string]interface{}{
			"time_type":      postedData["time_type"],
			"deposit_amount": postedData["deposit_amount"],
			"gift_bonus":     postedData["gift_bonus"],
			"interval":       postedData["interval"],
			"activity_img":   postedData["activity_img"],
			"activity_rule":  postedData["activity_rule"],
			"invite_again":   postedData["invite_again"],
			"updated":        tools.NowMicro(),
		}
		if umap["time_type"].(string) == "2" {
			umap["time_start"] = tools.GetTimeStampByString(postedData["time_start"].(string))
			umap["time_end"] = tools.GetTimeStampByString(postedData["time_end"].(string))
		}
		if _, err := engine.Table("invite_friends_rules").Where("id=?", id).Update(umap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "编辑规则失败")
			return
		}
		response.Ok(c)
	},
	Enable: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		umap := map[string]interface{}{
			"state":   postedData["state"],
			"updated": tools.NowMicro(),
		}
		if _, err := engine.Table("invite_friends_rules").Update(umap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "修改邀请好友活动失败")
			return
		}
		response.Ok(c)
	},
	Agree: func(c *gin.Context) {
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			log.Err("无法获取id信息!\n")
			return
		}
		viewData := pongo2.Context{"id": idStr}
		response.Render(c, "activities/agree.html", viewData)
	},
	Refuse: func(c *gin.Context) {
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			log.Err("无法获取id信息!\n")
			return
		}
		viewData := pongo2.Context{"id": idStr}
		response.Render(c, "activities/refuse.html", viewData)
	},
	AgreeDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		//防止多人同时更改
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		rKey := "invite_friends:" + "_" + idStr
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		if num > 1 {
			response.Err(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)

		inviteFriendsInfo := &models.InviteFriend{}
		if exists, err := models.InviteFriends.FindById(platform, id, inviteFriendsInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找这条记录失败")
			return
		}
		if inviteFriendsInfo.State != 1 {
			response.Err(c, "这条记录已经被处理了")
			return
		}
		userInfo := &models.User{}
		if exists, err := models.Users.Find(platform, userInfo, builder.NewCond().And(builder.Eq{"username": inviteFriendsInfo.IntroduceUsername})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户信息失败")
			return
		}
		accountInfo := &models.Account{}
		if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": userInfo.Id})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户账户信息失败")
			return
		}
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Err(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		administrator := GetLoginAdmin(c)
		umap := map[string]interface{}{
			"remark":     postedData["remark"],
			"status":     4,
			"updated":    tools.NowMicro(),
			"bonus_time": time.Now().Unix(),
			"operator":   administrator.Name,
		}
		if _, err := session.Table("invite_lists").Where("id=?", id).Update(umap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "通过失败")
			return
		}
		transAction := &models.Transaction{}
		billNo := tools.GetBillNo("inf", 5)
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"description":   "邀请好友-派奖",
			"administrator": administrator.Name,
			"admin_user_id": administrator.Id,
			"serial_number": billNo,
		}
		transType := consts.TransTypeInviteFriendsBonus
		if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, inviteFriendsInfo.BonusAmount, extraMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, err.Error())
			return
		}
		_ = session.Commit()
		response.Ok(c)
	},
	RefuseDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		//防止多人同时更改
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		rKey := "invite_friends:" + "_" + idStr
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		if num > 1 {
			response.Err(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)

		inviteFriends := &models.InviteFriend{}
		if exists, err := models.InviteFriends.FindById(platform, id, inviteFriends); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找这条记录失败")
			return
		}
		if inviteFriends.State != 1 {
			response.Err(c, "这条记录已经被处理了")
			return
		}
		engine := common.Mysql(platform)
		defer engine.Close()
		umap := map[string]interface{}{
			"remark":  postedData["remark"],
			"status":  postedData["cause_failure"],
			"updated": tools.NowMicro(),
		}
		if _, err := engine.Table("invite_lists").Where("id=?", id).Update(umap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "启用活动失败")
			return
		}
		response.Ok(c)
	},
}
