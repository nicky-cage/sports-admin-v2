package controllers

import (
	"encoding/json"
	"fmt"
	common "sports-common"
	"sports-common/log"
	"sports-common/redis"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// UserDepositCoinMatch 用户代币存款匹配
type UserDepositCoinMatch struct {
	Id                int32   `json:"id"`                  // 编号
	WalletId          int32   `json:"wallet_id"`           // 收款钱包编号
	Money             float64 `json:"money"`               // 金额 - 元
	VirtualCoin       float64 `json:"virtual_coin"`        // 金额 - 代币
	UserWalletAddress string  `json:"user_wallet_address"` // 用户钱包地址
	WalletAddress     string  `json:"wallet_address"`      // 收款钱包地址
}

// UserDepositCoinMatches 代币存款自动匹配
var UserDepositCoinMatches = struct {
	*ActionList
	*ActionCreate
	SaveRecharge func(c *gin.Context) // 保存自动匹配订单
}{
	ActionList: &ActionList{
		Model:    models.UserDepositCoinMatches,
		ViewFile: "user_deposit_coins/_matches.html",
		QueryCond: map[string]interface{}{
			"pay_user_name":   "%",
			"pay_address":     "%",
			"receive_address": "%",
			"order_number":    "%",
			"admin_name":      "%",
		},
		Rows: func() interface{} {
			return &[]models.UserDepositCoinMatch{}
		},
		OrderBy: func(C *gin.Context) string {
			return "created DESC"
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.UserCards,
		ViewFile: "user_deposit_coins/match_edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			rows := []UserDepositCoinMatch{}
			SQL := "SELECT ud.id, ud.money, ud.virtual_coin, ud.user_wallet_address, dv.wallet_address, ud.wallet_id " +
				"FROM user_deposits AS ud, deposit_virtuals AS dv " +
				"WHERE ud.status = 1 AND ud.type = 4 AND ud.wallet_id = dv.id"
			myClient := common.Mysql(platform)
			defer myClient.Close()
			_ = myClient.SQL(SQL).Find(&rows)
			retJson := ""
			if retBytes, err := json.Marshal(rows); err == nil {
				retJson = string(retBytes)
			}
			wallets := []models.DepositVirtual{}
			_ = models.DepositVirtuals.FindAllNoCount(platform, &wallets)
			return pongo2.Context{
				"deposits":        rows,
				"deposits_json":   retJson,
				"virtual_wallets": wallets,
			}
		},
	},
	SaveRecharge: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		rechargeID := 0                                     // 充值ID
		walletID := 0                                       // 收款钱包ID
		payAddress := ""                                    // 付款地址
		remark := ""                                        // 备注
		if v, exists := postedData["recharge_id"]; exists { // 获取充值记录
			if rv, err := strconv.Atoi(v.(string)); err == nil && rv > 0 {
				rechargeID = rv
			}
		}
		if v, exists := postedData["wallet_id"]; exists { // 收款钱包ID
			if rv, err := strconv.Atoi(v.(string)); err == nil && rv > 0 {
				walletID = rv
			}
		}
		if v, exists := postedData["pay_address"]; exists { // 付款地址
			payAddress = strings.TrimSpace(v.(string))
		}
		if v, exists := postedData["remark"]; exists { // 备注
			remark = strings.TrimSpace(v.(string))
		}
		if rechargeID <= 0 || walletID <= 0 || payAddress == "" || remark == "" {
			response.Err(c, "提交信息有误, 缺少必要字段")
			return
		}

		platform := request.GetPlatform(c) // 平台识别号
		// -- 防止多人同时更改
		rKey, err := redis.Lock(platform, fmt.Sprintf("confirm-deposit-%d", rechargeID))
		if err != nil {
			log.Err(err.Error())
			fmt.Println("缓存服务器加锁失败: ", err)
			response.Err(c, "请不要同一时间内多次提交")
			return
		}
		defer redis.Unlock(platform, rKey)

		row := models.UserDeposit{}
		cond := builder.NewCond().And(builder.Eq{"id": rechargeID}).And(builder.Eq{"user_wallet_address": payAddress}).And(builder.Eq{"wallet_id": walletID})
		mySQLClient := common.Mysql(platform)
		defer mySQLClient.Close()

		wRow := models.DepositVirtual{}
		if exists, err := mySQLClient.Table("deposit_virtuals").Where(builder.NewCond().And(builder.Eq{"id": walletID})).Get(&wRow); err != nil || !exists {
			fmt.Println("err: ", err)
			response.Err(c, "查找收款钱包失败")
			return
		}

		mySQLClient.Begin() // 开启事务
		if exists, err := mySQLClient.Where(cond).ForUpdate().Get(&row); err != nil || !exists {
			mySQLClient.Rollback()
			response.Err(c, "查找不到此充值信息")
			return
		}
		if row.UserWalletAddress != payAddress || row.Status != 1 { // 必须是未处理状态
			mySQLClient.Rollback()
			response.Err(c, "充值信息状态有误")
			return
		}

		admin := GetLoginAdmin(c)
		// 用于保存到财务日志表
		financeData := map[string]interface{}{
			"bill_no":   postedData["order_no"],
			"type":      0,
			"operating": "存款结束",
			"result":    "失败",
			"operator":  admin.Name,
			"consuming": time.Now().Unix() - int64(row.Created/1000000),
			"remark":    "代币充值订单自动匹配",
			"created":   tools.NowMicro(),
		}
		// 用于保存到用户存款记录表
		depositData := map[string]interface{}{
			"arrive_money":     row.Money,
			"confirm_money":    row.Money,
			"remark":           remark,
			"status":           3,
			"finance_admin":    admin.Name,
			"updated":          tools.NowMicro(),
			"is_first_deposit": 1,
		}

		// -- 以下, 统一处理保存问题
		if err := saveConfirmDeposit(platform, &row, depositData, financeData, c, mySQLClient); err != nil {
			response.Err(c, err.Error())
			return
		}

		// 保存相关记录信息
		record := models.UserDepositCoinMatch{
			OrderNumber:    row.OrderNo,
			PayUserId:      int(row.UserId),
			PayUserName:    row.Username,
			PayAddress:     payAddress,
			Amount:         row.Money,
			VirtualCoin:    row.VirtualCoin, // 代币金额
			ReceiveName:    wRow.Name,
			ReceiveAddress: wRow.WalletAddress,
			AdminId:        int(admin.Id),
			AdminName:      admin.Name,
			Remark:         remark,
			Created:        tools.NowMicro(),
			Updated:        tools.NowMicro(),
		}
		if lastID, err := mySQLClient.Table("user_deposit_coin_matches").Insert(record); err != nil || lastID <= 0 {
			mySQLClient.Rollback()
			fmt.Println("自动配单保存自动匹配信息失败:", err)
			response.Err(c, "保存自动匹配信息失败")
			return
		}

		if err := mySQLClient.Commit(); err != nil {
			mySQLClient.Rollback()
			fmt.Println("自动配单事务处理失败:", err)
			response.Err(c, "事务处理失败")
			return
		}

		response.Ok(c)
	},
}
