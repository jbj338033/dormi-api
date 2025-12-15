package middleware

import (
	"net/http"

	"dormi-api/internal/dto"
	"dormi-api/internal/model"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...model.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Error:   "unauthorized",
			})
			c.Abort()
			return
		}

		role := userRole.(model.Role)
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, dto.Response{
			Success: false,
			Error:   "insufficient permissions",
		})
		c.Abort()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return RequireRole(model.RoleAdmin)
}

func RequireAdminOrSupervisor() gin.HandlerFunc {
	return RequireRole(model.RoleAdmin, model.RoleSupervisor)
}
