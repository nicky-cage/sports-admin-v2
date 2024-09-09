package bti

import (
	"context"
	"encoding/json"
	"fmt"
	"sports-common/consts/game_detail"
	"sports-common/es"
	"sports-common/log"
	"sports-common/tools"
	models "sports-models"

	"github.com/olivere/elastic/v7"
)

type BtiDetail struct{}

func NewGameDetail() *BtiDetail {
	return new(BtiDetail)
}

func (ths *BtiDetail) Data(billNo string, platform string) (*models.GameDetail, []*models.GameDetail) {
	res := ths.GetBetInfo(billNo, platform)
	//var detail []models.GameDetail
	detail := make([]*models.GameDetail, 10)
	err := json.Unmarshal([]byte(res.ExtendStr), &detail)
	if err != nil {
		log.Err(err.Error())
	}
	commonDetail := new(models.GameDetail)
	commonDetail.PlayName = res.Playname
	commonDetail.ValidMoney = res.ValidMoney
	commonDetail.NetMoney = res.NetMoney
	commonDetail.RebateMoney = tools.ToFixed(res.RebateMoney, 2)
	commonDetail.RebateRate = res.RebateRatio
	commonDetail.CreateTime = res.CreatedAt
	commonDetail.UpdateTime = res.UpdatedAt
	commonDetail.BetMoney = res.BetMoney
	commonDetail.Status = res.Status
	commonDetail.VenueCode = "BTI"
	commonDetail.VenueName = "BTi体育"
	commonDetail.GameType = "体育"
	commonDetail.BillNo = billNo

	if detail[0].BetType != "组合投注" { //当时普通投注时
		commonDetail.Odds = detail[0].Odds
		commonDetail.MyWin = detail[0].MyWin
		detail[0].BetProject = ths.GetBetProject(detail[0].BetProject)
	} else { //当时串单时
		commonDetail.Odds = 1.0
		for _, v := range detail {
			commonDetail.Odds *= v.WagerOdds
			commonDetail.MyWin += v.MyWin
			v.BetProject = ths.GetBetProject(v.BetProject)
		}
	}
	return commonDetail, detail
	//status 需要转化。用fitter函数转化

	//拼凑

}

func (ths *BtiDetail) GetBetInfo(billNo string, platform string) models.WagerRecord {
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

func (ths *BtiDetail) GetBetProject(name string) string {
	var temp string
	temp, ok := game_detail.EventTypeName[name]
	if !ok {
		temp = name
	}
	return temp
}

func (ths *BtiDetail) GetBetDetail(name string) string {
	return name
}
