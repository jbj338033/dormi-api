package middleware

import (
	"net/http"
	"strings"

	"dormi-api/internal/dto"
	"dormi-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Error:   "authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Error:   "invalid authorization header format",
			})
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Error:   "invalid or expired token",
			})
			c.Abort()
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Error:   "invalid user id in token",
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}
