package controllers

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
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
)

// UserDepositOnlines 存款管理-存款列表
var UserDepositOnlines = struct {
	*ActionUpdate
	*ActionSave
	List        func(*gin.Context)
	ConfirmDo   func(*gin.Context) // 人工确认
	GetStatus   func(*gin.Context) // 获取状态
	AddSlip     func(*gin.Context) // 添加存款单页面
	AddSlipSave func(*gin.Context) // 保存存款单
	OrderInfo   func(*gin.Context) // 查询订单信息
	UserInfo    func(*gin.Context) // 用户信息
	*ActionExport
}{
	List: func(c *gin.Context) {
		c.Set("is_online", true)
		UserDeposits.List(c)
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.UserDeposits,
		ViewFile: "user_deposit_onlines/user_deposits_edit.html",
		Row: func() interface{} {
			return &models.UserDeposit{}
		},
		ExtendData: func(c *gin.Context) pongo2.Context {
			depositCards := []models.DepositCard{}
			platform := request.GetPlatform(c)
			engine := common.Mysql(platform)
			defer engine.Close()
			if err := engine.Table("deposit_cards").Where("status=2").Find(&depositCards); err != nil {
				log.Logger.Error(err.Error())
			}
			payments := make([]models.Payment, 0)
			if err := engine.Table("payments").Where("is_online=2").Find(&payments); err != nil {
				log.Logger.Error(err.Error())
			}
			userLabel := func() string {
				if rows, err := engine.QueryString("SELECT users.label FROM users, user_deposits WHERE users.id = user_deposits.user_id AND user_deposits.id = ?", c.DefaultQuery("id", "0")); err == nil && len(rows) > 0 {
					return rows[0]["label"]
				}
				return ""
			}()

			return pongo2.Context{
				"rows":       depositCards,
				"rs":         payments,
				"user_label": userLabel,
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.UserDeposits,
		SaveBefore: func(c *gin.Context, data *map[string]interface{}) error {
			var startAt int64
			var endAt int64
			if value, exists := (*data)["created"].(string); !exists {
				currentTime := time.Now().Unix()
				startAt = tools.SecondToMicro(currentTime - currentTime%86400)
				endAt = startAt + tools.SecondToMicro(86400)
			} else {
				areas := strings.Split(value, " - ")
				startAt = tools.GetTimeStampByString(areas[0] + " 00:00:00")
				endAt = tools.GetTimeStampByString(areas[1]+" 00:00:00") + 86400
			}
			(*data)["time_start"] = startAt
			(*data)["time_end"] = endAt
			delete((*data), "created")
			return nil
		},
	},
	AddSlip: func(c *gin.Context) {
		depositCards := make([]models.DepositCard, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if err := engine.Table("deposit_cards").Where("status=2").Find(&depositCards); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查村收款银行卡错误")
			return
		}
		viewData := pongo2.Context{
			"rows":         depositCards,
			"payTypes":     consts.PaymentTypes,
			"channelTypes": models.Payments.All(platform),
		}
		response.Render(c, "user_deposit_onlines/add_slips.html", viewData)
	},
	AddSlipSave: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		stype := postedData["type"].(string)
		username := postedData["username"].(string)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		userInfo := &models.User{}
		b, err := engine.Table("users").Where("username=?", username).Get(userInfo)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查询用户失败")
			return
		}
		if !b {
			response.Err(c, "用户不存在失败")
			return
		}
		administrator := GetLoginAdmin(c)
		imap := map[string]interface{}{
			"order_no":  tools.GetBillNo("D", 5),
			"user_id":   userInfo.Id,
			"top_code":  userInfo.TopCode,
			"top_name":  userInfo.TopName,
			"top_id":    userInfo.TopId,
			"type":      stype,
			"money":     postedData["money"],
			"username":  postedData["username"],
			"comment":   postedData["comment"],
			"status":    1,
			"applicant": administrator.Name,
			"created":   tools.NowMicro(),
		}
		if stype == "1" { //在线存款
			payCode := postedData["pay_code"].(string)
			payCodeInfo := &models.Payment{}
			b, err := engine.Table("payments").Where("code=?", payCode).Get(payCodeInfo)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "查询支付编码错误")
				return
			}
			if !b {
				response.Err(c, "支付编码不存在")
				return
			}
			imap["channel_type"] = postedData["channel_type"]
			imap["pay_code"] = payCode
			imap["business_id"] = payCodeInfo.Id
			imap["business_name"] = payCode
		} else { //离线存款
			if accountByNameVal, exists := postedData["account_by_name"]; exists { // 如果有存款卡号相关信息
				accountByName := accountByNameVal.(string)
				if accountByName != "" {
					tempStr := strings.Split(accountByName, "-")
					depositCardInfo := &models.DepositCard{}
					if _, err := engine.Table("deposit_cards").Where("bank_card=?", tempStr[2]).Get(depositCardInfo); err != nil {
						log.Logger.Error(err.Error())
						response.Err(c, "查询收款银行卡错误")
						return
					}
					imap["card_number_id"] = depositCardInfo.Id
					imap["card_number"] = tempStr[2]
				}
				imap["account_by_name"] = accountByName
			}
			if depositName, exists := postedData["deposit_name"]; exists { // 如果有存款姓名
				imap["deposit_name"] = depositName
			}
		}
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		if _, err := session.Table("user_deposits").Insert(imap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "添加失败")
			return
		}
		iMap := map[string]interface{}{
			"bill_no":   postedData["order_no"],
			"type":      0,
			"operating": "后台手动添加",
			"result":    "成功",
			"operator":  administrator.Name,
			"created":   tools.NowMicro(),
		}
		if v, exists := postedData["remark"]; exists { // 如果有备注, 再添加备注
			iMap["remark"] = v
		}
		if _, err := session.Table("finance_logs").Insert(iMap); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "操作失败")
			return
		}
		_ = session.Commit()
		response.Ok(c)
	},
	ConfirmDo: func(c *gin.Context) { //手动确认
		postedData := request.GetPostedData(c)
		// 关于基础数据的验证 ----------------------------------------------------------------------------------
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		submit := postedData["submit"].(string)

		arriveMoneyStr := postedData["arrive_money"].(string)
		arriveMoney, err := strconv.ParseFloat(arriveMoneyStr, 64)
		if err != nil || arriveMoney < 0 {
			response.Err(c, "到账金额有误: "+arriveMoneyStr)
			return
		}
		confirmMoneyStr := postedData["confirm_money"].(string)
		confirmMoney, err := strconv.ParseFloat(confirmMoneyStr, 64)
		if err != nil || confirmMoney < 0 {
			response.Err(c, "确认金误有误: "+confirmMoneyStr)
			return
		}

		platform := request.GetPlatform(c)
		//防止多人同时更改
		rKey, err := redis.Lock(platform, "confirm-deposit-"+idStr)
		if err != nil {
			log.Err(err.Error())
			fmt.Println("缓存服务器加锁失败: ", err)
			response.Err(c, "请不要同一时间内多次提交")
			return
		}
		defer redis.Unlock(platform, rKey)

		administrator := GetLoginAdmin(c)
		r := &models.UserDeposit{}
		if exists, err := models.UserDeposits.FindById(platform, id, r); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查询用户存款信息失败")
			return
		}
		if r.Status != 1 {
			response.Err(c, "该订单已经被操作过")
			return
		}

		// 用于保存到用户存款记录表
		depositData := map[string]interface{}{
			"arrive_money":     arriveMoney,
			"confirm_money":    confirmMoney,
			"remark":           postedData["remark"],
			"status":           3,
			"finance_admin":    administrator.Name,
			"updated":          tools.NowMicro(),
			"is_first_deposit": 1,
		}

		// 用于保存到财务日志表
		financeData := map[string]interface{}{
			"bill_no":   postedData["order_no"],
			"type":      0,
			"operating": "存款结束",
			"result":    "失败",
			"operator":  administrator.Name,
			"consuming": time.Now().Unix() - int64(r.Created),
			"remark":    postedData["remark"],
			"created":   tools.NowMicro(),
		}

		if submit == "2" { //失败按钮 -----------------------------
			session := common.Mysql(platform)
			defer session.Close()
			if err := session.Begin(); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "事务启动失败")
				return
			}
			depositData["status"] = 3
			if _, err := session.Table("user_deposits").Where("id=?", id).Update(depositData); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "操作失败")
				return
			}
			if _, err := session.Table("finance_logs").Insert(financeData); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "操作失败")
				return
			}

			if r.CardNumberId > 0 { // 如果是失败,则对相应的银行卡进行累减
				if err := models.DepositCards.ReduceUsedMoney(platform, int(r.CardNumberId), r.Money, int(r.Created), session); err != nil {
					log.Logger.Error(err)
					_ = session.Rollback()
					response.Err(c, "扣减银行卡相关信息失败")
				}
			}
			_ = session.Commit()
			response.Message(c, "已将订单状态设置为失败")
			return
		}

		// 以下处理成功按钮 ---------------------------------------------

		if err := saveConfirmDeposit(platform, r, depositData, financeData, c); err != nil {
			response.Err(c, err.Error())
			return
		}

		response.Message(c, "操作成功")
	},
	GetStatus: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		userDepositInfo := &models.UserDeposit{}
		platform := request.GetPlatform(c)
		if exists, err := models.UserDeposits.FindById(platform, id, userDepositInfo); !exists || err != nil {
			response.Err(c, "查找用户存款信息失败")
			return
		}
		data := map[string]interface{}{
			"status": userDepositInfo.Status,
		}
		response.Result(c, data)
	},
	OrderInfo: func(c *gin.Context) {
		orderNumber := c.DefaultQuery("order_number", "")
		if orderNumber == "" {
			response.Err(c, "缺少订单号码")
			return
		}
		result, err := models.PaymentThirds.OrderInfo(orderNumber)
		if err != nil {
			response.Err(c, err.Error())
			return
		}

		response.Result(c, result)
	},
	UserInfo: func(c *gin.Context) {
		Platforms := request.GetPlatform(c)
		db := common.Mysql(Platforms)
		defer db.Close()
		username := c.Query("username")
		sql := "select a.realname,b.balance from users a join accounts b on a.username=b.username where a.username='" + username + "'"
		res, _ := db.QueryString(sql)
		if len(res) == 0 {
			response.Err(c, "用户不存在")
		} else {
			data := map[string]interface{}{
				"realname": res[0]["realname"],
				"money":    res[0]["balance"],
			}
			response.Result(c, data)
		}
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
