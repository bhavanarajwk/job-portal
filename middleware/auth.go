package middleware

import (
	"strings"

	"job-portal/models"
	"job-portal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

// JWTAuth validates the Bearer token and injects user_id and role into context
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.Unauthorized(c, "authorization header format must be: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

// RequireRole restricts access to users with one of the specified roles
func RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextRole)
		if !exists {
			utils.Unauthorized(c, "unauthorized")
			c.Abort()
			return
		}

		userRole := models.Role(roleVal.(string))
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "you do not have permission to access this resource")
		c.Abort()
	}
}

// GetUserID extracts the authenticated user's UUID from context
func GetUserID(c *gin.Context) uuid.UUID {
	val, _ := c.Get(ContextUserID)
	return val.(uuid.UUID)
}

// GetRole extracts the authenticated user's role from context
func GetRole(c *gin.Context) models.Role {
	val, _ := c.Get(ContextRole)
	return models.Role(val.(string))
}
