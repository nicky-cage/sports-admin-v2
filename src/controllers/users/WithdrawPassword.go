package users

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

// WithdrawPassword 修改提现密码
func (ths *Users) WithdrawPassword(c *gin.Context) {
	//获取用户的Id 检验是否绑定手机  生成验证码。发送验证码。 保存验证码。  验证验证码， 重置资金密码。
	postData := request.GetPostedData(c)
	id := postData["id"]
	if id == "" {
		response.Err(c, "缺少必要参数")
		return
	}
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "SELECT phone, id, username FROM users WHERE id = ?"
	phone, err := dbSession.QueryString(sql, id)
	if err != nil {
		log.Err(err.Error())
	}
	if len(phone) == 0 {
		response.Err(c, "未绑定手机号码")
		return
	}
	resetSql := "update users set withdraw_password='' where id=?"
	_, err = dbSession.Exec(resetSql, phone[0]["id"])
	if err != nil {
		log.Err(err.Error())
		response.Err(c, "重置错误")
		return
	}
	response.Result(c, "已重置")
}
