package ag

import (
	"bufio"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sports-common/config"
	"sports-common/consts/game_detail"
	"sports-common/es"
	"sports-common/log"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req"
	"github.com/olivere/elastic/v7"
)

type AgDetail struct{}

func NewGameDetail() *AgDetail {
	return new(AgDetail)
}

func (ths *AgDetail) Data(billNo string, platform string) (*models.GameDetail, []*models.GameDetail) {
	res := ths.GetBetInfo(billNo, platform)
	detailList := make([]*models.GameDetail, 10)
	detail := new(models.GameDetail)
	err := json.Unmarshal([]byte(res.ExtendStr), &detailList)

	if err != nil {
		// 修正历史遗留问题
		var ext models.GameDetail
		err = json.Unmarshal([]byte(res.ExtendStr), &ext)
		if err == nil {
			detail = &ext
			if res.ExtendDetail != "" {
				detail.BetDetail = ths.GetDetail(res.ExtendDetail, platform)
			} else {
				detail.BetDetail = ths.GetBetDetail(detail.BetProject, detail.BetDetail)
			}
		}
	} else {
		if res.ExtendDetail != "" {
			detailList[0].BetDetail = ths.GetDetail(res.ExtendDetail, platform)
		} else {
			detailList[0].BetDetail = ths.GetBetDetail(detailList[0].BetProject, detailList[0].BetDetail)
		}
		// detailList[0].SportName = detailList[0].GameBillNo
		// detailList[0].Score = ths.GetBetScore(detailList[0].GameBillNo, detailList[0].BetProject, detailList[0].StartTime)
		// detailList[0].BetProject = ths.GetBetProject(detailList[0].BetProject)
		detail = detailList[0]
		detail.GameType = detailList[0].GameType
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
	detail.VenueCode = "AG"
	detail.VenueName = "AG国际厅"

	detail.BillNo = billNo

	return detail, detailList
}

func (ths *AgDetail) GetBetInfo(billNo string, platform string) models.WagerRecord {

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

func (ths *AgDetail) GetDetail(data string, platform string) string {
	var result string
	var err error
	var info response.RespInfo
	req.SetTimeout(50 * time.Second)
	req.Debug = true
	header := req.Header{
		"Accept": "application/json",
	}
	params := req.Param{
		"code": "AG",
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

func (ths *AgDetail) GetBetProject(name string) string {
	var temp string
	temp, ok := game_detail.AgGameName[name]
	if !ok {
		temp = name
	}
	return temp
}

func (ths *AgDetail) GetBetDetail(betType string, name string) string {
	return game_detail.AgGameType[betType][name]
}

func (ths *AgDetail) GetBetScore(gameBillNo string, projectName string, betTime string) string {
	//获取当局结果。
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, betTime, loc)
	tmpStr := tmp.Format("20060102") //某天的文件夹。
	tmpDay, _ := strconv.Atoi(tmp.Format("200601021504"))

	agXmlSavePath := config.Get("ag.xml") + tmpStr

	files, err := ioutil.ReadDir(agXmlSavePath)
	if err != nil {
		log.Logger.Errorf("ag err: %v", err)
		return ""
	}

	fileName := config.Get("ag.xml") + tmpStr + "/"

	fileNums := len(files)
	if fileNums == 1 {
		fileName = fileName + files[0].Name()
	}
	for k, v := range files {
		if strings.Index(fileName, ".xml") > 0 {
			break
		}
		for key, val := range files {
			name, _ := strconv.Atoi(strings.Replace(v.Name(), ".xml", "", 1))
			names, _ := strconv.Atoi(strings.Replace(val.Name(), ".xml", "", 1))
			//最开始一个。 和最后一个都要考虑

			if name+4 >= tmpDay && tmpDay < names {
				fileName = fileName + val.Name()
				break
			}
			if k+1 == fileNums && key+1 == fileNums && tmpDay > names {
				fileName = fileName + val.Name()
				break
			}
		}
	}

	//打开所有文件，获取具体的一个文件夹，
	//地址获取，
	file, err := os.Open(fileName)
	if err != nil {
		log.Err(err.Error())
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	type TempXml struct {
		XMLName xml.Name `xml:"docs"`
		Text    string   `xml:",chardata"`
		Row     []struct {
			Text     string `xml:",chardata"`
			CardList string `xml:"cardlist,attr"`
		} `xml:"row"`
	}
	var temp TempXml
	matchStr := `gmcode="%s"`
	matchStr = fmt.Sprintf(matchStr, gameBillNo)
	for scanner.Scan() { //每行匹配。
		lineTxt := scanner.Text()
		if strings.Index(lineTxt, matchStr) > 0 {
			tempStr := `<docs>%s</docs>`
			lineTxts := fmt.Sprintf(tempStr, lineTxt)
			err := xml.Unmarshal([]byte(lineTxts), &temp)
			if err != nil {
				log.Err(err.Error())
				return ""
			}

		}
	}
	var res string
	if len(temp.Row) == 0 {
		return res
	}
	switch projectName {
	case "BAC":
		res = ths.GetAGINResult(temp.Row[0].CardList)
	case "BJ":
		res = ths.GetBJResult(temp.Row[0].CardList)
	case "DT":
		res = ths.GetDtResult(temp.Row[0].CardList)
	case "SHB":
		res = temp.Row[0].CardList
	case "ULPK":
		res = ths.GetULPKResult(temp.Row[0].CardList)
	case "NN":
		res = ths.GetNNResult(temp.Row[0].CardList)
	case "ZJH":
		res = ths.GetZJHResult(temp.Row[0].CardList)
	case "BF":
		res = ths.GetBFResult(temp.Row[0].CardList)
	case "SG":
		res = ths.GetSGResult(temp.Row[0].CardList)
	case "ROU":
		res = temp.Row[0].CardList
	}

	return res
	//翻译
}

func (ths *AgDetail) GetAGINResult(score string) string {
	arr := strings.Split(score, ";")
	arr1 := strings.Split(arr[0], ",")
	arr2 := strings.Split(arr[1], ",")
	var sum int
	var sums int
	for _, v := range arr1 {
		temp := strings.Split(v, ".")
		nums, _ := strconv.Atoi(temp[1])
		if nums < 10 {
			sum += nums
		}
	}
	for _, v := range arr2 {
		temp := strings.Split(v, ".")
		nums, _ := strconv.Atoi(temp[1])
		if nums < 10 {
			sums += nums
		}
	}

	temp := "庄:" + strconv.Itoa(sum) + ", 闲:" + strconv.Itoa(sums)
	return temp
}

func (ths *AgDetail) GetDtResult(score string) string {
	arr := strings.Split(score, ";")
	arr1 := strings.Split(arr[0], ",")
	arr2 := strings.Split(arr[1], ",")
	var sum int
	var sums int
	for _, v := range arr1 {
		temp := strings.Split(v, ".")
		sum, _ = strconv.Atoi(temp[1])

	}
	for _, v := range arr2 {
		temp := strings.Split(v, ".")
		sums, _ = strconv.Atoi(temp[1])

	}

	temp := "庄:" + strconv.Itoa(sum) + ", 闲:" + strconv.Itoa(sums)
	return temp
}

func (ths *AgDetail) GetULPKResult(score string) string {
	arr := strings.Split(score, ";")
	arr1 := strings.Split(arr[0], ",")
	arr2 := strings.Split(arr[1], ",")
	sum := "公共牌:"
	sums := "玩家牌:"
	for _, v := range arr1 {
		temp := strings.Split(v, ".")
		sum = sum + temp[1] + ","

	}
	for _, v := range arr2 {
		temp := strings.Split(v, ".")
		sums = sums + temp[1] + ","

	}
	sums = strings.TrimRight(sums, ",")
	temp := sum + sums
	return temp
}

func (ths *AgDetail) GetBJResult(score string) string {
	arr := strings.Split(score, ";")
	sum := "庄家:"
	sums := "闲家:"
	for k, v := range arr {
		tempArr := strings.Split(v[2:], ".")
		for key, val := range tempArr {
			if key > 0 && k%2 > 0 && k == 0 {
				sum = sum + val + ","
				continue
			}
			if key > 0 && k%2 > 0 && k > 0 {
				sums = sums + val + ","
				continue
			}
		}

	}

	sums = strings.TrimRight(sums, ",")
	temp := sum + sums
	return temp
}

func (ths *AgDetail) GetNNResult(score string) string {
	sum := "庄家:"
	sums1 := "闲家一:"
	sums2 := "闲家二:"
	sums3 := "闲家三:"
	arr := strings.Split(score, ";")
	//S.6,C.2,C.9,H.8,S.7;H.7,C.7,S.11,C.6,S.5;D.10,S.10,D.5,H.11,C.3;C.12,H.4,S.4,D.9,C.4
	for k, v := range arr {
		if k > 0 && k == 1 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sum = sum + temp[1] + ","
			}
		}
		if k > 0 && k == 2 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sums1 = sums1 + temp[1] + ","
			}
		}
		if k > 0 && k == 3 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sums2 = sums2 + temp[1] + ","
			}
		}
		if k > 0 && k == 4 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sums3 = sums3 + temp[1] + ","
			}
		}
	}

	sum = strings.TrimRight(sum, ",")
	sum1 := strings.TrimRight(sums1, ",")
	sum2 := strings.TrimRight(sums2, ",")
	sum3 := strings.TrimRight(sums3, ",")
	temp := sum + sum1 + sum2 + sum3
	return temp
}

func (ths *AgDetail) GetZJHResult(score string) string {
	arr := strings.Split(score, ";")
	arr1 := strings.Split(arr[0], ",")
	arr2 := strings.Split(arr[1], ",")
	sum := "庄家:"
	sums := "闲家:"
	for _, v := range arr1 {
		temp := strings.Split(v, ".")
		sum = sum + temp[1] + ","

	}
	for _, v := range arr2 {
		temp := strings.Split(v, ".")
		sums = sums + temp[1] + ","

	}

	sums = strings.TrimRight(sums, ",")
	temp := sum + sums
	return temp
}

func (ths *AgDetail) GetBFResult(score string) string {
	//H.1,C.3,S.12,C.5,C.2;
	arr := strings.Split(score, ";")
	arr1 := strings.Split(arr[0], ",")
	arr2 := strings.Split(arr[1], ",")
	sum := "庄家:"
	sums := "闲家:"
	for _, v := range arr1 {
		temp := strings.Split(v, ".")
		sum = sum + temp[1] + ","

	}
	for _, v := range arr2 {
		temp := strings.Split(v, ".")
		sums = sums + temp[1] + ","

	}

	sums = strings.TrimRight(sums, ",")
	temp := sum + sums
	return temp
}

func (ths *AgDetail) GetSGResult(score string) string {
	////D.2;  H.1,S.13,H.4;D.5,S.3,S.4;
	sum := "庄家:"
	sums1 := "闲家一:"
	sums2 := "闲家二:"
	sums3 := "闲家三:"
	arr := strings.Split(score, ";")
	//S.7; ;C.6,S.9,C.10;S.13,C.9,D.8;H.2,D.5,H.8
	for k, v := range arr { //S.5,S.8,C.4
		if k > 0 && k == 1 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sum = sum + temp[1] + ","
			}
		}
		if k > 0 && k == 2 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sums1 = sums1 + temp[1] + ","
			}
		}
		if k > 0 && k == 3 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sums2 = sums2 + temp[1] + ","
			}
		}
		if k > 0 && k == 4 {
			temp := strings.Split(v, ",")
			for _, val := range temp {
				temp := strings.Split(val, ".")
				sums3 = sums3 + temp[1] + ","
			}
		}
	}

	sum = strings.TrimRight(sum, ",")
	sum1 := strings.TrimRight(sums1, ",")
	sum2 := strings.TrimRight(sums2, ",")
	sum3 := strings.TrimRight(sums3, ",")
	temp := sum + sum1 + sum2 + sum3
	return temp
}
