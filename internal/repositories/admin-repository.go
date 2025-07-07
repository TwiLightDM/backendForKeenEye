package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewAdminRepository(pool *pgxpool.Pool, builder squirrel.StatementBuilderType) *AdminRepository {
	return &AdminRepository{pool: pool, builder: builder}
}

func (repo *AdminRepository) Create(ctx context.Context, admin entities.Admin) (int, error) {

	sql, args, err := repo.builder.
		Insert("admins").
		Columns("fio", "phone_number", "account_id").
		Values(admin.Fio, admin.PhoneNumber, admin.AccountId).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, SqlStatementError
	}

	var newID int
	err = repo.pool.QueryRow(ctx, sql, args...).Scan(&newID)
	if err != nil {
		return 0, SqlInsertError
	}

	return newID, nil
}

func (repo *AdminRepository) ReadById(ctx context.Context, id int) (entities.Admin, error) {
	var fio, group, phoneNumber string
	var accountId int

	sql, args, err := repo.builder.
		Select("fio", "group_name", "phone_number", "account_id").
		From("admins").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Admin{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&group,
		&phoneNumber,
		&accountId,
	)
	if err != nil {
		return entities.Admin{}, SqlReadError
	}

	return entities.Admin{Id: id, Fio: fio, PhoneNumber: phoneNumber, AccountId: accountId}, nil
}

func (repo *AdminRepository) ReadByAccountId(ctx context.Context, accountId int) (entities.Admin, error) {
	var fio, group, phoneNumber string
	var id int

	sql, args, err := repo.builder.
		Select("id", "fio", "group_name", "phone_number").
		From("admins").
		Where(squirrel.Eq{"account_id": accountId}).
		ToSql()

	if err != nil {
		return entities.Admin{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&id,
		&fio,
		&group,
		&phoneNumber,
	)
	if err != nil {
		return entities.Admin{}, SqlReadError
	}

	return entities.Admin{Id: id, Fio: fio, PhoneNumber: phoneNumber, AccountId: accountId}, nil
}

func (repo *AdminRepository) Update(ctx context.Context, id int, updates map[string]any) (entities.Admin, error) {
	sql, args, err := repo.builder.
		Update("admins").
		Where(squirrel.Eq{"id": id}).
		SetMap(updates).
		Suffix("RETURNING id, fio, phone_number, account_id").
		ToSql()

	if err != nil {
		return entities.Admin{}, SqlStatementError
	}

	var student entities.Admin
	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&student.Id,
		&student.Fio,
		&student.PhoneNumber,
		&student.AccountId,
	)

	if err != nil {
		return entities.Admin{}, SqlUpdateError
	}

	return student, nil
}

func (repo *AdminRepository) Delete(ctx context.Context, id int) error {
	sql, args, err := repo.builder.
		Delete("admins").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return SqlStatementError
	}

	_, err = repo.pool.Exec(ctx, sql, args...)
	if err != nil {
		return SqlDeleteError
	}

	return nil
}
