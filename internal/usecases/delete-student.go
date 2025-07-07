package usecases

import (
	"context"
)

type DeleteStudentUsecase struct {
	StudentRepo DeleteStudentRepository
}

type DeleteStudentRequestDto struct {
	Id int
}

func NewDeleteStudentUsecase(StudentRepo DeleteStudentRepository) DeleteStudentUsecase {
	return DeleteStudentUsecase{StudentRepo: StudentRepo}
}

func (uc *DeleteStudentUsecase) DeleteStudent(ctx context.Context, request DeleteStudentRequestDto) error {

	err := uc.StudentRepo.Delete(ctx, request.Id)
	if err != nil {
		return DeleteError
	}

	return nil
}
