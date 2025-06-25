package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type CreateStudentRepository interface {
	Create(ctx context.Context, student entities.Student) (int, error)
}

type ReadAllStudentsRepository interface {
	Read(ctx context.Context) ([]entities.Student, error)
}

type ReadStudentRepository interface {
	ReadById(ctx context.Context, id int) (entities.Student, error)
}

type UpdateStudentRepository interface {
	Update(ctx context.Context, id int, updates map[string]any) (entities.Student, error)
}

type DeleteStudentRepository interface {
	DeleteById(ctx context.Context, id int) error
}
