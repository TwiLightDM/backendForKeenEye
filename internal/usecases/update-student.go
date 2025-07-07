package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type UpdateStudentUsecase struct {
	studentRepo UpdateStudentRepository
}

type UpdateStudentRequestDto struct {
	Id          int
	Fio         string
	GroupName   string
	PhoneNumber string
	GroupId     int
	AccountId   int
}

type UpdateStudentResponseDto struct {
	Student entities.Student `json:"student"`
}

func NewUpdateStudentUsecase(StudentRepo UpdateStudentRepository) UpdateStudentUsecase {
	return UpdateStudentUsecase{studentRepo: StudentRepo}
}

func (uc *UpdateStudentUsecase) UpdateStudent(ctx context.Context, request UpdateStudentRequestDto) (UpdateStudentResponseDto, error) {
	var response UpdateStudentResponseDto
	updates := make(map[string]any)

	if request.Id == 0 {
		return response, MissingIdError
	}
	if request.Fio != "" {
		updates["fio"] = request.Fio
	}
	if request.GroupName != "" {
		updates["group_name"] = request.GroupName
	}
	if request.PhoneNumber != "" {
		updates["phone_number"] = request.PhoneNumber
	}
	if request.GroupId != 0 {
		updates["group_id"] = request.GroupId
	}
	if request.AccountId != 0 {
		updates["account_id"] = request.AccountId
	}
	if len(updates) == 0 {
		return response, NoFieldsError
	}

	student, err := uc.studentRepo.Update(ctx, request.Id, updates)
	if err != nil {
		return response, UpdateError
	}
	response = UpdateStudentResponseDto{
		Student: student,
	}
	return response, nil
}
