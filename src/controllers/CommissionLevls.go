package controllers

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var DataVenueType = []map[string]string{
	{"id": "1", "name": "体育"},
	{"id": "2", "name": "电竞"},
	{"id": "3", "name": "真人"},
	{"id": "4", "name": "电游"},
	{"id": "5", "name": "捕鱼"},
	{"id": "6", "name": "彩票"},
	{"id": "7", "name": "棋牌"},
}

// CommissionLevls 返水等级
var CommissionLevls = struct {
	*ActionList
	Setup   func(*gin.Context)
	Details func(*gin.Context)
	SaveDo  func(*gin.Context)
}{
	ActionList: &ActionList{
		Model:    models.UserRebateVips,
		ViewFile: "commission_levels/commission_levels.html",
		OrderBy: func(*gin.Context) string {
			return "id ASC"
		},
		Rows: func() interface{} {
			return &[]models.UserRebateVip{}
		},
	},
	Setup: func(c *gin.Context) { //设置
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		vipStr, _ := c.GetQuery("vip")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sqlVip := "select * from user_rebate_vips where user_vip='" + vipStr + "'"
		dataVip, err := dbSession.QueryString(sqlVip)
		if err != nil {
			response.Err(c, "系统内部繁忙")
			log.Logger.Error(err.Error())
			return
		}
		tempSilce := make([]map[string]string, 0)
		//tempMap["ratio"] = strconv.FormatFloat(userRebateSetting.Ratio, 'f', -1, 64)
		for _, v := range DataVenueType {
			sqlVenue := "select id,name,venue_type from game_venues where venue_type=" + v["id"]
			dataVenue, _ := dbSession.QueryString(sqlVenue)
			for _, vv := range dataVenue {
				tempMap := make(map[string]string)
				userRebateSetting := &models.UserRebateSetting{}
				b, _ := models.UserRebateSettings.Find(platform, userRebateSetting, builder.NewCond().
					And(builder.Eq{"vip_id": idStr}).
					And(builder.Eq{"venue_id": vv["id"]}.
						And(builder.Eq{"type_id": vv["venue_type"]})))
				if b {
					tempMap["ratio"] = strconv.FormatFloat(userRebateSetting.Ratio, 'f', -1, 64)
					tempMap["type_id"] = strconv.Itoa(int(userRebateSetting.TypeId))
					tempMap["type_name"] = userRebateSetting.TypeName
					tempMap["venue_id"] = strconv.Itoa(int(userRebateSetting.VenueId))
					tempMap["vip_id"] = strconv.Itoa(int(userRebateSetting.VipId))
					tempMap["id"] = vv["id"]
					tempMap["name"] = vv["name"]
				} else {
					tempMap["ratio"] = "0.00"
					tempMap["type_id"] = vv["venue_type"]
					tempMap["venue_id"] = vv["id"]
					tempMap["type_name"] = ""
					tempMap["vip_id"] = idStr
					tempMap["id"] = vv["id"]
					tempMap["name"] = vv["name"]
				}
				tempSilce = append(tempSilce, tempMap)
			}
		}
		viewData := pongo2.Context{
			"child":         tempSilce,
			"venue":         DataVenueType,
			"id":            idStr,
			"day_max_water": dataVip[0]["day_max_water"],
		}
		response.Render(c, "commission_levels/setup.html", viewData)
	},
	Details: func(c *gin.Context) { //详情
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			response.Err(c, "无法获取id信息!\n")
			return
		}
		vipStr, _ := c.GetQuery("vip")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sqlVip := "select * from user_rebate_vips where user_vip='" + vipStr + "'"
		dataVip, err := dbSession.QueryString(sqlVip)
		if err != nil {
			response.Err(c, "系统内部繁忙")
			log.Logger.Error(err.Error())
			return
		}
		tempSilce := make([]map[string]string, 0)
		//tempMap["ratio"] = strconv.FormatFloat(userRebateSetting.Ratio, 'f', -1, 64)
		for _, v := range DataVenueType {
			sqlVenue := "select id,name,venue_type from game_venues where venue_type=" + v["id"]
			dataVenue, _ := dbSession.QueryString(sqlVenue)
			for _, vv := range dataVenue {
				tempMap := make(map[string]string)
				userRebateSetting := &models.UserRebateSetting{}
				b, _ := models.UserRebateSettings.Find(platform, userRebateSetting, builder.NewCond().
					And(builder.Eq{"vip_id": idStr}).
					And(builder.Eq{"venue_id": vv["id"]}.
						And(builder.Eq{"type_id": vv["venue_type"]})))
				if b {
					tempMap["ratio"] = strconv.FormatFloat(userRebateSetting.Ratio, 'f', -1, 64)
					tempMap["type_id"] = strconv.Itoa(int(userRebateSetting.TypeId))
					tempMap["type_name"] = userRebateSetting.TypeName
					tempMap["venue_id"] = strconv.Itoa(int(userRebateSetting.VenueId))
					tempMap["vip_id"] = strconv.Itoa(int(userRebateSetting.VipId))
					tempMap["id"] = vv["id"]
					tempMap["name"] = vv["name"]
				} else {
					tempMap["ratio"] = "0.00"
					tempMap["type_id"] = vv["venue_type"]
					tempMap["venue_id"] = vv["id"]
					tempMap["type_name"] = ""
					tempMap["vip_id"] = idStr
					tempMap["id"] = vv["id"]
					tempMap["name"] = vv["name"]
				}
				tempSilce = append(tempSilce, tempMap)
			}
		}
		viewData := pongo2.Context{
			"child":         tempSilce,
			"venue":         DataVenueType,
			"id":            idStr,
			"day_max_water": dataVip[0]["day_max_water"],
		}
		response.Render(c, "commission_levels/details.html", viewData)
	},
	SaveDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		vipIdStr, exists := postedData["vip_id"].(string)
		if !exists || exists && vipIdStr == "0" {
			response.Err(c, "id为空")
			return
		}
		vipId, _ := strconv.Atoi(vipIdStr)
		ratiaoSile := []string{}
		for k, v := range postedData {
			if strings.Contains(k, "ratio") {
				temp := strings.Split(k, "_")
				s := temp[2] + "," + v.(string)
				ratiaoSile = append(ratiaoSile, s)
			}
		}
		//事务操作
		platform := request.GetPlatform(c)
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		uMap := map[string]interface{}{
			"day_max_water": postedData["day_max_water"],
			"updated":       tools.NowMicro(),
		}
		if _, err := session.Table("user_rebate_vips").Where("id=?", vipId).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "更新失败")
			return
		}
		for _, v := range ratiaoSile {
			t := strings.Split(v, ",")
			venueId, _ := strconv.Atoi(t[0])
			gameVenueInfo := &models.GameVenue{}
			if b, err := session.Table("game_venues").Where("id=?", venueId).Get(gameVenueInfo); !b || err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "查询失败")
				return
			}
			typeId := gameVenueInfo.VenueType
			uuMap := map[string]interface{}{
				"ratio":   t[1],
				"updated": tools.NowMicro(),
			}
			userRebateSetting := &models.UserRebateSetting{}
			b, err := session.Table("user_rebate_settings").Where("vip_id=? and venue_id=?", vipId, venueId).Exist(userRebateSetting)
			if err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "查询失败")
				return
			}
			if b {
				if _, err := session.Table("user_rebate_settings").Where("vip_id=? and venue_id=? and type_id=?", vipId, venueId, typeId).Update(uuMap); err != nil {
					log.Logger.Error(err.Error())
					_ = session.Rollback()
					response.Err(c, "更新失败")
					return
				}
			} else {
				iMap := map[string]interface{}{
					"vip_id":    vipId,
					"venue_id":  venueId,
					"type_id":   typeId,
					"ratio":     t[1],
					"type_name": "",
					"created":   tools.NowMicro(),
				}
				if _, err := session.Table("user_rebate_settings").Insert(iMap); err != nil {
					log.Logger.Error(err.Error())
					_ = session.Rollback()
					response.Err(c, "更新失败")
					return
				}
			}
		}
		_ = session.Commit()
		response.Ok(c)
	},
}
