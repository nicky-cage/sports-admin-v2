package base_controller

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// ActionUpdate 修改记录
type ActionUpdate struct {
	Model      common.IModel                     // 模型 - 必须
	Row        func() interface{}                // 单条记录 - 必须
	ExtendData func(*gin.Context) pongo2.Context // 扩展数据
	ViewFile   string                            // 视图文件 - 必须
}

// Update 修改页面
func (ths *ActionUpdate) Update(c *gin.Context) {
	idStr, exists := c.GetQuery("id")
	if !exists || idStr == "" {
		log.Err("无法获取id信息!\n")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id == 0 {
		log.Err("ID转换失败: %v\n", err)
		return
	}

	row := ths.Row()
	cond := builder.NewCond().And(builder.Eq{"id": id})
	platform := request.GetPlatform(c)
	exists, err = ths.Model.Find(platform, row, cond)
	if !exists || err != nil {
		log.Err("获取详情失败: %v\n", err)
		return
	}

	viewData := pongo2.Context{"r": row}
	if ths.ExtendData != nil {
		data := ths.ExtendData(c)
		for k, v := range data {
			viewData[k] = v
		}
	}

	SetLoginAdmin(c)
	response.Render(c, ths.ViewFile, viewData)
}
