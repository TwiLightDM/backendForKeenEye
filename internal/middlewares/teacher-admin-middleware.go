package middlewares

import (
	"github.com/gin-gonic/gin"
)

func TeacherAdminMiddleware() gin.HandlerFunc {
	return checkingRoles("teacher", "admin")
}
