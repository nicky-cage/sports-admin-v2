package controllers

import (
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

type UserAccountSetAuditsStruct struct {
	models.UserAccountSet `xorm:"extends"`
	Vip                   int32
}

// UserAccountAudits 上下分-审核列表
var UserAccountAudits = struct {
	List   func(*gin.Context)
	Agree  func(*gin.Context) //同意页面
	Refuse func(*gin.Context) //拒绝页面
	SaveDo func(*gin.Context) //保存
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		var start_at int64
		var end_at int64
		if value, exists := c.GetQuery("created"); !exists {
			currentDayTime := time.Now().Format("2006-01-02")
			start_at = tools.GetMicroTimeStampByString(currentDayTime + " 00:00:00")
			end_at = tools.GetMicroTimeStampByString(currentDayTime + " 23:59:59")
			cond = cond.And(builder.Gte{"user_account_sets.created": start_at}).And(builder.Lte{"user_account_sets.created": end_at})
		} else {
			areas := strings.Split(value, " - ")
			start_at = tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
			end_at = tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
			cond = cond.And(builder.Gte{"user_account_sets.created": start_at}).And(builder.Lte{"user_account_sets.created": end_at})
		}
		cond = cond.And(builder.Eq{"user_account_sets.status": 1})
		if min, ok := c.GetQuery("money_min"); ok {
			cond = cond.And(builder.Gte{"user_account_sets.money": min})
		}
		if max, ok := c.GetQuery("money_max"); ok {
			cond = cond.And(builder.Lte{"user_account_sets.money": max})
		}
		username := c.DefaultQuery("username", "")
		applicant := c.DefaultQuery("applicant", "")
		sType := c.DefaultQuery("type", "")
		limit, offset := request.GetOffsets(c)
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"user_account_sets.username": username})
		}
		if len(applicant) > 0 {
			cond = cond.And(builder.Eq{"user_account_sets.applicant": applicant})
		}
		if len(sType) > 0 {
			cond = cond.And(builder.Eq{"user_account_sets.type": sType})
		}
		userAccountSets := make([]UserAccountSetAuditsStruct, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("user_account_sets").Join("LEFT OUTER", "users", "user_account_sets.user_id = users.id").Where(cond).OrderBy("user_account_sets.id DESC").Limit(limit, offset).FindAndCount(&userAccountSets)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ss := new(SumStruct)
		sumTotal, _ := engine.Table("user_account_sets").Where(cond).Sum(ss, "money")
		viewData := pongo2.Context{
			"rows":      userAccountSets,
			"total":     total,
			"sum_money": sumTotal,
		}
		viewFile := "user_account_sets/audits.html"
		if request.IsAjax(c) {
			viewFile = "user_account_sets/_audits.html"
		}
		SetLoginAdmin(c)
		response.Render(c, viewFile, viewData)
	},
	Agree: func(c *gin.Context) {
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		sql := "select a.*,b.vip from user_account_sets a left join users b on a.user_id=b.id where a.id=" + idStr
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		data, err := dbSession.QueryString(sql)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		cr, _ := strconv.Atoi(data[0]["created"])
		data[0]["created"] = time.UnixMicro(int64(cr)).Format("2006-01-02 15:04:05")
		viewData := pongo2.Context{"r": data[0]}
		response.Render(c, "user_account_sets/agree.html", viewData)
	},
	Refuse: func(c *gin.Context) { //下分
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		sql := "select a.*,b.vip from user_account_sets a left join users b on a.user_id=b.id where a.id=" + idStr
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		data, err := dbSession.QueryString(sql)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		cr, _ := strconv.Atoi(data[0]["created"])
		data[0]["created"] = time.UnixMicro(int64(cr)).Format("2006-01-02 15:04:05")
		viewData := pongo2.Context{"r": data[0]}
		response.Render(c, "user_account_sets/refuse.html", viewData)
	},
	SaveDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		//防止多人同时更改-开始
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		rKey := "ls:UserAccountAudits:" + idStr + ":" + postedData["username"].(string)
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		if num > 1 {
			response.Err(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)

		userAccountAuditsInfo := &models.UserAccountSet{}
		if exists, err = models.UserAccountSets.FindById(platform, id, userAccountAuditsInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "订单不存在")
			return
		}
		if userAccountAuditsInfo.Status != 1 {
			response.Err(c, "该订单已经被操作")
			return
		}
		accountInfo := &models.Account{}
		if exists, err = models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": userAccountAuditsInfo.UserId})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户账户信息失败")
			return
		}
		userInfo := &models.User{}
		if exists, err = models.Users.FindById(platform, int(userAccountAuditsInfo.UserId), userInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户信息失败")
			return
		}
		isFirstDeposit := 1
		if models.UserDeposits.IsFirstTime(platform, int(userInfo.Id), false) { // 判断是否首存
			isFirstDeposit = 2
		}
		//事务操作
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		uMap := map[string]interface{}{
			"user_vip":     postedData["user_vip"],
			"status":       postedData["status"],
			"audit_remark": postedData["audit_remark"],
			"audit":        GetLoginAdmin(c).Name,
			"updated":      tools.NowMicro(),
			"top_id":       userInfo.TopId,
		}
		if _, err := session.Table("user_account_sets").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "更新失败")
			return
		}

		administrator := GetLoginAdmin(c)
		transAction := &models.Transaction{}
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"administrator": administrator.Name,
			"admin_user_id": administrator.Id,
			"serial_number": userAccountAuditsInfo.BillNo,
		}
		//同意
		if postedData["status"] == "2" {
			depositMap := map[string]interface{}{
				"order_no":         userAccountAuditsInfo.BillNo,
				"user_id":          userAccountAuditsInfo.UserId,
				"username":         userAccountAuditsInfo.Username,
				"top_code":         userInfo.TopCode,
				"top_name":         userInfo.TopName,
				"top_id":           userInfo.TopId,
				"real_name":        userInfo.RealName,
				"deposit_name":     "代客充值",
				"type":             3,
				"money":            userAccountAuditsInfo.Money,
				"arrive_money":     userAccountAuditsInfo.Money,
				"confirm_money":    userAccountAuditsInfo.Money,
				"top_money":        userAccountAuditsInfo.Money,
				"discount":         0.00,
				"finance_admin":    administrator.Name,
				"status":           2,
				"is_first_deposit": isFirstDeposit,
				"created":          tools.NowMicro(),
				"deposit_at":       tools.Now(),
				"confirm_at":       tools.Now(),
				"updated":          tools.NowMicro(),
				"remark":           postedData["audit_remark"],
			}
			if _, err := session.Table("user_deposits").Insert(depositMap); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "同意失败")
				return
			}

			transType := consts.TransTypeManualPointsLess
			if postedData["stype"] == "1" {
				extraMap["description"] = postedData["applicant_remark"]
				transType = consts.TransTypeManualPointsPlus
			} else {
				extraMap["description"] = "手动下分"
			}
			platform := request.GetPlatform(c)
			if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, userAccountAuditsInfo.Money, extraMap); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}
		}
		_ = session.Commit()
		//覆盖用户钱包的数据
		if accountInfo.Id > 0 {
			_ = accountInfo.ResetCacheData(redis)
		}
		response.Ok(c)
		//事务操作
	},
}
