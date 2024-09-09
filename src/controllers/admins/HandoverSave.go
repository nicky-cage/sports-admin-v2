package admins

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

func (ths *Admins) HandoverSave(c *gin.Context) {
	platform := request.GetPlatform(c)
	postData := request.GetPostedData(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	role := postData["receiver_role"].(string)
	receiver := postData["receiver"].(string)
	sql := "update admins set role_name= ? where name=? "
	_, err := dbSession.Exec(sql, role, receiver)
	if err != nil {
		log.Err(err.Error())
		response.Err(c, "交接失败")
		return
	}
	response.Ok(c)
}
