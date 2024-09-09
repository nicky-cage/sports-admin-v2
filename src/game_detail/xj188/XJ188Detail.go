package xj188

import (
	"context"
	"encoding/json"
	"fmt"
	"sports-common/es"
	"sports-common/log"
	"sports-common/tools"
	models "sports-models"

	"github.com/olivere/elastic/v7"
)

type XJ188Detail struct{}

func NewGameDetail() *XJ188Detail {
	return new(XJ188Detail)
}

func (ths *XJ188Detail) Data(billNo string, platform string) (*models.GameDetail, []*models.GameDetail) {
	res := ths.GetBetInfo(billNo, platform)
	detailList := make([]*models.GameDetail, 10)
	detail := new(models.GameDetail)
	err := json.Unmarshal([]byte(res.ExtendStr), &detailList)
	if err != nil {
		log.Err(err.Error())
	}

	detail.Odds = 1.0
	for _, v := range detailList {
		detail.MyWin += v.MyWin
		if v.WagerOdds != 0 {
			detail.Odds *= v.WagerOdds
		}
	}

	detail.PlayName = res.Playname
	detail.ValidMoney = res.ValidMoney
	detail.NetMoney = res.NetMoney
	detail.RebateMoney = tools.ToFixed(res.RebateMoney, 2)
	detail.RebateRate = res.RebateRatio
	detail.CreateTime = res.CreatedAt
	detail.UpdateTime = res.UpdatedAt
	detail.BetMoney = res.BetMoney
	detail.Status = res.Status
	detail.VenueCode = "XJ188"
	// if platform == "tianji" || platform == "shipu" {
	// 	detail.VenueName = "天际体育"
	// 	res.GameCode = "天际体育"
	// } else if platform == "xingkong" {
	// 	detail.VenueName = "星空体育"
	// } else {
	// 	detail.VenueName = "小金体育"
	// }
	detail.VenueName = "小金体育"
	detail.GameType = "体育"
	detail.BillNo = billNo

	return detail, detailList
}

func (ths *XJ188Detail) GetBetInfo(billNo string, platform string) models.WagerRecord {
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
