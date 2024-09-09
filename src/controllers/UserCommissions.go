package controllers

import (
	"context"
	"encoding/json"
	common "sports-common"
	"sports-common/consts"
	"sports-common/es"
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
	"github.com/olivere/elastic/v7"
	"xorm.io/builder"
)

type RebateRes struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key           string `json:"key"`
		DocCount      int    `json:"doc_count"`
		GroupGameCode struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Key           int `json:"key"`
				DocCount      int `json:"doc_count"`
				GroupGameType struct {
					Value float64 `json:"value"`
				} `json:"group_game_type"`
			} `json:"buckets"`
		} `json:"group_game_code"`
	} `json:"buckets"`
}

type OutResType struct {
	Name        string  //大类AG
	Child       int     //子类捕鱼
	Value       float64 //流水
	Ratio       float64 //返水比例
	RebateMoney float64 //返水金额
}
type UserCommissionsStruct struct {
	models.UserRebateRecord `xorm:"extends"`
	Vip                     int32
}
type SumRecordStruct struct {
	Money float64
}

// 返水
var UserCommissions = struct {
	Record func(*gin.Context) //返水记录
	Manual func(*gin.Context) //手动返水
	Issue  func(*gin.Context) //手动返水-发放
}{
	Manual: func(c *gin.Context) {
		var startAt int64
		var endAt int64
		// if value, exists := c.GetQuery("commissons_time"); !exists {
		// 	currentTime := time.Now().Unix()
		// 	startAt = currentTime - currentTime%86400
		// 	endAt = startAt + 86400
		// } else {
		// 	startAt = tools.GetTimeStampByString(value + " 00:00:00")
		// 	endAt = tools.GetTimeStampByString(value+" 00:00:00") + 86400
		// }
		day := c.DefaultQuery("commissons_time", time.Now().Format("2006-01-02"))
		startAt = tools.GetTimeStampByString(day + " 00:00:00")
		endAt = tools.GetTimeStampByString(day+" 00:00:00") + 86400
		username := c.DefaultQuery("username", "jos123")
		platform := request.GetPlatform(c)
		esIndexName := platform + "_wagers"
		esClient, err := es.GetClientByPlatform(platform)
		if err != nil {
			response.Err(c, "es连接无响应")
			return
		}
		defer esClient.Stop()
		isRebate := false       //是否应该返水
		hasRebate := 0.00       //已返水
		remainingRebate := 0.00 //剩余未返水
		userInfo := &models.User{}
		b, err := models.Users.Find(platform, userInfo, builder.NewCond().And(builder.Eq{"username": username}))
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		if !b {
			response.Err(c, "用户不存在")
			return
		}
		boolQuery := elastic.NewBoolQuery()
		boolQuery.Must(elastic.NewMatchQuery("username", username))
		boolQuery.Filter(elastic.NewRangeQuery("created_at").Gte(startAt).Lte(endAt))
		sumMoneyAggs := elastic.NewSumAggregation().Field("valid_money")
		aggs := elastic.NewTermsAggregation().Field("game_type").SubAggregation("group_game_type", sumMoneyAggs)
		aggss := elastic.NewTermsAggregation().Field("game_code").SubAggregation("group_game_code", aggs)
		sumRes, err := esClient.Search(esIndexName).Query(boolQuery).
			Size(0).Aggregation("group_sum_money", aggss).
			Do(context.Background())
		if err != nil {
			response.Err(c, "es获取聚合数据无响应")
			return
		}
		buckets := &RebateRes{}
		outRes := make([]OutResType, 0)
		sumRebateMoney := 0.00
		sumNeedRebate := 0.00
		if b, ok := sumRes.Aggregations["group_sum_money"]; ok {
			_ = json.Unmarshal(b, buckets)
			for _, v := range buckets.Buckets {
				for _, vv := range v.GroupGameCode.Buckets {
					s := OutResType{}
					s.Name = v.Key
					s.Child = vv.Key
					s.Value = vv.GroupGameType.Value
					sumRebateMoney += s.Value
					outRes = append(outRes, s)
				}
			}
			if len(outRes) > 0 {
				for i, v := range outRes {
					gameVenueInfo := &models.GameVenue{}
					_, _ = models.GameVenues.Find(platform, gameVenueInfo, builder.NewCond().And(builder.Eq{"pid": 0, "code": v.Name}))
					outRes[i].Name = gameVenueInfo.Name
					rebateSettingInfo := &models.UserRebateSetting{}
					_, _ = models.UserRebateSettings.Find(platform, rebateSettingInfo, builder.NewCond().
						And(builder.Eq{"venue_id": gameVenueInfo.Id}).
						And(builder.Eq{"type_id": v.Child}).
						And(builder.Eq{"vip_id": userInfo.Vip}))
					outRes[i].Ratio = rebateSettingInfo.Ratio
					outRes[i].RebateMoney = v.Value * rebateSettingInfo.Ratio / 100
					sumNeedRebate += outRes[i].RebateMoney
				}
			}
		}
		userRebateInfo := &models.UserRebateRecord{}
		exist, err := models.UserRebateRecords.Find(platform, userRebateInfo, builder.NewCond().
			And(builder.Eq{"username": username}).
			And(builder.Eq{"day": day}))
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统繁忙")
			return
		}
		if !exist {
			if sumNeedRebate > 0 {
				isRebate = true
			}
		} else {
			hasRebate = userRebateInfo.Money
		}
		remainingRebate = sumNeedRebate - hasRebate
		if remainingRebate > 0 {
			isRebate = true
		}
		viewData := pongo2.Context{
			"userInfo":        userInfo,
			"outRes":          outRes,
			"isRebate":        isRebate,
			"hasRebate":       hasRebate,
			"remainingRebate": remainingRebate,
			"sumRebateMoney":  sumRebateMoney,
			"sumNeedRebate":   sumNeedRebate,
		}
		viewFile := "user_commissons/manuals.html"
		if request.IsAjax(c) {
			viewFile = "user_commissons/_manuals.html"
		}
		response.Render(c, viewFile, viewData)
	},
	Issue: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		day := time.Now().Format("2006-01-02")
		rKey := "rebate:issue:" + day + postedData["username"].(string)
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

		userInfo := &models.User{}
		_, _ = models.Users.Find(platform, userInfo, builder.NewCond().And(builder.Eq{"username": postedData["username"]}))
		accountInfo := &models.Account{}
		_, _ = models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": userInfo.Id}))
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		administrator := GetLoginAdmin(c)
		rebateMoney, _ := strconv.ParseFloat(postedData["rebate_money"].(string), 64)
		imap := map[string]interface{}{
			"username":   postedData["username"],
			"user_id":    userInfo.Id,
			"issue_type": 2,
			"operator":   administrator.Name,
			"money":      rebateMoney,
			"day":        postedData["commissons_day"],
			"created":    tools.NowMicro(),
		}
		if _, err := session.Table("user_rebate_record").Insert(imap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "发放失败")
			return
		}
		transAction := &models.Transaction{}
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            c.ClientIP(),
			"description":   "手动返水",
			"administrator": administrator.Name,
			"admin_user_id": administrator.Id,
			"serial_number": tools.GetBillNo("FS", 0),
		}
		transType := consts.TransTypeRebateBonus
		if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, rebateMoney, extraMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, err.Error())
			return
		}
		_ = session.Commit()
		response.Ok(c)
	},
	Record: func(c *gin.Context) {
		var start_at int64
		var end_at int64
		if value, exists := c.GetQuery("created"); !exists {
			current_time := time.Now().Unix()
			start_at = current_time - current_time%86400
			end_at = start_at + 86400
		} else {
			areas := strings.Split(value, " - ")
			start_at = tools.GetTimeStampByString(areas[0] + " 00:00:00")
			end_at = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
		}
		username := c.DefaultQuery("username", "")
		cond := builder.NewCond()
		cond = cond.And(builder.Gte{"user_rebate_records.created": start_at}).And(builder.Lte{"user_rebate_records.created": end_at})
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"user_rebate_records.username": username})
		}
		limit, offset := request.GetOffsets(c)
		userRebateRecords := make([]UserCommissionsStruct, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("user_rebate_records").Join("LEFT OUTER", "users", "user_rebate_records.user_id = users.id").Where(cond).OrderBy("user_rebate_records.id DESC").Limit(limit, offset).FindAndCount(&userRebateRecords)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ss := new(SumRecordStruct)
		sumTotal, _ := engine.Table("user_rebate_records").Where(cond).Sum(ss, "money")
		viewData := pongo2.Context{
			"rows":      userRebateRecords,
			"total":     total,
			"sum_money": sumTotal,
		}
		viewFile := "user_commissons/records.html"
		if request.IsAjax(c) {
			viewFile = "user_commissons/_records.html"
		}
		response.Render(c, viewFile, viewData)
	},
}
