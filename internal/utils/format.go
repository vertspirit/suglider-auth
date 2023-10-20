package utils

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, errCode int64) (responseData gin.H) {
	responseData = gin.H{
		"status":    "success",
		"code":      errCode,
		"message":   CodeMap[errCode],
	}

	return
}

func ErrorResponse(c *gin.Context, errCode int64, errInfo ...interface{}) (responseData gin.H) {
	responseData = gin.H{
		"status":    "error",
		"code":      errCode,
		"message":   CodeMap[errCode],
		"error_msg":  errInfo,
	}

	return
}