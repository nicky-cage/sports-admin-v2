package user_detail_accounts

import (
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

func (ths *UserDetailAccounts) TransferOut(g *gin.Context) {
	platform := request.GetPlatform(g)
	postedData := request.GetPostedData(g)
	idStr, exists := postedData["id"].(string)
	if !exists || exists && idStr == "0" {
		response.Err(g, "id为空")
		return
	}
	money, merr := strconv.ParseFloat(postedData["money"].(string), 64)
	if merr != nil {
		response.Err(g, "id为空")
		return
	}
	//防止多人同时更改
	//redis := common.Redis(platform)
	//rKey := "_" + idStr + "_" + postedData["code"].(string)
	//num, err := redis.Incr(rKey).Result()
	//if err != nil {
	//	log.Err(err.Error())
	//	response.Err(g, "系统繁忙")
	//	return
	//}
	//if num > 1 {
	//	response.Err(g, "请稍等片刻再试")
	//	return
	//}
	//defer redis.Del(rKey)
	//场馆扣钱
	req.SetTimeout(50 * time.Second)
	req.Debug = true
	header := req.Header{
		"Accept": "application/json",
	}
	param := req.Param{
		"out_code": postedData["code"],
		"in_code":  "CENTERWALLET",
		"all_in":   2,
		"money":    money,
	}
	baseTransferUrl := config.Get("internal.internal_game_service") + config.Get("internal_api.transfer_url")
	TransferUrl := baseTransferUrl + "?user_id=" + idStr + "&platform=" + platform
	r, err := req.Post(TransferUrl, header, param)
	if err != nil {
		log.Logger.Error(err.Error())
		response.Err(g, "系统异常")
		return
	}
	res := TransferState{}
	_ = r.ToJSON(&res)
	if res.Errcode != 0 {
		response.Err(g, "转出错误")
		return
	}
	response.Ok(g)
}
