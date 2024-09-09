package users

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.Users,
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		id := (*m)["id"].(string)
		admin := base_controller.GetLoginAdmin(c)
		platform := request.GetPlatform(c)
		if id == "" {
			return nil
		}

		vip := (*m)["vip"].(string)
		name := (*m)["realname"].(string)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select vip,realname,vip_high from users where id=" + id
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return nil
		}
		//vip不同  需要记录
		if vip != res[0]["vip"] {
			p, err := strconv.Atoi(vip)
			if err != nil {
				log.Err(err.Error())
			}
			rvip, err := strconv.Atoi(res[0]["vip"])
			if err != nil {
				log.Err(err.Error())
			}
			var vipLogs = new(models.UserVipLog)
			vipLogs.BeforeVip = int32(rvip)
			vipLogs.AfterVip = int32(p)
			vipLogs.Username = (*m)["username"].(string)

			idint, err := strconv.Atoi(id)
			if err != nil {
				log.Err(err.Error())
			}
			vipLogs.UserId = idint
			//累计存款
			todayDateStr := time.Now().Format("2006-01-02")
			start := tools.GetTimeStampByString(todayDateStr + " 00:00:00")
			sumd := new(models.UserDeposit)
			d, err := dbSession.In("user_id", id).And(builder.Lt{"confirm_at": start}).And(builder.Eq{"status": 2}).Sum(sumd, "money")
			if err != nil {
				log.Err(err.Error())
			}

			//如果存入的vip 大于数据库里的vip，
			if p > rvip {
				vipLogs.AdjustType = 1
				vipHigh, _ := strconv.Atoi(res[0]["vip_high"])
				if p > vipHigh {
					(*m)["vip_high"] = vip
					vipInt, _ := strconv.Atoi(vip)
					var userLevel models.UserLevel
					_, err := dbSession.Table("user_levels").Where("id=?", vipInt).Get(&userLevel)
					if err != nil {
						log.Err(err.Error())
					}
					//发放升级礼金
					redis := common.Redis(platform)
					defer common.RedisRestore(platform, redis)
					accountInfo := &models.Account{}
					if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": idint})); !exists || err != nil {
						if err != nil {
							log.Logger.Error(err.Error())
						}

					}
					userInfo := &models.User{}
					if exists, err := models.Users.FindById(platform, idint, userInfo); !exists || err != nil {
						if err != nil {
							log.Logger.Error(err.Error())
						}

					}

					transType := consts.TransTypeAdjustmentDividendPlus

					// 事务操作
					session := common.Mysql(platform)
					defer session.Close()
					if err := session.Begin(); err != nil {
						log.Logger.Error(err.Error())
						_ = session.Rollback()

					}

					transAction := &models.Transaction{}
					extraMap := map[string]interface{}{
						"proxy_ip":      "",
						"ip":            "",
						"description":   "发放VIP" + strconv.Itoa(rvip) + "升级礼金",
						"administrator": "system",
						"admin_user_id": "system",
						"serial_number": tools.GetBillNo("v", 5),
					}
					if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, userLevel.UpgradeBonus, extraMap); err != nil {
						log.Logger.Error(err.Error())
						_ = session.Rollback()

					}
					_ = session.Commit()

					//覆盖用户钱包的数据
					if accountInfo.Id > 0 {
						_ = accountInfo.ResetCacheData(redis)
					}

					billNo := tools.GetBillNo("h", 5)
					// 红利记录
					diSql := "insert into user_dividends(bill_no,username,user_id,top_name,top_id,type,is_automatic,money_type,flow_limit,money,created,reviewer,reviewer_remark,state,vip) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
					_, err = dbSession.Exec(diSql, billNo, userInfo.Username, id, userInfo.TopName, userInfo.TopId, 3, 2, 1, 1, userLevel.UpgradeBonus, time.Now().Unix(), "system", "发放升级礼金", 2, userInfo.Vip+1)
					if err != nil {
						log.Err(err.Error())
					}
					//稽核记录
				}
			} else if p < rvip {
				vipLogs.AdjustType = 2
			} else {
				vipLogs.AdjustType = 0
			}
			vipLogs.Updated = tools.NowMicro()
			vipLogs.Admin = admin.Name
			vipLogs.Deposits = d
			// 累计有效下注
			integralsql := "select sum(valid_money) as money from user_daily_reports where user_id =%s and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + todayDateStr + "','%s') and game_code='0'"
			integralsqll := fmt.Sprintf(integralsql, id, "%y-%m-%d ", "%y-%m-%d")
			integralRes, err := dbSession.QueryString(integralsqll)
			if err != nil {
				log.Err(err.Error())
			}
			integral, _ := strconv.ParseFloat(integralRes[0]["money"], 64)
			sportIntegralsql := "select sum(valid_money) as money from user_daily_reports where user_id =%s and  DATE_FORMAT(day,'%s')< DATE_FORMAT('" + todayDateStr + "','%s')  and game_type=1 "
			sportIntegralsqll := fmt.Sprintf(sportIntegralsql, id, "%y-%m-%d", "%y-%m-%d")
			sportIntegralRes, err := dbSession.QueryString(sportIntegralsqll)
			sportIntegral, _ := strconv.ParseFloat(sportIntegralRes[0]["money"], 64)

			vipLogs.ValidBet = integral + sportIntegral
			if err != nil {
				log.Err(err.Error())
			}
			_, ierr := dbSession.Insert(vipLogs)
			if ierr != nil {
				log.Err(ierr.Error())
			}

		}
		//real 不同需要记录，并且
		if name != res[0]["realname"] {
			var realName = new(models.UserBanksNameLog)
			realName.Updated = tools.NowMicro()
			realName.UserId = id
			realName.Admin = admin.Name
			realName.BeforeName = res[0]["realname"]
			realName.AfterName = name
			res, err := dbSession.Insert(realName)
			if err != nil {
				log.Err(err.Error())
			}
			if res > 0 {
				sql := "update user_cards set real_name=? where user_id=?"
				_, eerr := dbSession.Exec(sql, name, id)
				if eerr != nil {
					log.Err(eerr.Error())
				}
			}
		}
		return nil
	},
	SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
		id := (*m)["id"].(string)
		ids, _ := strconv.Atoi(id)
		platform := request.GetPlatform(c)
		_ = models.DelUserByIdHandler(platform, ids)
	},
}
