package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewStudentRepository(pool *pgxpool.Pool, builder squirrel.StatementBuilderType) *StudentRepository {
	return &StudentRepository{pool: pool, builder: builder}
}

func (repo *StudentRepository) Read(ctx context.Context) ([]entities.Student, error) {
	var id int
	var fio, phoneNumber sql.NullString
	var groupId sql.NullInt32
	sql, args, err := repo.builder.
		Select("id", "fio", "phone_number", "group_id").
		From("students").
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

	var students []entities.Student
	for rows.Next() {
		err = rows.Scan(
			&id,
			&fio,
			&phoneNumber,
			&groupId,
		)
		if err != nil {
			return nil, SqlScanError
		}

		var student = entities.Student{
			Id:          id,
			Fio:         validateString(fio),
			PhoneNumber: validateString(phoneNumber),
			GroupId:     validateInt(groupId),
		}

		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return students, nil
}

func (repo *StudentRepository) ReadById(ctx context.Context, id int) (entities.Student, error) {
	var fio, phoneNumber sql.NullString
	var groupId sql.NullInt32

	sql, args, err := repo.builder.
		Select("fio", "phone_number", "group_id").
		From("students").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Student{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&phoneNumber,
		&groupId,
	)
	if err != nil {
		return entities.Student{}, SqlReadError
	}

	return entities.Student{
		Id:          id,
		Fio:         validateString(fio),
		PhoneNumber: validateString(phoneNumber),
		GroupId:     validateInt(groupId),
	}, nil
}

func (repo *StudentRepository) ReadByGroupId(ctx context.Context, groupId int) ([]entities.Student, error) {
	var id int
	var fio, phoneNumber sql.NullString

	sql, args, err := repo.builder.
		Select("id, fio, phone_number").
		From("students").
		Where(squirrel.Eq{"group_id": groupId}).
		ToSql()

	if err != nil {
		return nil, SqlStatementError
	}

	rows, err := repo.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, SqlReadError
	}
	defer rows.Close()

	var students []entities.Student
	for rows.Next() {

		err = rows.Scan(&id, &fio, &phoneNumber)
		if err != nil {
			return nil, SqlScanError
		}

		students = append(students, entities.Student{
			Id:          id,
			Fio:         validateString(fio),
			PhoneNumber: validateString(phoneNumber),
			GroupId:     groupId,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return students, nil
}

func (repo *StudentRepository) Update(ctx context.Context, id int, updates map[string]any) (entities.Student, error) {
	var fio, phoneNumber sql.NullString
	var groupId sql.NullInt32

	sql, args, err := repo.builder.
		Update("students").
		Where(squirrel.Eq{"id": id}).
		SetMap(updates).
		Suffix("RETURNING fio, phone_number, group_id").
		ToSql()

	if err != nil {
		return entities.Student{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&phoneNumber,
		&groupId,
	)

	if err != nil {
		return entities.Student{}, SqlUpdateError
	}

	return entities.Student{
		Id:          id,
		Fio:         validateString(fio),
		PhoneNumber: validateString(phoneNumber),
		GroupId:     validateInt(groupId),
	}, nil
}

func (repo *StudentRepository) SoftDelete(ctx context.Context, id int) error {
	sql, args, err := repo.builder.
		Update("students").
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
