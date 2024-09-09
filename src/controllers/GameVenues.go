package controllers

import (
	"fmt"
	"sports-admin/caches"
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

// GameVenues 场馆列表
var GameVenues = struct {
	//*ActionList
	*ActionCreate
	*ActionUpdate
	Saved     func(c *gin.Context)
	Stated    func(c *gin.Context)
	List      func(c *gin.Context)
	Games     func(c *gin.Context)
	Add       func(c *gin.Context)
	StateSave func(c *gin.Context)
	*ActionState
	*ActionSave
}{
	//ActionList: &ActionList{
	//	Model:    models.GameVenues,
	//	ViewFile: "game_venues/list.html",
	//	OrderBy: func(*gin.Context) string {
	//		return "id DESC"
	//	},
	//	Rows: func() interface{} {
	//		return &[]models.GameVenue{}
	//	},
	//	QueryCond: map[string]interface{}{
	//		"ename":    "%",
	//		"name":     "%",
	//		"type":     "=",
	//		"maintain": "=",
	//	},
	//	GetQueryCond: func(C *gin.Context) builder.Cond {
	//		cond := builder.NewCond()
	//		cond = cond.And(builder.Eq{"pid": 0})
	//		return cond
	//	},
	//	ExtendData: func(*gin.Context) pongo2.Context {
	//		return pongo2.Context{
	//			"venue_types":      consts.GameVenueTypes,
	//			"game_venues":      caches.GameVenues.All(),
	//			"suport_platforms": consts.GamePlatformTypes,
	//		}
	//	},
	//},
	ActionCreate: &ActionCreate{
		Model:    models.GameVenues,
		ViewFile: "game_venues/edit.html",
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.GameVenues,
		ViewFile: "game_venues/edit.html",
		Row: func() interface{} {
			return &models.GameVenue{}
		},
	},
	Saved: func(c *gin.Context) {
		var id string
		dataPost := request.GetPostedData(c)
		platform := request.GetPlatform(c)
		id = dataPost["id"].(string)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		var content string
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		admin := GetLoginAdmin(c)
		fsql := "select platform_rate,wallet,pid,name,is_online,maintain,wallet_name,code from game_venues where id='%s'"
		ffsql := fmt.Sprintf(fsql, id)
		res, err := dbSession.QueryString(ffsql)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "修改错误")
			return
		}

		if dataPost["platform_rate"] != nil {
			sort := dataPost["sort"].(string)
			rate := dataPost["platform_rate"].(string)
			remark := dataPost["remark"].(string)
			isOnline := dataPost["is_online"].(string)
			n_rate, ferr := strconv.ParseFloat(rate, 64)
			if ferr != nil {
				log.Err(ferr.Error())
				response.Err(c, "修改错误")
				return
			}
			updateTime := time.Now().Unix()
			real_rate := n_rate / 100
			str_rate := strconv.FormatFloat(real_rate, 'f', -1, 64)
			if isOnline == "1" {
				sql := "update game_venues set sort=?,platform_rate=? ,remark=? ,wallet=?, is_online=?,maintain=? ,updated=? where id=?"
				_, qerr := dbSession.Exec(sql, sort, real_rate, remark, 1, 1, 1, updateTime, id)
				if qerr != nil {
					log.Err(qerr.Error())
					response.Err(c, "修改错误")
					return
				}
				lsql := "update game_venues set wallet=?,is_online=?,maintain=?,platform_rate=? where pid=?"

				_, lerr := dbSession.Exec(lsql, 1, 1, 1, real_rate, id)
				if lerr != nil {
					log.Err(lerr.Error())
					response.Err(c, "修改错误")
					return
				}
				//修改电子场馆
				gsql := "update games set is_online=? where game_code=?"
				_, gerr := dbSession.Exec(gsql, 2, res[0]["code"])
				if gerr != nil {
					log.Err(gerr.Error())
				}
				if res[0]["maintain"] == "2" {
					content = res[0]["wallet_name"] + "状态由锁定变成正常," + dataPost["ename"].(string) + "场馆平台费由" + res[0]["platform_rate"] + "调整为" + str_rate
				} else {
					content = res[0]["wallet_name"] + "状态由下线变成正常," + dataPost["ename"].(string) + "场馆平台费由" + res[0]["platform_rate"] + "调整为" + str_rate
				}

			} else if isOnline == "2" {
				sql := "update game_venues set sort=?,platform_rate=? ,remark=? ,wallet=?, is_online=?,maintain=?, updated=? where id=?"
				_, qerr := dbSession.Exec(sql, sort, real_rate, remark, 2, res[0]["is_online"], 2, updateTime, id)
				if qerr != nil {
					log.Err(qerr.Error())
					response.Err(c, "修改错误")
					return
				}
				lsql := "update game_venues set wallet=?,maintain=?,platform_rate=? where pid=?"
				_, lerr := dbSession.Exec(lsql, 2, 2, real_rate, id)
				if lerr != nil {
					log.Err(lerr.Error())
					response.Err(c, "修改错误")
					return
				}
				//修改电子场馆
				gsql := "update games set is_online=? where game_code=?"
				_, gerr := dbSession.Exec(gsql, 1, res[0]["code"])
				if gerr != nil {
					log.Err(gerr.Error())
				}
				if res[0]["is_online"] == "1" {
					content = res[0]["wallet_name"] + "状态由正常变成锁定," + dataPost["ename"].(string) + "场馆平台费由" + res[0]["platform_rate"] + "调整为" + str_rate
				} else {
					content = res[0]["wallet_name"] + "状态由下线变成锁定," + dataPost["ename"].(string) + "场馆平台费由" + res[0]["platform_rate"] + "调整为" + str_rate
				}
			} else if isOnline == "3" {
				sql := "update game_venues set sort=?,platform_rate=? ,remark=? ,wallet=?, is_online=?,maintain=?, updated=? where id=?"
				_, qerr := dbSession.Exec(sql, sort, real_rate, remark, 2, 2, 1, updateTime, id)
				if qerr != nil {
					log.Err(qerr.Error())
					response.Err(c, "修改错误")
					return
				}

				lsql := "update game_venues set is_online=?,wallet=?,maintain=?,platform_rate=? where pid=?"

				_, lerr := dbSession.Exec(lsql, 2, 2, 1, real_rate, id)
				if lerr != nil {
					log.Err(lerr.Error())
					response.Err(c, "修改错误")
					return
				}
				//修改电子场馆
				gsql := "update games set is_online=? where game_code=?"
				_, gerr := dbSession.Exec(gsql, 1, res[0]["code"])
				if gerr != nil {
					log.Err(gerr.Error())
				}
				if res[0]["maintain"] == "2" {
					content = res[0]["wallet_name"] + "状态由锁定变成下线," + dataPost["ename"].(string) + "场馆平台费由" + res[0]["platform_rate"] + "调整为" + str_rate
				} else {
					content = res[0]["wallet_name"] + "状态由正常变成下线," + dataPost["ename"].(string) + "场馆平台费由" + res[0]["platform_rate"] + "调整为" + str_rate
				}
			}
			//删除电子游戏的缓存。
			codeKey := consts.EgameCacheKey + res[0]["code"]
			redis.Del(codeKey)

			idint, _ := strconv.Atoi(id)
			var venue_log models.GameVenueLog
			venue_log.Admin = admin.Name
			venue_log.Content = content
			venue_log.VenueName = res[0]["name"]
			venue_log.VenueId = uint8(idint)
			venue_log.Created = tools.NowMicro()
			venue_log.Remark = remark
			_, err4 := dbSession.Insert(&venue_log)
			if err4 != nil {
				log.Err(err4.Error())
				response.Err(c, "更新失败 "+err4.Error())
				return
			}

		} else {
			//场馆排序修改
			sort := dataPost["sort"].(string)
			sql := "update game_venues set sort=?  where id=?"
			_, err := dbSession.Exec(sql, sort, id)
			if err != nil {
				log.Err(err.Error())
				response.Err(c, "更新改错误")
			}
			//场馆排序也要删除缓存，
		}
		_, _ = redis.Del(consts.CkeyGameList).Result()
		//读取一次数据，填充数据
		_, _ = models.GetGameListAllHandler(platform)
		response.Ok(c)
	},
	Stated: func(c *gin.Context) {
		id := c.Query("id")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select id,name,ename,is_online,maintain from game_venues where id= " + id
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		response.Render(c, "game_venues/updated.html", pongo2.Context{"r": res[0]})
	},
	Games: func(c *gin.Context) {
		cond := builder.NewCond()
		name := c.DefaultQuery("name", "")
		ename := c.DefaultQuery("ename", "")
		online := c.DefaultQuery("is_online", "")
		venueType := c.DefaultQuery("venue_type", "")

		if len(name) > 0 {
			cond = cond.And(builder.Eq{"name": name})
		}
		if len(ename) > 0 {
			cond = cond.And(builder.Eq{"ename": ename})
		}
		if len(online) > 0 {
			switch online {
			case "1":
				cond = cond.And(builder.Eq{"is_online": 1}).And(builder.Neq{"maintain": 2})
			case "2":
				cond = cond.And(builder.Eq{"maintain": 2})
			case "3":
				cond = cond.And(builder.Eq{"is_online": 2}).And(builder.Neq{"maintain": 2})
			}

		}
		if len(venueType) > 0 {
			cond = cond.And(builder.Eq{"venue_type": venueType})
		}
		res := make([]models.GameVenue, 0)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		cond = cond.And(builder.Neq{"pid": 0})

		err := dbSession.Table("game_venues").Where(cond).OrderBy("sort ASC").Find(&res)
		if err != nil {
			log.Err(err.Error())
			return
		}
		SetLoginAdmin(c)
		response.Render(c, "game_venues/_games.html", pongo2.Context{"rows": res})
	},
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		name := c.DefaultQuery("name", "")
		ename := c.DefaultQuery("ename", "")
		online := c.DefaultQuery("is_online", "")
		//venueType := c.DefaultQuery("venue_type", "")

		if len(name) > 0 {
			cond = cond.And(builder.Eq{"name": name})
		}
		if len(ename) > 0 {
			cond = cond.And(builder.Eq{"ename": ename})
		}
		platform := request.GetPlatform(c)
		if len(online) > 0 {
			switch online {
			case "1":
				cond = cond.And(builder.Eq{"is_online": 1}).And(builder.Neq{"maintain": 2})
			case "2":
				cond = cond.And(builder.Eq{"maintain": 2})
			case "3":
				cond = cond.And(builder.Eq{"is_online": 2}).And(builder.Neq{"maintain": 2})
			}
		}
		//if len(venueType) > 0 {
		//	if venueTypeId, err := strconv.Atoi(venueType); err != nil {
		//		log.Logger.Error("格式化场馆类型出错: ", err)
		//	} else {
		//		cond = cond.And(builder.Eq{"venue_type": venueTypeId})
		//	}
		//}

		//cond = cond.And(builder.Eq{"venue_type": 2})
		res := make([]models.GameVenue, 0)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		cond = cond.And(builder.Eq{"pid": 0})
		err := dbSession.Table("game_venues").Where(cond).OrderBy("sort ASC").Find(&res)
		if err != nil {
			log.Err(err.Error())
			return
		}
		viewData := pongo2.Context{
			"rows":             res,
			"suport_platforms": consts.GamePlatformTypes,
			"game_venues":      caches.GameVenues.All(platform),
			"venue_types":      consts.GameVenueTypes,
		}
		SetLoginAdmin(c)
		if request.IsAjax(c) {
			response.Render(c, "game_venues/_ilist.html", viewData)
		} else {
			response.Render(c, "game_venues/list.html", viewData)
		}
	},
	Add: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		games := caches.GameVenues.All(platform)
		response.Render(c, "game_venues/created.html", pongo2.Context{"games": games})
	},
	StateSave: func(c *gin.Context) {
		var linef int32
		var start string
		var end string
		var lowerSql string
		var content string
		postData := request.GetPostedData(c)
		platform := request.GetPlatform(c)
		id := postData["id"].(string)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		remark := postData["remark"].(string)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		online := postData["is_online"].(string)
		name := postData["name"].(string)

		rowGameVenue := models.GameVenue{}
		bOk, rerr := dbSession.Table("game_venues").ID(id).Get(&rowGameVenue)
		if rerr != nil || !bOk {
			errMsg := "更新失败 "
			if rerr != nil {
				errMsg += rerr.Error()
			}
			log.Err(errMsg)
			response.Err(c, errMsg)
			return
		}

		//当is_online 等于2时，维护
		if online == "2" {
			if postData["created"].(string) != "" {
				created := postData["created"].(string)
				areas := strings.Split(created, " - ")
				start = areas[0]
				end = areas[1]
			}

			lowerSql = "update game_venues set is_online=? ,wallet=?,maintain = ?,remark=?,maintain_start_at=?,maintain_end_at=? where id = ?"
			_, err := dbSession.Exec(lowerSql, rowGameVenue.IsOnline, 2, "2", remark, start, end, id)
			if err != nil {
				log.Err(err.Error())
				response.Err(c, "更新失败 "+err.Error())
				return
			}

			if rowGameVenue.IsOnline == 1 {
				content = name + "状态由正常修改为维护"
			} else {
				content = name + "状态由下线修改为维护"
			}
			//子动父动,查询相同pid下的维护状态，都为维护则，钱包维护
			allsql := "select maintain from game_venues where pid=%d"
			allsqll := fmt.Sprintf(allsql, rowGameVenue.Pid)
			aRes, _ := dbSession.QueryString(allsqll)
			num := len(aRes)

			var maintainLen = 1

			for _, v := range aRes {
				if v["maintain"] == "2" {
					if maintainLen < 3 {
						maintainLen++
					}
				}
			}
			if num == maintainLen {
				psql := "update game_venues set maintain=2,is_online=1,wallet=2 where id=?"
				_, err := dbSession.Exec(psql, rowGameVenue.Pid)
				if err != nil {
					log.Err(err.Error())
				}
			}
			//关联电子游戏
			gint := strings.Index(rowGameVenue.Ename, "electronic")
			if gint > 0 {
				gsql := "update games set is_online=? where game_code=?"
				_, err := dbSession.Exec(gsql, 1, rowGameVenue.Code)
				if err != nil {
					log.Err(err.Error())
				}
			}
		} else {
			var statusName string
			if online == "3" {
				linef = 2
				statusName = "下线"
			} else {
				statusName = "正常"
				linef = 1
				//下级场馆想要恢复为正常的时候

			}

			lowerSql = "update game_venues set is_online=?,wallet=?,maintain = ?,remark=? where id = ?"

			_, err := dbSession.Exec(lowerSql, linef, 1, 1, remark, id)
			if err != nil {
				log.Err(err.Error())
				response.Err(c, "更新失败 "+err.Error())
				return
			}
			if rowGameVenue.Maintain == 2 {
				content = name + "状态由维护修改为" + statusName
			} else {
				if rowGameVenue.IsOnline == 1 {
					content = name + "状态由正常修改为" + statusName
				} else {
					content = name + "状态由下线修改为" + statusName
				}
			}

			allsql := "select is_online from game_venues where pid=%d"
			allsqll := fmt.Sprintf(allsql, rowGameVenue.Pid)
			aRes, _ := dbSession.QueryString(allsqll)
			num := len(aRes)
			var maintainLen int = 1
			// maintainLen = 1

			if online == "3" {
				//子动父动,查询相同pid下的下线状态，都为下线则，钱包下线
				for _, v := range aRes {
					if v["is_online"] == "2" {
						if maintainLen < 3 {
							maintainLen++
						}
					}
				}

				if num == maintainLen {
					psql := "update game_venues set is_online=2,wallet=2,maintain=1 where id=?"
					_, err := dbSession.Exec(psql, rowGameVenue.Pid)
					if err != nil {
						log.Err(err.Error())
					}
				}
				//关联电子游戏
				gint := strings.Index(rowGameVenue.Ename, "electronic")
				if gint > 0 {
					gsql := "update games set is_online=? where game_code=?"
					_, err := dbSession.Exec(gsql, 1, rowGameVenue.Code)
					if err != nil {
						log.Err(err.Error())
					}
				}
			} else {
				for _, v := range aRes {
					if v["is_online"] == "1" {
						if maintainLen < 3 {
							maintainLen++
						}
					}
				}
				if num == maintainLen {
					psql := "update game_venues set is_online=1,wallet=1,maintain=1 where id=?"
					_, err := dbSession.Exec(psql, rowGameVenue.Pid)
					if err != nil {
						log.Err(err.Error())
					}
				}
				//关联电子游戏
				gint := strings.Index(rowGameVenue.Ename, "electronic")
				if gint > 0 {
					gsql := "update games set is_online=? where game_code=?"
					_, err := dbSession.Exec(gsql, 2, rowGameVenue.Code)
					if err != nil {
						log.Err(err.Error())
					}
				}
			}
		}

		//删除缓存key
		_, _ = redis.Del(consts.CkeyGameList).Result()
		//读取一次数据，填充数据
		_, _ = models.GetGameListAllHandler(platform)
		//删除电子游戏常
		codeKey := consts.EgameCacheKey + rowGameVenue.Code
		redis.Del(codeKey)
		//将数据插入日志
		idint, _ := strconv.Atoi(id)
		var venue_log models.GameVenueLog
		admin := GetLoginAdmin(c)
		venue_log.Admin = admin.Name
		venue_log.Content = content
		venue_log.VenueName = name
		venue_log.VenueId = uint8(idint)
		venue_log.Created = tools.NowMicro()
		_, err4 := dbSession.Insert(&venue_log)
		if err4 != nil {
			log.Err(err4.Error())
			response.Err(c, "更新失败 "+err4.Error())
			return
		}
		response.Ok(c)
	},
	ActionState: &ActionState{
		Model: models.GameVenues,
		Field: "is_online",
		StateAfter: func(c *gin.Context) {
			platform := request.GetPlatform(c)
			redis := common.Redis(platform)
			defer common.RedisRestore(platform, redis)
			_, _ = redis.Del(consts.CkeyGameList).Result() //删除缓存key
			_, _ = models.GetGameListAllHandler(platform)  //读取一次数据，填充数据
		},
	},
	ActionSave: &ActionSave{
		Model: models.GameVenues,
	},
}
