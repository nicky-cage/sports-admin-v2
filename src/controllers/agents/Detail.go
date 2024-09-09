package agents

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

func (ths *Agents) Detail(c *gin.Context) {
	name := c.Query("username")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select top_name,username from users where username='" + name + "'"
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		return
	}
	if len(res) > 0 {
		response.Result(c, res[0])
	} else {
		response.Err(c, "用户不存在")
	}
}
