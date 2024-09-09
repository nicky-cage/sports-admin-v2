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

// UpgradeVip 升级vip信息
type UpgradeVip struct {
	UserID    int `form:"user_id" json:"user_id"`       // 用户编号
	FromLevel int `form:"from_level" json:"from_level"` // 从vip
	ToLevel   int `form:"to_level" json:"to_level"`     // 至vip
}

// UpVip 升级用户vip
func UpVip(c *gin.Context) {
	var query UpgradeVip
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Err(c, "提交数据格式有误")
		return
	}

	if query.UserID <= 0 || query.FromLevel <= 0 || query.ToLevel >= 10 || // 各个数值必须正确
		query.FromLevel >= query.ToLevel || // 从等级必须小于到等级
		query.FromLevel+1 != query.ToLevel { // 只能小于1
		response.Err(c, "提交数据范围有误")
		return
	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()

	// 写入vip升级日志 并赠送俸禄
	userID := query.UserID
	if err := myClient.Begin(); err != nil {
		log.Logger.Error("启动用户vip升级操作出错:", err)
		_ = myClient.Rollback()
		response.Err(c, "启动事务发生错误")
		return
	}

	userLevel := &models.UserLevel{}
	if exists, err := myClient.Table("user_levels").Where("id = ?", query.ToLevel).Get(userLevel); !exists || err != nil {
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

	userInfo := &models.User{}
	if exists, err := myClient.Table("users").Where("id = ?", userID).Get(userInfo); !exists || err != nil {
		log.Logger.Error("获取用户信息出错:", err, ", id:", userID, ", exists:", exists)
		_ = myClient.Rollback()
		response.Err(c, "获取用户信息出错")
		return
	}

	// 是否应该赠送升级礼金 - 如果以前有送过, 则不再送了
	shouldGiveUpMoney := func() bool {
		sRow := models.TotalInt{}
		vSQL := fmt.Sprintf("SELECT COUNT(*) AS total FROM user_vip_logs WHERE user_id = %d AND after_vip = %d AND adjust_type = 1 LIMIT 1", userID, query.ToLevel)
		if exists, err := myClient.SQL(vSQL).Get(&sRow); err != nil || !exists {
			log.Logger.Info("获取用户VIP历史:", vSQL)
			log.Logger.Error("获取用户vip积分历史出错:", err, ", exists:", exists, ", userID:", userID)
			return false
		} else if sRow.Total > 0 {
			log.Logger.Info("没用用户VIP升级历史, from", query.FromLevel, ", to", query.ToLevel, ", userID: ", userID)
			return false
		}
		return true
	}()
	currentTime := tools.NowMicro() // 当前时间
	if shouldGiveUpMoney {
		billNo := tools.GetBillNo("h", 5) // 生成订单号码
		// 红利记录
		diData := map[string]interface{}{
			"bill_no":         billNo,
			"user_id":         userID,
			"username":        userInfo.Username,
			"top_id":          userInfo.TopId,
			"top_name":        userInfo.TopName,
			"type":            3,                      // 类型
			"is_automatic":    2,                      // 是否自动
			"money_type":      1,                      // 类型
			"flow_limit":      2,                      // 流水限制
			"flow_multiple":   1,                      // 只要一倍流水, 会员等级表的那个是首存的流水
			"money":           userLevel.UpgradeBonus, // 升级礼金
			"created":         currentTime,
			"updated":         currentTime,
			"reviewer":        "system",
			"reviewer_remark": fmt.Sprintf("发放升级礼金 VIP(%d-%d)", userInfo.Vip-1, userInfo.Vip),
			"state":           2, // 是否自动
			"vip":             userInfo.Vip + 1,
		}
		diArr := utils.MapToInsertSQL("user_dividends", diData)
		if _, err := myClient.Exec(diArr...); err != nil {
			_ = myClient.Rollback()
			log.Logger.Error("设置用户红利记录出错: ", err.Error())
			response.Err(c, "发放升级礼金出错")
			return
		}
	}

	type TotalRow struct {
		Total      float64 `json:"total"`
		TotalSport float64 `json:"total_sport"`
	}
	betRow := TotalRow{}
	totalSQL := fmt.Sprintf("SELECT user_id, "+
		"SUM(IF(game_code = '0', valid_money, 0 )) AS total, "+ // -- 总的积分
		"SUM(IF(game_type = '1' OR game_type = '2', valid_money, 0)) AS total_sport "+ // -- 体育和电竞积分
		"FROM user_daily_reports WHERE user_id = %d", userID)
	if exists, err := myClient.SQL(totalSQL).Get(&betRow); err != nil { // 获取统计信息出错
		log.Logger.Error("获取用户投注统计信息出错:", totalSQL, ", error:", err.Error())
		_ = myClient.Rollback()
		response.Err(c, "获取用户投注信息出错")
		return
	} else if !exists {
		log.Logger.Error("获取用户投注统计信息出错, 没有获取任何数据: ", totalSQL)
		_ = myClient.Rollback()
		response.Err(c, "获取用户投注信息出错: 缺少记录")
		return
	}

	// 记录用户vip升级
	updateData := map[string]interface{}{
		"user_id":     userInfo.Id,
		"username":    userInfo.Username,
		"valid_bet":   betRow.Total + betRow.TotalSport,
		"deposits":    0,
		"before_vip":  userInfo.Vip,
		"after_vip":   userInfo.Vip + 1,
		"adjust_type": 1,
		"admin":       "system",
		"created":     currentTime,
		"updated":     currentTime,
	}
	upArr := utils.MapToInsertSQL("user_vip_logs", updateData)
	if _, err := myClient.Exec(upArr...); err != nil {
		log.Logger.Error("记录用户VIP升级日志失败:", err.Error())
		_ = myClient.Rollback()
		response.Err(c, "记录用户vip升级日志失败")
		return
	}

	upSQL := fmt.Sprintf("UPDATE users SET vip = %d, vip_high = %d, last_vip_update_at = '%s' WHERE id = %d LIMIT 1",
		userInfo.Vip+1, userInfo.VipHigh+1, tools.CurrentTime().Format("2006-01-02 15:04:05"), userInfo.Id)
	log.Logger.Info(upSQL)
	if _, err := myClient.Exec(upSQL); err != nil {
		log.Logger.Error("更新用户表修改vip信息出错:", err)
		_ = myClient.Rollback()
		response.Err(c, "记录用户vip升级日志失败")
		return
	}

	// 如果需要升级礼金, 则需要记录账变
	if shouldGiveUpMoney { // 记录账变
		rdClient := common.Redis(platform)
		defer rdClient.Close()
		extraMap := map[string]interface{}{
			"proxy_ip":      "",
			"ip":            "",
			"description":   "发放VIP" + strconv.Itoa(int(userInfo.Vip)) + "晋级礼金",
			"administrator": "system",
			"admin_user_id": "system",
			"serial_number": tools.GetBillNo("v", 5),
		}
		transType := consts.TransTypeAdjustmentDividendPlus // 调整红利 - 升级vip
		transAction := &models.Transaction{}
		if _, err := transAction.AddTransaction(platform, myClient, rdClient, userInfo, accountInfo, transType, userLevel.UpgradeBonus, extraMap); err != nil {
			log.Logger.Error("增加用户存款事务时出现错误:", err)
			_ = myClient.Rollback()
			response.Err(c, "发放升级礼金时出错")
			return
		}
	}

	if err := myClient.Commit(); err != nil {
		log.Logger.Error("提交事务日志出错:", err)
		_ = myClient.Rollback()
		response.Err(c, "发放升级礼金时出错")
		return
	}

	response.Ok(c)
}
