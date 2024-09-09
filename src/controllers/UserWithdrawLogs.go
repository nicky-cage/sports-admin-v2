package controllers

import (
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// UserWithdrawLogs 提款管理-日志记录
var UserWithdrawLogs = struct {
	*ActionList
}{
	ActionList: &ActionList{
		Model:    models.FinanceLogs,
		ViewFile: "user_withdraws/logs.html",
		OrderBy: func(*gin.Context) string {
			return "id DESC"
		},
		Rows: func() interface{} {
			return &[]models.FinanceLog{}
		},
		ProcessRow: func(c *gin.Context, rs interface{}) {
			list := rs.(*[]models.FinanceLog)
			for i, v := range *list {
				if v.Consuming != "" {
					HourStr := ""
					MinuteStr := ""
					SecondStr := ""
					ConsumingInt, _ := strconv.Atoi(v.Consuming)
					Hour := ConsumingInt / 3600
					Minute := (ConsumingInt - Hour*3600) / 60
					Second := ConsumingInt - Hour*3600 - Minute*60
					if Hour < 10 {
						HourStr = "0" + strconv.Itoa(Hour)
					} else {
						HourStr = strconv.Itoa(Hour)
					}
					if Minute < 10 {
						MinuteStr = "0" + strconv.Itoa(Minute)
					} else {
						MinuteStr = strconv.Itoa(Minute)
					}
					if Second < 10 {
						SecondStr = "0" + strconv.Itoa(Second)
					} else {
						SecondStr = strconv.Itoa(Second)
					}
					(*list)[i].Consuming = HourStr + ":" + MinuteStr + ":" + SecondStr
				}
			}
		},
		GetQueryCond: func(c *gin.Context) builder.Cond {
			billNo, _ := c.GetQuery("bill_no")
			cond := builder.NewCond()
			cond = cond.And(builder.Eq{"type": 1})
			cond = cond.And(builder.Eq{"bill_no": billNo})
			return cond
		},
	},
}
