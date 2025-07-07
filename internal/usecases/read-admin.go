package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadAdminUsecase struct {
	AdminRepo ReadAdminRepository
}

type ReadAdminRequestDto struct {
	Id int
}

type ReadAdminResponseDto struct {
	Admin entities.Admin `json:"admin"`
}

func NewReadAdminUsecase(AdminRepo ReadAdminRepository) ReadAdminUsecase {
	return ReadAdminUsecase{AdminRepo: AdminRepo}
}

func (uc *ReadAdminUsecase) ReadAdmin(ctx context.Context, request ReadAdminRequestDto) (ReadAdminResponseDto, error) {
	var response ReadAdminResponseDto

	admin, err := uc.AdminRepo.ReadById(ctx, request.Id)
	if err != nil {
		return response, ReadError
	}

	response = ReadAdminResponseDto{
		Admin: admin,
	}
	return response, nil
}
