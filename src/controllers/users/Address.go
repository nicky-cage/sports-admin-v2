package users

import (
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/gin-gonic/gin"
)

// Address 用户地址
// 点击收货省份，城市，返回对应下的省份下城市，城市下的线程，
// 现点击会自动到省份，城市的附近，如运维没有需求在删掉
func (ths *Users) Address(c *gin.Context) {
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	Type := c.Query("type")
	var list []map[string]string

	switch Type {
	case "1":
		code := c.Query("code")
		csql := "select * from cities where province_code=" + code
		res, err := dbSession.QueryString(csql)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "程序错误")
		}
		list = res
	case "2":
		code := c.Query("code")
		dsql := "select * from districts where city_code= " + code
		res, err := dbSession.QueryString(dsql)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "程序错误")
		}
		list = res
	}
	response.Result(c, list)
}
