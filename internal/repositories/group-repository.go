package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewGroupRepository(pool *pgxpool.Pool, builder squirrel.StatementBuilderType) *GroupRepository {
	return &GroupRepository{pool: pool, builder: builder}
}

func (repo *GroupRepository) Create(ctx context.Context, group entities.Group) (int, error) {
	var teacherId any
	if group.TeacherId == 0 {
		teacherId = nil
	} else {
		teacherId = group.TeacherId
	}

	sql, args, err := repo.builder.
		Insert("groups").
		Columns("name", "teacher_id").
		Values(group.Name, teacherId).
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

func (repo *GroupRepository) Read(ctx context.Context) ([]entities.Group, error) {
	var id int
	var name sql.NullString
	var teacherId sql.NullInt32
	sql, args, err := repo.builder.
		Select("id", "name", "teacher_id").
		From("groups").
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

	var groups []entities.Group
	for rows.Next() {
		err = rows.Scan(
			&id,
			&name,
			&teacherId,
		)
		if err != nil {
			return nil, SqlScanError
		}

		group := entities.Group{
			Id:        id,
			Name:      validateString(name),
			TeacherId: validateInt(teacherId),
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return groups, nil
}

func (repo *GroupRepository) ReadById(ctx context.Context, id int) (entities.Group, error) {
	var name sql.NullString
	var teacherId sql.NullInt32

	sql, args, err := repo.builder.
		Select("name", "teacher_id").
		From("groups").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Group{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&name,
		&teacherId,
	)
	if err != nil {
		return entities.Group{}, SqlReadError
	}

	return entities.Group{Id: id, Name: validateString(name), TeacherId: validateInt(teacherId)}, nil
}

func (repo *GroupRepository) Update(ctx context.Context, id int, updates map[string]any) (entities.Group, error) {
	var name sql.NullString
	var teacherId sql.NullInt32
	sql, args, err := repo.builder.
		Update("groups").
		Where(squirrel.Eq{"id": id}).
		SetMap(updates).
		Suffix("RETURNING name, teacher_id").
		ToSql()

	if err != nil {
		return entities.Group{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&name,
		&teacherId,
	)

	if err != nil {
		return entities.Group{}, SqlUpdateError
	}

	return entities.Group{Id: id, Name: validateString(name), TeacherId: validateInt(teacherId)}, nil
}

func (repo *GroupRepository) SoftDelete(ctx context.Context, id int) error {
	sql, args, err := repo.builder.
		Update("groups").
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
