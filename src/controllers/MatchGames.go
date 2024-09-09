package controllers

import (
	"fmt"
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"xorm.io/builder"
)

type MatchGameTemp struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct {
		CurrentPage int `json:"current_page"`
		List        []struct {
			MatchId   int    `json:"match_id" `   //编号
			MatchTime string `json:"match_time"`  //开始时间
			MatchName string `json:"match_name"`  //联赛名称
			HomeName  string `json:"home_name"`   //主队名称
			AwayName  string `json:"away_name"`   //客队名称
			SportType int    `json:"sport_type" ` //运动类型 1 足球，2篮球
		} `json:"list"`
		Total int `json:"total"`
	} `json:"data"`
}

var MatchGames = struct {
	Created func(c *gin.Context)
	*ActionList
	*ActionDelete
	*ActionSave
}{
	Created: func(c *gin.Context) {
		created := c.Query("start_time")
		param := req.Param{}
		starts := time.Now().Format("2006-01-02")
		param["start_time"] = starts + " 00:00:00"
		param["end_time"] = starts + " 23:59:59"

		//当有时间搜索
		if created != "" {
			areas := strings.Split(created, " - ")
			param["start_time"] = areas[0] + " 00:00:00"
			param["end_time"] = areas[1] + " 23:59:59"
		}
		Nami := make([]interface{}, 0)
		Bti := make([]interface{}, 0)

		// 从接口获取数据
		header := req.Header{
			"Accept": "application/json",
		}
		//获取已匹配赛事的id
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		start := time.Now().Format("2006-01-02")
		startTime, _ := time.Parse("2006-01-02 15:04:05", start)
		sql := "select nami_id,bti_id from match_games where nami_start_time>%d"
		sqll := fmt.Sprintf(sql, startTime.Unix())
		matchRes, err := dbSession.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
		}
		var btiString string
		var namiString string
		for _, v := range matchRes {
			btiString = btiString + v["bti_id"] + "."
			namiString = namiString + v["nami_id"] + "."
		}
		req.Debug = true
		req.SetTimeout(50 * time.Second)
		baseRegisterUrl := config.Get("internal.internal_member_service") + config.Get("internal_api.match_game_url")
		registerUrl := baseRegisterUrl

		sportType := c.Query("type")
		var num int
		var account int
		//体育类型请求
		if sportType != "" {
			//当有参数时, 循环1次   2 <=2 sport_id=2
			temp, _ := strconv.Atoi(sportType)
			num = temp
			account = temp
		} else {
			//当没有请求参数时，循环2次，
			num = 2
			account = 1
		}
		var namiPage string
		var btiPage string
		//只查询某个请求。
		namiPage = c.Query("nami_page")
		if namiPage != "" {
			param["current_page"] = namiPage
			//当时nami分页,mami 请求一次。 bti 循环应该不动。或者 都请求了。 但是值不插入。也没毛病。

		}
		btiPage = c.Query("bti_page")
		if btiPage != "" {
			param["current_page"] = btiPage
		}
		param["page_size"] = 15
		var namiTotal int
		var btiTotal int
		for i := account; i <= num; i++ {
			param["data_type"] = 1
			param["sport_id"] = i
			re, err := req.Post(registerUrl, header, param)
			if err != nil {
				response.Err(c, "系统异常")
				log.Err(err.Error())
				return
			}
			res := MatchGameTemp{}
			_ = re.ToJSON(&res)
			namiTotal = namiTotal + res.Data.Total
			for _, v := range res.Data.List {
				v.SportType = i
				//还要判断id是否存在，存在不插入，不存在插入。
				if !strings.Contains(namiString, strconv.Itoa(v.MatchId)) && len(Nami) < 15 {
					Nami = append(Nami, v)
				}
			}
			//多个nami的篮球结果， 多个足球结果
		}

		for o := account; o <= num; o++ {
			param["data_type"] = 2
			param["sport_id"] = o
			re, err := req.Post(registerUrl, header, param)
			if err != nil {
				response.Err(c, "系统异常")
				log.Err(err.Error())
				return
			}
			res := MatchGameTemp{}
			err = re.ToJSON(&res)
			if err != nil {
				log.Err(err.Error())
			}
			btiTotal = btiTotal + res.Data.Total
			for _, v := range res.Data.List {
				v.SportType = o
				if !strings.Contains(btiString, strconv.Itoa(v.MatchId)) && len(Bti) < 15 {
					Bti = append(Bti, v)
				}
			}
		}

		if created != "" {
			temp := map[string]interface{}{}
			temp["nami"] = Nami
			temp["nami_total"] = namiTotal
			temp["bti"] = Bti
			temp["bti_total"] = btiTotal
			temp["nami_page"] = namiPage
			temp["bti_page"] = btiPage
			response.Result(c, temp)
			return
		} else {
			response.Render(c, "match_games/list.html", pongo2.Context{"nami": Nami, "bti": Bti, "nami_total": namiTotal, "bti_total": btiTotal})
		}

	},
	ActionList: &ActionList{
		Model:    models.MatchGames,
		ViewFile: "match_games/i_list.html",
		Rows: func() interface{} {
			return &[]models.MatchGame{}
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			cond := builder.NewCond()
			if league_name, ok := c.GetQuery("league_name"); ok {
				cond = cond.And(builder.Like{"nami_league_name", league_name})
			}
			//if start_time, ok := c.GetQuery("start_time"); ok {
			//	value := strings.Split(start_time, " - ")
			//	start, _ := time.Parse("2006-01-02 15:04:05", value[0])
			//	end, _ := time.Parse("2006-01-02 15:04:05", value[1])
			//	if start != end {
			//		cond = cond.And(builder.Gte{"nami_start_time": start}).And(builder.Lte{"nami_start_time": end})
			//	}
			//}
			if team_name, ok := c.GetQuery("team_name"); ok {
				cond = cond.And(builder.Like{"nami_home_team", team_name})
			}
			return cond
		},
		QueryCond: map[string]interface{}{
			"type":             "=",
			"nami_league_name": "%",
		},
	},
	ActionSave: &ActionSave{
		Model: models.MatchGames,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			admin := GetLoginAdmin(c)
			(*m)["admin"] = admin.Name
			(*m)["risk_time"] = time.Now().Format("2006-01-02 15:04:05")
			if (*m)["type"].(string) == "篮球" {
				(*m)["type"] = "2"
			} else {
				(*m)["type"] = "1"
			}
			nami_time, _ := time.Parse("2006-01-02 15:04:05", (*m)["nami_start_time"].(string))
			(*m)["nami_start_time"] = nami_time.Unix()
			bti_time, _ := time.Parse("2006-01-02 15:04:05", (*m)["bti_start_time"].(string))
			(*m)["bti_start_time"] = bti_time.Unix()
			return nil
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.MatchGames,
	},
}
