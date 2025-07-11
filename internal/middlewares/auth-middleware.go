package middlewares

import (
	"backendForKeenEye/internal/entities"
	"backendForKeenEye/internal/usecases"
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(ctx context.Context, authService *usecases.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		switch {
		case strings.HasPrefix(authHeader, "Basic "):
			BasicAuthMiddleware(ctx, authService)(c)

		case strings.HasPrefix(authHeader, "Bearer "):
			JWTAuthMiddleware(ctx, authService)(c)

		default:
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unsupported or missing Authorization header"})
		}
	}
}

func BasicAuthMiddleware(ctx context.Context, authService *usecases.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.TrimPrefix(c.GetHeader("Authorization"), "Basic ")
		decoded, err := base64.StdEncoding.DecodeString(auth)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid base64 encoding"})
			return
		}

		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Basic auth format"})
			return
		}

		login, password := parts[0], parts[1]
		user, err := authService.GetUserByLoginAndPassword(ctx, login, password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Basic credentials"})
			return
		}

		c.Set("user", user)
		AttachUserRoleData(ctx, c, authService, user)
		c.Next()
	}
}

func JWTAuthMiddleware(ctx context.Context, authService *usecases.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		user, err := authService.GetUserByAccessToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			return
		}

		c.Set("user", user)
		AttachUserRoleData(ctx, c, authService, user)
		c.Next()
	}
}

func AttachUserRoleData(ctx context.Context, c *gin.Context, authService *usecases.AuthService, user entities.User) {
	switch user.Role {
	case "student":
		if student, err := authService.GetStudentById(ctx, user.Id); err == nil {
			c.Set("student", student)
		}
	case "teacher":
		if teacher, err := authService.GetTeacherById(ctx, user.Id); err == nil {
			c.Set("teacher", teacher)
		}
	case "admin":
		if admin, err := authService.GetAdminById(ctx, user.Id); err == nil {
			c.Set("admin", admin)
		}
	}
}
