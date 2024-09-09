package agents

import (
	"errors"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.Users,
	SaveBefore: func(context *gin.Context, m *map[string]interface{}) error {
		phone := (*m)["phone"].(string)
		if len(phone) > 11 {
			return errors.New("手机长度不能大于11")
		}
		commission := (*m)["agent_commission"].(string)
		if commission != "" {
			arr := strings.Split(commission, "-")
			(*m)["agent_commission"] = arr[0]
		}
		return nil
	},
	SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
		status := (*m)["type"].(string)
		id := (*m)["id"].(string)
		rate := (*m)["rate"].(string)
		commission := (*m)["agent_commission"].(string)
		arr := strings.Split(commission, "-")
		temp, _ := strconv.ParseFloat(rate, 64)
		rates := temp / 100
		admin := base_controller.GetLoginAdmin(c)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "select * from agent_commission_plans where user_id=" + id
		res, _ := dbSession.QueryString(sql)
		num := len(res)
		//1 是普通模式， 2是占成模式。
		if status == "2" {
			//普通转占成， 占成转占成
			if num == 1 {
				sql := "update agent_commission_plans set rate=? where id=?"
				dbSession.Exec(sql, rates, res[0]["id"])
			} else {
				//先删除后插入
				delSql := "delete from agent_commission_plans where user_id=" + id
				dbSession.Exec(delSql)
				sql := "insert into agent_commission_plans(user_id,agent_commission,type,rate,creat_admin,created) values(?,?,?,?,?,?)"
				_, err := dbSession.Exec(sql, id, arr[0], "2", rates, admin.Name, time.Now().Unix())
				if err != nil {
					log.Err(err.Error())
				}
			}
		} else {
			sql := "delete from agent_commission_plans where user_id=" + id
			dbSession.Exec(sql)
			var list []models.AgentCommissionPlan
			dbSession.Table("agent_commission_plans").Where("agent_commission=? and user_id=0", arr[0]).Find(&list)
			newCommission := make([]models.AgentCommissionPlan, len(list))
			idInt, _ := strconv.Atoi(id)
			for k, v := range list {
				newCommission[k].UserId = uint64(idInt)
				newCommission[k].Created = tools.NowMicro()
				newCommission[k].CreatAdmin = admin.Name
				newCommission[k].AgentCommission = v.AgentCommission
				newCommission[k].Type = v.Type
				newCommission[k].Rate = v.Rate
				newCommission[k].NegativeProfit = v.NegativeProfit
				newCommission[k].Level = v.Level
				newCommission[k].LevelId = v.LevelId
				newCommission[k].ActiveNum = v.ActiveNum
			}
			_, err := dbSession.Table("agent_commission_plans").Insert(&newCommission)
			//插入新的
			if err != nil {
				log.Err(err.Error())
			}
		}
	},
}
