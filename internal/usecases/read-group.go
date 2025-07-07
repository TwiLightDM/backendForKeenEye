package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadGroupUsecase struct {
	GroupRepo ReadGroupRepository
}

type ReadGroupRequestDto struct {
	Id int
}

type ReadGroupResponseDto struct {
	Group entities.Group `json:"group"`
}

func NewReadGroupUsecase(GroupRepo ReadGroupRepository) ReadGroupUsecase {
	return ReadGroupUsecase{GroupRepo: GroupRepo}
}

func (uc *ReadGroupUsecase) ReadGroup(ctx context.Context, request ReadGroupRequestDto) (ReadGroupResponseDto, error) {
	var response ReadGroupResponseDto

	group, err := uc.GroupRepo.ReadById(ctx, request.Id)
	if err != nil {
		return response, ReadError
	}

	response = ReadGroupResponseDto{
		Group: group,
	}
	return response, nil
}
