package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ricky97gr/lpcenter/server/utils"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		host := ctx.ClientIP()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		ctx.Next()
		// 不打印 OPTIONS 请求的日志
		if method == "OPTIONS" {
			return
		}
		raw := ctx.Request.URL.RawQuery
		status := ctx.Writer.Status()
		utils.Logger.Infof("| %d | %s | '%s' | %s | %+v | \t%s\t \n", status, host, path, method, time.Since(start), raw)

	}

}
