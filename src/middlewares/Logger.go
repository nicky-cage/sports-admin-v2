package middlewares

import (
	"sports-common/log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()                                // 开始时间
		c.Next()                                           // 处理请求
		end := time.Now()                                  // 结束时间
		latency := end.Sub(start)                          // 执行时间
		log.Logger.Infof("| %3d | %13v | %15s | %s  %s |", // 指定日志打印出来的格式。分别是状态码，执行时间,请求ip,请求方法,请求路由
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}
