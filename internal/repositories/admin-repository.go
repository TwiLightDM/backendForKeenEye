package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
	"database/sql"
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

func (repo *AdminRepository) ReadById(ctx context.Context, id int) (entities.Admin, error) {
	var fio, phoneNumber sql.NullString

	sql, args, err := repo.builder.
		Select("fio", "phone_number").
		From("admins").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Admin{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&phoneNumber,
	)
	if err != nil {
		return entities.Admin{}, SqlReadError
	}

	return entities.Admin{Id: id, Fio: validateString(fio), PhoneNumber: validateString(phoneNumber)}, nil
}

func (repo *AdminRepository) Update(ctx context.Context, id int, updates map[string]any) (entities.Admin, error) {
	var fio, phoneNumber sql.NullString
	sql, args, err := repo.builder.
		Update("admins").
		Where(squirrel.Eq{"id": id}).
		SetMap(updates).
		Suffix("RETURNING fio, phone_number").
		ToSql()

	if err != nil {
		return entities.Admin{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&phoneNumber,
	)

	if err != nil {
		return entities.Admin{}, SqlUpdateError
	}

	return entities.Admin{
		Id:          id,
		Fio:         validateString(fio),
		PhoneNumber: validateString(phoneNumber),
	}, nil
}

func (repo *AdminRepository) SoftDelete(ctx context.Context, id int) error {
	sql, args, err := repo.builder.
		Update("admins").
		Where(squirrel.Eq{"id": id}).
		Set("is_deleted", true).
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
