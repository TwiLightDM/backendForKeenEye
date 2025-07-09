package usecases

import (
	"backendForKeenEye/internal/entities"
	"context"
)

type Cryptographer interface {
	HashPassword(password string) (string, string, error)
	PasswordComparison(password, hashedPassword, salt string) (bool, error)
}

type JWTGenerator interface {
	GenerateAccessJWT(data map[string]any) (string, error)
	GenerateRefreshJWT(data map[string]any) (string, error)
	ParseJWT(tokenString string) (map[string]any, error)
}

type ReadAllStudentsRepository interface {
	Read(ctx context.Context) ([]entities.Student, error)
}

type ReadAllStudentsByGroupIdRepository interface {
	ReadByGroupId(ctx context.Context, groupId int) ([]entities.Student, error)
}

type ReadStudentRepository interface {
	ReadById(ctx context.Context, id int) (entities.Student, error)
}

type UpdateStudentRepository interface {
	Update(ctx context.Context, id int, updates map[string]any) (entities.Student, error)
}

type DeleteStudentRepository interface {
	SoftDelete(ctx context.Context, id int) error
}

type CreateUserRepository interface {
	Create(ctx context.Context, student entities.User) (int, error)
}

type ReadUserRepository interface {
	ReadByLogin(ctx context.Context, login string) (entities.User, error)
	ReadById(ctx context.Context, id int) (entities.User, error)
}

type ReadAllTeachersRepository interface {
	Read(ctx context.Context) ([]entities.Teacher, error)
}

type ReadTeacherRepository interface {
	ReadById(ctx context.Context, id int) (entities.Teacher, error)
}

type UpdateTeacherRepository interface {
	Update(ctx context.Context, id int, updates map[string]any) (entities.Teacher, error)
}

type DeleteTeacherRepository interface {
	SoftDelete(ctx context.Context, id int) error
}

type CreateGroupRepository interface {
	Create(ctx context.Context, teacher entities.Group) (int, error)
}

type ReadAllGroupsRepository interface {
	Read(ctx context.Context) ([]entities.Group, error)
}

type ReadGroupRepository interface {
	ReadById(ctx context.Context, id int) (entities.Group, error)
}

type UpdateGroupRepository interface {
	Update(ctx context.Context, id int, updates map[string]any) (entities.Group, error)
}

type DeleteGroupRepository interface {
	SoftDelete(ctx context.Context, id int) error
}

type ReadAdminRepository interface {
	ReadById(ctx context.Context, id int) (entities.Admin, error)
}

type UpdateAdminRepository interface {
	Update(ctx context.Context, id int, updates map[string]any) (entities.Admin, error)
}

type DeleteAdminRepository interface {
	SoftDelete(ctx context.Context, id int) error
}
