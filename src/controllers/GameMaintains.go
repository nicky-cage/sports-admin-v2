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

// GameMaintains 维护设置
var GameMaintains = struct {
	List func(c *gin.Context)
	*ActionCreate
	Update func(c *gin.Context)
	Save   func(c *gin.Context)
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		sate := c.DefaultQuery("state", "")
		if len(sate) > 0 {
			cond = cond.And(builder.Eq{"state": sate})
		}
		venueId := c.DefaultQuery("venue_id", "")
		if len(venueId) > 0 {
			cond = cond.And(builder.Eq{"venue_id": venueId})
			cond = cond.And(builder.Eq{"pid": 0})
		}
		var list []models.GameVenueMaintain
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		err := dbSession.Table("game_venue_maintains").Join("INNER", "game_venues", "game_venues.id=game_venue_maintains.venue_id").Where(cond).Find(&list)
		if err != nil {
			log.Err(err.Error())
			return
		}
		viewDate := pongo2.Context{
			"rows": list,
		}
		response.Render(c, "game_maintains/_list.html", viewDate)
	},
	ActionCreate: &ActionCreate{
		Model:    models.GameVenueMaintains,
		ViewFile: "game_maintains/edit.html",
	},
	Update: func(c *gin.Context) {
		id := c.Query("id")
		platform := request.GetPlatform(c)
		var list []models.GameVenueMaintain
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		err := dbSession.Table("game_venue_maintains").Join("INNER", "game_venues", "game_venues.id=game_venue_maintains.venue_id").Where("venue_id=?", id).Find(&list)
		if err != nil {
			log.Err(err.Error())
			return
		}
		admin := GetLoginAdmin(c)

		response.Render(c, "game_maintains/edit.html", pongo2.Context{"r": list[0], "adimin": admin.Name})
	},
	Save: func(c *gin.Context) {

		data := request.GetPostedData(c)
		maintainTime := data["maintain_time"].(string)
		times := strings.Split(maintainTime, " - ")
		strStart := times[0]
		strEnd := times[1]
		timeStart := tools.GetTimeStampByString(strStart)
		timeEnd := tools.GetTimeStampByString(strEnd)

		platform := request.GetPlatform(c)
		admin := GetLoginAdmin(c)
		//(*data)["admin_id"] = admin.Id

		state := data["state"].(string)
		venueState, _ := strconv.Atoi(state)
		gameId := data["venue_id"].(string)
		venueId, _ := strconv.Atoi(gameId)

		var list = make(map[string]interface{})
		list["admin_name"] = admin.Name
		list["created"] = tools.NowMicro()
		list["Remark"] = data["remark"].(string)
		list["time_start"] = uint32(timeStart)
		list["time_end"] = uint32(timeEnd)
		list["state"] = uint8(venueState)
		list["venue_id"] = uint32(venueId)

		_, _ = models.GameVenueMaintainLogs.Create(platform, list) //创建一条记录

		//还要将维护状态之类的 存到redis
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)

		var wallet string
		var maintain string
		if state == "1" {
			//维护
			wallet = "1"
			maintain = "1"
		} else {
			wallet = "2"
			maintain = "2"
		}

		// 需要将场馆下的所有游戏的维护状态打开和 关闭钱包，和场馆的状态
		vSql := "update game_venues set maintain=?,wallet=?,updated=? where id=?"
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		_, err := dbSession.Exec(vSql, maintain, wallet, time.Now().Unix(), gameId)
		if err != nil {
			log.Err(err.Error())
			return
		}
		lSql := "update game_venues set  maintain=?,wallet=? where pid=?"
		_, lerr := dbSession.Exec(lSql, maintain, wallet, gameId)
		if lerr != nil {
			log.Err(lerr.Error())
			return
		}
		//更改视图表的另一个
		mSql := "update game_venue_maintains set time_start=?,time_end=?,state=?,admin_name=?,remark=? where venue_id=?"
		_, verr := dbSession.Exec(mSql, uint32(timeStart), uint32(timeEnd), state, admin.Name, data["remark"].(string), gameId)
		if verr != nil {
			log.Err(verr.Error())
			return
		}

		_, _ = redis.Del(consts.CkeyGameList).Result()
		//读取一次数据，填充数据
		_, _ = models.GetGameListAllHandler(platform)
		response.Ok(c)
	},
}
