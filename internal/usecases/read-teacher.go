package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type ReadTeacherUsecase struct {
	TeacherRepo ReadTeacherRepository
}

type ReadTeacherRequestDto struct {
	Id int
}

type ReadTeacherResponseDto struct {
	Teacher entities.Teacher `json:"teacher"`
}

func NewReadTeacherUsecase(TeacherRepo ReadTeacherRepository) ReadTeacherUsecase {
	return ReadTeacherUsecase{TeacherRepo: TeacherRepo}
}

func (uc *ReadTeacherUsecase) ReadTeacher(ctx context.Context, request ReadTeacherRequestDto) (ReadTeacherResponseDto, error) {
	var response ReadTeacherResponseDto

	teacher, err := uc.TeacherRepo.ReadById(ctx, request.Id)
	if err != nil {
		return response, ReadError
	}

	response = ReadTeacherResponseDto{
		Teacher: teacher,
	}
	return response, nil
}
