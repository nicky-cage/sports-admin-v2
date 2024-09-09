package controllers

import (
	"encoding/json"
	"net/url"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var RedEnvelopes = struct {
	List func(c *gin.Context)
	Save func(c *gin.Context)
}{
	List: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		row := models.GetActivityRed(platform)
		response.Render(c, "activities/_activity_red.html", ViewData{"red": row})
	},
	Save: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		jsonData, _ := url.QueryUnescape(postedData["data"].(string))
		row := models.RedEnvelopeRain{}
		err := json.Unmarshal([]byte(jsonData), &row)
		if err != nil {
			response.Err(c, "保存信息失败: "+err.Error())
			return
		}

		if len(row.GrabCondRows) > 0 {
			res, err := json.Marshal(row.GrabCondRows)
			if err == nil {
				row.GrabConds = string(res)
			}
		}
		if len(row.GrabTimeRows) > 0 {
			res, err := json.Marshal(row.GrabTimeRows)
			if err == nil {
				row.GrabTimes = string(res)
			}
		}
		if len(row.RankingListRows) > 0 {
			res, err := json.Marshal(row.RankingListRows)
			if err == nil {
				row.RankingList = string(res)
			}
		}
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		total, err := func() (int64, error) {
			if row.Id > 0 {
				return dbSession.Table("red_envelope_rains").ID(row.Id).Update(&row)
			}
			return dbSession.Table("red_envelope_rains").Insert(&row)
		}()

		if total == 0 {
			response.Err(c, "保存信息失败: ")
			return
		}
		if err != nil {
			response.Err(c, "保存信息失败: "+err.Error())
			return
		}

		response.Ok(c)
	},
}
