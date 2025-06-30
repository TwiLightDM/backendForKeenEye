package main

import (
	"backendForKeenEye/config"
	_ "backendForKeenEye/docs"
	"backendForKeenEye/internal/controllers"
	"backendForKeenEye/internal/middleware"
	"backendForKeenEye/internal/repositories"
	"backendForKeenEye/internal/usecases"
	encryption_service "backendForKeenEye/pkg/encryption-service"
	"backendForKeenEye/pkg/postgres"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

// @title Backend for KeenEye
// @version 1.0.0
// @description Backend for KeenEye

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey BasicAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	pgClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	err = pgClient.MigrateUp()
	if err != nil {
		fmt.Println("failed to migrate:", err)
	}

	ctx := context.Background()

	encryption := encryption_service.NewEncryptionService(cfg.Salt)

	studentsRepo := repositories.NewStudentRepository(pgClient.Pool, pgClient.Builder)
	accountsRepo := repositories.NewAccountRepository(pgClient.Pool, pgClient.Builder)

	authService := usecases.NewAuthService(accountsRepo, encryption)

	createStudentUsecase := usecases.NewCreateStudentUsecase(studentsRepo)
	createAccountUsecase := usecases.NewCreateAccountUsecase(accountsRepo, encryption)

	readAllStudentsUsecase := usecases.NewReadAllStudentsUsecase(studentsRepo)
	readStudentUsecase := usecases.NewReadStudentUsecase(studentsRepo)

	updateStudentUsecase := usecases.NewUpdateStudentUsecase(studentsRepo)

	deleteStudentUsecase := usecases.NewDeleteStudentUsecase(studentsRepo)

	studentController := controllers.NewStudentController(&createStudentUsecase, &readAllStudentsUsecase, &readStudentUsecase, &updateStudentUsecase, &deleteStudentUsecase)
	accountController := controllers.NewAccountController(&createAccountUsecase)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authMiddleware := middleware.AuthMiddleware(ctx, authService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/api/create-student", authMiddleware, studentController.CreateStudent)
	router.POST("/api/create-account", accountController.CreateAccount)

	router.GET("/api/read-all-students", studentController.ReadAllStudents)
	router.GET("/api/read-student", studentController.ReadStudent)

	router.PUT("/api/update-student", authMiddleware, studentController.UpdateStudent)

	router.DELETE("/api/delete-student", authMiddleware, studentController.DeleteStudent)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router)
	fmt.Println(err)
}
