package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
	"fmt"
)

type CreateAccountUsecase struct {
	accountRepo CreateAccountRepository
	crypto      Cryptographer
	jwt         JWTGenerator
}

type CreateAccountRequestDto struct {
	Login    string
	Password string
	Role     string
}

type CreateAccountResponseDto struct {
	Id           int    `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewCreateAccountUsecase(accountRepo CreateAccountRepository, crypto Cryptographer, jwt JWTGenerator) CreateAccountUsecase {
	return CreateAccountUsecase{accountRepo: accountRepo, crypto: crypto, jwt: jwt}
}

func (uc *CreateAccountUsecase) CreateAccount(ctx context.Context, request CreateAccountRequestDto) (CreateAccountResponseDto, error) {
	var response CreateAccountResponseDto
	hashedPassword, salt, err := uc.crypto.HashPassword(request.Password)
	if err != nil {
		return response, HashPasswordError
	}

	Account := entities.Account{Login: request.Login, Password: hashedPassword, Salt: salt, Role: request.Role}

	_, err = Account.Validate()
	if err != nil {
		return response, ValidationError
	}

	id, err := uc.accountRepo.Create(ctx, Account)
	if err != nil {
		return response, CreateError
	}

	var data = make(map[string]any)
	data["id"] = id

	accessToken, err := uc.jwt.GenerateAccessJWT(data)
	if err != nil {
		return response, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := uc.jwt.GenerateRefreshJWT(data)
	if err != nil {
		return response, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	response = CreateAccountResponseDto{
		Id:           id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}
