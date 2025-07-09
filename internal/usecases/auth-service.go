package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
	"fmt"
)

type AuthService struct {
	userRepo    ReadUserRepository
	studentRepo ReadStudentRepository
	teacherRepo ReadTeacherRepository
	adminRepo   ReadAdminRepository
	encryption  Cryptographer
	jwt         JWTGenerator
}

func NewAuthService(userRepo ReadUserRepository, studentRepo ReadStudentRepository, teacherRepo ReadTeacherRepository, adminRepo ReadAdminRepository, encryption Cryptographer, jwt JWTGenerator) *AuthService {
	return &AuthService{userRepo: userRepo, studentRepo: studentRepo, teacherRepo: teacherRepo, adminRepo: adminRepo, encryption: encryption, jwt: jwt}
}

func (a *AuthService) GetUserByLoginAndPassword(ctx context.Context, login, password string) (entities.User, error) {
	user, err := a.userRepo.ReadByLogin(ctx, login)
	if err != nil {
		return entities.User{}, UserNotFoundError
	}

	_, err = a.encryption.PasswordComparison(user.Password, password, user.Salt)
	if err != nil {
		return entities.User{}, DifferentPasswordError
	}

	return user, nil
}

func (a *AuthService) GetUserByAccessToken(ctx context.Context, token string) (entities.User, error) {
	dataFromToken, err := a.jwt.ParseJWT(token)
	if err != nil {
		return entities.User{}, fmt.Errorf("invalid token: %w", err)
	}

	id, ok := dataFromToken["sub"].(float64)
	if !ok {
		return entities.User{}, fmt.Errorf("invalid token payload")
	}

	user, err := a.userRepo.ReadById(ctx, int(id))
	if err != nil {
		return entities.User{}, UserNotFoundError
	}

	return user, nil
}

func (a *AuthService) GetStudentById(ctx context.Context, id int) (entities.Student, error) {
	student, err := a.studentRepo.ReadById(ctx, id)
	if err != nil {
		return entities.Student{}, UserAccountNotFoundError
	}

	return student, nil
}

func (a *AuthService) GetTeacherById(ctx context.Context, id int) (entities.Teacher, error) {
	teacher, err := a.teacherRepo.ReadById(ctx, id)
	if err != nil {
		return entities.Teacher{}, UserAccountNotFoundError
	}

	return teacher, nil
}

func (a *AuthService) GetAdminById(ctx context.Context, id int) (entities.Admin, error) {
	admin, err := a.adminRepo.ReadById(ctx, id)
	if err != nil {
		return entities.Admin{}, UserAccountNotFoundError
	}

	return admin, nil
}
