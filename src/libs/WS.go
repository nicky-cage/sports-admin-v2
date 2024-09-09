package libs

import (
	"errors"
	"fmt"
	"sports-common/mapping"
	"sports-common/ws"
	models "sports-models"
	"strconv"
)

// WebSocketConfig 设置websocket所需ID
var WebSocketConfig = struct {
	GetConnectID func(string, interface{}) (string, error)
}{
	GetConnectID: func(platform string, data interface{}) (string, error) {

		r := ws.RequestLoginByAdmin{}
		realData := data.(map[string]interface{})
		if err := mapping.MapToStruct(realData, &r); err != nil {
			return "", err
		}

		adminID, err := strconv.Atoi(r.AdminID)
		if err != nil {
			return "", err
		}

		user := models.LoginAdmins.GetLogin(platform, int(adminID))
		if user == nil {
			fmt.Println("[用户]可能并未登录")
			return "", errors.New("[后台]用户还没有登录")
		}
		return r.AdminID, nil
	},
}
