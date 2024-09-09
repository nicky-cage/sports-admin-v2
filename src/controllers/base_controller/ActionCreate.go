package base_controller

import (
	common "sports-common"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// ActionCreate 添加记录
type ActionCreate struct {
	Model      common.IModel                     // 模型
	ExtendData func(*gin.Context) pongo2.Context // 扩展数据
	ViewFile   string                            // 视图文件 - 必须
}

// Create 添加页面
func (ths *ActionCreate) Create(c *gin.Context) {
	viewData := pongo2.Context{}
	if ths.ExtendData != nil {
		data := ths.ExtendData(c)
		for k, v := range data {
			viewData[k] = v
		}
	}
	if ths.Model != nil {
		viewData["r"] = ths.Model
		response.Render(c, ths.ViewFile, viewData)
		return
	}

	response.Render(c, ths.ViewFile, viewData)
}
