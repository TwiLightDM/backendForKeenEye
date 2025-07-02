package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
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
		return entities.Account{}, AccountNotFoundError
	}

	_, err = a.encryption.PasswordComparison(account.Password, password, account.Salt)
	if err != nil {
		return entities.Account{}, DifferentPasswordError
	}

	return account, nil
}
