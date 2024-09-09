package agent_commission_plans

import (
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (ths *AgentCommissionPlans) Add(c *gin.Context) {
	dataPost := request.GetPostedData(c)
	admin := base_controller.GetLoginAdmin(c)
	times := time.Now().Unix()
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()

	//是否存在相同的佣金名。
	mode := dataPost["mode"]
	if mode == "1" {
		sql := "insert  agent_commission_plans(level_id, level,active_num,negative_profit,rate,agent_commission,created,creat_admin,type) values(?,?,?,?,?,?,?,?,?)"

		for _, v := range dataPost {
			var rate float64
			res, _ := v.(map[string]interface{})
			var Num int
			var po int

			if len(res) == 0 {
				continue
			}
			if res["active_num"].(string) != "" && res["active_num"] != nil {
				activeNum, _ := strconv.Atoi(res["active_num"].(string))

				Num = activeNum
			}
			if res["negative_profit"].(string) != "" {
				pro, _ := strconv.Atoi(res["negative_profit"].(string))
				po = pro
			}
			if res["rate"].(string) != "" {
				pro, _ := strconv.ParseFloat(res["rate"].(string), 64)
				rate = tools.ToFixed(pro, 0) / 100
			}
			dbSession.Exec(sql, res["level_id"].(string), res["level"].(string), Num, po, rate, res["agent_commission"].(string), times, admin.Name, mode)

		}
	} else {
		var rate float64
		if dataPost["rate"] == "" {
			response.Err(c, "请输入占成比例")
			return
		}
		sql := "insert  agent_commission_plans(rate,agent_commission,created,creat_admin,type) values(?,?,?,?,?)"
		if dataPost["rate"].(string) != "" {
			pro, _ := strconv.ParseFloat(dataPost["rate"].(string), 64)
			rate = tools.ToFixed(pro, 0) / 100
		}
		dbSession.Exec(sql, rate, dataPost["agent_commission"], time.Now().Unix(), admin.Name, 2)
	}

	response.Ok(c)
}
