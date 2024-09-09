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
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var Recommends = struct {
	*ActionCreate
	*ActionDelete
	*ActionUpdate
	*ActionList
	*ActionSave
	//Update  func(c *gin.Context)
	Created func(c *gin.Context)
	Saved   func(c *gin.Context)
}{
	ActionList: &ActionList{
		Model:    models.Recommends,
		ViewFile: "recommend/list.html",
		QueryCond: map[string]interface{}{
			"league_name": "%",
			"play_team":   "%",
			"play_item":   "=",
			"sort":        "=",
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			cond := builder.NewCond()
			//比赛时间，
			//发布昵称
			//发布时间
			return cond
		},
		Rows: func() interface{} {
			return &[]models.Recommend{}
		},
	},
	ActionCreate: &ActionCreate{
		ViewFile: "recommend/created.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			//昵称
			platform := request.GetPlatform(c)
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			var res []models.RecommendNickName
			err := dbSession.Table("recommend_nicknames").Find(&res)
			if err != nil {
				log.Err(err.Error())
			}
			return pongo2.Context{"nickname": res}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.Recommends,
		ViewFile: "recommend/created.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			//昵称
			platform := request.GetPlatform(c)
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			var res []models.RecommendNickName
			err := dbSession.Table("recommend_nicknames").Find(&res)
			if err != nil {
				log.Err(err.Error())
			}
			return pongo2.Context{"nickname": res}
			//投注选项
			//投注赔率，

			//推荐该该项

		},
		Row: func() interface{} {
			return &models.Recommend{}
		},
	},
	//Update: func(c *gin.Context) {
	//	dbSession := common.Mysql(platform)
	//  defer dbSession.Close()
	//	lastDay := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//	var userArr []models.UserDailyReport
	//	//查询昨天有投注过的会员，减少升级压力
	//	_ = dbSession.Cols("user_id").Where("day=?", lastDay).GroupBy("user_id").Find(platform, &userArr)
	//	//查询所有的会员。
	//	tempArr := make([]uint64, 0)
	//	for _, v := range userArr {
	//		tempArr = append(tempArr, v.UserId)
	//	}
	//	var users []models.User
	//	_ = dbSession.Cols("id", "vip_high", "vip", "username").In("id", tempArr).Find(platform, &users)
	//
	//	start := time.Now().Format("2006-01-02")
	//	//昨天时间转化为时间戳
	//	uStart := tools.GetTimeStampByString(start + " 00:00:00")
	//	var level []models.UserLevel
	//	_ = dbSession.Find(platform, &level)
	//
	//	//统计会员投注积分  game_type=1 体育 1元2积分，
	//	for _, v := range users {
	//		integralsql := "select sum(valid_money) as money from user_daily_reports where user_id =%d and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + start + "','%s')"
	//		integralsqll := fmt.Sprintf(integralsql, v.Id, "%y-%m-%d ", "%y-%m-%d")
	//		integralRes, err := dbSession.QueryString(integralsqll)
	//		if err != nil {
	//			log.Err(err.Error())
	//		}
	//		integral, _ := strconv.ParseFloat(integralRes[0]["money"], 64)
	//		sportIntegralsql := "select sum(valid_money) as money from user_daily_reports where user_id =%d and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + start + "','%s')  and (game_type=1 or game_type= 2)"
	//		sportIntegralsqll := fmt.Sprintf(sportIntegralsql, v.Id, "%y-%m-%d", "%y-%m-%d")
	//		sportIntegralRes, err := dbSession.QueryString(sportIntegralsqll)
	//		if err != nil {
	//			log.Err(err.Error())
	//		}
	//		sportIntegral, _ := strconv.ParseFloat(sportIntegralRes[0]["money"], 64)
	//		levelNum := len(level)
	//		for k, val := range level {
	//			//查找目前存款
	//			dSql := "select sum(money)  as money from user_deposits  where user_id=%d and status=2 and confirm_at< %d "
	//			dSqll := fmt.Sprintf(dSql, v.Id, uStart)
	//			dRes, err := dbSession.QueryString(dSqll)
	//			if err != nil {
	//				log.Err(err.Error())
	//			}
	//			dMoney, _ := strconv.Atoi(dRes[0]["money"])
	//			if integral+sportIntegral < val.UpgradeDeposit {
	//
	//				if v.Vip+1 < int32(val.Id) && v.Vip+1 > v.VipHigh {
	//
	//					num := int(int32(val.Digit) - v.Vip)
	//
	//					for i := 1; i <= num; i++ {
	//						Grand(v.Id, level[v.Vip+int32(i)].UpgradeBonus, dMoney, (integral + sportIntegral))
	//					}
	//				}
	//				//程序应该只执行一次。
	//				goto Loop
	//			} else {
	//				if levelNum == k+1 && v.VipHigh < int32(k+1) && integral+sportIntegral >= val.UpgradeDeposit {
	//
	//					Grand(v.Id, level[k].UpgradeBonus, dMoney, (integral + sportIntegral))
	//
	//				}
	//			}
	//		}
	//	Loop:
	//	}
	//},
	ActionSave: &ActionSave{
		Model: models.Recommends,
	},
	ActionDelete: &ActionDelete{
		Model: models.Recommends,
	},
	Created: func(c *gin.Context) {
		response.Render(c, "recommend/nickname_create.html")
		//dbSession := common.Mysql(platform)
		// defer dbSession.Close()
		////查询所有的会员。
		//var users []models.User
		//_ = dbSession.Cols("id", "vip_high", "vip").Find(platform, &users)
		//
		////上月时间,查询日报表，
		//year, months, _ := time.Now().Date()
		//thisMonth := time.Date(year, months, 1, 0, 0, 0, 0, time.Local)
		//start := thisMonth.AddDate(0, -1, 0).Format("2006-01-02")
		////定时任务时1号，使用当前日期，如变动会不准
		//end := time.Now().Format("2006-01-02")
		//uStart := thisMonth.AddDate(0, -1, 0).Unix()
		//uEnd, _ := time.Parse("2006-01-02", end)
		//
		//var level []models.UserLevel
		//var vipLogs models.UserVipLog
		////	sumDay := new(models.UserDailyReport)
		//
		//_ = dbSession.Find(platform, &level)
		//
		////统计会员投注积分  game_type=1 体育 1元2积分，
		//for _, v := range users {
		//	integralsql := "select sum(valid_money) as money from user_daily_reports where user_id =%d and  DATE_FORMAT(day,'%s')> DATE_FORMAT('" + start + "','%s') and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + end + "','%s') "
		//	integralsqll := fmt.Sprintf(integralsql, v.Id, "%y-%m-%d %h:%i%s", "%y-%m-%d", "%y-%m-%d %h:%i%s", "%y-%m-%d")
		//	integralRes, err := dbSession.QueryString(integralsqll)
		//	if err != nil {
		//		log.Err(err.Error())
		//	}
		//	integral, _ := strconv.ParseFloat(integralRes[0]["money"], 64)
		//	sportIntegralsql := "select sum(valid_money) as money from user_daily_reports where user_id =%d and  DATE_FORMAT(day,'%s')> DATE_FORMAT('" + start + "','%s') and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + end + "','%s') and game_type=1 "
		//	sportIntegralsqll := fmt.Sprintf(sportIntegralsql, v.Id, "%y-%m-%d %h:%i%s", "%y-%m-%d", "%y-%m-%d %h:%i%s", "%y-%m-%d")
		//	sportIntegralRes, err := dbSession.QueryString(sportIntegralsqll)
		//	if err != nil {
		//		log.Err(err.Error())
		//	}
		//	sportIntegral, _ := strconv.ParseFloat(sportIntegralRes[0]["money"], 64)
		//	//integral, _ := dbSession.Table("user_daily_reports").Where("day>? and day<? and user_id=?", start, end, v.Id).Sum(sumDay, "valid_money")
		//	//sportIntegral, _ := dbSession.Table("user_daily_reports").Where("day>? and day<? and game_type=? and user_id=?", start, end, 1, v.Id).Sum(sumDay, "valid_money")
		//	//vip记录所需总积分，
		//	logSql := "select sum(valid_money) as money from user_daily_reports where user_id =%d and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + end + "','%s')"
		//	logSqll := fmt.Sprintf(logSql, v.Id, "%y-%m-%d ", "%y-%m-%d")
		//	logRes, err := dbSession.QueryString(logSqll)
		//	if err != nil {
		//		log.Err(err.Error())
		//	}
		//	integralF, _ := strconv.ParseFloat(logRes[0]["money"], 64)
		//
		//	sporlogsql := "select sum(valid_money) as money from user_daily_reports where user_id =%d and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + end + "','%s')  and game_type=1 or game_type= 2 "
		//	sporlogsqll := fmt.Sprintf(sporlogsql, v.Id, "%y-%m-%d", "%y-%m-%d")
		//	sportloglRes, err := dbSession.QueryString(sporlogsqll)
		//	if err != nil {
		//		log.Err(err.Error())
		//	}
		//	integralS, _ := strconv.ParseFloat(sportloglRes[0]["money"], 64)
		//
		//	//判断是否当月晋级
		//	_, verr := dbSession.Table("user_vip_logs").Where("user_id=? and created>? and created<?", v.Id, uStart, uEnd.Unix()).Get(&vipLogs)
		//	if verr != nil {
		//		log.Logger.Infof("user_vip_logs error")
		//	}
		//	if vipLogs.AdjustType == 1 {
		//		//上月有晋级，直接发放俸禄。
		//		for _, val := range level {
		//			if v.Vip == int32(val.Id) {
		//				Gran(v.Id, val.MonthBonus)
		//			}
		//		}
		//	} else {
		//		//判断是否够积分保级
		//		for k, val := range level {
		//			if v.Vip == int32(val.Id) {
		//				if val.HoldStream > (integral + float64(sportIntegral)) {
		//					//降级  2才是vip1
		//					if v.Vip >= 2 {
		//						sql := "update users set vip=? where id=?"
		//						_, err := dbSession.Exec(sql, v.Vip-1, v.Id)
		//						if err != nil {
		//							log.Err(err.Error())
		//						}
		//						//vip 降级记录
		//						dSql := "select sum(money)  as money from user_deposits  where user_id=%d and status=2 and confirm_at< %d "
		//						dSqll := fmt.Sprintf(dSql, v.Id, time.Now().Unix())
		//						dRes, err := dbSession.QueryString(dSqll)
		//						if err != nil {
		//							log.Err(err.Error())
		//						}
		//
		//						changeSql := "insert into user_vip_logs('username','user_id','valid_bet','deposits','before_vip','" +
		//							"after_vip','adjust_type','admin','created') values(?,?,?,?,?,?,?,?,?)"
		//						_, err = dbSession.Exec(changeSql, v.Username, v.Id, integralF+integralS, dRes[0]["money"], v.Vip, v.Vip-1, 2, "system", time.Now().Unix())
		//						if err != nil {
		//							log.Err(err.Error())
		//						}
		//
		//						//积分不够保级时 ，当vip>1才有下级
		//						if level[k-1].HoldStream <= (integral + float64(sportIntegral)) {
		//							//大于下级的投注积分，发放下级的俸禄
		//							Gran(v.Id, level[k-1].MonthBonus)
		//						}
		//					}
		//
		//				} else {
		//					//够保级，发放
		//					Gran(v.Id, val.MonthBonus)
		//
		//				}
		//			}
		//		}
		//	}
		//}
	},
	Saved: func(c *gin.Context) {
		postData := request.GetPostedData(c)
		delete(postData, "file")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		_, err := dbSession.Table("recommend_nicknames").Insert(postData)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "添加失败")
			return
		}
		response.Ok(c)
	},
}

func Grand(id uint64, money float64, dMoney int, iMoney float64, platform string) {
	redis := common.Redis(platform)
	defer common.RedisRestore(platform, redis)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	accountInfo := &models.Account{}

	if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": id})); !exists || err != nil {
		if err != nil {
			log.Err(err.Error())
		}
		return
	}
	userInfo := &models.User{}
	if exists, err := models.Users.FindById(platform, int(id), userInfo); !exists || err != nil {
		if err != nil {
			log.Err(err.Error())
		}
		return
	}

	transType := consts.TransTypeAdjustmentDividendPlus

	//事务操作
	session := common.Mysql(platform)
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()
		return
	}

	transAction := &models.Transaction{}

	extraMap := map[string]interface{}{
		"proxy_ip":      "",
		"ip":            "",
		"description":   "发放VIP" + strconv.Itoa(int(userInfo.Vip)) + "晋级礼金",
		"administrator": "system",
		"admin_user_id": "system",
		"serial_number": tools.GetBillNo("v", 5),
	}

	nows := time.Now().Unix()
	billNo := tools.GetBillNo("h", 5)

	//红利记录
	diSql := "insert into user_dividends(bill_no,username,user_id,top_name,top_id,type,is_automatic,money_type,flow_limit,money,created,updated,reviewer,reviewer_remark,state,vip) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := dbSession.Exec(diSql, billNo, userInfo.Username, id, userInfo.TopName, userInfo.TopId, 3, 2, 1, 1, money, nows, nows, "system", "发放升级礼金", 2, userInfo.Vip+1)
	if err != nil {
		log.Err(err.Error())
		_ = session.Rollback()
	}

	//记录用户vip升级
	changeSql := "insert into user_vip_logs(username,user_id,valid_bet,deposits,before_vip," +
		"after_vip,adjust_type,admin,created) values(?,?,?,?,?,?,?,?,?)"

	_, err = dbSession.Exec(changeSql, userInfo.Username, userInfo.Id, iMoney, dMoney, userInfo.Vip, userInfo.Vip+1, 1, "system", nows)
	if err != nil {
		log.Err(err.Error())
		_ = session.Rollback()
	}

	sql := "update users set vip=vip+1,vip_high=vip_high+1 where id=?"
	_, err = dbSession.Exec(sql, userInfo.Id)
	if err != nil {
		log.Err(err.Error())
		_ = session.Rollback()
	}

	// 稽核记录
	asql := "insert into audits(user_id,type,status,dividend_money,multiple,need_flow) values(?,?,?,?,?,?)"
	_, err = dbSession.Exec(asql, userInfo.Id, 1, 1, money, 1, money)
	if err != nil {
		log.Err(err.Error())
	}
	if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, money, extraMap); err != nil {
		log.Logger.Error(err.Error())
		_ = session.Rollback()

		return
	}

	_ = session.Commit()
	//覆盖用户钱包的数据
	if accountInfo.Id > 0 {
		_ = accountInfo.ResetCacheData(redis)
	}

}
