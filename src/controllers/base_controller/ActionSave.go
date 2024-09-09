package base_controller

import (
	common "sports-common"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

// ActionSave 保存数据
type ActionSave struct {
	Model           common.IModel                                     // model - 必须
	Validator       func(string, map[string]interface{}) error        //校验器
	CreateValidator func(string, map[string]interface{}) error        //校验器
	CreateBefore    func(*gin.Context, *map[string]interface{}) error //添加数据之前处理, 可以中断
	CreateAfter     func(*gin.Context, *map[string]interface{})       //添加数据之后处理, 不可中断
	UpdateValidator func(string, map[string]interface{}) error        //校验器
	UpdateBefore    func(*gin.Context, *map[string]interface{}) error //修改数据之前处理, 可以中断
	UpdateAfter     func(*gin.Context, *map[string]interface{})       //修改数据之后处理, 不可中断
	SaveBefore      func(*gin.Context, *map[string]interface{}) error //保存数据之前处理, 可以中断
	SaveAfter       func(*gin.Context, *map[string]interface{})       //保存数据之后处理, 不可中断
}

// Save 保存数据 - ajax
func (ths *ActionSave) Save(c *gin.Context) {
	postedData := request.GetPostedData(c)
	isCreate := false
	validator := ths.Validator                                                //优化选择before/after验证器
	if idStr, exists := postedData["id"]; !exists || exists && idStr == "0" { //id不能为0
		isCreate = true
		if ths.CreateAfter != nil {
			validator = ths.CreateValidator
		}
	} else {
		if ths.UpdateAfter != nil {
			validator = ths.UpdateValidator
		}
	}
	platform := request.GetPlatform(c)
	if validator != nil {
		if err := ths.Validator(platform, postedData); err != nil {
			response.Err(c, err.Error())
			return
		}
	}
	if isCreate && ths.CreateBefore != nil { //创建前验证
		if err := ths.CreateBefore(c, &postedData); err != nil { //可以取消
			response.Err(c, err.Error())
			return
		}
	} else if !isCreate && ths.UpdateBefore != nil { //修改前验证
		if err := ths.UpdateBefore(c, &postedData); err != nil { //可以取消
			response.Err(c, err.Error())
			return
		}
	}
	if ths.SaveBefore != nil { //保存前验证
		if err := ths.SaveBefore(c, &postedData); err != nil {
			response.Err(c, err.Error())
			return
		}
	}
	var id uint64 = 0
	if isCreate {
		createdId, err := ths.Model.Create(platform, postedData)
		if err != nil {
			response.Err(c, err.Error())
			return
		}
		id = createdId
	} else {
		err := ths.Model.Update(platform, postedData)
		if err != nil {
			response.Err(c, err.Error())
			return
		}
	}

	if isCreate && ths.CreateAfter != nil {
		ths.CreateAfter(c, &postedData) //不可取消
	} else if !isCreate && ths.UpdateAfter != nil {
		ths.UpdateAfter(c, &postedData)
	}
	if ths.SaveAfter != nil {
		ths.SaveAfter(c, &postedData)
	}

	if isCreate {
		response.Result(c, struct {
			ID uint64 `json:"id"`
		}{
			ID: id,
		})
	} else {
		response.Ok(c)
	}
}
