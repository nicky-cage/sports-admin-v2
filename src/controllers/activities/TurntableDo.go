package activities

import (
	"encoding/json"
	"fmt"
	"net/url"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

func (*Activities) TurntableDo(c *gin.Context) { // 参数获取
	postedData := request.GetPostedData(c)
	jsonData, _ := url.QueryUnescape(postedData["data"].(string))
	row := models.TurntableDetail{}
	err := json.Unmarshal([]byte(jsonData), &row)
	if err != nil {
		fmt.Println("保存转般失败: ", err)
		response.Err(c, "保存幸运转盘信息失败")
		return
	}
	if len(row.ConditionSettingRows) > 0 {
		res, err := json.Marshal(row.ConditionSettingRows)
		if err == nil {
			row.ConditionSetting = string(res)
		}
	}
	if len(row.PrizeSettingsRows) > 0 {
		res, err := json.Marshal(row.PrizeSettingsRows)
		if err == nil {
			row.PrizeSettings = string(res)
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
	if total, err := func() (int64, error) {
		if row.Id > 0 {
			return dbSession.Table("turntables").ID(row.Id).Update(&row)
		}
		return dbSession.Table("turntables").Insert(&row)
	}(); total == 0 {
		response.Err(c, "保存信息无效")
		return
	} else if err != nil {
		log.Err(err.Error())
		fmt.Println("保存转盘信息失败: ", err)
		response.Err(c, "保存信息失败")
		return
	}

	response.Ok(c)
}
