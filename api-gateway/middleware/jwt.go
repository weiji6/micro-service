package middleware

import (
	"etcd/api-gateway/pkg/e"
	"etcd/api-gateway/pkg/util"
	"github.com/gin-gonic/gin"
	"time"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claim, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail
			} else if time.Now().Unix() > claim.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeOut
			}
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(uint(code)),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
