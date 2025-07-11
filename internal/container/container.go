package container

import (
	"backendForKeenEye/config"
	"backendForKeenEye/internal/controllers"
	"backendForKeenEye/internal/middlewares"
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

	UserController    controllers.UserController
	StudentController controllers.StudentController
	TeacherController controllers.TeacherController
	AdminController   controllers.AdminController
	GroupController   controllers.GroupController

	AuthMiddleware         func() func(c *gin.Context)
	AdminMiddleware        func() func(c *gin.Context)
	TeacherAdminMiddleware func() func(c *gin.Context)
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
	userRepo := repositories.NewUserRepository(pgClient.Pool, pgClient.Builder)
	teacherRepo := repositories.NewTeacherRepository(pgClient.Pool, pgClient.Builder)
	adminRepo := repositories.NewAdminRepository(pgClient.Pool, pgClient.Builder)
	groupRepo := repositories.NewGroupRepository(pgClient.Pool, pgClient.Builder)

	authService := usecases.NewAuthService(userRepo, studentRepo, teacherRepo, adminRepo, encryption, jwt)

	readAllStudents := usecases.NewReadAllStudentsUsecase(studentRepo)
	readAllStudentsByGroupId := usecases.NewReadAllStudentsByGroupIdUsecase(studentRepo)
	readStudent := usecases.NewReadStudentUsecase(studentRepo)
	updateStudent := usecases.NewUpdateStudentUsecase(studentRepo)
	deleteStudent := usecases.NewDeleteStudentUsecase(studentRepo)

	createUser := usecases.NewCreateUserUsecase(userRepo, encryption, jwt)

	readAllTeachers := usecases.NewReadAllTeachersUsecase(teacherRepo)
	readTeacher := usecases.NewReadTeacherUsecase(teacherRepo)
	updateTeacher := usecases.NewUpdateTeacherUsecase(teacherRepo)
	deleteTeacher := usecases.NewDeleteTeacherUsecase(teacherRepo)

	readAdmin := usecases.NewReadAdminUsecase(adminRepo)
	updateAdmin := usecases.NewUpdateAdminUsecase(adminRepo)
	deleteAdmin := usecases.NewDeleteAdminUsecase(adminRepo)

	createGroup := usecases.NewCreateGroupUsecase(groupRepo)
	readAllGroups := usecases.NewReadAllGroupsUsecase(groupRepo)
	readGroup := usecases.NewReadGroupUsecase(groupRepo)
	updateGroup := usecases.NewUpdateGroupUsecase(groupRepo)
	deleteGroup := usecases.NewDeleteGroupUsecase(groupRepo)

	accountController := controllers.NewUserController(&createUser)

	studentController := controllers.NewStudentController(
		&readGroup,
		&readAllStudents,
		&readAllStudentsByGroupId,
		&readStudent,
		&updateStudent,
		&deleteStudent,
	)

	teacherController := controllers.NewTeacherController(
		&readAllTeachers,
		&readTeacher,
		&updateTeacher,
		&deleteTeacher,
	)

	adminController := controllers.NewAdminController(
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
		Cfg:                    *cfg,
		Ctx:                    ctx,
		PGClient:               pgClient,
		UserController:         accountController,
		StudentController:      studentController,
		TeacherController:      teacherController,
		AdminController:        adminController,
		GroupController:        groupController,
		AuthMiddleware:         func() func(c *gin.Context) { return middlewares.AuthMiddleware(ctx, authService) },
		AdminMiddleware:        func() func(c *gin.Context) { return middlewares.AdminMiddleware() },
		TeacherAdminMiddleware: func() func(c *gin.Context) { return middlewares.TeacherAdminMiddleware() },
	}
}
