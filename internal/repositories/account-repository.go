package repositories

import (
	"backendForKeenEye/internal/entities"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	pool    *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewAccountRepository(pool *pgxpool.Pool, builder squirrel.StatementBuilderType) *AccountRepository {
	return &AccountRepository{pool: pool, builder: builder}
}

func (repo *AccountRepository) Create(ctx context.Context, Account entities.Account) (int, error) {
	sql, args, err := repo.builder.
		Insert("accounts").
		Columns("login", "password", "salt").
		Values(Account.Login, Account.Password, Account.Salt).
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

func (repo *AccountRepository) ReadByLogin(ctx context.Context, login string) (entities.Account, error) {
	var id int
	var password, salt string
	sql, args, err := repo.builder.
		Select("id", "password", "salt").
		From("accounts").
		Where(squirrel.Eq{"login": login}).
		ToSql()

	if err != nil {
		return entities.Account{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&id,
		&password,
		&salt,
	)
	if err != nil {
		return entities.Account{}, SqlReadError
	}

	return entities.Account{Id: id, Login: login, Password: password, Salt: salt}, nil
}

func (repo *AccountRepository) ReadById(ctx context.Context, id int) (entities.Account, error) {
	var login, password, salt string
	sql, args, err := repo.builder.
		Select("login", "password", "salt").
		From("accounts").
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		return entities.Account{}, SqlStatementError
	}

	err = repo.pool.QueryRow(ctx, sql, args...).Scan(
		&login,
		&password,
		&salt,
	)
	if err != nil {
		return entities.Account{}, SqlReadError
	}

	return entities.Account{Id: id, Login: login, Password: password, Salt: salt}, nil
}
