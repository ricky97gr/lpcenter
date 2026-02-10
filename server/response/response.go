package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, result interface{}, total int64) {
	ctx.JSON(
		http.StatusOK,
		gin.H{
			"code":   SuccessCode,
			"msg":    errCodeMap[SuccessCode].msgCn,
			"result": result,
			"total":  total,
		},
	)
}

func Failed(ctx *gin.Context, errCode int32, msg string) {
	if msg == "" {
		msg = errCodeMap[errCode].msgCn
	}
	ctx.JSON(
		http.StatusServiceUnavailable,
		gin.H{
			"code": errCode,
			"msg":  msg,
		},
	)
}
