package router

import (
	_ "backendForKeenEye/docs"
	"backendForKeenEye/internal/container"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func NewRouter(c *container.Container) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := c.AuthMiddleware()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/api/create-student", auth, c.StudentController.CreateStudent)
	router.POST("/api/create-account", c.AccountController.CreateAccount)

	router.GET("/api/read-all-students", c.StudentController.ReadAllStudents)
	router.GET("/api/read-student", c.StudentController.ReadStudent)

	router.PUT("/api/update-student", auth, c.StudentController.UpdateStudent)
	router.DELETE("/api/delete-student", auth, c.StudentController.DeleteStudent)

	return router
}
