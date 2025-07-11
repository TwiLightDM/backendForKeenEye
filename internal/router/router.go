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
	admin := c.AdminMiddleware()
	teacherAdmin := c.TeacherAdminMiddleware()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/api/create-user", auth, admin, c.UserController.CreateUser)

	router.GET("/api/read-all-students", auth, admin, c.StudentController.ReadAllStudents)
	router.GET("/api/read-all-students-by-group-id", auth, c.StudentController.ReadAllStudentsByGroupId)
	router.GET("/api/read-student", auth, c.StudentController.ReadStudent)
	router.PUT("/api/update-student", auth, c.StudentController.UpdateStudent)
	router.DELETE("/api/delete-student", auth, admin, c.StudentController.DeleteStudent)

	router.GET("/api/read-all-teachers", auth, admin, c.TeacherController.ReadAllTeachers)
	router.GET("/api/read-teacher", auth, teacherAdmin, c.TeacherController.ReadTeacher)
	router.PUT("/api/update-teacher", auth, teacherAdmin, c.TeacherController.UpdateTeacher)
	router.DELETE("/api/delete-teacher", auth, admin, c.TeacherController.DeleteTeacher)

	router.GET("/api/read-admin", auth, admin, c.AdminController.ReadAdmin)
	router.PUT("/api/update-admin", auth, admin, c.AdminController.UpdateAdmin)
	router.DELETE("/api/delete-admin", auth, admin, c.AdminController.DeleteAdmin)

	router.POST("/api/create-group", auth, admin, c.GroupController.CreateGroup)
	router.GET("/api/read-all-groups", auth, admin, c.GroupController.ReadAllGroups)
	router.GET("/api/read-group", auth, c.GroupController.ReadGroup)
	router.PUT("/api/update-group", auth, admin, c.GroupController.UpdateGroup)
	router.DELETE("/api/delete-group", auth, admin, c.GroupController.DeleteGroup)

	return router
}
