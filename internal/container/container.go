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

	PGClient *postgres.Client

	AccountController controllers.AccountController
	StudentController controllers.StudentController
	TeacherController controllers.TeacherController
	AdminController   controllers.AdminController
	GroupController   controllers.GroupController

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
	teacherRepo := repositories.NewTeacherRepository(pgClient.Pool, pgClient.Builder)
	adminRepo := repositories.NewAdminRepository(pgClient.Pool, pgClient.Builder)
	groupRepo := repositories.NewGroupRepository(pgClient.Pool, pgClient.Builder)

	authService := usecases.NewAuthService(accountRepo, studentRepo, teacherRepo, adminRepo, encryption, jwt)

	createStudent := usecases.NewCreateStudentUsecase(studentRepo)
	readAllStudents := usecases.NewReadAllStudentsUsecase(studentRepo)
	readAllStudentsByGroupId := usecases.NewReadAllStudentsByGroupIdUsecase(studentRepo)
	readStudent := usecases.NewReadStudentUsecase(studentRepo)
	updateStudent := usecases.NewUpdateStudentUsecase(studentRepo)
	deleteStudent := usecases.NewDeleteStudentUsecase(studentRepo)

	createAccount := usecases.NewCreateAccountUsecase(accountRepo, encryption, jwt)

	createTeacher := usecases.NewCreateTeacherUsecase(teacherRepo)
	readAllTeachers := usecases.NewReadAllTeachersUsecase(teacherRepo)
	readTeacher := usecases.NewReadTeacherUsecase(teacherRepo)
	updateTeacher := usecases.NewUpdateTeacherUsecase(teacherRepo)
	deleteTeacher := usecases.NewDeleteTeacherUsecase(teacherRepo)

	createAdmin := usecases.NewCreateAdminUsecase(adminRepo)
	readAdmin := usecases.NewReadAdminUsecase(adminRepo)
	updateAdmin := usecases.NewUpdateAdminUsecase(adminRepo)
	deleteAdmin := usecases.NewDeleteAdminUsecase(adminRepo)

	createGroup := usecases.NewCreateGroupUsecase(groupRepo)
	readAllGroups := usecases.NewReadAllGroupsUsecase(groupRepo)
	readGroup := usecases.NewReadGroupUsecase(groupRepo)
	updateGroup := usecases.NewUpdateGroupUsecase(groupRepo)
	deleteGroup := usecases.NewDeleteGroupUsecase(groupRepo)

	accountController := controllers.NewAccountController(&createAccount)

	studentController := controllers.NewStudentController(
		&readGroup,
		&createStudent,
		&readAllStudents,
		&readAllStudentsByGroupId,
		&readStudent,
		&updateStudent,
		&deleteStudent,
	)

	teacherController := controllers.NewTeacherController(
		&createTeacher,
		&readAllTeachers,
		&readTeacher,
		&updateTeacher,
		&deleteTeacher,
	)

	adminController := controllers.NewAdminController(
		&createAdmin,
		&readAdmin,
		&updateAdmin,
		&deleteAdmin,
	)

	groupController := controllers.NewGroupController(
		&createGroup,
		&readAllGroups,
		&readGroup,
		&updateGroup,
		&deleteGroup,
	)

	return &Container{
		Cfg:               *cfg,
		Ctx:               ctx,
		PGClient:          pgClient,
		AccountController: accountController,
		StudentController: studentController,
		TeacherController: teacherController,
		AdminController:   adminController,
		GroupController:   groupController,
		AuthMiddleware:    func() func(c *gin.Context) { return middleware.AuthMiddleware(ctx, authService) },
	}
}
