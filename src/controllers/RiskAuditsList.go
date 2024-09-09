package controllers

import (
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

var RiskAuditsList = struct {
	*ActionList
	*ActionSave
	Lists func(c *gin.Context)
}{
	ActionSave: &ActionSave{
		Model: models.UserWithdraws,
	},
	Lists: func(c *gin.Context) {
		admin := GetLoginAdmin(c)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := fmt.Sprintf("select * from user_withdraws where status=1 and risk_admin='%s' and process_step = 2 and type = 1 limit 50", admin.Name)
		csql := "select count(*) as total from user_withdraws where status=1 and risk_admin='" + admin.Name + "' and process_step=2 and type=1 "
		sqll := sql
		res, err := dbSession.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
			return
		}
		cRes, _ := dbSession.QueryString(csql)
		vSql := "select vip, label from users where id='%s' "
		lastSql := "select money,confirm_at from user_deposits where user_id='%s' and status=2 order by confirm_at DESC limit 1 "
		aSql := "select money,updated from user_account_sets where user_id=%s and status=2 and type= 1 order by updated DESC limit 1"

		for _, v := range res {
			vsqll := fmt.Sprintf(vSql, v["user_id"])
			vRes, err := dbSession.QueryString(vsqll)
			if err != nil {
				log.Err(err.Error())
				return
			}

			vipInt, _ := strconv.Atoi(vRes[0]["vip"])
			v["vip"] = strconv.Itoa(vipInt - 1)
			v["label"] = vRes[0]["label"]

			lastsqll := fmt.Sprintf(lastSql, v["user_id"])
			lastRes, err := dbSession.QueryString(lastsqll)
			if err != nil {
				log.Err(err.Error())
				return
			}
			aSqll := fmt.Sprintf(aSql, v["user_id"])
			aRes, err := dbSession.QueryString(aSqll)
			if err != nil {
				log.Err(err.Error())
				return
			}
			var atime int
			if aRes != nil {
				atime, _ = strconv.Atoi(aRes[0]["updated"])
			}
			var wtime int
			if lastRes != nil {
				wtime, _ = strconv.Atoi(lastRes[0]["confirm_at"])
			}

			if atime > wtime {
				v["last_money"] = aRes[0]["money"]
			} else {
				if wtime > 0 {
					v["last_money"] = lastRes[0]["money"]
				} else {
					v["last_money"] = "0"
				}
			}
			if v["bank_name"] == "其他银行" {
				sql_card := "select bank_name from user_cards where card_number= " + v["bank_card"]
				data_card, _ := dbSession.QueryString(sql_card)
				if data_card[0]["bank_name"] != "" {
					v["bank_name"] = data_card[0]["bank_name"]
				}
			}
		}

		SetLoginAdmin(c)
		Viewdata := pongo2.Context{
			"rows":      res,
			"total":     cRes[0]["total"],
			"vipLevels": caches.UserLevels.All(platform),
		}
		if request.IsAjax(c) {
			response.Render(c, "risk_audits/_hands.html", Viewdata)
		} else {
			response.Render(c, "risk_audits/hands.html", Viewdata)
		}
	},
}
