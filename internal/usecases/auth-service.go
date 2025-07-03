package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
	"fmt"
)

type AuthService struct {
	accountRepo ReadAccountRepository
	encryption  Cryptographer
	jwt         JWTGenerator
}

func NewAuthService(accountRepo ReadAccountRepository, encryption Cryptographer, jwt JWTGenerator) *AuthService {
	return &AuthService{accountRepo: accountRepo, encryption: encryption, jwt: jwt}
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
		return entities.Account{}, fmt.Errorf("account not found: %w", err)
	}
	return account, nil
}
