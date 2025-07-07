package usecases

import (
	"context"
)

type DeleteTeacherUsecase struct {
	TeacherRepo DeleteTeacherRepository
}

type DeleteTeacherRequestDto struct {
	Id int
}

func NewDeleteTeacherUsecase(TeacherRepo DeleteTeacherRepository) DeleteTeacherUsecase {
	return DeleteTeacherUsecase{TeacherRepo: TeacherRepo}
}

func (uc *DeleteTeacherUsecase) DeleteTeacher(ctx context.Context, request DeleteTeacherRequestDto) error {

	err := uc.TeacherRepo.SoftDelete(ctx, request.Id)
	if err != nil {
		return DeleteError
	}

	return nil
}
