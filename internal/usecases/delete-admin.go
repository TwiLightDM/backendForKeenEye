package usecases

import (
	"context"
)

type DeleteAdminUsecase struct {
	AdminRepo DeleteAdminRepository
}

type DeleteAdminRequestDto struct {
	Id int
}

func NewDeleteAdminUsecase(AdminRepo DeleteAdminRepository) DeleteAdminUsecase {
	return DeleteAdminUsecase{AdminRepo: AdminRepo}
}

func (uc *DeleteAdminUsecase) DeleteAdmin(ctx context.Context, request DeleteAdminRequestDto) error {

	err := uc.AdminRepo.SoftDelete(ctx, request.Id)
	if err != nil {
		return DeleteError
	}

	return nil
}
