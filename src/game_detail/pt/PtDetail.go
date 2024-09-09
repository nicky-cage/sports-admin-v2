package pt

import (
	"context"
	"encoding/json"
	"fmt"
	"sports-common/config"
	"sports-common/es"
	"sports-common/log"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"time"

	"github.com/imroc/req"
	"github.com/olivere/elastic/v7"
)

type PtDetail struct{}

func NewGameDetail() *PtDetail {
	return new(PtDetail)
}

func (ths *PtDetail) Data(billNo string, platform string) (*models.GameDetail, []*models.GameDetail) {
	res := ths.GetBetInfo(billNo, platform)
	detailList := make([]*models.GameDetail, 1)
	//detail := new(models.GameDetail)
	var ext models.GameDetail
	err := json.Unmarshal([]byte(res.ExtendStr), &ext)
	if err != nil {
		log.Err(err.Error())
	}
	if res.ExtendDetail != "" {
		ext.BetDetail = ths.GetDetail(res.ExtendDetail, platform)
	}
	detail := &ext
	detail.PlayName = res.Playname
	detail.ValidMoney = res.ValidMoney
	detail.NetMoney = res.NetMoney
	detail.RebateMoney = tools.ToFixed(res.RebateMoney, 2)
	detail.RebateRate = res.RebateRatio
	detail.CreateTime = res.CreatedAt
	detail.UpdateTime = res.UpdatedAt
	detail.BetMoney = res.BetMoney
	detail.Status = res.Status
	detail.VenueCode = ext.VenueCode
	detail.VenueName = ext.VenueName
	detail.GameType = ext.GameType
	detail.BillNo = billNo

	detailList[0] = &ext

	return detail, detailList
}

func (ths *PtDetail) GetBetInfo(billNo string, platform string) models.WagerRecord {
	esIndexName := platform + "_wagers"
	esClient, err := es.GetClientByPlatform(platform)
	if err != nil {
		fmt.Println(err)
		//return
	}
	defer esClient.Stop()
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery("bill_no", billNo))
	res, err := esClient.Search(esIndexName).Query(boolQuery).Do(context.Background())
	if err != nil {
		//response.Err(c, "es获取列表数据无响应: "+err.Error())
		//return
		log.Err(err.Error())
	}
	temp := models.WagerRecord{}
	if res.Hits.TotalHits.Value > 0 {
		for _, v := range res.Hits.Hits {
			err := json.Unmarshal(v.Source, &temp)
			if err != nil {
				log.Err(err.Error())
			}
		}
	}
	return temp
}

func (ths *PtDetail) GetDetail(data string, platform string) string {
	var result string
	var err error
	var info response.RespInfo
	req.SetTimeout(50 * time.Second)
	req.Debug = true
	header := req.Header{
		"Accept": "application/json",
	}
	params := req.Param{
		"code": "PT",
		"data": data,
	}

	baseGameUrl := config.Get("internal.internal_game_service", "")
	GameUrl := baseGameUrl + "/game/v1/internal/get_detail?platform=" + platform
	resp, err := req.Post(GameUrl, header, req.BodyJSON(params))
	if err != nil {
		log.Logger.Error(err.Error())
		return ""
	}
	err = resp.ToJSON(&info)
	if err == nil {
		if info.Data != nil {
			result = info.Data.(string)
		}
	}
	return result
}
