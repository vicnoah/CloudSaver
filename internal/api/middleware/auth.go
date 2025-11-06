package middleware

import (
	"net/http"
	"strings"

	"cloudsaver/internal/api/response"
	"cloudsaver/internal/pkg/utils"
	"cloudsaver/internal/repository"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(jwtSecret string, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过认证的路径
		skipPaths := []string{
			"/api/user/login",
			"/api/user/register",
		}

		path := c.Request.URL.Path
		for _, skipPath := range skipPaths {
			if path == skipPath {
				c.Next()
				return
			}
		}

		// 提取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Error("未提供Token"))
			c.Abort()
			return
		}

		// 解析Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.Error("Token格式错误"))
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析Token
		claims, err := utils.ParseToken(tokenString, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("Token无效或已过期"))
			c.Abort()
			return
		}

		// 验证用户是否存在
		user, err := userRepo.FindByUUID(claims.UserUUID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error("用户不存在"))
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_uuid", user.UUID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}
