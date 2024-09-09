package admin_tools

import (
	"fmt"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"sports-common/utils"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DownVipQuery 升级vip信息
type DownVipQuery struct {
	UserID    int     `form:"user_id" json:"user_id"`       // 用户编号
	FromLevel int     `form:"from_level" json:"from_level"` // 从vip
	ToLevel   int     `form:"to_level" json:"to_level"`     // 至vip
	Reduct    int     `form:"reduct" json:"reduct"`         // 是否扣款 1:是 0:否
	Valid     float64 `form:"valid" json:"valid"`           // 有效投注
}

// DownVip 降级vip - 针对某一统计区间进行降级
func DownVip(c *gin.Context) {
	var query DownVipQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Err(c, "提交数据格式有误")
	}

	if query.UserID <= 0 || query.FromLevel <= 0 || query.ToLevel >= 10 || // 各个数值必须正确
		query.FromLevel <= query.ToLevel || // 从等级必须小于到等级
		query.FromLevel-1 != query.ToLevel { // 只能小于1
		response.Err(c, "提交数据范围有误")
		return
	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()

	userID := query.UserID
	if err := myClient.Begin(); err != nil {
		log.Logger.Error("启动用户vip降级操作出错:", err)
		_ = myClient.Rollback()
		response.Err(c, "启动事务发生错误")
		return
	}

	userInfo := &models.User{}
	if exists, err := myClient.Table("users").Where("id = ?", userID).Get(userInfo); !exists || err != nil {
		log.Logger.Error("获取用户信息出错:", err, ", id:", userID, ", exists:", exists)
		_ = myClient.Rollback()
		response.Err(c, "获取用户信息出错")
		return
	}
	if int(userInfo.Vip) != query.FromLevel {
		response.Err(c, "用户当前VIP等级有误")
		return
	}

	userLevel := &models.UserLevel{}
	if exists, err := myClient.Table("user_levels").Where("id = ?", query.FromLevel).Get(userLevel); !exists || err != nil {
		log.Logger.Error("获取用户等级信息出错:", err)
		_ = myClient.Rollback()
		response.Err(c, "获取用户等级信息出错")
		return
	}

	accountInfo := &models.Account{}
	if exists, err := myClient.Table("accounts").Where("user_id = ?", userID).Get(accountInfo); !exists || err != nil {
		log.Logger.Error("获取用户账户信息出错:", err, ", user_id:", userID, ", exists:", exists)
		_ = myClient.Rollback()
		response.Err(c, "获取用户账户信息出错")
		return
	}

	currentTime := tools.NowMicro() // 当前时间

	// 记录用户vip升级
	updateData := map[string]interface{}{
		"user_id":     userInfo.Id,
		"username":    userInfo.Username,
		"valid_bet":   query.Valid,
		"deposits":    0,
		"before_vip":  query.FromLevel,
		"after_vip":   query.ToLevel,
		"adjust_type": 2, // 降级
		"admin":       "system",
		"created":     currentTime,
		"updated":     currentTime,
	}
	upArr := utils.MapToInsertSQL("user_vip_logs", updateData)
	if _, err := myClient.Exec(upArr...); err != nil {
		log.Logger.Error("记录用户VIP降级日志失败:", err.Error())
		_ = myClient.Rollback()
		response.Err(c, "记录用户VIP降级日志失败")
		return
	}

	downSQL := fmt.Sprintf("UPDATE users SET vip = %d, last_vip_update_at = '%s' WHERE id = %d LIMIT 1",
		query.ToLevel, tools.CurrentTime().Format("2006-01-02 15:04:05"), userInfo.Id)
	log.Logger.Info(downSQL)
	if _, err := myClient.Exec(downSQL); err != nil {
		log.Logger.Error("更新用户表修改vip信息出错:", err)
		_ = myClient.Rollback()
		response.Err(c, "修改用户VIP等级失败")
		return
	}

	if query.Reduct == 1 { // 需要扣除月俸禄
		rdClient := common.Redis(platform)
		defer rdClient.Close()
		extraMap := map[string]interface{}{
			"proxy_ip":      "127.0.0.1",
			"ip":            "127.0.0.1",
			"description":   "扣除" + strconv.Itoa(int(userInfo.Vip)) + "本月误发每月俸禄",
			"administrator": "system",
			"admin_user_id": "system",
			"serial_number": tools.GetBillNo("v", 5),
		}
		transType := consts.TransTypeAdjustmentLess // 降级vip - 调整 - 减少
		transAction := &models.Transaction{}
		if _, err := transAction.AddTransaction(platform, myClient, rdClient, userInfo, accountInfo, transType, userLevel.MonthBonus, extraMap); err != nil {
			log.Logger.Error("扣除每月俸禄出现错误:", err)
			_ = myClient.Rollback()
			response.Err(c, "扣除多发每月俸禄出错: "+err.Error())
			return
		}
	}

	if err := myClient.Commit(); err != nil {
		log.Logger.Error("写入日志信息失败:", err)
		response.Err(c, "")
		return
	}

	response.Ok(c)
}
