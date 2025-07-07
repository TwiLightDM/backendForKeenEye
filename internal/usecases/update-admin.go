package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type UpdateAdminUsecase struct {
	adminRepo UpdateAdminRepository
}

type UpdateAdminRequestDto struct {
	Id          int
	Fio         string
	PhoneNumber string
	AccountId   int
}

type UpdateAdminResponseDto struct {
	Admin entities.Admin `json:"admin"`
}

func NewUpdateAdminUsecase(AdminRepo UpdateAdminRepository) UpdateAdminUsecase {
	return UpdateAdminUsecase{adminRepo: AdminRepo}
}

func (uc *UpdateAdminUsecase) UpdateAdmin(ctx context.Context, request UpdateAdminRequestDto) (UpdateAdminResponseDto, error) {
	var response UpdateAdminResponseDto
	updates := make(map[string]any)

	if request.Id == 0 {
		return response, MissingIdError
	}
	if request.Fio != "" {
		updates["fio"] = request.Fio
	}
	if request.PhoneNumber != "" {
		updates["phone_number"] = request.PhoneNumber
	}
	if request.AccountId != 0 {
		updates["account_id"] = request.AccountId
	}
	if len(updates) == 0 {
		return response, NoFieldsError
	}

	admin, err := uc.adminRepo.Update(ctx, request.Id, updates)
	if err != nil {
		return response, UpdateError
	}
	response = UpdateAdminResponseDto{
		Admin: admin,
	}
	return response, nil
}
