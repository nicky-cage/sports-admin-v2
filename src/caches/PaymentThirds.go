package caches

import (
	common "sports-common"
)

const keyPaymentThirds = "payment_thirds"

// PaymentThirds HelpCategories 配置信息
var PaymentThirds = struct {
	Load func(string)
	Get  func(string, string) string
	All  func(string) map[string]string
}{
	Load: func(platform string) {
		sql := "SELECT code, name FROM payments"
		dbSession := common.Mysql(platform)
		defer dbSession.Close()

		rArr := map[string]string{}
		rows, err := dbSession.QueryString(sql)
		if err != nil {
			panic("获取三方支付方式出错")
		}

		for _, r := range rows {
			rArr[r["code"]] = r["name"]
		}

		_ = setCache(platform, keyPaymentThirds, rArr)
	},
	Get: func(platform string, code string) string {
		rArr := map[string]string{}
		_ = getCache(platform, keyPaymentThirds, &rArr)
		if v, exists := rArr[code]; exists {
			return v
		}
		return ""
	},
	All: func(platform string) map[string]string {
		rArr := map[string]string{}
		_ = getCache(platform, keyPaymentThirds, &rArr)
		return rArr
	},
}
