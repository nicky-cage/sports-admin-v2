package users

import (
	"fmt"
	"sports-admin/caches"
	"sports-admin/controllers/base_controller"
	"sports-common/request"
	"sports-common/tools"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var ActionExport = &base_controller.ActionExport{
	Columns: []base_controller.ExportHeader{
		{"用户编号", "id"},
		{"会员账号", "username"},
		{"中心钱包", "available"}, // 中心钱包余额
		{"会员等级", "vip"},
		{"真实姓名", "realname"},
		{"是否代理", "is_agent"},
		{"上级代理", "top_name"},
		{"电话", "phone"},
		{"电子邮件", "email"},
		{"QQ", "qq"},
		{"微信", "we_chat"},
		{"注册IP/地区", "register_ip"},
		{"注册时间", "created"},
		{"最后登录时间", "last_login_at"},
		{"最后登录IP/地区", "last_login_ip"},
		{"状态", "status"},
	},
	ProcessRow: func(m *map[string]interface{}, c *gin.Context) {
		(*m)["is_agent"] = base_controller.FieldToYesNo(fmt.Sprintf("%v", (*m)["is_agent"]))
		field := "last_login_at"
		theTime := int((*m)[field].(float64))
		(*m)[field] = base_controller.FieldToDateTime(fmt.Sprintf("%d", theTime))
		field = "created"
		theTime = int((*m)[field].(float64))
		(*m)[field] = base_controller.FieldToDateTime(fmt.Sprintf("%d", theTime))
		(*m)["status"] = base_controller.FieldToStatus(fmt.Sprintf("%v", (*m)["status"]))
		val, err := strconv.Atoi(fmt.Sprintf("%v", (*m)["vip"]))
		(*m)["available"] = fmt.Sprintf("%.2f", (*m)["available"].(float64))
		realVal := val - 1
		platform := request.GetPlatform(c)
		userLevel := caches.UserLevels.Get(platform, realVal)
		if err == nil {
			(*m)["vip"] = userLevel.Name
		}
		getIP := func(field string) {
			val := fmt.Sprintf("%v", (*m)[field])
			(*m)[field] = val + " | " + func() string {
				areas := strings.Trim(tools.GetAreaByIp(val), "[]")
				areaArr := strings.Split(areas, ",")
				area := ""
				if len(areaArr) <= 1 {
					if areaArr[0] != "" {
						return areaArr[0]
					}
					return "*未知地区*"
				}
				if areaArr[0] == "本机地址" {
					return "本机地址"
				}
				area += areaArr[0]
				if len(areaArr) <= 2 || areaArr[1] == "" {
					return area
				}
				area += "-" + areaArr[1]
				if len(areaArr) <= 3 || areaArr[2] == "" {
					return area
				}
				area += "-" + areaArr[2]
				return area
			}()
		}
		getIP("register_ip")
		getIP("last_login_ip")
	},
}
