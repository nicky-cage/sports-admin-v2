package controllers

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DepositLimits 站点设置-存款限制
var DepositLimits = struct {
	Index      func(*gin.Context) // 数据列表
	Remind     func(*gin.Context) // 提示
	RemindSave func(*gin.Context) // 保存提示相关
	Allow      func(*gin.Context) // 是否允许提交多笔未支付订单
}{
	Index: func(c *gin.Context) { //默认首页
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		level := []models.DepositLimit{}
		_ = dbSession.Table("deposit_limits").Find(&level)
		sql := "select * from deposit_limits_set"
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		row := res[0]
		response.Render(c, "deposit_limits/_index.html", ViewData{
			"rows": level,
			"res":  row,
		})
	},
	Remind: func(c *gin.Context) { //时间提醒设置
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from deposit_limits_set"
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		response.Render(c, "deposit_limits/remind.html", ViewData{"r": res[0]})
	},
	RemindSave: func(c *gin.Context) {
		response.Ok(c)
	},
	Allow: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		if valueStr, exists := postedData["value"]; !exists {
			response.Err(c, "提交的值有误")
		} else if value, err := strconv.Atoi(valueStr.(string)); err != nil {
			response.Err(c, "")
		} else if value != 1 && value != 2 {
			response.Err(c, "")
		} else {
			platform := request.GetPlatform(c)
			if value == 2 {
				models.DepositLimitsSets.TurnOn(platform)
			} else {
				models.DepositLimitsSets.TurnOn(platform, false)
			}
			response.Ok(c)
		}
	},
}
