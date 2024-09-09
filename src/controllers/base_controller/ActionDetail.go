package base_controller

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// ActionDetail 记录详情
type ActionDetail struct {
	Model common.IModel
	Row   func() interface{} //单条记录
}

// Detail 详情 - get
func (ths *ActionDetail) Detail(c *gin.Context) {
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
	response.Result(c, row)
}
