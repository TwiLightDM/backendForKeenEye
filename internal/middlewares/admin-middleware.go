package middlewares

import (
	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return checkingRoles("admin")
}
