package middleware

import (
	"strings"

	"fail2ban-dashboard/internal/auth"
	"fail2ban-dashboard/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens from the Authorization header or cookie.
func AuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// Try Authorization header first
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Fallback to cookie
		if tokenStr == "" {
			if cookie, err := c.Cookie("token"); err == nil {
				tokenStr = cookie
			}
		}

		if tokenStr == "" {
			response.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(tokenStr)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
