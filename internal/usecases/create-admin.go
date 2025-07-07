package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type CreateAdminUsecase struct {
	AdminRepo CreateAdminRepository
}

type CreateAdminRequestDto struct {
	Fio         string
	PhoneNumber string
	AccountId   int
}

type CreateAdminResponseDto struct {
	Id int `json:"id"`
}

func NewCreateAdminUsecase(AdminRepo CreateAdminRepository) CreateAdminUsecase {
	return CreateAdminUsecase{AdminRepo: AdminRepo}
}

func (uc *CreateAdminUsecase) CreateAdmin(ctx context.Context, request CreateAdminRequestDto) (CreateAdminResponseDto, error) {
	var response CreateAdminResponseDto
	student := entities.Admin{Fio: request.Fio, PhoneNumber: request.PhoneNumber, AccountId: request.AccountId}

	id, err := uc.AdminRepo.Create(ctx, student)
	if err != nil {
		return response, CreateError
	}

	response = CreateAdminResponseDto{
		Id: id,
	}
	return response, nil
}
