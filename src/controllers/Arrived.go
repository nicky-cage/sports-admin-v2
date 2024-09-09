package controllers

import (
	"encoding/json"
	"fmt"
	"net/url"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// Arrived 会员签到
var Arrived = struct {
	Update func(*gin.Context)
	Save   func(*gin.Context)
}{
	Update: func(c *gin.Context) {
		platform := request.GetPlatform(c)
		row := models.GetActivityArrived(platform) // 获取签到活动配置
		viewData := ViewData{"arrived": row}
		response.Render(c, "activities/_activity_arrived.html", viewData)
	},
	Save: func(c *gin.Context) { // 参数获取
		postedData := request.GetPostedData(c)
		jsonData, _ := url.QueryUnescape(postedData["data"].(string))
		row := models.ArrivedActivity{}
		err := json.Unmarshal([]byte(jsonData), &row)
		if err != nil {
			fmt.Println("保存签到配置失败: ", err)
			response.Err(c, "保存签到信息失败")
			return
		}
		if len(row.SignRewardRows) > 0 {
			bytes, err := json.Marshal(row.SignRewardRows)
			if err == nil {
				row.SignRewards = string(bytes)
			}
		}

		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		total, err := func() (int64, error) {
			if row.Id > 0 {
				return dbSession.Table("arrived_activities").ID(row.Id).Update(&row)
			}
			return dbSession.Table("arrived_activities").Insert(&row)
		}()
		if total == 0 || err != nil {
			fmt.Println("保存签到信息失败: ", err)
			response.Err(c, "保存信息失败")
			return
		}

		response.Ok(c)
	},
}
