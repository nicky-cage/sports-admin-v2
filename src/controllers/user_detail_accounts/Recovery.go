package user_detail_accounts

import (
	"fmt"
	common "sports-common"
	"sports-common/config"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/redis"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
)

// Recovery 钱包回收
func (ths *UserDetailAccounts) Recovery(c *gin.Context) {
	idStr := c.DefaultQuery("id", "0") // 用户id
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		response.Err(c, "用户编号错误")
		return
	}
	gameID, gameErr := strconv.Atoi(c.DefaultQuery("game_id", "0"))
	if gameErr != nil || id < 0 {
		response.Err(c, "游戏编号错误")
		return
	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()

	// 检查是否有未结算注单, 如果有, 则不能回收
	uRow := struct {
		RegisterTime int `json:"register_time" xorm:"created"` // 注册时间
		LastLoginAt  int `json:"last_login_at"`                // 最后登录时间
	}{}
	uSQL := fmt.Sprintf("SELECT last_login_at, created FROM users WHERE id = %d", id)
	if exists, err := myClient.SQL(uSQL).Get(&uRow); err != nil {
		log.Logger.Error("查找用户相关记录错误: ", err)
		response.Err(c, "查找用户相关记录出错")
		return
	} else if !exists {
		log.Logger.Error("查找不到用户相关记录错误: ")
		response.Err(c, "查找用户相关记录失败")
		return
	}

	dRow := struct {
		LastDepositTime int `json:"last_deposit_time" xorm:"last_deposit_time"` // 最后存款时间
	}{}
	dSQL := fmt.Sprintf("SELECT MAX(created) AS last_deposit_time "+
		"FROM user_deposits WHERE user_id = %d AND created > %d AND status = 2", id, uRow.RegisterTime)
	if exists, err := myClient.SQL(dSQL).Get(&dRow); err != nil {
		log.Logger.Error("查找用户充值记录有误: ", err)
		response.Err(c, "查找用户充值记录有误")
		return
	} else if !exists {
		response.Err(c, "用户没有充值记录, 不需要回收")
		return
	}

	// gameCode := ""
	// if gameID > 0 {
	// 	gRow := models.GameVenue{}
	// 	if exists, err := myClient.SQL("SELECT * FROM game_venues WHERE id = ?", gameID).Get(&gRow); err != nil || !exists {
	// 		log.Logger.Error("查找游戏信息有误:", err)
	// 		response.Err(c, "游戏信息查找失败")
	// 		return
	// 	}
	// 	gameCode = gRow.Code
	// }
	// pgClient := pgsql.GetConnForReading(platform)
	// defer pgClient.Close()
	// currentTime := tools.Now()
	// wRow := struct {
	// 	Total int `json:"total"`
	// }{}
	// wSQL := fmt.Sprintf("SELECT COUNT(*) AS total "+
	// 	"FROM wager_records "+
	// 	"WHERE user_id = %d AND status = 0 AND created_at < %d AND created_at > %d", id, currentTime, dRow.LastDepositTime)
	// if gameCode != "" {
	// 	wSQL += fmt.Sprintf(" AND game_code = '%s'", gameCode)
	// }
	// log.Logger.Info(wSQL)
	// if _, err := pgClient.QueryOne(&wRow, wSQL); err != nil {
	// 	log.Logger.Error("查找用记投注记录失败: ", err)
	// 	response.Err(c, "查找用户投注记录错误")
	// 	return
	// } else if wRow.Total > 0 {
	// 	response.Err(c, "用户 "+idStr+" 有未结算注单, 不能回收")
	// 	return
	// }
	// // 最后存款 30 分钟后, 才能回收, 以防止有未拉回注单
	// if currentTime-int64(dRow.LastDepositTime) < 60*30 {
	// 	response.Err(c, "用户存款30分钟内不能回收")
	// 	return
	// }

	// 以下, 强制用户回收余额
	if c.DefaultQuery("code", "normal") == "force" {
		rClient := common.Redis(platform)
		defer common.RedisRestore(platform, rClient)
		lockKey, err := redis.Lock(platform, "user_reback_wallet:"+idStr)
		defer redis.Unlock(platform, lockKey)
		if err != nil {
			response.Err(c, "锁定用户操作失败")
			return
		}
		// 先处理 user_games_v1 里面的东西
		sql := fmt.Sprintf("SELECT SUM(money_v1) AS total FROM user_games WHERE game_code <> 'IME' AND user_id = %d AND `version` = 1", id)
		if err := myClient.Begin(); err != nil {
			response.Err(c, "事务处理失败")
			return
		}
		rows, err := myClient.QueryString(sql)
		if err == nil && len(rows) > 0 {
			total, err := strconv.ParseFloat(rows[0]["total"], 64)
			//fmt.Println("total = ", total)
			if err == nil && total > 0 { // 有金额才会处理
				user := models.User{}
				exists, err := myClient.Table("users").SQL("SELECT * FROM users WHERE id = ?", id).Get(&user)
				if err != nil || !exists {
					myClient.Rollback()
					//fmt.Println("获取用户信息出现错误: ", err, ", exists = ", exists)
					response.Err(c, "获取用户信息用误")
					return
				}
				userAccount := models.Account{}
				exists, err = myClient.Table("accounts").SQL("SELECT * FROM accounts WHERE user_id = ?", id).Get(&userAccount)
				if err != nil || !exists {
					myClient.Rollback()
					//fmt.Println("获取用户账户信息有误: ", err, ", exists = ", exists)
					response.Err(c, "获取账户信息用误")
					return
				}

				// 调整信息
				userId := user.Id
				userName := user.Username
				centerMoney := userAccount.Available // 中心钱包
				adjustType := 2                      // 调整类型 - 系统调整
				adjustMethod := 1                    // 调整方法 - 上分
				flowLimit := 0                       // 流水限制
				flowMultiple := 0                    // 流水倍数
				adjustMoney := total                 // 调整金额
				remark := "三方平台维护,强制回收钱包"            // 备注
				adminName := "admin"                 // 管理员名称
				billNo := tools.GetBillNo("hl", 5)   // 订单号
				vip := user.Vip                      // VIP
				topId := user.TopId                  // 上级编号
				vSQL := " INSERT INTO user_resets " +
					"(user_id, username, center_money, adjust_type, adjust_method, flow_limit, flow_multiple, adjust_money, " +
					"remark, updated, admin_name, bill_no, vip, top_id, created, status) " +
					"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 2) "
				currentTime := tools.NowMicro()
				res, err := myClient.Exec(vSQL,
					userId, userName, centerMoney, adjustType, adjustMethod,
					flowLimit, flowMultiple, adjustMoney, remark, currentTime, adminName, billNo, vip, topId, currentTime)
				if err != nil {
					myClient.Rollback()
					response.Err(c, "保存调整信息出错")
					return
				}

				effected, err := res.LastInsertId()
				if err != nil || effected <= 0 {
					myClient.Rollback()
					response.Err(c, "写入数据库时出错")
					return
				}

				ua := models.UserAccountSet{ // 账户信记录
					UserId:          user.Id,
					Username:        user.Username,
					Applicant:       "admin",
					Audit:           "admin",
					Created:         currentTime,
					Updated:         currentTime,
					ApplicantRemark: "强制回收游戏钱包余额",
					AuditRemark:     "三方平台维护",
					BillNo:          billNo,
					Type:            uint64(adjustMethod),
					Money:           total,
					UserVip:         user.Vip,
					Status:          2, // status
				}
				transType := consts.TransTypeAdjustmentPlus
				transAction := &models.Transaction{}

				extraMap := map[string]interface{}{
					"proxy_ip":      "",
					"ip":            c.ClientIP(),
					"description":   "系统调整,强制回收钱包余额",
					"administrator": "admin",
					"admin_user_id": 1001,
					"serial_number": billNo,
				}
				if _, err := transAction.AddTransaction(platform, myClient, rClient, &user, &userAccount, transType, total, extraMap); err != nil {
					myClient.Rollback()
					response.Err(c, err.Error())
					return
				}
				if _, err = myClient.Table("user_account_sets").Insert(ua); err != nil {
					myClient.Rollback()
					response.Err(c, err.Error())
					return
				}

				// 以后升级为2.0之后,将不再访问此数据库
				res, err = myClient.Exec("UPDATE user_games SET `version` = 2 WHERE user_id = ? AND `version` = 1", id)
				if err != nil {
					myClient.Rollback()
					fmt.Println(" -- 修改版本号出错: ", err)
					response.Err(c, err.Error())
					return
				}
				_, err = res.LastInsertId()
				if err != nil {
					myClient.Rollback()
					fmt.Println("err: ", err)
					response.Err(c, err.Error())
					return
				}
				if err = myClient.Commit(); err != nil {
					myClient.Rollback()
					fmt.Println("提交事务出错: ", err)
					response.Err(c, err.Error())
					return
				}
			}
		}
	}

	// 以下, 再回收用户余额
	req.SetTimeout(30 * time.Second)
	req.Debug = true
	header := req.Header{"Accept": "application/json"}
	//url 是game.ip.vhost
	baseTransferUrl := config.Get("internal.internal_game_service") + config.Get("internal_api.recovery_url")
	TransferUrl := baseTransferUrl + "?user_id=" + idStr + "&platform=" + platform + "&game_id=" + strconv.Itoa(gameID)
	log.Logger.Info("transfer url = ", TransferUrl)
	r, err := req.Post(TransferUrl, header)
	if err != nil {
		log.Logger.Error("回收发生错误: ", err.Error())
		response.Err(c, "系统异常")
		return
	}
	res := TransferState{}
	_ = r.ToJSON(&res)
	if res.Errcode != 0 {
		log.Logger.Error("回收错误: ", res.Errcode, ", ", res.Message)
		response.Err(c, "回收错误")
		return
	}
	response.Ok(c)
}
