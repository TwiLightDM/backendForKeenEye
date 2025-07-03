package container

import (
	"backendForKeenEye/config"
	"backendForKeenEye/internal/controllers"
	"backendForKeenEye/internal/middleware"
	"backendForKeenEye/internal/repositories"
	"backendForKeenEye/internal/usecases"
	encryptionService "backendForKeenEye/pkg/encryption-service"
	jwtService "backendForKeenEye/pkg/jwt-service"
	"backendForKeenEye/pkg/postgres"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

type Container struct {
	Cfg config.Config
	Ctx context.Context

	StudentController controllers.StudentController
	AccountController controllers.AccountController

	AuthMiddleware func() func(c *gin.Context)
}

func NewContainer() *Container {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pgClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to create DB client: %v", err)
	}

	if err = pgClient.MigrateUp(); err != nil {
		fmt.Printf("failed to migrate: %v\n", err)
	}

	ctx := context.Background()
	encryption := encryptionService.NewEncryptionService(cfg.Salt)
	jwt := jwtService.NewJWTService(cfg.Key, cfg.AccessTime, cfg.RefreshTime)

	studentRepo := repositories.NewStudentRepository(pgClient.Pool, pgClient.Builder)
	accountRepo := repositories.NewAccountRepository(pgClient.Pool, pgClient.Builder)

	authService := usecases.NewAuthService(accountRepo, encryption, jwt)

	createStudent := usecases.NewCreateStudentUsecase(studentRepo)
	readAllStudents := usecases.NewReadAllStudentsUsecase(studentRepo)
	readStudent := usecases.NewReadStudentUsecase(studentRepo)
	updateStudent := usecases.NewUpdateStudentUsecase(studentRepo)
	deleteStudent := usecases.NewDeleteStudentUsecase(studentRepo)
	createAccount := usecases.NewCreateAccountUsecase(accountRepo, encryption, jwt)

	accountController := controllers.NewAccountController(&createAccount)
	studentController := controllers.NewStudentController(
		&createStudent,
		&readAllStudents,
		&readStudent,
		&updateStudent,
		&deleteStudent,
	)

	return &Container{
		Cfg:               *cfg,
		Ctx:               ctx,
		StudentController: studentController,
		AccountController: accountController,
		AuthMiddleware:    func() func(c *gin.Context) { return middleware.AuthMiddleware(ctx, authService) },
	}
}
