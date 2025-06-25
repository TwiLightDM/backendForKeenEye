package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
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

func (repo *StudentRepository) Create(ctx context.Context, Student entities.Student) (int, error) {

	sql, args, err := repo.builder.
		Insert("students").
		Columns("fio", "group_name", "phone_number").
		Values(Student.Fio, Student.GroupName, Student.PhoneNumber).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, SqlStatementError
	}

	var newID int
	err = repo.pool.QueryRow(ctx, sql, args...).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert Student: %w", err)
	}

	return newID, nil
}

func (repo *StudentRepository) Read(ctx context.Context) ([]entities.Student, error) {

	sql, args, err := repo.builder.
		Select("id, fio, group_name, phone_number").
		From("students").
		ToSql()

	if err != nil {
		return nil, SqlStatementError
	}

	rows, err := repo.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query students: %w", err)
	}
	defer rows.Close()

	var students []entities.Student
	for rows.Next() {
		var student entities.Student
		err = rows.Scan(
			&student.Id,
			&student.Fio,
			&student.GroupName,
			&student.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan Student: %w", err)
		}

		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return students, nil
}

func (repo *StudentRepository) ReadById(ctx context.Context, id int) (entities.Student, error) {
	var fio, group, phoneNumber string

	sql, args, err := repo.builder.
		Select("fio", "group_name", "phone_number").
		From("students").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Student{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&group,
		&phoneNumber,
	)
	if err != nil {
		return entities.Student{}, fmt.Errorf("failed to read Student: %w", err)
	}

	return entities.Student{Id: id, Fio: fio, GroupName: group, PhoneNumber: phoneNumber}, nil
}

func (repo *StudentRepository) Update(ctx context.Context, id int, updates map[string]any) (entities.Student, error) {
	sql, args, err := repo.builder.
		Update("students").
		Where(squirrel.Eq{"id": id}).
		SetMap(updates).
		Suffix("RETURNING id, fio, group_name, phone_number").
		ToSql()

	if err != nil {
		return entities.Student{}, SqlStatementError
	}

	var student entities.Student
	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&student.Id,
		&student.Fio,
		&student.GroupName,
		&student.PhoneNumber,
	)

	if err != nil {
		return entities.Student{}, fmt.Errorf("failed to update Student: %w", err)
	}

	return student, nil
}

func (repo *StudentRepository) DeleteById(ctx context.Context, id int) error {
	sql, args, err := repo.builder.
		Delete("students").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return SqlStatementError
	}

	_, err = repo.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete Student: %w", err)
	}

	return nil
}
