package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadAllStudentsByGroupIdUsecase struct {
	StudentRepo ReadAllStudentsByGroupIdRepository
}

type ReadAllStudentsByGroupIdResponseDto struct {
	Students []entities.Student `json:"students"`
}

type ReadAllStudentsByGroupIdRequestDto struct {
	GroupId int
}

func NewReadAllStudentsByGroupIdUsecase(studentRepo ReadAllStudentsByGroupIdRepository) ReadAllStudentsByGroupIdUsecase {
	return ReadAllStudentsByGroupIdUsecase{StudentRepo: studentRepo}
}

func (uc *ReadAllStudentsByGroupIdUsecase) ReadAllStudentsByGroupId(ctx context.Context, request ReadAllStudentsByGroupIdRequestDto) (ReadAllStudentsByGroupIdResponseDto, error) {
	var response ReadAllStudentsByGroupIdResponseDto

	students, err := uc.StudentRepo.ReadByGroupId(ctx, request.GroupId)
	if err != nil {
		return response, ReadError
	}

	response = ReadAllStudentsByGroupIdResponseDto{
		Students: students,
	}
	return response, nil
}
