package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
					"path":  c.Request.URL.Path,
				}).Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
