package middlewares

import (
	"sports-common/config"
	"sports-common/consts"
	"strings"

	"github.com/gin-gonic/gin"
)

// 检测白名单当中的域名/IP
var checkWhiteList = func(key, checkValue string) bool {
	whites := config.Get(key, "")
	if whites == "" {
		return false
	}
	iArr := strings.Split(whites, ",") // 以,号分隔
	for _, white := range iArr {
		if white == "" {
			continue
		}
		realWhite := strings.TrimSpace(white) //去掉左右两边空格
		if realWhite == checkValue {
			return true
		}
	}
	return false
}

// DomainCheck 跨域请求
func DomainCheck() gin.HandlerFunc {

	return func(c *gin.Context) {
		hostName := c.Request.Host // 主机域名 - 当前来访

		// 先检测平台预设的盘口/站点url在不在范围之内
		for url := range consts.PlatformUrls {
			if hostName == url {
				c.Next()
				return
			}
		}

		// 再检测域名白名单 - 设置过ip白名单之后, 不再需要域名白名单
		// 这个地方主要处理 websocket 连接的白名单问题, 同时防止外部人员利用修改hosts来访问
		ip := c.ClientIP()                       // 来访ip
		if checkWhiteList("white_list.ip", ip) { // 如果ip白名单也满足条件, 则
			c.Next()
			return
		} else {
			Render(c, "unauthorized.html", "域名未经授权: "+hostName)
			return
		}
	}
}
