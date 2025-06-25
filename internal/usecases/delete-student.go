package usecases

import (
	"context"
	"fmt"
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

	err := uc.StudentRepo.DeleteById(ctx, request.Id)
	if err != nil {
		return fmt.Errorf("failed to delete Student record: %w", err)
	}

	return nil
}
