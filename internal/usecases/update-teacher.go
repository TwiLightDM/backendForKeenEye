package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type UpdateTeacherUsecase struct {
	teacherRepo UpdateTeacherRepository
}

type UpdateTeacherRequestDto struct {
	Id          int
	Fio         string
	PhoneNumber string
	AccountId   int
}

type UpdateTeacherResponseDto struct {
	Teacher entities.Teacher `json:"teacher"`
}

func NewUpdateTeacherUsecase(TeacherRepo UpdateTeacherRepository) UpdateTeacherUsecase {
	return UpdateTeacherUsecase{teacherRepo: TeacherRepo}
}

func (uc *UpdateTeacherUsecase) UpdateTeacher(ctx context.Context, request UpdateTeacherRequestDto) (UpdateTeacherResponseDto, error) {
	var response UpdateTeacherResponseDto
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

	teacher, err := uc.teacherRepo.Update(ctx, request.Id, updates)
	if err != nil {
		return response, UpdateError
	}
	response = UpdateTeacherResponseDto{
		Teacher: teacher,
	}
	return response, nil
}
