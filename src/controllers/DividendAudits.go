package controllers

import (
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"xorm.io/builder"
)

type TransferResp struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct {
		Money       string `json:"money"`
		GameBalance []struct {
			GameCode string `json:"game_code"`
			Balance  string `json:"balance"`
		} `json:"game_balance"`
	} `json:"data"`
}

type BalanceResp struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct {
		Balance  float64 `json:"balance"`
		GameCode string  `json:"game_code"`
		State    int     `json:"state"`
	} `json:"data"`
}

type SumStruct struct {
	Money float64
}

// 红利-审核列表
var DividendAudits = struct {
	List        func(*gin.Context)
	EditView    func(*gin.Context)
	Agree       func(*gin.Context) //同意
	Refuse      func(*gin.Context) //拒绝
	BatchAgree  func(*gin.Context) //批量同意
	BatchRefuse func(*gin.Context) //批量拒绝
}{
	List: func(c *gin.Context) {
		cond := builder.NewCond()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); !exists {
			currentTime := time.Now().Unix()
			startAt = tools.SecondToMicro(currentTime - currentTime%86400)
			endAt = startAt + tools.SecondToMicro(86400)
		} else {
			areas := strings.Split(value, " - ")
			startAt = tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
			endAt = tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
		}
		cond = cond.And(builder.Gte{"created": startAt}).And(builder.Lte{"created": endAt})
		cond = cond.And(builder.Eq{"state": 1})
		username := c.DefaultQuery("username", "")
		billNo := c.DefaultQuery("bill_no", "")
		applicant := c.DefaultQuery("applicant", "")
		sType := c.DefaultQuery("type", "")
		moneyType := c.DefaultQuery("money_type", "")
		flowLimit := c.DefaultQuery("flow_limit", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"username": username})
		}
		if len(applicant) > 0 {
			cond = cond.And(builder.Eq{"applicant": applicant})
		}
		if len(sType) > 0 {
			cond = cond.And(builder.Eq{"type": sType})
		}
		if len(billNo) > 0 {
			cond = cond.And(builder.Eq{"bill_no": billNo})
		}
		if len(moneyType) > 0 {
			cond = cond.And(builder.Eq{"money_type": moneyType})
		}
		if len(flowLimit) > 0 {
			cond = cond.And(builder.Eq{"flow_limit": flowLimit})
		}
		limit, offset := request.GetOffsets(c)
		userDividends := make([]models.UserDividend, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("user_dividends").Where(cond).OrderBy("id DESC").Limit(limit, offset).FindAndCount(&userDividends)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		ss := new(SumStruct)
		sumTotal, _ := engine.Table("user_dividends").Where(cond).Sum(ss, "money")
		viewData := pongo2.Context{
			"rows":      userDividends,
			"total":     total,
			"sum_money": sumTotal,
		}
		viewFile := "dividend_managements/applies.html"
		if request.IsAjax(c) {
			viewFile = "dividend_managements/_applies.html"
		}
		response.Render(c, viewFile, viewData)
	},
	EditView: func(c *gin.Context) {
		idStr, exists := c.GetQuery("id")
		if !exists || idStr == "" {
			log.Err("无法获取id信息!\n")
			return
		}
		typeStr, _ := c.GetQuery("type")
		data := make(map[string]string)
		data["id"] = idStr
		viewData := pongo2.Context{"r": data}
		if typeStr == "agree" {
			response.Render(c, "dividend_managements/view_agree.html", viewData)
		} else if typeStr == "refuse" {
			response.Render(c, "dividend_managements/view_refuse.html", viewData)
		} else if typeStr == "batch_agree" {
			response.Render(c, "dividend_managements/view_batch_agree.html", viewData)
		} else {
			response.Render(c, "dividend_managements/view_batch_refuse.html", viewData)
		}
	},
	Agree: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		platform := request.GetPlatform(c)
		dividend := &models.UserDividend{}
		b, err := models.UserDividends.FindById(platform, id, dividend)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查询错误")
			return
		}
		if !b {
			response.Err(c, "找不这条记录")
			return
		}
		if dividend.State != 1 {
			response.Err(c, "这条红利申请已经被处理")
			return
		}
		//防止多人同时更改
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		rKey := dividend.BillNo + "_" + idStr + "_" + dividend.Username
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "redis系统繁忙")
			return
		}
		if num > 1 {
			response.Err(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)

		uMap := make(map[string]interface{})
		userInfo := &models.User{}
		if exists, err := models.Users.FindById(platform, int(dividend.UserId), userInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户信息失败")
			return
		}
		accountInfo := &models.Account{}
		if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": dividend.UserId})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			response.Err(c, "查找用户账户信息失败")
			return
		}
		//事务操作
		session := common.Mysql(platform)
		defer session.Close()
		if err := session.Begin(); err != nil {
			log.Logger.Error(err.Error())
			_ = session.Rollback()
			response.Err(c, "事务启动失败")
			return
		}
		administrator := GetLoginAdmin(c)
		successDo := false
		if dividend.MoneyType == 1 { //中心钱包
			transAction := &models.Transaction{}
			extraMap := map[string]interface{}{
				"proxy_ip":      "",
				"ip":            c.ClientIP(),
				"description":   "红利审核-同意",
				"administrator": administrator.Name,
				"admin_user_id": administrator.Id,
				"serial_number": dividend.BillNo,
			}
			transType := consts.TransTypeAdjustmentDividendPlus
			if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, dividend.Money, extraMap); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}
			uMap["before_money"] = accountInfo.Available
			uMap["money"] = dividend.Money
			uMap["after_money"] = accountInfo.Available + dividend.Money
			successDo = true
		} else { //场馆钱包
			//如果外接请求成功successDo = true
			req.SetTimeout(50 * time.Second)
			req.Debug = true
			//先查余额
			headerB := req.Header{
				"Accept": "application/json",
			}
			paramB := req.Param{
				"game_code": dividend.Venue,
			}
			baseBalanceUrl := consts.InternalGameServUrl
			BalanceUrl := baseBalanceUrl + "/game/v1/internal/balance?user_id=" + strconv.Itoa(int(dividend.UserId))
			rB, err := req.Post(BalanceUrl, headerB, paramB)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, err.Error())
				return
			}
			balanceResp := BalanceResp{}
			BalanceJsonErr := rB.ToJSON(&balanceResp)
			if BalanceJsonErr != nil {
				log.Logger.Error(BalanceJsonErr.Error())
				response.Err(c, "同意失败")
				return
			}
			if balanceResp.Errcode != 0 { //查询余额失败
				response.Err(c, "同意失败")
				return
			}
			//再转账
			header := req.Header{
				"Accept": "application/json",
			}
			param := req.Param{
				"out_code": "CENTERWALLET",
				"in_code":  dividend.Venue,
				"all_in":   2,
				"money":    dividend.Money,
			}
			baseTransferUrl := consts.InternalGameServUrl
			TransferUrl := baseTransferUrl + "/game/v1/internal/transfer?user_id=" + strconv.Itoa(int(dividend.UserId))
			r, err := req.Post(TransferUrl, header, param)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, err.Error())
				return
			}
			transferResp := TransferResp{}
			transferJsonErr := r.ToJSON(&transferResp)
			if transferJsonErr != nil {
				log.Logger.Error(transferJsonErr.Error())
				response.Err(c, "同意失败")
				return
			}
			if transferResp.Errcode != 0 {
				response.Err(c, "同意失败")
				return
			} else {
				uMap["before_money"] = balanceResp.Data.Balance
				uMap["money"] = dividend.Money
				beforeMoney := balanceResp.Data.Balance
				uMap["after_money"] = beforeMoney + dividend.Money
				successDo = true
			}
		}
		if successDo {
			uMap["reviewer_remark"] = postedData["reviewer_remark"]
			uMap["state"] = 2
			uMap["reviewer"] = GetLoginAdmin(c).Name
			uMap["updated"] = time.Now().UnixMicro()
			if _, err := session.Table("user_dividends").Where("id=?", id).Update(uMap); err != nil {
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "操作失败")
				return
			}
		}
		_ = session.Commit()
		//覆盖用户钱包的数据
		if accountInfo.Id > 0 {
			_ = accountInfo.ResetCacheData(redis)
		}
		response.Ok(c)
	},
	Refuse: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		dividend := &models.UserDividend{}
		platform := request.GetPlatform(c)
		b, err := models.UserDividends.FindById(platform, id, dividend)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "查询错误")
			return
		}
		if !b {
			response.Err(c, "找不这条红利申请记录")
			return
		}
		if dividend.State != 1 {
			response.Err(c, "这条红利申请已经被处理")
			return
		}
		//防止多人同时更改
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		rKey := dividend.BillNo + "_" + idStr + "_" + dividend.Username
		num, err := redis.Incr(rKey).Result()
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "redis系统繁忙")
			return
		}
		if num > 1 {
			response.Err(c, "请稍等片刻再试")
			return
		}
		defer redis.Del(rKey)
		//更改状态
		engine := common.Mysql(platform)
		defer engine.Close()
		uMap := map[string]interface{}{
			"reviewer_remark": postedData["reviewer_remark"],
			"state":           3,
			"reviewer":        GetLoginAdmin(c).Name,
			"updated":         tools.NowMicro(),
		}
		if dividend.MoneyType == 1 { //中心钱包
			accountInfo := &models.Account{}
			_, _ = models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": dividend.UserId}))
			uMap["before_money"] = accountInfo.Available
			uMap["money"] = dividend.Money
			uMap["after_money"] = accountInfo.Available
		} else { //场馆钱包
			//查询余额
			req.SetTimeout(50 * time.Second)
			req.Debug = true
			header := req.Header{
				"Accept": "application/json",
			}
			param := req.Param{
				"game_code": dividend.Venue,
			}
			baseBalanceUrl := consts.InternalGameServUrl
			BalanceUrl := baseBalanceUrl + "/game/v1/internal/balance?user_id=" + strconv.Itoa(int(dividend.UserId))
			r, err := req.Post(BalanceUrl, header, param)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, err.Error())
				return
			}
			balanceResp := BalanceResp{}
			_ = r.ToJSON(&balanceResp)
			if balanceResp.Errcode != 0 { //查询余额失败
				response.Err(c, "拒绝失败")
				return
			} else {
				uMap["before_money"] = balanceResp.Data.Balance
				uMap["money"] = dividend.Money
				beforeMoney := balanceResp.Data.Balance
				uMap["after_money"] = beforeMoney + dividend.Money
			}
		}
		if _, err := engine.Table("user_dividends").Where("id=?", id).Update(uMap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "操作失败")
			return
		}
		response.Ok(c)
	},
	BatchAgree: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		idSilce := strings.Split(idStr, ",")
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		administrator := GetLoginAdmin(c)
		for _, v := range idSilce {
			uMap := make(map[string]interface{})
			successDo := false
			dividend := &models.UserDividend{}
			id, _ := strconv.Atoi(v)
			b, err := models.UserDividends.FindById(platform, id, dividend)
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			//记录不存在或已经被处理过
			if !b || dividend.State != 1 {
				continue
			}
			rKey := dividend.BillNo + "_" + idStr + "_" + dividend.Username
			num, err := redis.Incr(rKey).Result()
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			//这条记录正在被处理的过滤
			if num > 1 {
				continue
			}
			session := common.Mysql(platform)
			defer session.Close()
			if err := session.Begin(); err != nil {
				log.Err(err.Error())
				_ = session.Rollback()
				continue
			}
			userInfo := &models.User{}
			if exists, err := models.Users.FindById(platform, int(dividend.UserId), userInfo); !exists || err != nil {
				if err != nil {
					log.Logger.Error(err.Error())
				}
				response.Err(c, "查找用户信息失败")
				return
			}
			accountInfo := &models.Account{}
			if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": dividend.UserId})); !exists || err != nil {
				if err != nil {
					log.Logger.Error(err.Error())
				}
				response.Err(c, "查找用户账户信息失败")
				return
			}
			if dividend.MoneyType == 1 { //中心钱包
				transAction := &models.Transaction{}
				extraMap := map[string]interface{}{
					"proxy_ip":      "",
					"ip":            c.ClientIP(),
					"description":   "红利批量审核-同意",
					"administrator": administrator.Name,
					"admin_user_id": administrator.Id,
					"serial_number": dividend.BillNo,
				}
				transType := consts.TransTypeAdjustmentDividendPlus
				if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, dividend.Money, extraMap); err != nil {
					log.Logger.Error(err.Error())
					_ = session.Rollback()
					continue
				}
				uMap["before_money"] = accountInfo.Available
				uMap["money"] = dividend.Money
				uMap["after_money"] = accountInfo.Available + dividend.Money
				successDo = true
			} else { //场馆钱包
				//如果外接请求成功successDo = true
				req.SetTimeout(50 * time.Second)
				req.Debug = true
				//先查余额
				headerB := req.Header{
					"Accept": "application/json",
				}
				paramB := req.Param{
					"game_code": dividend.Venue,
				}
				baseBalanceUrl := consts.InternalGameServUrl
				BalanceUrl := baseBalanceUrl + "/game/v1/internal/balance?user_id=" + strconv.Itoa(int(dividend.UserId))
				rB, err := req.Post(BalanceUrl, headerB, paramB)
				if err != nil {
					log.Logger.Error(err.Error())
					response.Err(c, err.Error())
					return
				}
				balanceResp := BalanceResp{}
				_ = rB.ToJSON(&balanceResp)
				if balanceResp.Errcode != 0 { //查询余额失败
					continue
				}
				//再转账
				header := req.Header{
					"Accept": "application/json",
				}
				param := req.Param{
					"out_code": "CENTERWALLET",
					"in_code":  dividend.Venue,
					"all_in":   2,
					"money":    dividend.Money,
				}
				baseTransferUrl := consts.InternalGameServUrl
				TransferUrl := baseTransferUrl + "/game/v1/internal/transfer?user_id=" + strconv.Itoa(int(dividend.UserId))
				r, err := req.Post(TransferUrl, header, param)
				if err != nil {
					log.Logger.Error(err.Error())
					continue
				}
				transferResp := TransferResp{}
				_ = r.ToJSON(&transferResp)
				if transferResp.Errcode != 0 {
					continue
				} else {
					uMap["before_money"] = balanceResp.Data.Balance
					uMap["money"] = dividend.Money
					beforeMoney := balanceResp.Data.Balance
					uMap["after_money"] = beforeMoney + dividend.Money
					successDo = true
				}
			}
			if successDo {
				uMap["reviewer_remark"] = postedData["reviewer_remark"]
				uMap["state"] = 2
				uMap["reviewer"] = GetLoginAdmin(c).Name
				uMap["updated"] = time.Now().UnixMicro()
				if _, err := session.Table("user_dividends").Where("id=?", v).Update(uMap); err != nil {
					log.Logger.Error(err.Error())
					_ = session.Rollback()
					continue
				}
			}
			_ = session.Commit()
			//覆盖用户钱包的数据
			if accountInfo.Id > 0 {
				_ = accountInfo.ResetCacheData(redis)
			}
			redis.Del(rKey)
		}
		response.Ok(c)
	},
	BatchRefuse: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		idSilce := strings.Split(idStr, ",")
		platform := request.GetPlatform(c)
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		engine := common.Mysql(platform)
		defer engine.Close()
		administrator := GetLoginAdmin(c)
		for _, v := range idSilce {
			dividend := &models.UserDividend{}
			vid, _ := strconv.Atoi(v)
			b, err := models.UserDividends.FindById(platform, vid, dividend)
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			//记录不存在或已经被处理过
			if !b || dividend.State != 1 {
				continue
			}
			rKey := dividend.BillNo + "_" + idStr + "_" + dividend.Username
			num, err := redis.Incr(rKey).Result()
			if err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			//这条记录正在被处理的过滤
			if num > 1 {
				continue
			}
			uMap := map[string]interface{}{
				"reviewer_remark": postedData["reviewer_remark"],
				"state":           3,
				"reviewer":        administrator.Name,
				"updated":         tools.NowMicro(),
			}
			if dividend.MoneyType == 1 { //中心钱包
				accountInfo := &models.Account{}
				if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": dividend.UserId})); !exists || err != nil {
					if err != nil {
						log.Logger.Error(err.Error())
					}
					continue
				}
				uMap["before_money"] = accountInfo.Available
				uMap["money"] = dividend.Money
				uMap["after_money"] = accountInfo.Available
			} else { //场馆钱包
				//查询余额
				req.SetTimeout(50 * time.Second)
				req.Debug = true
				header := req.Header{
					"Accept": "application/json",
				}
				param := req.Param{
					"game_code": dividend.Venue,
				}
				baseBalanceUrl := consts.InternalGameServUrl
				BalanceUrl := baseBalanceUrl + "/game/v1/internal/balance?user_id=" + strconv.Itoa(int(dividend.UserId))
				r, err := req.Post(BalanceUrl, header, param)
				if err != nil {
					log.Logger.Error(err.Error())
					continue
				}
				balanceResp := BalanceResp{}
				_ = r.ToJSON(&balanceResp)
				if balanceResp.Errcode != 0 { //查询余额失败
					continue
				} else {
					uMap["before_money"] = balanceResp.Data.Balance
					uMap["money"] = dividend.Money
					beforeMoney := balanceResp.Data.Balance
					uMap["after_money"] = beforeMoney + dividend.Money
				}
			}
			if _, err := engine.Table("user_dividends").Where("id=?", v).Update(uMap); err != nil {
				log.Logger.Error(err.Error())
				continue
			}
			redis.Del(rKey)
		}
		response.Ok(c)
	},
}
