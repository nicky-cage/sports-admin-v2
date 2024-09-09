package base_controller

import (
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// ActionDelete 删除记录
type ActionDelete struct {
	Model        common.IModel                         // model - 必须
	Row          func() interface{}                    // 用于保存记叟信息 - 可选 - 用于指定要处理的数据类型
	DeleteBefore func(*gin.Context, interface{}) error // 删除之前处理, 可以中断
	DeleteAfter  func(*gin.Context, interface{})       // 删除之后处理, 不可中断
}

// Delete 删除记录 - ajax
func (ths *ActionDelete) Delete(c *gin.Context) {
	idStr := c.DefaultQuery("id", "0")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	platform := request.GetPlatform(c)
	// 先查找相应的记录
	cond := builder.And(builder.Eq{"id": id})
	var row interface{}
	if ths.Row != nil {
		row = ths.Row()
	} else {
		row = &IdRecord{}
	}
	exists, err := ths.Model.Find(platform, row, cond)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	if !exists {
		response.Err(c, "删除失败: 此记录不存在")
		return
	}
	if ths.DeleteBefore != nil {
		if err := ths.DeleteBefore(c, row); err != nil {
			response.Err(c, err.Error())
			return
		}
	}
	err = ths.Model.Delete(platform, idStr)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	if ths.DeleteAfter != nil {
		ths.DeleteAfter(c, row)
	}
	response.Ok(c)
}
