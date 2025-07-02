package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type CreateAccountUsecase struct {
	accountRepo CreateAccountRepository
	crypto      Cryptographer
}

type CreateAccountRequestDto struct {
	Login    string
	Password string
}

type CreateAccountResponseDto struct {
	Id int `json:"id"`
}

func NewCreateAccountUsecase(accountRepo CreateAccountRepository, crypto Cryptographer) CreateAccountUsecase {
	return CreateAccountUsecase{accountRepo: accountRepo, crypto: crypto}
}

func (uc *CreateAccountUsecase) CreateAccount(ctx context.Context, request CreateAccountRequestDto) (CreateAccountResponseDto, error) {
	var response CreateAccountResponseDto
	hashedPassword, salt, err := uc.crypto.HashPassword(request.Password)
	if err != nil {
		return response, HashPasswordError
	}

	Account := entities.Account{Login: request.Login, Password: hashedPassword, Salt: salt}

	id, err := uc.accountRepo.Create(ctx, Account)
	if err != nil {
		return response, CreateError
	}

	response = CreateAccountResponseDto{
		Id: id,
	}
	return response, nil
}
