package usecases

import (
	"context"
)

type DeleteGroupUsecase struct {
	GroupRepo DeleteGroupRepository
}

type DeleteGroupRequestDto struct {
	Id int
}

func NewDeleteGroupUsecase(GroupRepo DeleteGroupRepository) DeleteGroupUsecase {
	return DeleteGroupUsecase{GroupRepo: GroupRepo}
}

func (uc *DeleteGroupUsecase) DeleteGroup(ctx context.Context, request DeleteGroupRequestDto) error {

	err := uc.GroupRepo.SoftDelete(ctx, request.Id)
	if err != nil {
		return DeleteError
	}

	return nil
}
