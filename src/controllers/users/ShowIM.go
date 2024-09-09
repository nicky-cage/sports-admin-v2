package users

import (
	"fmt"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ths *Users) ShowIM(c *gin.Context) {
	typeStr := c.DefaultQuery("type", "")
	if typeStr == "" {
		response.Err(c, "错误的数据信息, 也有可能缺少权限")
		return
	}
	userIDStr := c.DefaultQuery("id", "")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		response.Err(c, "缺少必要用户编号信息")
		return
	}

	typeMap := map[string]string{
		"im_qq":     "qq",
		"im_wechat": "we_chat",
		"im_phone":  "phone",
		"im_email":  "email",
	}
	field, exists := typeMap[typeStr]
	if !exists {
		response.Err(c, "缺少相关数据信息, 或者缺少权限")
		return
	}

	platform := request.GetPlatform(c)
	sql := fmt.Sprintf("SELECT IFNULL(%s, '') AS %s  FROM users WHERE id = %d", field, field, userID)
	mConn := common.Mysql(platform)
	defer mConn.Close()
	rows, err := mConn.QueryString(sql)
	if err != nil || len(rows) == 0 {
		response.Err(c, "查询数据有误: "+err.Error())
		return
	}

	message := "[" + strings.ToUpper(field) + "]<br />" + func() string {
		msg := rows[0][field]
		if msg == "" {
			return "暂未设置"
		}
		return msg
	}()
	response.Result(c, message)
}
