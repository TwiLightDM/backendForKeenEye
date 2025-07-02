package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadStudentUsecase struct {
	StudentRepo ReadStudentRepository
}

type ReadStudentRequestDto struct {
	Id int
}

type ReadStudentResponseDto struct {
	Student entities.Student `json:"student"`
}

func NewReadStudentUsecase(StudentRepo ReadStudentRepository) ReadStudentUsecase {
	return ReadStudentUsecase{StudentRepo: StudentRepo}
}

func (uc *ReadStudentUsecase) ReadStudent(ctx context.Context, request ReadStudentRequestDto) (ReadStudentResponseDto, error) {
	var response ReadStudentResponseDto

	student, err := uc.StudentRepo.ReadById(ctx, request.Id)
	if err != nil {
		return response, ReadError
	}

	response = ReadStudentResponseDto{
		Student: student,
	}
	return response, nil
}
