package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type CreateGroupUsecase struct {
	GroupRepo CreateGroupRepository
}

type CreateGroupRequestDto struct {
	Name      string
	TeacherId int
}

type CreateGroupResponseDto struct {
	Id int `json:"id"`
}

func NewCreateGroupUsecase(GroupRepo CreateGroupRepository) CreateGroupUsecase {
	return CreateGroupUsecase{GroupRepo: GroupRepo}
}

func (uc *CreateGroupUsecase) CreateGroup(ctx context.Context, request CreateGroupRequestDto) (CreateGroupResponseDto, error) {
	var response CreateGroupResponseDto
	student := entities.Group{Name: request.Name, TeacherId: request.TeacherId}

	id, err := uc.GroupRepo.Create(ctx, student)
	if err != nil {
		return response, CreateError
	}

	response = CreateGroupResponseDto{
		Id: id,
	}
	return response, nil
}
