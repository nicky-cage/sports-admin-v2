package base_controller

import (
	"encoding/json"
	common "sports-common"
	"sports-common/es"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

// ActionEsList 数据列表
type ActionEsList struct {
	ModelEs      common.IModelEs
	QueryCond    map[string]interface{}                //查询条件
	GetQueryCond func(*gin.Context) *elastic.BoolQuery //获得查询条件, 此条件会与 elastic.BoolQuery条件进行累加
	ProcessRow   func(interface{})                     //默认处理函数
	ViewFile     string                                //视图文件 - 必须
	//Rows               func() interface{}                    //多条记录 - 不需要
	Row                func() interface{}                //一条记录 - 必须 返回指针类型
	ExtendData         func(*gin.Context) pongo2.Context //扩展数据
	OrderBy            func(*gin.Context) string         //获取排序
	SearchPageCriteria common.SearchPageCriteria
}

// List 记录列表 - get
func (ths *ActionEsList) List(c *gin.Context) {
	row := ths.Row()
	//rows := ths.Rows()
	rows := make([]interface{}, 0)

	limit, offset := request.GetOffsets(c)
	searchPageCriteria := common.SearchPageCriteria{
		Limit:         limit,
		Offset:        offset,
		OrderType:     false,
		SortFieldName: "created_at",
	}
	//platform := request.GetPlatform(c)
	var (
		esResp *elastic.SearchResult
		err    error
	)
	cond := request.GetEsQueryCond(c, ths.QueryCond)
	if ths.GetQueryCond != nil { //将条件进行合并
		condTemp := ths.GetQueryCond(c)
		if condTemp != nil { //确保附加条件不为空
			if cond != nil { // 如果有条件, 则合并条件
				cond = cond.Filter(condTemp)
			} else { //如果没有条件, 则生成条件
				cond = condTemp
			}
		}
	}
	platform := request.GetPlatform(c)
	client, err := es.GetClientByPlatform(platform)
	if err != nil {
		log.Err("获取列表信息出错: %v\n", err)
		response.Err(c, "获取列表错误")
		return
	}
	defer client.Stop()

	// 关于 order by 的判断
	esResp, err = ths.ModelEs.Search(platform, client, cond, &searchPageCriteria)

	if err != nil {
		log.Err("获取列表信息出错: %v\n", err)
		response.Err(c, "获取列表错误")
		return
	}

	realViewFile := ths.ViewFile
	if request.IsAjax(c) && !strings.Contains(ths.ViewFile, "/_") { //如果是ajax请求
		realViewFile = strings.ReplaceAll(ths.ViewFile, "/", "/_")
	}
	total := 0
	if esResp.Hits.TotalHits.Value > 0 {
		total = int(esResp.Hits.TotalHits.Value) //获取总分页数目
		for _, v := range esResp.Hits.Hits {
			//fmt.Printf("v.Source %s\n", string(v.Source))
			//err := json.Unmarshal(v.Source, &rows) //另外一种取数据的方法
			//fmt.Printf("v.Source: %#v\n", string(v.Source))
			err := json.Unmarshal(v.Source, row) //另外一种取数据的方法
			if err != nil {
				log.Err("json.Unmarshal %v", err)
				continue
			}

			//这里必须断言取值，否则，数据都是最后一条
			if data, ok := row.(*models.EsLoginLogs); ok {
				rows = append(rows, *data)
			} else if data, ok := row.(*models.EsLoginLogs); ok {
				rows = append(rows, *data)
			}
			//fmt.Printf("row: %#v\n", row)
			//fmt.Printf("rows: %#v\n", rows)

		}

	}
	if ths.ProcessRow != nil {
		ths.ProcessRow(rows)
	}
	viewData := pongo2.Context{
		"rows":  rows,
		"total": total,
	}
	if ths.ExtendData != nil { //如果有附加数据, 则进行追加
		data := ths.ExtendData(c)
		for k, v := range data {
			viewData[k] = v
		}
	}

	response.Render(c, realViewFile, viewData)
}
