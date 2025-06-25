package controllers

import (
	"backendForKeenEye/internal/usecases"
	"context"
)

type CreateStudentUsecase interface {
	CreateStudent(context.Context, usecases.CreateStudentRequestDto) (usecases.CreateStudentResponseDto, error)
}

type ReadAllStudentsUsecase interface {
	ReadAllStudents(context.Context) (usecases.ReadAllStudentsResponseDto, error)
}

type ReadStudentUsecase interface {
	ReadStudent(context.Context, usecases.ReadStudentRequestDto) (usecases.ReadStudentResponseDto, error)
}

type UpdateStudentUsecase interface {
	UpdateStudent(context.Context, usecases.UpdateStudentRequestDto) (usecases.UpdateStudentResponseDto, error)
}

type DeleteStudentUsecase interface {
	DeleteStudent(context.Context, usecases.DeleteStudentRequestDto) error
}
