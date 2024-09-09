package base_controller

import (
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// ActionState 自动修改状态
type ActionState struct {
	Model       common.IModel            //模型
	Field       string                   //字段名称
	StateBefore func(*gin.Context) error //修改状态前处理
	StateAfter  func(*gin.Context)       //修改状态后处理
}

// State 修改状态 - ajax
// 需要提交2个参数
// id: 记录编号; to_state: 要更改到的状态值
func (ths *ActionState) State(c *gin.Context) {
	idStr := c.DefaultQuery("id", "0") //检测ID
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	if ths.Field == "" {
		ths.Field = "state" //如果没有设置状态字段, 则默认置为 state
	}
	toStateStr, exists := c.GetQuery("to_" + ths.Field)
	if !exists {
		response.Err(c, "必须提供要处理的状态")
		return
	}
	toState, err := strconv.Atoi(toStateStr)
	if err != nil {
		response.Err(c, "状态类型有误")
		return
	}
	platform := request.GetPlatform(c)
	cond := builder.And(builder.Eq{"id": id}).And(builder.Neq{ths.Field: toState}) // 先查找相应的记录
	row := &struct {
		Id uint64 `json:"id"`
	}{}

	exists, err = ths.Model.Find(platform, row, cond)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	if !exists {
		response.Err(c, "修改状态失败: 此记录不存在")
		return
	}
	if ths.StateBefore != nil {
		if err := ths.StateBefore(c); err != nil {
			response.Err(c, err.Error())
			return
		}
	}
	data := map[string]interface{}{
		"id":      id,
		ths.Field: toState,
	}
	err = ths.Model.Update(platform, data)
	if err != nil {
		response.Err(c, err.Error())
		return
	}
	if ths.StateAfter != nil {
		ths.StateAfter(c)
	}
	response.Ok(c)
}
