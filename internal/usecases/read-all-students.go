package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadAllStudentsUsecase struct {
	StudentRepo ReadAllStudentsRepository
}

type ReadAllStudentsResponseDto struct {
	Students []entities.Student `json:"students"`
}

func NewReadAllStudentsUsecase(StudentRepo ReadAllStudentsRepository) ReadAllStudentsUsecase {
	return ReadAllStudentsUsecase{StudentRepo: StudentRepo}
}

func (uc *ReadAllStudentsUsecase) ReadAllStudents(ctx context.Context) (ReadAllStudentsResponseDto, error) {
	var response ReadAllStudentsResponseDto

	students, err := uc.StudentRepo.Read(ctx)
	if err != nil {
		return response, ReadError
	}

	response = ReadAllStudentsResponseDto{
		Students: students,
	}
	return response, nil
}
