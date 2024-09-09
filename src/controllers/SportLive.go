package controllers

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"time"

	"github.com/flosch/pongo2"
	gin "github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var SportLives = struct {
	List            func(*gin.Context)
	Stop            func(*gin.Context)
	Manage          func(*gin.Context)
	ManageWords     func(*gin.Context)
	ManageWordsSave func(*gin.Context)
}{
	List: func(c *gin.Context) {
		username := c.Query("username")
		cond := builder.NewCond()
		platform := request.GetPlatform(c)
		onlineUserNames := GetAllOnlineUser(platform)

		if len(username) != 0 {
			test := false
			var name string
			for _, v := range onlineUserNames {
				if v == username {
					test = true
					name = v
				}
			}
			if test {
				cond = cond.And(builder.Eq{"username": name})
			} else {
				cond = cond.And(builder.Eq{"username": ""})
			}
		}
		nickName := c.Query("nickname")
		if len(nickName) != 0 {
			cond = cond.And(builder.Eq{"nickname": nickName})
		}
		//搜索，展示

		if len(nickName) == 0 && len(username) == 0 {

			cond = cond.And(builder.In("username", onlineUserNames))
		}
		limit, offset := request.GetOffsets(c)
		db := common.Mysql(platform)
		defer db.Close()
		var list []models.User
		db.Table("users").Where(cond).Cols("username,id,nickname").Limit(limit, offset).Find(&list)

		// redis 禁言列表
		redis := common.Redis(platform)
		redisMap, _ := redis.HGetAll(platform + "sport_live:gag").Result()

		var data []models.User
		for _, v := range list {
			// 禁言。
			if timeMap, ok := redisMap[v.Username]; ok { // 2否，1是
				//需要处理过去时间。 删掉之前的。
				temp := tools.GetTimeStampByString(timeMap)

				if time.Now().Unix() > temp {
					v.Gag = 2
					redis.HDel(platform+"sport_live:gag", v.Username)
				} else {
					v.Gag = 1 //是再禁言中
				}

			} else {
				v.Gag = 2 //未禁言
			}
			data = append(data, v)
		}

		viewData := pongo2.Context{
			"rows":  data,
			"total": len(onlineUserNames),
		}
		viewFile := "sport_live/live.html"
		if request.IsAjax(c) {
			viewFile = "sport_live/_list.html"
		}
		response.Render(c, viewFile, viewData)
	},
	Stop: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		id := c.Query("id")
		value := c.Query("value")
		redis := common.Redis(platform)
		if value == "1" { //原来是禁言状态
			redis.HDel(platform+"sport_live:gag", id)
		} else { //原来是未禁言状态
			temp := time.Now().AddDate(0, 0, 1).Format("2006-01-02 15:04:05")
			redis.HSet(platform+"sport_live:gag", id, temp)
		}
		response.Ok(c)
	},
	Manage: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		db := common.Mysql(platform)
		defer db.Close()
		var list []models.FilterPhrase
		limit, offset := request.GetOffsets(c)
		err := db.Table("filter_phrases").Limit(limit, offset).Find(&list)
		if err != nil {
			log.Err(err.Error())
		}
		var lists []models.FilterPhrase
		db.Table("filter_phrases").Find(&lists)
		viewData := pongo2.Context{
			"row":   list,
			"total": len(lists),
		}
		viewFile := "sport_live/manage.html"
		if request.IsAjax(c) {
			viewFile = "sport_live/_manage.html"
		}
		response.Render(c, viewFile, viewData)
	},
	ManageWords: func(c *gin.Context) {
		//修改
		platform := request.GetPlatform(c)
		id := c.Query("id")
		var data models.FilterPhrase
		if id != "" {
			db := common.Mysql(platform)
			defer db.Close()
			db.Table("filter_phrases").Where("id=?", id).Get(&data)
		}

		response.Render(c, "sport_live/created.html", pongo2.Context{"row": data})
	},
	ManageWordsSave: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		data := request.GetPostedData(c)
		num := data["method"].(string)
		db := common.Mysql(platform)
		defer db.Close()
		switch num {
		case "1": //插入
			types := data["type"].(string)
			name := data["name"].(string)
			sql := "insert into filter_phrases(type,name) values(?,?)"
			db.Exec(sql, types, name)
		case "2": //修改
			types := data["type"].(string)
			name := data["name"].(string)
			id := data["id"].(string)
			sql := "update filter_phrases set type=?,name=? where id=?"
			db.Exec(sql, types, name, id)
		case "3": //删除
			id := data["id"].(string)
			sql := "delete from filter_phrases where id=? "
			db.Exec(sql, id)
		}
		response.Ok(c)
	},
}
