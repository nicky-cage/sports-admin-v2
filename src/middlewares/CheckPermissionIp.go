package middlewares

import (
	"sports-admin/caches"
	"sports-common/config"
	"sports-common/request"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// CheckPermissionIp 检测是否有IP授权
func CheckPermissionIp() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 先要检测白名单 - 如果有白名单, 则直接返回
		whites := config.Get("white_list.ip", "")
		if whites != "" { // 如果白名单不为空
			iArr := strings.Split(whites, ",") // 以,号分隔
			for _, IP := range iArr {
				whiteIP := strings.TrimSpace(IP) //去掉左右两边空格
				if whiteIP == clientIP {
					c.Set("white_list_ip", whiteIP)
					c.Next()
					return
				}
			}
		}

		platform := request.GetPlatform(c)
		hasPermission := caches.PermissionIps.HasPermission(platform, clientIP) // 是否有授权
		if !hasPermission {                                                     // 如果没有权限
			//Render(c, "404.html", "操作失败, 你的IP没有权限", pongo2.Context{"ip": ip})
			Render(c, "404.html", "Operation failed, your IP does not have permission", pongo2.Context{"ip": clientIP})
			return
		}
		c.Next()
	}
}
