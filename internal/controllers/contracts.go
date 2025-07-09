package controllers

import (
	"backendForKeenEye/internal/usecases"
	"context"
)

type CreateUserUsecase interface {
	CreateUser(context.Context, usecases.CreateUserRequestDto) (usecases.CreateUserResponseDto, error)
}
type ReadAllStudentsUsecase interface {
	ReadAllStudents(context.Context) (usecases.ReadAllStudentsResponseDto, error)
}

type ReadAllStudentsByGroupIdUsecase interface {
	ReadAllStudentsByGroupId(context.Context, usecases.ReadAllStudentsByGroupIdRequestDto) (usecases.ReadAllStudentsByGroupIdResponseDto, error)
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

type ReadAllTeachersUsecase interface {
	ReadAllTeachers(context.Context) (usecases.ReadAllTeachersResponseDto, error)
}

type ReadTeacherUsecase interface {
	ReadTeacher(context.Context, usecases.ReadTeacherRequestDto) (usecases.ReadTeacherResponseDto, error)
}

type UpdateTeacherUsecase interface {
	UpdateTeacher(context.Context, usecases.UpdateTeacherRequestDto) (usecases.UpdateTeacherResponseDto, error)
}

type DeleteTeacherUsecase interface {
	DeleteTeacher(context.Context, usecases.DeleteTeacherRequestDto) error
}

type ReadAdminUsecase interface {
	ReadAdmin(context.Context, usecases.ReadAdminRequestDto) (usecases.ReadAdminResponseDto, error)
}

type UpdateAdminUsecase interface {
	UpdateAdmin(context.Context, usecases.UpdateAdminRequestDto) (usecases.UpdateAdminResponseDto, error)
}

type DeleteAdminUsecase interface {
	DeleteAdmin(context.Context, usecases.DeleteAdminRequestDto) error
}

type CreateGroupUsecase interface {
	CreateGroup(context.Context, usecases.CreateGroupRequestDto) (usecases.CreateGroupResponseDto, error)
}

type ReadAllGroupsUsecase interface {
	ReadAllGroups(context.Context) (usecases.ReadAllGroupsResponseDto, error)
}

type ReadGroupUsecase interface {
	ReadGroup(context.Context, usecases.ReadGroupRequestDto) (usecases.ReadGroupResponseDto, error)
}

type UpdateGroupUsecase interface {
	UpdateGroup(context.Context, usecases.UpdateGroupRequestDto) (usecases.UpdateGroupResponseDto, error)
}

type DeleteGroupUsecase interface {
	DeleteGroup(context.Context, usecases.DeleteGroupRequestDto) error
}
