package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
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

func (repo *TeacherRepository) Create(ctx context.Context, teacher entities.Teacher) (int, error) {

	sql, args, err := repo.builder.
		Insert("teachers").
		Columns("fio", "phone_number", "account_id").
		Values(teacher.Fio, teacher.PhoneNumber, teacher.AccountId).
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

func (repo *TeacherRepository) Read(ctx context.Context) ([]entities.Teacher, error) {
	sql, args, err := repo.builder.
		Select("id, fio, phone_number", "account_id").
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
		var student entities.Teacher
		err = rows.Scan(
			&student.Id,
			&student.Fio,
			&student.PhoneNumber,
			&student.AccountId,
		)
		if err != nil {
			return nil, SqlScanError
		}

		teachers = append(teachers, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return teachers, nil
}

func (repo *TeacherRepository) ReadById(ctx context.Context, id int) (entities.Teacher, error) {
	var fio, group, phoneNumber string
	var accountId int

	sql, args, err := repo.builder.
		Select("fio", "group_name", "phone_number", "account_id").
		From("teachers").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Teacher{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&fio,
		&group,
		&phoneNumber,
		&accountId,
	)
	if err != nil {
		return entities.Teacher{}, SqlReadError
	}

	return entities.Teacher{Id: id, Fio: fio, PhoneNumber: phoneNumber, AccountId: accountId}, nil
}

func (repo *TeacherRepository) ReadByAccountId(ctx context.Context, accountId int) (entities.Teacher, error) {
	var fio, phoneNumber string
	var id int

	sql, args, err := repo.builder.
		Select("id", "fio", "phone_number").
		From("teachers").
		Where(squirrel.Eq{"account_id": accountId}).
		ToSql()

	if err != nil {
		return entities.Teacher{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&id,
		&fio,
		&phoneNumber,
	)
	if err != nil {
		return entities.Teacher{}, SqlReadError
	}

	return entities.Teacher{Id: id, Fio: fio, PhoneNumber: phoneNumber, AccountId: accountId}, nil
}

func (repo *TeacherRepository) Update(ctx context.Context, id int, updates map[string]any) (entities.Teacher, error) {
	sql, args, err := repo.builder.
		Update("teachers").
		Where(squirrel.Eq{"id": id}).
		SetMap(updates).
		Suffix("RETURNING id, fio, phone_number, account_id").
		ToSql()

	if err != nil {
		return entities.Teacher{}, SqlStatementError
	}

	var student entities.Teacher
	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&student.Id,
		&student.Fio,
		&student.PhoneNumber,
		&student.AccountId,
	)

	if err != nil {
		fmt.Println(err)
		return entities.Teacher{}, SqlUpdateError
	}

	return student, nil
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
