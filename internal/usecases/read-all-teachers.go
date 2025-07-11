package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadAllTeachersUsecase struct {
	TeacherRepo ReadAllTeachersRepository
}

type ReadAllTeachersResponseDto struct {
	Teachers []entities.Teacher `json:"teachers"`
}

func NewReadAllTeachersUsecase(TeacherRepo ReadAllTeachersRepository) ReadAllTeachersUsecase {
	return ReadAllTeachersUsecase{TeacherRepo: TeacherRepo}
}

func (uc *ReadAllTeachersUsecase) ReadAllTeachers(ctx context.Context) (ReadAllTeachersResponseDto, error) {
	var response ReadAllTeachersResponseDto

	teachers, err := uc.TeacherRepo.Read(ctx)
	if err != nil {
		return response, ReadError
	}

	response = ReadAllTeachersResponseDto{
		Teachers: teachers,
	}
	return response, nil
}
