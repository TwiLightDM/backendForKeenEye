package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
	"fmt"
)

type CreateUserUsecase struct {
	userRepo CreateUserRepository
	crypto   Cryptographer
	jwt      JWTGenerator
}

type CreateUserRequestDto struct {
	Login    string
	Password string
	Role     string
}

type CreateUserResponseDto struct {
	Id           int    `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewCreateUserUsecase(userRepo CreateUserRepository, crypto Cryptographer, jwt JWTGenerator) CreateUserUsecase {
	return CreateUserUsecase{userRepo: userRepo, crypto: crypto, jwt: jwt}
}

func (uc *CreateUserUsecase) CreateUser(ctx context.Context, request CreateUserRequestDto) (CreateUserResponseDto, error) {
	var response CreateUserResponseDto
	hashedPassword, salt, err := uc.crypto.HashPassword(request.Password)
	if err != nil {
		return response, HashPasswordError
	}

	user := entities.User{Login: request.Login, Password: hashedPassword, Salt: salt, Role: request.Role}

	_, err = user.Validate()
	if err != nil {
		return response, ValidationError
	}

	id, err := uc.userRepo.Create(ctx, user)
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

	response = CreateUserResponseDto{
		Id:           id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}
