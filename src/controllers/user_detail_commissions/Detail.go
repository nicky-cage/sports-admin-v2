package user_detail_commissions

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

func (ths *UserDetailCommissions) Detail(c *gin.Context) {
	id := c.Query("user_id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select balance from accounts where user_id=" + id
	res, err := dbSession.QueryString(sql)
	if err != nil {
		log.Err(err.Error())
		return
	}
	response.Result(c, res[0]["balance"])
}
