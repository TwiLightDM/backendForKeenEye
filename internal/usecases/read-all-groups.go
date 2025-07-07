package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadAllGroupsUsecase struct {
	GroupRepo ReadAllGroupsRepository
}

type ReadAllGroupsResponseDto struct {
	Groups []entities.Group `json:"groups"`
}

func NewReadAllGroupsUsecase(GroupRepo ReadAllGroupsRepository) ReadAllGroupsUsecase {
	return ReadAllGroupsUsecase{GroupRepo: GroupRepo}
}

func (uc *ReadAllGroupsUsecase) ReadAllGroups(ctx context.Context) (ReadAllGroupsResponseDto, error) {
	var response ReadAllGroupsResponseDto

	groups, err := uc.GroupRepo.Read(ctx)
	if err != nil {
		return response, ReadError
	}

	response = ReadAllGroupsResponseDto{
		Groups: groups,
	}
	return response, nil
}
