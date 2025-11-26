package middleware

import (
	"net/http"

	"cloudsaver/internal/api/response"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 异常恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, response.Error("服务器内部错误"))
				c.Abort()
			}
		}()
		c.Next()
	}
}
