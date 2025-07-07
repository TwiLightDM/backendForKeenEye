package middleware

import (
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

		account, err := authService.GetAccountByLoginAndPassword(ctx, login, password)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Basic credentials"})
			return
		}

		switch account.Role {
		case "student":
			student, err := authService.GetStudentByAccountId(ctx, account.Id)
			if err == nil {
				c.Set("student", student)
			}
		case "teacher":
			teacher, err := authService.GetTeacherByAccountId(ctx, account.Id)
			if err == nil {
				c.Set("teacher", teacher)
			}
		case "admin":
			admin, err := authService.GetAdminByAccountId(ctx, account.Id)
			if err == nil {
				c.Set("admin", admin)
			}
		}

		c.Set("account", account)
		c.Next()
	}
}

func JWTAuthMiddleware(ctx context.Context, jwtAuthService *usecases.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		account, err := jwtAuthService.GetAccountByAccessToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			return
		}

		switch account.Role {
		case "student":
			student, err := jwtAuthService.GetStudentByAccountId(ctx, account.Id)
			if err == nil {
				c.Set("student", student)
			}
		case "teacher":
			teacher, err := jwtAuthService.GetTeacherByAccountId(ctx, account.Id)
			if err == nil {
				c.Set("teacher", teacher)
			}
		case "admin":
			admin, err := jwtAuthService.GetAdminByAccountId(ctx, account.Id)
			if err == nil {
				c.Set("admin", admin)
			}
		}

		c.Set("account", account)
		c.Next()
	}
}
