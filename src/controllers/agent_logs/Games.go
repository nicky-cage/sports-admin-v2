package agent_logs

import (
	"context"
	"fmt"
	"reflect"
	"sports-admin/controllers/base_controller"
	common "sports-common"
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
)

func (ths *AgentLogs) Games(c *gin.Context) {
	var startAt int64
	var endAt int64
	if value, exists := c.GetQuery("created"); !exists {
		currentTime := time.Now().Unix()
		startAt = currentTime - currentTime%86400
		endAt = startAt + 86400
	} else {
		areas := strings.Split(value, " - ")
		startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
		endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
	}

	venueType := c.DefaultQuery("venue_type", "")
	username := c.DefaultQuery("username", "")
	topName := c.DefaultQuery("top_name", "")
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)
	pageSize := 15

	platform := request.GetPlatform(c)
	esIndexName := platform + "_wagers"
	esClient, err := es.GetClientByPlatform(platform)
	//	res, err := esClient.Search(esIndexName).Do(context.Background())
	if err != nil {
		response.Err(c, "es获取列表数据无响应")
		return
	}
	defer esClient.Stop()

	boolQuery := elastic.NewBoolQuery()
	if len(username) > 0 {
		boolQuery.Must(elastic.NewMatchQuery("username", username))
	}
	if len(venueType) > 0 {
		arr := strings.Split(venueType, "-")
		boolQuery.Must(elastic.NewMatchQuery("game_code", arr[0]))
		boolQuery.Must(elastic.NewMatchQuery("game_type", arr[1]))
	}
	if len(topName) > 0 {
		boolQuery.Must(elastic.NewMatchQuery("top_name", topName))
	}

	boolQuery.Filter(elastic.NewRangeQuery("created_at").Gte(startAt).Lte(endAt))
	res, err := esClient.Search(esIndexName).Query(boolQuery).From((page-1)*pageSize).Size(pageSize).Sort("created_at", false).Do(context.Background())
	if err != nil {
		response.Err(c, "es获取列表数据无响应")
		return
	}

	total := res.Hits.TotalHits.Value
	var typ models.WagerRecord
	var rows []models.WagerRecord
	if total > 0 {
		for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
			t := item.(models.WagerRecord)
			rows = append(rows, t)
		}
	}
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select top_name from users where id=%d"
	for k, res := range rows {
		sqll := fmt.Sprintf(sql, res.UserId)
		sres, err := dbSession.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
			return
		}
		if len(sres) > 0 {
			rows[k].TopName = sres[0]["top_name"]
		}
	}

	vsql := "select code,venue_type,name from game_venues where pid!=0"
	vRes, err := dbSession.QueryString(vsql)
	if err != nil {
		log.Err(err.Error())
	}
	tSql := "select * from user_tags"
	tRes, err := dbSession.QueryString(tSql)
	if err != nil {
		log.Err(err.Error())
	}
	base_controller.SetLoginAdmin(c)
	if request.IsAjax(c) {
		response.Render(c, "agent_logs/_games.html", pongo2.Context{"res": rows, "total": total, "venue_type": vRes, "tags": tRes})
	} else {
		response.Render(c, "agent_logs/index.html", pongo2.Context{"res": rows, "total": total, "venue_type": vRes, "tags": tRes})
	}
}
