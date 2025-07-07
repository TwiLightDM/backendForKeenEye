package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type UpdateGroupUsecase struct {
	groupRepo UpdateGroupRepository
}

type UpdateGroupRequestDto struct {
	Id        int
	Name      string
	TeacherId int
}

type UpdateGroupResponseDto struct {
	Group entities.Group `json:"group"`
}

func NewUpdateGroupUsecase(GroupRepo UpdateGroupRepository) UpdateGroupUsecase {
	return UpdateGroupUsecase{groupRepo: GroupRepo}
}

func (uc *UpdateGroupUsecase) UpdateGroup(ctx context.Context, request UpdateGroupRequestDto) (UpdateGroupResponseDto, error) {
	var response UpdateGroupResponseDto
	updates := make(map[string]any)

	if request.Id == 0 {
		return response, MissingIdError
	}
	if request.Name != "" {
		updates["name"] = request.Name
	}
	if request.TeacherId != 0 {
		updates["teacher_id"] = request.TeacherId
	}
	if len(updates) == 0 {
		return response, NoFieldsError
	}

	group, err := uc.groupRepo.Update(ctx, request.Id, updates)
	if err != nil {
		return response, UpdateError
	}
	response = UpdateGroupResponseDto{
		Group: group,
	}
	return response, nil
}
