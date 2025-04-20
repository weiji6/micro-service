package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		c.Keys["user"] = service[0]
		c.Keys["task"] = service[1]
		c.Next()
	}
}

// 错误处理中间件
func ErrorMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				context.JSON(200, gin.H{
					"code": 404,
					"msg":  fmt.Sprintf("%s", r),
				})
				context.Abort()
			}
		}()
		context.Next()
	}
}
