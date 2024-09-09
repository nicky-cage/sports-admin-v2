package agent_domains

import (
	"errors"
	"fmt"
	"sports-admin/controllers/base_controller"
	common "sports-common"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

var ActionSave = &base_controller.ActionSave{
	Model: models.AgentDomains,
	SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
		// 自动设定用户ID
		if userName, exists := (*m)["username"]; !exists {
			return errors.New("缺少代理用户名称")
		} else if userName == "" {
			return errors.New("代理用户名称为空")
		} else {

			platform := request.GetPlatform(c)
			myClient := common.Mysql(platform)
			defer myClient.Close()

			row := struct {
				Id int `json:"id"`
			}{}
			if exists, err := myClient.SQL(fmt.Sprintf("SELECT id FROM users WHERE username = '%s' LIMIT 1", userName)).Get(&row); err != nil {
				return errors.New("获取代理用户信息失败:" + err.Error())
			} else if !exists {
				return errors.New("查找用户信息失败")
			} else {
				(*m)["user_id"] = row.Id
			}
		}
		return nil
	},
	SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
		// 重写 配置文件
		// 重启 nginx
	},
}
