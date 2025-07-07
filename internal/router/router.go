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

	router.POST("/api/create-account", c.AccountController.CreateAccount)

	router.POST("/api/create-student", auth, c.StudentController.CreateStudent)
	router.GET("/api/read-all-students", auth, c.StudentController.ReadAllStudents)
	router.GET("/api/read-all-students-from-group", auth, c.StudentController.ReadAllStudentsByGroupId)
	router.GET("/api/read-student", auth, c.StudentController.ReadStudent)
	router.PUT("/api/update-student", auth, c.StudentController.UpdateStudent)
	router.DELETE("/api/delete-student", auth, c.StudentController.DeleteStudent)

	router.POST("/api/create-teacher", auth, c.TeacherController.CreateTeacher)
	router.GET("/api/read-all-teachers", auth, c.TeacherController.ReadAllTeachers)
	router.GET("/api/read-teacher", auth, c.TeacherController.ReadTeacher)
	router.PUT("/api/update-teacher", auth, c.TeacherController.UpdateTeacher)
	router.DELETE("/api/delete-teacher", auth, c.TeacherController.DeleteTeacher)

	router.POST("/api/create-admin", auth, c.AdminController.CreateAdmin)
	router.GET("/api/read-admin", auth, c.AdminController.ReadAdmin)
	router.PUT("/api/update-admin", auth, c.AdminController.UpdateAdmin)
	router.DELETE("/api/delete-admin", auth, c.AdminController.DeleteAdmin)

	router.POST("/api/create-group", auth, c.GroupController.CreateGroup)
	router.GET("/api/read-all-groups", auth, c.GroupController.ReadAllGroups)
	router.GET("/api/read-group", auth, c.GroupController.ReadGroup)
	router.PUT("/api/update-group", auth, c.GroupController.UpdateGroup)
	router.DELETE("/api/delete-group", auth, c.GroupController.DeleteGroup)

	return router
}
