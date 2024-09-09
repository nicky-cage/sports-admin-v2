package agent_withdraw_hrs

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ActionExport = &base_controller.ActionExport{
	Columns: []base_controller.ExportHeader{
		{Name: "序号", Field: "id"},
		{Name: "订单编号", Field: "bill_no"},
		{Name: "会员账号", Field: "username"},
		{Name: "会员等级", Field: "vip"},
		{Name: "订单金额", Field: "money"},
		{Name: "银行卡信息", Field: "bank_realname"},
		{Name: "申请时间", Field: "created"},
		{Name: "是否代付", Field: "business_type"},
		{Name: "风控审核时间", Field: "risk_process_at"},
		{Name: "风控审核人", Field: "risk_admin"},
		{Name: "订单状态", Field: "status"},
		{Name: "完成时间", Field: "updated"},
		{Name: "出款人", Field: "finance_admin"},
		{Name: "商户名称", Field: "status"},
		{Name: "出款卡号", Field: "card_number"},
	},
	ProcessRow: func(m *map[string]interface{}, c *gin.Context) {
		setTime := func(field string) {
			theTime := int((*m)[field].(float64))
			(*m)[field] = base_controller.FieldToDateTime(fmt.Sprintf("%d", theTime))
		}
		setTime("updated")
		setTime("created")
		setTime("risk_process_at")

		val, err := strconv.Atoi(fmt.Sprintf("%v", (*m)["user_vip"]))
		realVal := val - 1
		platform := request.GetPlatform(c)
		userLevel := caches.UserLevels.Get(platform, realVal)
		if err == nil {
			(*m)["user_vip"] = userLevel.Name
		}
		tval := fmt.Sprintf("%v/%v/%v/%v", (*m)["bank_realname"], (*m)["bank_name"], (*m)["bank_card"], (*m)["bank_address"])
		(*m)["bank_realname"] = tval
		nval := (*m)["status"].(float64)
		if nval == 1 {
			(*m)["status"] = "处理中"
		} else if nval == 2 {
			(*m)["status"] = "成功"
		} else {
			(*m)["status"] = "失败"
		}
	},
}
