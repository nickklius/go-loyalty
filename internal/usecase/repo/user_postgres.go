package repo

import (
	"context"
	"fmt"

	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/storage/postgres"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

type UserRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) Store(ctx context.Context, u entity.User) error {
	sql, args, err := r.Builder.
		Insert("gophermart.users").
		Columns("login, password").
		Values(u.Login, u.Password).
		ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo - Store - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return usecase.ParseError(err)
	}

	return nil
}

func (r *UserRepo) Check(ctx context.Context, u entity.User) (entity.User, error) {
	var user entity.User

	sql, args, err := r.Builder.
		Select("login, password").
		From("gophermart.users").
		Where("login = ? AND password = ?", u.Login, u.Password).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("UserRepo - Check - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&user.Login, &user.Password)
	if err != nil {
		return user, usecase.ParseError(err)
	}

	return user, nil

}
