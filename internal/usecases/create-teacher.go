package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type CreateTeacherUsecase struct {
	TeacherRepo CreateTeacherRepository
}

type CreateTeacherRequestDto struct {
	Fio         string
	PhoneNumber string
	AccountId   int
}

type CreateTeacherResponseDto struct {
	Id int `json:"id"`
}

func NewCreateTeacherUsecase(TeacherRepo CreateTeacherRepository) CreateTeacherUsecase {
	return CreateTeacherUsecase{TeacherRepo: TeacherRepo}
}

func (uc *CreateTeacherUsecase) CreateTeacher(ctx context.Context, request CreateTeacherRequestDto) (CreateTeacherResponseDto, error) {
	var response CreateTeacherResponseDto
	student := entities.Teacher{Fio: request.Fio, PhoneNumber: request.PhoneNumber, AccountId: request.AccountId}

	id, err := uc.TeacherRepo.Create(ctx, student)
	if err != nil {
		return response, CreateError
	}

	response = CreateTeacherResponseDto{
		Id: id,
	}
	return response, nil
}
