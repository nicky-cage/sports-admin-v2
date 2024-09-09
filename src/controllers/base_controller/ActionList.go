package base_controller

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"xorm.io/builder"
)

// ActionList 数据列表
type ActionList struct {
	Model             common.IModel                          // model - 必须
	QueryCond         map[string]interface{}                 // 查询条件
	GetQueryCond      func(*gin.Context) builder.Cond        // 获得查询条件, 此条件会与 QueryCond条件进行累加
	ProcessRow        func(*gin.Context, interface{})        // 默认处理函数
	ViewFile          string                                 // 视图文件 - 必须
	Rows              func() interface{}                     // 多条记录 - 必须
	ExtendData        func(*gin.Context) pongo2.Context      // 扩展数据
	OrderBy           func(*gin.Context) string              // 获取排序
	GetQuerySession   func(*gin.Context) *xorm.Session       // 得到自定义的查询session
	RequireParameters bool                                   // 必须带参数访问 - 如果此值为true, 则访问必须带参数否则将返回空白
	AfterAction       func(*gin.Context, *response.ViewData) // 处理之后
}

// List 记录列表 - get
func (ths *ActionList) List(c *gin.Context) {
	realViewFile := ths.ViewFile
	if request.IsAjax(c) && !strings.Contains(ths.ViewFile, "/_") { //如果是ajax请求
		realViewFile = strings.ReplaceAll(ths.ViewFile, "/", "/_")
	}

	tools.TimeDebugEnable = false
	t1 := tools.TimeDebugBegin("执行主体之前")
	SetLoginAdmin(c)
	if ths.RequireParameters {
		requestData := c.Request.URL.Query()
		if len(requestData) == 0 {
			response.Render(c, realViewFile, response.ViewData{})
			return
		}
		if _, exists := requestData["page"]; exists {
			delete(requestData, "page")
			if len(requestData) == 0 {
				response.Render(c, realViewFile, response.ViewData{})
				return
			}
		}
	}

	t2 := tools.TimeDebugAt(t1, "")
	rows := ths.Rows()
	limit, offset := request.GetOffsets(c)
	platform := request.GetPlatform(c)
	var total uint64
	var err error
	cond := request.GetQueryCond(c, ths.QueryCond)

	if ths.GetQueryCond != nil { //将条件进行合并
		condTemp := ths.GetQueryCond(c)
		if condTemp != nil { //确保附加条件不为空
			if cond != nil { // 如果有条件, 则合并条件
				cond = cond.And(condTemp)
			} else { //如果没有条件, 则生成条件
				cond = condTemp
			}
		}
	}

	t3 := tools.TimeDebugAt(t2, "执行判断是否 Export 之前")
	// 如果是导出
	if exportExcel := c.DefaultQuery("export_excel", ""); exportExcel != "" {
		err := func() error {
			if ths.OrderBy != nil {
				return ths.Model.FindAllNoCount(platform, rows, cond, ths.OrderBy(c))
			} else {
				return ths.Model.FindAllNoCount(platform, rows, cond)
			}
		}()
		if err != nil {
			log.Logger.Error("执行导出数据时出错: ", err)
			response.Err(c, "导出数据时出错")
			return
		}
		response.Result(c, rows)
		return
	}

	t4 := tools.TimeDebugAt(t3, "执行 FindAll 之前")
	// 关于 order by 的判断
	if ths.OrderBy != nil {
		total, err = ths.Model.FindAll(platform, rows, cond, limit, offset, ths.OrderBy(c))
	} else {
		total, err = ths.Model.FindAll(platform, rows, cond, limit, offset)
	}

	t5 := tools.TimeDebugAt(t4, "执行 ProcessRow 之前")
	if ths.ProcessRow != nil {
		ths.ProcessRow(c, rows)
	}
	if err != nil {
		log.Err("获取列表信息出错: %v\n", err)
		response.Err(c, "获取列表错误")
		return
	}

	viewData := pongo2.Context{
		"rows":  rows,
		"total": total,
	}

	t6 := tools.TimeDebugAt(t5, "执行 ExtendData 之前")
	if ths.ExtendData != nil { //如果有附加数据, 则进行追加
		data := ths.ExtendData(c)
		for k, v := range data {
			viewData[k] = v
		}
	}

	if ths.AfterAction != nil {
		ths.AfterAction(c, &viewData)
	}

	t7 := tools.TimeDebugAt(t6, "执行 response.Render 之前")
	response.Render(c, realViewFile, viewData)
	_ = tools.TimeDebugAt(t7, "执行 response.Render 之后")
}
