package index

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

var Right = func(c *gin.Context) {
	admin := base_controller.GetLoginAdmin(c)
	row := models.Admin{}
	platform := request.GetPlatform(c)
	exists, err := models.Admins.FindById(platform, int(admin.Id), &row)
	if err != nil || !exists {
		response.ErrorHTML(c, "管理员信息查找失败")
		return
	}
	role := caches.AdminRoles.Get(platform, int(admin.RoleId))
	if role == nil {
		response.ErrorHTML(c, "角色信息查找失败")
		return
	}

	lastRegCount, lastDepositCount, lastDepositTotal, lastWithdrawTotal := func() (int, int, float64, float64) {
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		todayBegin := tools.GetTodayBegin()
		lastBegin := todayBegin - 86400
		lastEnd := todayBegin - 1
		regCount, depositCount := func() (int, int) {
			sql := fmt.Sprintf("(SELECT COUNT(*) AS total, 'u' AS type  FROM users WHERE created >= %d AND created <= %d)"+
				" UNION ALL "+
				"(SELECT COUNT(*) AS total, 'd' AS type FROM user_deposits "+
				"WHERE status = 2 AND created >= %d AND created <= %d AND user_id IN "+
				"(SELECT id FROM users WHERE created >= %d AND created <= %d))",
				lastBegin, lastEnd, lastBegin, lastEnd, lastBegin, lastEnd)
			rows, err := dbSession.QueryString(sql)
			if err == nil || len(rows) > 0 {
				rCount := 0
				dCount := 0
				for _, r := range rows {
					if r["type"] == "u" {
						rCount, _ = strconv.Atoi(r["total"])
					} else if r["type"] == "d" {
						dCount, _ = strconv.Atoi(r["total"])
					}
				}
				return rCount, dCount
			}
			return 0, 0
		}()
		depositTotal, withdrawTotal := func() (float64, float64) {
			sql := "(SELECT SUM(arrive_money) AS total, 'd' AS type FROM user_deposits WHERE status = 2 AND created >= ? AND created <= ?)" +
				" UNION ALL " +
				"(SELECT SUM(actual_money) AS total, 'w' AS type FROM user_withdraws WHERE status = 2 AND created >= ? AND created <= ?)"
			rows, err := dbSession.QueryString(sql, lastBegin, lastEnd, lastBegin, lastEnd)
			if err == nil || len(rows) > 0 {
				dTotal := 0.0
				wTotal := 0.0
				for _, r := range rows {
					if r["type"] == "w" {
						wTotal, _ = strconv.ParseFloat(r["total"], 64)
					} else if r["type"] == "d" {
						dTotal, _ = strconv.ParseFloat(r["total"], 64)
					}
				}
				return dTotal, wTotal
			}
			return 0.0, 0.0
		}()
		return regCount, depositCount, depositTotal, withdrawTotal
	}()
	transRate := func() string {
		if lastRegCount == 0 {
			return "0.00"
		}
		return fmt.Sprintf("%.2f", float64(lastDepositCount)/float64(lastRegCount)*100.0)
	}()
	dwRate := func() string {
		if lastDepositTotal == 0.0 {
			return "0.00"
		}
		return fmt.Sprintf("%.2f", (lastWithdrawTotal/lastDepositTotal)*100.0)
	}()
	response.Render(c, "index/right.html", response.ViewData{
		"f":      models.FinanceMessages.Statistics(platform),
		"userID": admin.Id,
		"admin":  admin,
		"row":    row,
		"role": map[string]interface{}{
			"id":   role.Id,
			"name": role.Name,
		},
		"currentIP":      c.ClientIP(),
		"deposit_total":  lastDepositTotal,
		"withdraw_total": lastWithdrawTotal,
		"reg_count":      lastRegCount,
		"deposit_count":  lastDepositCount,
		"trans_rate":     transRate,
		"dw_rate":        dwRate,
	})
}
