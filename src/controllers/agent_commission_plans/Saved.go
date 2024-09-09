package agent_commission_plans

import (
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissionPlans) Saved(c *gin.Context) {
	postData := request.GetPostedData(c)
	admin := base_controller.GetLoginAdmin(c)
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	mode := postData["mode"]
	if mode == "1" {
		sql := "UPDATE agent_commission_plans SET level_id = ?, level = ?, active_num = ?, negative_profit = ?, rate = ?, agent_commission = ?, creat_admin = ? where id = ?"
		insertSQL := "INSERT INTO agent_commission_plans(level_id, level, active_num, negative_profit, rate, agent_commission, creat_admin) VALUES (?, ?, ?, ?, ?, ?, ?)"
		for _, v := range postData {
			res, _ := v.(map[string]interface{})
			var Num int
			var po int
			if len(res) == 0 {
				continue
			}
			val, exists := res["active_num"]
			if !exists || val == "" {
				response.Err(c, "活跃会员错误: ")
				return
			}
			activeNum, err := strconv.Atoi(fmt.Sprintf("%v", val))
			if err != nil {
				log.Err(err.Error())
				response.Err(c, "活跃会员错误: "+err.Error())
				return
			}
			Num = activeNum
			val, exists = res["negative_profit"]
			if !exists || val == "" {
				response.Err(c, "活跃会员错误: ")
				return
			}
			pro, err := strconv.Atoi(fmt.Sprintf("%v", val))
			if err != nil {
				log.Err(err.Error())
				response.Err(c, "佣金比例错误: "+err.Error())
				return
			}
			po = pro
			rates, _ := strconv.ParseFloat(res["rate"].(string), 64)
			rate := tools.ToFixed(rates, 0)
			if res["id"] == "0" { // 如果添加
				_, err := dbSession.Exec(insertSQL, res["level_id"].(string), res["level"].(string), Num, po, rate/100, res["agent_commission"].(string), admin.Name)
				if err != nil {
					log.Err(err.Error())
					response.Err(c, err.Error())
					return
				}
			} else { // 否则修改
				_, err := dbSession.Exec(sql, res["level_id"].(string), res["level"].(string), Num, po, rate/100, res["agent_commission"].(string), admin.Name, res["id"].(string))
				if err != nil {
					log.Err(err.Error())
					response.Err(c, err.Error())
					return
				}
			}
		}
	} else {
		//mode=2
		var rate float64
		sql := "UPDATE agent_commission_plans SET  rate = ?, agent_commission = ?, creat_admin = ? where id = ?"
		if postData["rate"].(string) != "" {
			pro, _ := strconv.ParseFloat(postData["rate"].(string), 64)
			rate = tools.ToFixed(pro, 0) / 100
		} else {
			response.Err(c, "请输入占成比例")
			return
		}
		dbSession.Exec(sql, rate, postData["agent_commission"], admin.Name, postData["id"])
	}

	response.Result(c, "修改成功")
}
