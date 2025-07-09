package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeacherRepository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewTeacherRepository(pool *pgxpool.Pool, builder squirrel.StatementBuilderType) *TeacherRepository {
	return &TeacherRepository{pool: pool, builder: builder}
}

func (repo *TeacherRepository) Read(ctx context.Context) ([]entities.Teacher, error) {
	var id sql.NullInt32
	var fio, phoneNumber sql.NullString
	sql, args, err := repo.builder.
		Select("id, fio, phone_number").
		From("teachers").
		Where(squirrel.Eq{"is_deleted": false}).
		ToSql()

	if err != nil {
		return nil, SqlStatementError
	}

	rows, err := repo.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, SqlReadError
	}
	defer rows.Close()

	var teachers []entities.Teacher
	for rows.Next() {
		err = rows.Scan(
			&id,
			&fio,
			&phoneNumber,
		)
		if err != nil {
			return nil, SqlScanError
		}

		teacher := entities.Teacher{
			Id:          validateInt(id),
			Fio:         validateString(fio),
			PhoneNumber: validateString(phoneNumber),
		}
		teachers = append(teachers, teacher)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return teachers, nil
}

func (repo *TeacherRepository) ReadById(ctx context.Context, id int) (entities.Teacher, error) {
	var fio, phoneNumber sql.NullString

	sql, args, err := repo.builder.
		Select("fio", "phone_number").
		From("teachers").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Teacher{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&phoneNumber,
	)
	if err != nil {
		return entities.Teacher{}, SqlReadError
	}

	return entities.Teacher{Id: id, Fio: validateString(fio), PhoneNumber: validateString(phoneNumber)}, nil
}

func (repo *TeacherRepository) Update(ctx context.Context, id int, updates map[string]any) (entities.Teacher, error) {
	var fio, phoneNumber sql.NullString
	sql, args, err := repo.builder.
		Update("teachers").
		Where(squirrel.Eq{"id": id}).
		SetMap(updates).
		Suffix("RETURNING fio, phone_number").
		ToSql()

	if err != nil {
		return entities.Teacher{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&phoneNumber,
	)

	if err != nil {
		fmt.Println(err)
		return entities.Teacher{}, SqlUpdateError
	}

	return entities.Teacher{
		Id:          id,
		Fio:         validateString(fio),
		PhoneNumber: validateString(phoneNumber),
	}, nil
}

func (repo *TeacherRepository) SoftDelete(ctx context.Context, id int) error {
	sql, args, err := repo.builder.
		Update("teachers").
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
