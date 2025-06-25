package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
	"fmt"
)

type CreateStudentUsecase struct {
	StudentRepo CreateStudentRepository
}

type CreateStudentRequestDto struct {
	Fio         string
	GroupName   string
	PhoneNumber string
}

type CreateStudentResponseDto struct {
	Id int `json:"id"`
}

func NewCreateStudentUsecase(StudentRepo CreateStudentRepository) CreateStudentUsecase {
	return CreateStudentUsecase{StudentRepo: StudentRepo}
}

func (uc *CreateStudentUsecase) CreateStudent(ctx context.Context, request CreateStudentRequestDto) (CreateStudentResponseDto, error) {
	var response CreateStudentResponseDto
	student := entities.Student{Fio: request.Fio, GroupName: request.GroupName, PhoneNumber: request.PhoneNumber}

	id, err := uc.StudentRepo.Create(ctx, student)
	if err != nil {
		return response, fmt.Errorf("failed to create Student record: %w", err)
	}

	response = CreateStudentResponseDto{
		Id: id,
	}
	return response, nil
}
