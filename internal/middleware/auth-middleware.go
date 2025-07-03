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

		case strings.HasPrefix(authHeader, "Jwt "):
			JWTAuthMiddleware(authService)(c)

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

		c.Set("account", account)
		c.Next()
	}
}

func JWTAuthMiddleware(jwtAuthService *usecases.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Jwt ")

		account, err := jwtAuthService.GetAccountByAccessToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
			return
		}

		c.Set("account", account)
		c.Next()
	}
}
