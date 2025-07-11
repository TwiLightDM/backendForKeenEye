package middlewares

import (
	"backendForKeenEye/internal/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

func checkingRoles(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		userRaw, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		user, ok := userRaw.(entities.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "User data is corrupted"})
			return
		}

		if _, ok = roleSet[user.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient role"})
			return
		}

		c.Next()
	}
}
