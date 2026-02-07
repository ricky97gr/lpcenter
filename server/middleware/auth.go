package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Failed(c, 401, "未授权，请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Failed(c, 401, "授权格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			response.Failed(c, 401, "令牌无效或已过期")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
