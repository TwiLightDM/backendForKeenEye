package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
	"fmt"
)

type AuthService struct {
	accountRepo ReadAccountRepository
	encryption  Cryptographer
}

func NewAuthService(accountRepo ReadAccountRepository, encryption Cryptographer) *AuthService {
	return &AuthService{accountRepo: accountRepo, encryption: encryption}
}

func (a *AuthService) GetAccountByLoginAndPassword(ctx context.Context, login, password string) (entities.Account, error) {
	account, err := a.accountRepo.ReadByLogin(ctx, login)
	if err != nil {
		return entities.Account{}, fmt.Errorf("account not found: %w", err)
	}

	_, err = a.encryption.PasswordComparison(account.Password, password, account.Salt)
	if err != nil {
		return entities.Account{}, fmt.Errorf("passwords are not similar: %w", err)
	}

	return account, nil
}
