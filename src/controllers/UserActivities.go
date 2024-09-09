package controllers

import (
	"encoding/json"
	"fmt"
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

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"xorm.io/builder"
)

// UserActivities 用户活动记录
var UserActivities = struct {
	Index  func(*gin.Context)
	Cancel func(*gin.Context)
	*ActionSave
	*ActionCreate
	Agree func(*gin.Context)
}{
	Index: func(c *gin.Context) { //默认首页
		limit, offset := request.GetOffsets(c)
		var other string
		var part string
		if offset == 0 {
			other = " order by created DESC limit 15"
		} else {
			temp := "order by created DESC  limit %d,%d "
			other = fmt.Sprintf(temp, limit, offset)
		}

		username := c.Query("username")
		if username != "" {
			part = part + " and username='" + username + "'"
		}
		topName := c.Query("top_name")
		if topName != "" {
			part = part + " and top_name='" + topName + "'"
		}
		status := c.Query("status")
		if status != "" {
			part = part + " and status='" + status + "'"
		}
		created := c.Query("created")
		if created != "" {
			areas := strings.Split(created, " - ")
			start := tools.GetMicroTimeStampByString(areas[0] + " 00:00:00")
			end := tools.GetMicroTimeStampByString(areas[1] + " 23:59:59")
			part = part + " and created>=" + strconv.Itoa(int(start)) + " and created<=" + strconv.Itoa(int(end)) + " "
		}
		Title := c.Query("activity_title")
		if Title != "" {
			part = part + " and activity_title='" + Title + "'"
		}
		money := c.Query("money")
		if money != "" {
			areas := strings.Split(money, "~")
			part = part + " and money>=" + areas[0] + " and money<=" + areas[1] + " "
		}

		platform := request.GetPlatform(c)
		myClient := common.Mysql(platform)
		defer myClient.Close()
		aSql := " select * from user_activities where applicant = 2 "
		cSql := " select count(*) as total ,sum(money) as money from user_activities where applicant = 2" + part

		aRes, err := myClient.QueryString(aSql + part + other)
		if err != nil {
			log.Err(err.Error())
		}
		total, err := myClient.QueryString(cSql)
		if err != nil {
			log.Err(err.Error())
		}
		for _, v := range aRes {
			var user models.User
			_, err := myClient.Table("users").Cols("vip", "is_agent", "top_name").Where("username=?", v["username"]).Get(&user)
			if err != nil {
				log.Err(err.Error())
			}
			v["vip"] = strconv.Itoa(int(user.Vip - 1))
			v["is_agent"] = strconv.Itoa(int(user.IsAgent))
			v["top_name"] = user.TopName
			sql := "select sum(valid_money) as money from user_daily_reports where user_id=%s and created>%s"
			sqll := fmt.Sprintf(sql, v["user_id"], v["created"])
			dRes, err := myClient.QueryString(sqll)
			if err != nil {
				log.Err(err.Error())
			}
			MultipleFinish, _ := strconv.ParseFloat(dRes[0]["money"], 64)
			MultipleRequirement, _ := strconv.ParseFloat(v["multiple_requirement"], 64)
			v["multiple_finish"] = "0.00"
			if MultipleRequirement <= MultipleFinish {
				if moneyFinished, err := strconv.ParseFloat(v["multiple_requirement"], 64); err == nil {
					v["multiple_finish"] = fmt.Sprintf("%.2f", moneyFinished)
				} else {
					v["multiple_finish"] = v["multiple_requirement"]
				}
			} else {
				if moneyFinished, err := strconv.ParseFloat(dRes[0]["money"], 64); err == nil {
					v["multiple_finish"] = fmt.Sprintf("%.2f", moneyFinished)
				} else {
					v["multiple_finish"] = dRes[0]["money"]
				}
			}
		}

		response.Render(c, "user_changes/_user_activities.html", ViewData{
			"rows":           aRes,
			"total":          total[0]["total"],
			"activity_money": total[0]["money"],
		})
	},
	Cancel: func(c *gin.Context) { //取消活动
		ids := c.Query("id")
		sql := "update user_activities set risk_admin=? ,updated=?,status=? where id=" + ids
		admin := GetLoginAdmin(c)
		platform := request.GetPlatform(c)
		myClient := common.Mysql(platform)
		defer myClient.Close()
		_, err := myClient.Exec(sql, admin.Name, time.Now().Unix(), 3)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "程序错误")
			return
		}
		response.Ok(c)
	},
	Agree: func(c *gin.Context) { //同意派发
		platform := request.GetPlatform(c)
		myClient := common.Mysql(platform)
		defer myClient.Close()
		ids := c.Query("id") //活动表id
		jsql := "select a.username,a.multiple,a.money ,b.id,a.activity_title from user_activities a join users b on a.username=b.username where a.id=" + ids
		jRes, err := myClient.QueryString(jsql)
		if err != nil {
			log.Err(err.Error())
		}

		id, _ := strconv.Atoi(jRes[0]["id"])
		//moneys := c.Query("money")
		money, _ := strconv.ParseFloat(jRes[0]["money"], 64)
		admin := GetLoginAdmin(c)

		sql := "update user_activities set risk_admin=? ,updated=?,status=? where id=" + ids
		_, err = myClient.Exec(sql, admin.Name, time.Now().Unix(), 2)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "程序错误")
			return
		}
		redis := common.Redis(platform)
		defer common.RedisRestore(platform, redis)
		accountInfo := &models.Account{}
		if exists, err := models.Accounts.Find(platform, accountInfo, builder.NewCond().And(builder.Eq{"user_id": id})); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			return
		}

		userInfo := &models.User{}
		if exists, err := models.Users.FindById(platform, id, userInfo); !exists || err != nil {
			if err != nil {
				log.Logger.Error(err.Error())
			}
			return
		}

		transType := consts.TransTypeAdjustmentActivityPlus

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
			"description":   jRes[0]["activity_title"],
			"administrator": admin.Name,
			"admin_user_id": admin.Id,
			"serial_number": tools.GetBillNo("v", 5),
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

		response.Ok(c)
	},
	ActionCreate: &ActionCreate{
		ViewFile: "user_changes/activities_create.html",
		ExtendData: func(c *gin.Context) response.ViewData {
			platform := request.GetPlatform(c)
			activities := []struct {
				Id           int    `json:"id" xorm:"id"`  // 编号
				Title        string `json:"title"`         // 活动标题
				ActivityType int    `json:"activity_type"` // 活动类型
			}{}
			activityTitle := ""
			myClient := common.Mysql(platform)
			defer myClient.Close()
			if err := myClient.SQL("SELECT id, title, activity_type FROM activities WHERE state = 2 ORDER BY id DESC").Find(&activities); err != nil {
				log.Logger.Error("获取活动信息列表失败", err)
			} else if len(activities) > 0 {
				activityTitle = activities[0].Title
			}
			activityListString := ""
			if bytes, err := json.Marshal(activities); err == nil {
				activityListString = string(bytes)
			}
			return response.ViewData{
				"activities":       activities,
				"activityListJSON": activityListString,
				"activityTitle":    activityTitle,
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.UserActivities,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			platform := request.GetPlatform(c)
			myClient := common.Mysql(platform)
			defer myClient.Close()

			//风控标签
			username := (*m)["username"].(string)
			str := ""
			userID := ""
			sql := "select id, label from users where username = '" + username + "'"
			if res, err := myClient.QueryString(sql); err != nil {
				log.Err(err.Error())
			} else if len(res) == 0 {
				return errors.New("用户不存在")
			} else {
				userID = res[0]["id"]
				if res[0]["label"] != "" {
					str = res[0]["label"] + ";3|9 "
				} else {
					str = "3|9 "
				}
			}

			uSQL := "update users set label='" + str + "' where username='" + username + "'"
			if _, err := myClient.Exec(uSQL); err != nil {
				log.Err(err.Error())
			}
			billNo := tools.GetBillNo("a", 5)
			admin := GetLoginAdmin(c)
			(*m)["bill_no"] = billNo
			(*m)["applicant"] = 2
			(*m)["proposer"] = admin.Name
			(*m)["status"] = 1
			(*m)["created"] = time.Now().UnixMicro()
			(*m)["user_id"] = userID
			return nil
		},
	},
}
