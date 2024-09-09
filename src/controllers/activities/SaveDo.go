package activities

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (*Activities) SaveDo(c *gin.Context) {
	postedData := request.GetPostedData(c)
	myClient := common.Mysql(request.GetPlatform(c))
	defer myClient.Close()

	activityTimeSlice := strings.Split(postedData["activity_time"].(string), " - ")
	activityShowTimeSlice := strings.Split(postedData["activity_show_time"].(string), " - ")
	gameCodeList := make([]string, 0)
	gameCodeListLen := len(postedData)
	for i := 0; i < gameCodeListLen; i++ {
		tempI := strconv.Itoa(i)
		v, ok := postedData["game_code_list["+tempI+"]"].(string)
		if ok {
			gameCodeList = append(gameCodeList, v)
		}
	}
	if postedData["content_form"].(string) == "2" {
		if postedData["join_type"].(string) == "1" { //手动参与
			if len(gameCodeList) == 0 {
				response.Err(c, "常规内容-手动参与-请选择场馆")
				return
			}
		}
	}
	newGameCodeList := make([]string, 0)
	for _, v := range gameCodeList {
		gameVenusInfo := &models.GameVenue{}
		_, _ = myClient.Table("game_venues").Where("id=?", v).Get(gameVenusInfo)
		temp := gameVenusInfo.Code + "-" + strconv.Itoa(int(gameVenusInfo.VenueType))
		newGameCodeList = append(newGameCodeList, temp)
	}
	gameCodeListStr := strings.Join(newGameCodeList, ",")
	aMap := map[string]interface{}{
		"activity_type":      postedData["activity_type"],
		"join_type":          postedData["join_type"],
		"activity_label":     postedData["activity_label"],
		"content_form":       postedData["content_form"],
		"special_offer":      postedData["special_offer"],
		"game_code_list":     gameCodeListStr,
		"title":              postedData["title"],
		"web_topic_url":      postedData["web_topic_url"],
		"mobile_topic_url":   postedData["mobile_topic_url"],
		"application_cycle":  postedData["application_cycle"],
		"lowest_level":       postedData["lowest_level"],
		"sort":               postedData["sort"],
		"activity_amount":    postedData["activity_amount"],
		"run_water":          postedData["run_water"],
		"give_rate":          postedData["give_rate"],
		"give_money_max":     postedData["give_money_max"],
		"start_at":           tools.GetTimeStampByString(activityTimeSlice[0]),
		"end_at":             tools.GetTimeStampByString(activityTimeSlice[1]),
		"show_time_start":    tools.GetTimeStampByString(activityShowTimeSlice[0]),
		"show_time_end":      tools.GetTimeStampByString(activityShowTimeSlice[1]),
		"activity_pic":       postedData["activity_pic"],
		"web_list_pic":       postedData["web_list_pic"],
		"app_h5_list_pic":    postedData["app_h5_list_pic"],
		"web_main_pic":       postedData["web_main_pic"],
		"web_background_pic": postedData["web_background_pic"],
		"app_h5_main_pic":    postedData["app_h5_main_pic"],
		"details":            postedData["details"],
		"mobile_details":     postedData["mobile_details"],
		"code":               postedData["code"],
	}
	if postedData["content_form"].(string) == "1" { //专题内容
		if postedData["activity_pic"].(string) == "" {
			delete(aMap, "activity_pic")
		}
		if postedData["web_list_pic"].(string) == "" {
			delete(aMap, "web_list_pic")
		}
		if postedData["app_h5_list_pic"].(string) == "" {
			delete(aMap, "app_h5_list_pic")
		}
		delete(aMap, "join_type")
		delete(aMap, "game_code_list")
		delete(aMap, "application_cycle")
		delete(aMap, "lowest_level")
		delete(aMap, "web_main_pic")
		delete(aMap, "web_background_pic")
		delete(aMap, "app_h5_main_pic")
		delete(aMap, "activity_time")
		delete(aMap, "activity_amount")
		delete(aMap, "give_rate")
		delete(aMap, "give_money_max")
		delete(aMap, "run_water")
		delete(aMap, "details")
		delete(aMap, "mobile_details")
	} else { //常规内容
		if postedData["activity_pic"].(string) == "" {
			delete(aMap, "activity_pic")
		}
		if postedData["web_list_pic"].(string) == "" {
			delete(aMap, "web_list_pic")
		}
		if postedData["app_h5_list_pic"].(string) == "" {
			delete(aMap, "app_h5_list_pic")
		}
		if postedData["web_main_pic"].(string) == "" {
			delete(aMap, "web_main_pic")
		}
		if postedData["web_background_pic"].(string) == "" {
			delete(aMap, "web_background_pic")
		}
		if postedData["app_h5_main_pic"].(string) == "" {
			delete(aMap, "app_h5_main_pic")
		}
		if postedData["join_type"].(string) == "1" { //手动参与
			delete(aMap, "web_topic_url")
			delete(aMap, "mobile_topic_url")
		} else { //自动参与
			delete(aMap, "web_topic_url")
			delete(aMap, "mobile_topic_url")
			delete(aMap, "game_code_list")
			delete(aMap, "application_cycle")
			delete(aMap, "lowest_level")
			delete(aMap, "activity_time")
			delete(aMap, "activity_amount")
			delete(aMap, "give_rate")
			delete(aMap, "give_money_max")
			delete(aMap, "run_water")
		}
	}
	if postedData["method"].(string) == "add" {
		aMap["created"] = tools.NowMicro()
		aMap["updated"] = tools.NowMicro()
		aMap["activity_id"] = tools.RandInt64(10000, 99999)
		if _, err := myClient.Table("activities").Insert(aMap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "新增活动失败")
			return
		}
	} else {
		idStr := postedData["id"].(string)
		id, _ := strconv.Atoi(idStr)
		aMap["updated"] = tools.NowMicro()
		if _, err := myClient.Table("activities").Where("id=?", id).Update(aMap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "编辑活动失败")
			return
		}
	}
	response.Ok(c)
}
