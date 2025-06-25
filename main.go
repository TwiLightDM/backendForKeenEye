package main

import (
	"backendForKeenEye/config"
	_ "backendForKeenEye/docs"
	"backendForKeenEye/internal/controllers"
	"backendForKeenEye/internal/repositories"
	"backendForKeenEye/internal/usecases"
	"backendForKeenEye/pkg/postgres"
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

	studentsRepo := repositories.NewStudentRepository(pgClient.Pool, pgClient.Builder)

	createStudentUsecase := usecases.NewCreateStudentUsecase(studentsRepo)
	readAllStudentsUsecase := usecases.NewReadAllStudentsUsecase(studentsRepo)
	readStudentUsecase := usecases.NewReadStudentUsecase(studentsRepo)
	updateStudentUsecase := usecases.NewUpdateStudentUsecase(studentsRepo)
	deleteStudentUsecase := usecases.NewDeleteStudentUsecase(studentsRepo)

	studentController := controllers.NewStudentController(&createStudentUsecase, &readAllStudentsUsecase, &readStudentUsecase, &updateStudentUsecase, &deleteStudentUsecase)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/api/create-student", studentController.CreateStudent)

	router.GET("/api/read-all-students", studentController.ReadAllStudents)
	router.GET("/api/read-student", studentController.ReadStudent)

	router.PUT("/api/update-student", studentController.UpdateStudent)

	router.DELETE("/api/delete-student", studentController.DeleteStudent)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router)
	fmt.Println(err)
}
