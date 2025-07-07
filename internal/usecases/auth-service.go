package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
	"fmt"
)

type AuthService struct {
	accountRepo ReadAccountRepository
	studentRepo ReadStudentRepository
	teacherRepo ReadTeacherRepository
	adminRepo   ReadAdminRepository
	encryption  Cryptographer
	jwt         JWTGenerator
}

func NewAuthService(accountRepo ReadAccountRepository, studentRepo ReadStudentRepository, teacherRepo ReadTeacherRepository, adminRepo ReadAdminRepository, encryption Cryptographer, jwt JWTGenerator) *AuthService {
	return &AuthService{accountRepo: accountRepo, studentRepo: studentRepo, teacherRepo: teacherRepo, adminRepo: adminRepo, encryption: encryption, jwt: jwt}
}

func (a *AuthService) GetAccountByLoginAndPassword(ctx context.Context, login, password string) (entities.Account, error) {
	account, err := a.accountRepo.ReadByLogin(ctx, login)
	if err != nil {
		return entities.Account{}, AccountNotFoundError
	}

	_, err = a.encryption.PasswordComparison(account.Password, password, account.Salt)
	if err != nil {
		return entities.Account{}, DifferentPasswordError
	}

	return account, nil
}

func (a *AuthService) GetAccountByAccessToken(ctx context.Context, token string) (entities.Account, error) {
	dataFromToken, err := a.jwt.ParseJWT(token)
	if err != nil {
		return entities.Account{}, fmt.Errorf("invalid token: %w", err)
	}

	accountID, ok := dataFromToken["sub"].(float64)
	if !ok {
		return entities.Account{}, fmt.Errorf("invalid token payload")
	}

	account, err := a.accountRepo.ReadById(ctx, int(accountID))
	if err != nil {
		return entities.Account{}, AccountNotFoundError
	}

	return account, nil
}

func (a *AuthService) GetStudentByAccountId(ctx context.Context, id int) (entities.Student, error) {
	student, err := a.studentRepo.ReadByAccountId(ctx, id)
	if err != nil {
		return entities.Student{}, UserAccountNotFoundError
	}

	return student, nil
}

func (a *AuthService) GetTeacherByAccountId(ctx context.Context, id int) (entities.Teacher, error) {
	teacher, err := a.teacherRepo.ReadByAccountId(ctx, id)
	if err != nil {
		return entities.Teacher{}, UserAccountNotFoundError
	}

	return teacher, nil
}

func (a *AuthService) GetAdminByAccountId(ctx context.Context, id int) (entities.Admin, error) {
	admin, err := a.adminRepo.ReadByAccountId(ctx, id)
	if err != nil {
		return entities.Admin{}, UserAccountNotFoundError
	}

	return admin, nil
}
