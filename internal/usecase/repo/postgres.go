package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/storage/postgres"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

type Repository struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repository {
	return &Repository{pg}
}

func (r *Repository) StoreUser(ctx context.Context, u entity.User) error {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("login, password").
		Values(u.Login, u.Password).
		ToSql()
	if err != nil {
		return fmt.Errorf("Repo - StoreUser - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return usecase.ParseError(err)
	}

	return nil
}

func (r *Repository) CheckUser(ctx context.Context, user entity.User) (entity.User, error) {
	var u entity.User

	sql, args, err := r.Builder.
		Select("id, login, password").
		From("users").
		Where("login = ? AND password = ?", user.Login, user.Password).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("Repo - CheckUser - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return u, usecase.ParseError(err)
	}

	return u, nil
}

func (r *Repository) StoreOrder(ctx context.Context, order entity.Order) error {
	sql, args, err := r.Builder.
		Insert("orders").
		Columns("user_id, number, status").
		Values(order.UserID, order.Number, order.Status).
		ToSql()

	if err != nil {
		return fmt.Errorf("Repo - StoreOrder - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			existing, err := r.getOrder(ctx, order.Number)
			if err != nil {
				return err
			}
			if existing.UserID == order.UserID {
				return usecase.ErrDBDuplicateOrderByUserItself
			}
			return usecase.ErrDBDuplicateOrder
		}
		return err
	}

	return nil
}

func (r *Repository) getOrder(ctx context.Context, number string) (entity.Order, error) {
	var order entity.Order

	sql, args, err := r.Builder.
		Select("id, user_id, number, status, uploaded_at, accrual").
		From("orders").
		Where("number = ?", number).
		ToSql()
	if err != nil {
		return order, fmt.Errorf("Repo - getOrder - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt, &order.Accrual)
	if err != nil {
		return order, err
	}

	return order, nil
}
