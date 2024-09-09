package agents

import (
	"fmt"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ths *Agents) LowerAdd(c *gin.Context) {
	data := request.GetPostedData(c)
	username := data["username"].(string)
	top_name := data["top_name"].(string)
	beforeAgent := data["before_agent"].(string)
	top_id := data["top_id"].(string)
	remark := data["remark"].(string)
	arr := strings.Split(username, ",")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform, true)
	defer dbSession.Close()

	sql := "update users set top_name='%s', top_id='%s',remark='%s',trans_before_agent='%s' where username='%s'"
	for _, v := range arr {
		csql := "select is_agent,top_name from users where username='" + v + "'"
		cRes, _ := dbSession.QueryString(csql)
		if cRes[0]["is_agent"] == "1" {
			response.Err(c, "该用户是代理")
			return
		}
		//if cRes[0]["top_name"] != "" {
		//	response.Err(c, "该用户"+v+"已有代理")
		//	return
		//}
		sqll := fmt.Sprintf(sql, top_name, top_id, remark, beforeAgent, v)
		_, err := dbSession.QueryString(sqll)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "转代错误")
			return
		}
	}
	response.Result(c, "ok")
}
