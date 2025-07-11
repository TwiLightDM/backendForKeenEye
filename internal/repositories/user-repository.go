package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewUserRepository(pool *pgxpool.Pool, builder squirrel.StatementBuilderType) *UserRepository {
	return &UserRepository{pool: pool, builder: builder}
}

func (repo *UserRepository) Create(ctx context.Context, user entities.User) (int, error) {
	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	sql, args, err := repo.builder.
		Insert("users").
		Columns("login", "password", "salt", "role").
		Values(user.Login, user.Password, user.Salt, user.Role).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, SqlStatementError
	}

	var newID int
	err = tx.QueryRow(ctx, sql, args...).Scan(&newID)
	if err != nil {
		return 0, SqlInsertError
	}

	var table string
	switch user.Role {
	case "student":
		table = "students"
	case "teacher":
		table = "teachers"
	case "admin":
		table = "admins"
	}

	sql, args, err = repo.builder.
		Insert(table).
		Columns("id").
		Values(newID).
		ToSql()

	if err != nil {
		return 0, SqlStatementError
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, SqlInsertError
	}

	return newID, nil
}

func (repo *UserRepository) ReadByLogin(ctx context.Context, login string) (entities.User, error) {
	var id int
	var password, salt, role string
	sql, args, err := repo.builder.
		Select("id", "password", "salt", "role").
		From("users").
		Where(squirrel.Eq{"login": login}).
		ToSql()

	if err != nil {
		return entities.User{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&id,
		&password,
		&salt,
		&role,
	)
	if err != nil {
		return entities.User{}, SqlReadError
	}

	return entities.User{Id: id, Login: login, Password: password, Salt: salt, Role: role}, nil
}

func (repo *UserRepository) ReadById(ctx context.Context, id int) (entities.User, error) {
	var login, password, salt, role string
	sql, args, err := repo.builder.
		Select("login", "password", "salt", "role").
		From("users").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.User{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&login,
		&password,
		&salt,
		&role,
	)
	if err != nil {
		return entities.User{}, SqlReadError
	}

	return entities.User{Id: id, Login: login, Password: password, Salt: salt, Role: role}, nil
}
