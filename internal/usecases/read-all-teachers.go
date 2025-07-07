package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadAllTeachersUsecase struct {
	TeacherRepo ReadAllTeachersRepository
}

type ReadAllTeachersResponseDto struct {
	Teachers []entities.Teacher `json:"students"`
}

func NewReadAllTeachersUsecase(TeacherRepo ReadAllTeachersRepository) ReadAllTeachersUsecase {
	return ReadAllTeachersUsecase{TeacherRepo: TeacherRepo}
}

func (uc *ReadAllTeachersUsecase) ReadAllTeachers(ctx context.Context) (ReadAllTeachersResponseDto, error) {
	var response ReadAllTeachersResponseDto

	students, err := uc.TeacherRepo.Read(ctx)
	if err != nil {
		return response, ReadError
	}

	response = ReadAllTeachersResponseDto{
		Teachers: students,
	}
	return response, nil
}
