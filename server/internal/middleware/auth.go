package middleware

import (
	"strings"

	"store-ops-server/internal/pkg/auth"
	"store-ops-server/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthRequired JWT 认证中间件
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.AuthFailed(c, "未提供认证信息")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.AuthFailed(c, "认证格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			response.AuthFailed(c, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}