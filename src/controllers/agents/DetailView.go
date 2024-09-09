package agents

import (
	"fmt"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *Agents) DetailView(c *gin.Context) {
	id := c.Query("id")
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	sql := "select * from users where id=%s"
	sqll := fmt.Sprintf(sql, id)
	res, err := dbSession.QueryString(sqll)
	if err != nil {
		log.Err(err.Error())
		return
	}
	response.Render(c, "agents/agents_detail.html", pongo2.Context{"r": res[0]})
}
