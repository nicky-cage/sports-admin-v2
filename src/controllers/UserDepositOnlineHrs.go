package controllers

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	sredis "sports-common/redis"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"xorm.io/builder"
)

// UserDepositOnlineHrs 存款管理-历史记录
var UserDepositOnlineHrs = struct {
	List      func(*gin.Context)
	Mistake   func(*gin.Context) //失误反转
	MistakeDo func(*gin.Context) //失误反转保存
	Fix       func(*gin.Context) // 补单
	*ActionExport
}{
	List: func(c *gin.Context) {
		c.Set("is_online", true)
		UserDepositHrs.List(c)
	},
	Mistake: func(c *gin.Context) {
		idStr := c.Query("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.ErrorHTML(c, err.Error())
			return
		}
		userDepositsInfo := &models.UserDeposit{}
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if _, err := engine.Table("user_deposits").Where("id=?", id).Get(userDepositsInfo); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "信息不存在")
			return
		}
		depositCards := make([]models.DepositCard, 0)
		if err := engine.Table("deposit_cards").Where("status=2").Find(&depositCards); err != nil {
			log.Logger.Error(err.Error())
		}
		viewData := pongo2.Context{
			"dr": depositCards,
			"r":  userDepositsInfo,
		}
		viewFile := "user_deposits/mistakes.html"
		response.Render(c, viewFile, viewData)
	},
	MistakeDo: func(c *gin.Context) { // -- 检测基本数据信息
		postedData := request.GetPostedData(c)
		id, err := strconv.Atoi(postedData["id"].(string)) // 编号
		if err != nil {
			response.Err(c, err.Error())
			return
		}
		remark := postedData["remark"].(string)
		if remark == "" {
			response.Err(c, "必须输入备注")
			return
		}
		cardNumberId, err := strconv.Atoi(postedData["card_number_id"].(string))
		if err != nil {
			response.Err(c, "提交的银行卡相关信息有误")
			return
		}
		submit := postedData["submit"].(string) // 提交成功/失败/失败不扣款

		// 获取存款订单信息
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		r := &models.UserDeposit{} // 订单信息
		if exists, err := engine.Table("user_deposits").Where("id=?", id).Get(r); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查询记录错误")
			return
		} else if !exists {
			response.Err(c, "无法查找到此存款订单信息")
			return
		}
		status := r.Status
		if status != 2 && status != 3 { // 此时订单必须处理成功/失败状态
			response.Err(c, "订单状态异异常")
			return
		}

		// 获取用户信息
		userInfo := &models.User{}
		if exists, err := models.Users.FindById(platform, int(r.UserId), userInfo); !exists || err != nil {
			log.Logger.Error(err)
			response.Err(c, "查找用户信息失败")
			return
		}

		rKey, err := sredis.Lock(platform, fmt.Sprintf("mistake-%d", id))
		if err != nil { //防止多人同时更改
			response.Err(c, err.Error())
			return
		}
		defer sredis.Unlock(platform, rKey)

		admin := GetLoginAdmin(c)
		card := caches.DepositCards.Get(platform, cardNumberId)
		if card == nil {
			response.Err(c, "无法获取银行卡信息")
			return
		}
		AccountByName := fmt.Sprintf("%s-%s-%s", card.BankCode, card.BankRealname, card.BankCard)
		saveData := map[string]interface{}{
			"updated":          tools.NowMicro(),
			"account_by_name":  AccountByName,
			"card_number_id":   cardNumberId,
			"card_number":      card.BankCard,
			"business_name":    card.BankRealname,
			"is_first_deposit": 1,
			"finance_admin":    admin.Name,
			"remark":           remark,
		}
		financeData := map[string]interface{}{
			"type":      0,
			"bill_no":   r.OrderNo,
			"result":    "",
			"operating": "",
			"operator":  admin.Name,
			"consuming": tools.Now() - int64(r.Updated/1000000),
			"remark":    remark,
			"created":   tools.NowMicro(),
		}

		switch submit {
		case "2_2": // 保存成功信息
			financeData["operating"] = "修改存款信息"
			financeData["result"] = "成功"
		case "2_3": // 将成功修改为失败, 需要减钱
			saveData["status"] = 3
			financeData["operating"] = "由成功修改为失败"
			financeData["result"] = "失败"
			saveData["confirm_money"] = 0
		case "3_2": // 将失败修改为成功, 需要加钱
			saveData["status"] = 2
			saveData["confirm_at"] = tools.Now()
			financeData["operating"] = "由失败修改为成功"
			financeData["result"] = "成功"
			// 需要判断是否首存
			if models.UserDeposits.IsFirstTime(platform, int(r.UserId)) {
				saveData["is_first_time"] = 2
			}
			//discount := models.Activities.GetActivityAmount(int(r.Id)) +
			//	r.ConfirmMoney*models.UserDepositDiscounts.GetDiscount(int(r.ChannelType), int(userInfo.Vip)) // 优惠 = 活动优惠 + 存款优惠
			saveData["top_money"] = r.Money + r.Discount // 上分金额 = 存款金额(实收) + 存款优惠
		case "3_3": // 保存失败信息
			financeData["operating"] = "修改存款信息"
			financeData["result"] = "失败"
		case "2_30": // 将成功修改为失败, 但不扣钱
			saveData["status"] = 3
			financeData["operating"] = "由成功修改为失败但不扣款"
			financeData["result"] = "失败"
		default:
			response.Err(c, "提交的状态变更信息有误")
			return
		}

		session := engine
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		if affected, err := session.Table("user_deposits").Where("id = ? ", r.Id).Update(saveData); err != nil || affected <= 0 {
			log.Logger.Error(err)
			_ = session.Rollback()
			response.Err(c, "保存存款订单信息失败")
			return
		}
		if id, err := session.Table("finance_logs").Insert(financeData); err != nil {
			_ = session.Rollback()
			response.Err(c, "保存财务日志失败: "+err.Error())
			return
		} else if id == 0 {
			_ = session.Rollback()
			response.Err(c, "保存财务日志出错")
			return
		}

		// 以下, 写账变信息
		if submit == "2_3" || submit == "3_2" { // 需要钱方面的加/减
			// 获取用户账户信息
			accountInfo := &models.Account{}
			if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().
				And(builder.Eq{"user_id": r.UserId})); !exists || err != nil {
				_ = session.Rollback
				log.Logger.Error(err)
				response.Err(c, "查找用户账户信息失败")
				return
			}

			// 如果是由成功修改为失败, 则不需要判断首存, 直接减去上分金额即可
			var money = r.ConfirmMoney
			if submit == "2_3" { // 如果是由成功 => 失败, 需要扣的总额 = 上分金额, 上分金额 = 到账金额 + 优惠金额
				money = 0 - r.ConfirmMoney - r.Discount
			}

			transAction := &models.Transaction{}
			extraMap := map[string]interface{}{
				"proxy_ip":      "",
				"ip":            c.ClientIP(),
				"description":   "存款-订单状态反转",
				"administrator": admin.Name,
				"admin_user_id": admin.Id,
				"serial_number": r.OrderNo,
			}
			transType := consts.TransTypeRechargeOffline
			if submit == "2_3" {
				transType = consts.TransTypeRechargeDeduct
			}
			platform := request.GetPlatform(c)
			redisSession := common.Redis(platform)
			defer common.RedisRestore(platform, redisSession)
			if _, err := transAction.AddTransaction(platform, session, redisSession, userInfo, accountInfo, transType, money, extraMap); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}

			// 需要写银行卡限额相关信息
			if r.Type == 2 { // 如果是离线存款
				// 如果是由成功=>失败, 并且是今天的存款,则减去银行卡的累计
				if submit == "2_3" && int64(r.Created) > tools.GetTodayBegin() { // 由成功 => 失败, 需要减去限额
					if err := models.DepositCards.ReduceUsedMoney(platform, int(r.CardNumberId), r.ConfirmMoney, int(r.Created), session); err != nil {
						_ = session.Rollback()
						response.Err(c, err.Error())
						return
					}
				}
				if submit == "3_2" { // 如果有失败=>成功, 则需要加上存款额度
					if err := models.DepositCards.AddUsedMoney(platform, int(r.CardNumberId), r.ConfirmMoney, session); err != nil {
						_ = session.Rollback()
						response.Err(c, err.Error())
						return
					}
				}
			}

			if accountInfo.Id > 0 { //覆盖用户钱包的数据
				_ = accountInfo.ResetCacheData(redisSession)
			}
		}

		_ = session.Commit()
		response.Message(c, "操作成功")
	},
	Fix: func(c *gin.Context) {
		orderNum := c.DefaultQuery("order_no", "")
		if orderNum == "" {
			response.Err(c, "错误的订单号码")
			return
		}

		platform := request.GetPlatform(c)
		rKey, err := sredis.Lock(platform, "fix-"+orderNum) // 加锁
		if err != nil {
			fmt.Println("缓存服务器加锁失败: ", err)
			response.Err(c, err.Error())
			return
		}
		defer sredis.Unlock(platform, rKey)

		order := models.UserDeposit{}
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		exists, err := dbSession.Table("user_deposits").
			Where("order_no = ?", orderNum).
			Get(&order)
		if !exists { // 订单不存在的情况
			response.Err(c, "订单记录查询失败")
			return
		}
		if err != nil { // 表示出错
			response.Err(c, "订单查询出错: "+err.Error())
			return
		}
		if tools.NowMicro()-order.Created > 86400*1000000 {
			response.Err(c, "无法补单: 订单过期已经超过1天")
			return
		}

		if order.Status == 2 { // 已经成功
			response.Err(c, "订单已经是成功状态, 不需要再次补单")
			return
		}

		rd := common.Redis(platform)
		defer common.RedisRestore(platform, rd)
		key := "fix:" + orderNum

		// 检查是否已经缓存此订单成功信息
		// 将此订单号信息保存到缓存当中, 1 分钟内再次提交时则更新并为用户上分
		isConfirm := c.DefaultQuery("confirm", "1")
		if isConfirm == "2" {
			statusStr, err := rd.Get(key).Result()
			if err != nil && err != redis.Nil {
				response.Err(c, err.Error())
				return
			}
			if err == redis.Nil {
				response.Err(c, "三方查单已经超时,请再试尝试")
				return
			}
			val, err := strconv.Atoi(statusStr)
			if err != nil {
				response.Err(c, "账号金额转换出错: "+err.Error())
				return
			}
			if val != models.ORDER_STATUS_SUCCESS { // 如查状态异常
				response.Err(c, "补单状态异常: 三方平台状态不匹配")
				return
			}

			// 各方确定正常, 开始进行补单程序计算
			administrator := GetLoginAdmin(c)
			depositData := map[string]interface{}{ // 用于保存到用户存款记录表
				"arrive_money":     order.Money,
				"confirm_money":    order.Money,
				"remark":           "后台财务补单",
				"status":           3,
				"finance_admin":    administrator.Name,
				"updated":          tools.NowMicro(),
				"is_first_deposit": 1,
			}
			financeData := map[string]interface{}{ // 用于保存到财务日志表
				"bill_no":   order.OrderNo,
				"type":      0,
				"operating": "存款结束",
				"result":    "成功",
				"operator":  administrator.Name,
				"consuming": time.Now().Unix() - int64(order.Created/1000000),
				"remark":    "后台财务补单",
				"created":   tools.NowMicro(),
			}

			if err := saveConfirmDeposit(platform, &order, depositData, financeData, c); err != nil {
				response.Err(c, err.Error())
				return
			}

			_, _ = rd.Del(key).Result() // 删除缓存
			response.Ok(c)              // 表示补单成功
			return
		}

		status, err := models.PaymentThirds.QueryOrder(orderNum)
		if err != nil {
			log.Logger.Error("三方查单失败: ", err)
			response.Err(c, "三方平台查单失败: "+err.Error())
			return
		}

		// 如果本地不是成功状态, 三方平台是成功状态 -> 则本地转换为成功状态 -> 为用户上分
		if status != models.ORDER_STATUS_SUCCESS {
			fmt.Println("手动补单失败: status = ", status)
			response.Err(c, "手动补单成功!<br />三方平台支付状态: 失败")
			return
		}

		rd.Set(key, status, time.Second*60) // 保存订单信息到redis

		response.Result(c, map[string]interface{}{
			"status":   status,
			"amount":   fmt.Sprintf("%.2f", order.Money),    // 存款金额
			"discount": fmt.Sprintf("%.2f", order.Discount), // 优惠金额
		})
	},
	ActionExport: &ActionExport{
		Columns: []ExportHeader{
			{"序号", "id"},
			{"订单编号", "order_no"},
			{"会员编号", "user_id"},
			{"会员名称", "username"},
			{"会员等级", "vip"},
			{"订单金额", "money"},
			{"到账金额", "arrive_money"},
			{"上分金额", "top_money"},
			{"存款优惠", "discount"},
			{"支付方式", "channel_type"},
			{"订单时间", "created"},
			{"完成时间", "updated"},
			{"支付渠道/编码", "account_by_name"},
			{"操作人", "finance_admin"},
			{"状态", "status"},
		},
		ProcessRawData: func(fields []string, rArr *[]map[string]interface{}, c *gin.Context) {
			base_controller.ExportRawDataReset(rArr)
		},
		ProcessRow: func(m *map[string]interface{}, c *gin.Context) {
			(*m)["id"] = (*m)["user_deposit.id"]
			(*m)["order_no"] = (*m)["user_deposit.order_no"]
			(*m)["username"] = (*m)["user_deposit.username"]
			(*m)["vip"] = base_controller.FieldToUserVip(c, (*m)["user.vip"])
			(*m)["money"] = fmt.Sprintf("%.2f", (*m)["user_deposit.money"].(float64))
			(*m)["arrive_money"] = fmt.Sprintf("%.2f", (*m)["user_deposit.arrive_money"].(float64))
			(*m)["top_money"] = fmt.Sprintf("%.2f", (*m)["user_deposit.top_money"].(float64))
			(*m)["discount"] = fmt.Sprintf("%.2f", (*m)["user_deposit.discount"].(float64))
			(*m)["created"] = base_controller.FieldToDateTime(fmt.Sprintf("%d", int((*m)["user_deposit.created"].(float64))))
			(*m)["updated"] = func() string {
				if (*m)["user_deposit.updated"] == nil {
					return ""
				}
				return base_controller.FieldToDateTime(fmt.Sprintf("%d", int((*m)["user_deposit.updated"].(float64))))
			}()
			(*m)["finance_admin"] = (*m)["user_deposit.finance_admin"]
			(*m)["channel_type"] = func() string {
				currentType := int((*m)["user_deposit.channel_type"].(float64))
				switch currentType {
				case 0:
					return "银行转账"
				case 1:
					return "网银转账"
				case 2:
					return "支付宝"
				case 3:
					return "微信"
				case 4:
					return "QQ钱包"
				case 5:
					return "快捷支付"
				case 6:
					return "京东"
				case 7:
					return "银行扫码"
				case 8:
					return "虚拟币"
				case 9:
					return "云闪付"
				default:
					return "未知渠道"
				}
			}()
			(*m)["account_by_name"] = func() string {
				name := base_controller.FieldToPaymentType(c, (*m)["user_deposit.account_by_name"])
				return fmt.Sprintf("%v-%v", name, (*m)["user_deposit.account_by_name"])
			}()
			(*m)["status"] = func() string {
				switch int((*m)["user_deposit.status"].(float64)) {
				case 1:
					return "待确认"
				case 2:
					return "成功"
				case 3:
					return "失败"
				}
				return "未知"
			}()
		},
	},
}
