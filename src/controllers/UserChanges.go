package controllers

import (
	"sports-admin/caches"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"xorm.io/builder"
)

// UserChanges 记录管理
var UserChanges = struct {
	List   func(*gin.Context)
	Agree  func(*gin.Context)
	Refuse func(*gin.Context)
	*ActionSave
}{
	List: func(c *gin.Context) { //默认首页
		cond := builder.NewCond()
		Type := c.DefaultQuery("adjust_type", "")
		if len(Type) > 0 {
			cond = cond.And(builder.Eq{"adjust_type": Type})
		}
		username := c.DefaultQuery("username", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"username": username})
		}
		//status := c.DefaultQuery("status", "")
		//if len(status) > 0 {
		//	cond = cond.And(builder.Eq{"status": status})
		//}
		billNo := c.DefaultQuery("bill_no", "")
		if len(billNo) > 0 {
			cond = cond.And(builder.Eq{"bill_no": billNo})
		}
		admin := c.DefaultQuery("admin", "")
		if len(admin) > 0 {
			cond = cond.And(builder.Eq{"admin": admin})
		}
		method := c.DefaultQuery("adjust_method", "")
		if len(method) > 0 {
			cond = cond.And(builder.Eq{"adjust_method": method})
		}
		create := c.Query("updated")
		if create != "" { //对时间进行处理
			areas := strings.Split(create, " - ")
			start := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
			end := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
			cond = cond.And(builder.Gte{"updated": start}).And(builder.Lte{"updated": end})
		}
		limit, offset := request.GetOffsets(c)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		rows := []models.UserReset{}
		total, ferr := dbSession.Table("user_resets").Where(cond).Limit(limit, offset).OrderBy("updated DESC").FindAndCount(&rows)
		if ferr != nil {
			log.Err(ferr.Error())
			return
		}
		var sumReset models.UserReset
		totals, err := dbSession.Table("user_resets").Where(cond).Sum(&sumReset, "adjust_money")
		if err != nil {
			log.Err(ferr.Error())
			return
		}
		var viewFile string
		if request.IsAjax(c) {
			viewFile = "user_changes/_user_changes.html"
		} else {
			viewFile = "user_changes/user_changes.html"
		}
		SetLoginAdmin(c)
		viewData := ViewData{
			"rows":          rows,
			"total":         total,
			"venue":         caches.GameVenues.All(platform),
			"dividendTypes": consts.DividendType,
			"sum_reset":     tools.ToFixed(totals, 0),
		}
		response.Render(c, viewFile, viewData)
	},
	Agree: func(c *gin.Context) { //账户调整-同意
		ids := c.Query("id")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from user_resets where id=" + ids
		res, err := dbSession.QueryString(sql)
		if err != nil {
			response.Err(c, err.Error())
			log.Err(err.Error())
			return
		}
		response.Render(c, "user_changes/agree.html", pongo2.Context{"r": res[0]})
	},
	Refuse: func(c *gin.Context) { //账户调整-拒绝
		ids := c.Query("id")
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from user_resets where id=" + ids
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
			return
		}
		response.Render(c, "user_changes/cancel.html", pongo2.Context{"r": res[0]})
	},
	//ActionList: &ActionList{
	//	Model:    models.UserResets,
	//	ViewFile: "user_changes/user_changes.html",
	//	Rows: func() interface{} {
	//		return []models.UserReset{}
	//	},
	//	QueryCond: map[string]interface{}{
	//		"adjust_type": "=",
	//		"username":    "%",
	//		"status":      "=",
	//		"bill_no":     "=",
	//		"admin":       "%",
	//	},
	//	GetQueryCond: func(c *gin.Context) builder.Cond {
	//		cond := builder.NewCond()
	//		cond.And(builder.Eq{"status": 1})
	//		return cond
	//	},
	//},
	ActionSave: &ActionSave{
		Model: models.UserResets,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			if (*m)["adjust_method"].(string) == "上分" {
				(*m)["adjust_method"] = "1"
			} else {
				(*m)["adjust_method"] = "2"
				if (*m)["status"].(string) != "3" {
					platform := request.GetPlatform(c)
					dbSession := common.Mysql(platform)
					defer dbSession.Close()
					var userMoney models.Account
					_, err := dbSession.Table("accounts").Where("user_id=?", (*m)["user_id"].(string)).Cols("available").Get(&userMoney)
					if err != nil {
						return errors.New("账户不存在")
					}
					mo, _ := strconv.ParseFloat((*m)["adjust_money"].(string), 64)
					if userMoney.Available-mo < 0 {
						return errors.New("金额不足")
					}
				}
			}
			return nil
		},
		SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
			// 涉及账变表，
			var lo models.UserAccountSet
			id := (*m)["user_id"].(string)
			ids, _ := strconv.Atoi(id)
			administrator := GetLoginAdmin(c)
			lo.UserId = uint64(ids)
			lo.Username = (*m)["username"].(string)
			lo.Applicant = (*m)["admin_name"].(string)
			lo.Audit = administrator.Name

			cr := tools.GetMicroTimeStampByString((*m)["updated"].(string))

			lo.Created = int64(cr)
			lo.Updated = tools.NowMicro()
			if (*m)["applicant_remark"] != nil {
				lo.ApplicantRemark = (*m)["applicant_remark"].(string)
			} else {
				lo.ApplicantRemark = ""
			}

			lo.AuditRemark = (*m)["remark"].(string)
			lo.BillNo = (*m)["bill_no"].(string)

			if (*m)["adjust_method"].(string) == "1" {
				lo.Type = 1
			} else {
				lo.Type = 2
			}
			platform := request.GetPlatform(c)
			mo, err := strconv.ParseFloat((*m)["adjust_money"].(string), 64)
			if err != nil {
				DeleteUserResetBillNo(platform, id)
				response.Err(c, "程序错误")
				log.Err(err.Error())
			}
			lo.Money = mo
			s, err := strconv.Atoi((*m)["vip"].(string))
			if err != nil {
				DeleteUserResetBillNo(platform, id)
				response.Err(c, "程序错误")
				log.Err(err.Error())
			}
			lo.UserVip = int32(s)
			o, verr := strconv.Atoi((*m)["status"].(string))
			if verr != nil {
				DeleteUserResetBillNo(platform, id)
				response.Err(c, "程序错误")
				log.Err(verr.Error())
			}
			dbSession := common.Mysql(platform)
			defer dbSession.Close()

			if (*m)["status"] == "3" {
				lo.Status = uint64(o)
				_, err := dbSession.Table("user_account_sets").Insert(lo)
				if err != nil {
					DeleteUserResetBillNo(platform, id)
					log.Err(err.Error())
					response.Err(c, "程序错误")
				}
				return
			}

			redis := common.Redis(platform)
			defer common.RedisRestore(platform, redis)
			rKey := (*m)["bill_no"].(string) + "_" + id + "_" + (*m)["username"].(string)
			num, err := redis.Incr(rKey).Result()
			if err != nil {
				DeleteUserResetBillNo(platform, id)
				log.Err(err.Error())
				response.Err(c, "系统繁忙")
				return
			}
			if num > 1 {
				DeleteUserResetBillNo(platform, id)
				response.Err(c, "请稍等片刻再试")
				return
			}
			defer redis.Del(rKey)

			money, _ := strconv.ParseFloat((*m)["adjust_money"].(string), 64)

			accountInfo := &models.Account{}
			if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": id})); !exists || err != nil {
				if err != nil {
					DeleteUserResetBillNo(platform, id)
					log.Logger.Error(err.Error())
				}
				response.Err(c, "查找用户账户信息失败")
				return
			}
			userInfo := &models.User{}
			if exists, err := models.Users.FindById(platform, ids, userInfo); !exists || err != nil {
				if err != nil {
					DeleteUserResetBillNo(platform, id)
					log.Logger.Error(err.Error())
				}
				response.Err(c, "查找用户信息失败")
				return
			}
			var transType int
			rs := (*m)["adjust_method"].(string)
			if rs == "1" {
				transType = consts.TransTypeAdjustmentPlus
			} else {
				transType = consts.TransTypeAdjustmentLess
			}

			//事务操作
			session := common.Mysql(platform)
			defer session.Close()
			if err := session.Begin(); err != nil {
				DeleteUserResetBillNo(platform, id)
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, "事务启动失败")
				return
			}

			transAction := &models.Transaction{}
			var desc string

			if (*m)["adjust_type"].(string) == "1" {
				desc = "红利补发"
			} else if (*m)["adjust_type"].(string) == "2" {
				desc = "系统调整"
			} else {
				desc = "输赢调整"
			}
			extraMap := map[string]interface{}{
				"proxy_ip":      "",
				"ip":            c.ClientIP(),
				"description":   desc,
				"administrator": administrator.Name,
				"admin_user_id": administrator.Id,
				"serial_number": (*m)["bill_no"],
			}

			if _, err := transAction.AddTransaction(platform, session, redis, userInfo, accountInfo, transType, money, extraMap); err != nil {
				DeleteUserResetBillNo(platform, id)
				log.Logger.Error(err.Error())
				_ = session.Rollback()
				response.Err(c, err.Error())
				return
			}
			lo.Status = uint64(o)

			_, eerr := session.Table("user_account_sets").Insert(lo)
			if eerr != nil {
				DeleteUserResetBillNo(platform, id)
				log.Err(eerr.Error())
				_ = session.Rollback()
				response.Err(c, "程序错误")
			}
			_ = session.Commit()
			//覆盖用户钱包的数据
			if accountInfo.Id > 0 {
				_ = accountInfo.ResetCacheData(redis)
			}

			key := consts.CtxKeyLoginUser + id
			redis.Del(key)
		},
	},
}

func DeleteUserResetBillNo(platform string, id string) {
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "delete from user_resets where id =" + id
	dbSession.Exec(sql)
}
