package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/nickklius/go-loyalty/internal/entity"
	"github.com/nickklius/go-loyalty/internal/storage/postgres"
	"github.com/nickklius/go-loyalty/internal/usecase"
)

type Repository struct {
	*postgres.Postgres
}

func NewPostgresRepository(pg *postgres.Postgres) *Repository {
	return &Repository{pg}
}

func (r *Repository) StoreUser(ctx context.Context, user entity.User) error {
	sql, args, err := r.Builder.
		Insert("users").
		Columns("login, password").
		Values(user.Login, user.Password).
		ToSql()
	if err != nil {
		return fmt.Errorf("repo - StoreUser - r.Builder: %w", err)
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
		return user, fmt.Errorf("repo - CheckUser - r.Builder: %w", err)
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
		return fmt.Errorf("repo - StoreOrder - r.Builder: %w", err)
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

func (r *Repository) GetOrders(ctx context.Context, userID string) ([]entity.Order, error) {
	var orders []entity.Order

	sql, args, err := r.Builder.
		Select("number, status, accrual, uploaded_at").
		From("orders").
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return orders, fmt.Errorf("repo - getOrders - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return orders, err
	}

	for rows.Next() {
		var order entity.Order

		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return orders, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return orders, err
	}

	return orders, nil
}

func (r *Repository) getOrder(ctx context.Context, number string) (entity.Order, error) {
	var order entity.Order

	sql, args, err := r.Builder.
		Select("id, user_id, number, status, uploaded_at, accrual").
		From("orders").
		Where("number = ?", number).
		ToSql()
	if err != nil {
		return order, fmt.Errorf("repo - getOrder - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt, &order.Accrual)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r *Repository) GetBalance(ctx context.Context, userID string) (entity.UserBalance, error) {
	var ub entity.UserBalance

	sql, args, err := r.Builder.
		Select("balance, spent").
		From("users").
		Where("id = ?", userID).
		ToSql()
	if err != nil {
		return ub, fmt.Errorf("repo - GetBalance - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&ub.Balance, &ub.Spent)
	if err != nil {
		return ub, err
	}

	return ub, nil
}

func (r *Repository) Withdraw(ctx context.Context, withdraw entity.Withdraw) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Select("balance, spent").
		From("users").
		Where("id = ?", withdraw.UserID).
		ToSql()
	if err != nil {
		return fmt.Errorf("repo - GetCurrentBalance - r.Builder: %w", err)
	}

	row := tx.QueryRow(ctx, sql, args...)

	var balance, spent float64

	err = row.Scan(&balance, &spent)
	if err != nil {
		return fmt.Errorf("no row error %w", err)
	}

	if withdraw.Sum > balance {
		return usecase.ErrDBNotEnoughBalanceForWithdrawal
	}

	sql, args, err = r.Builder.
		Update("users").
		Set("balance", balance-withdraw.Sum).
		Set("spent", spent+withdraw.Sum).
		Where("id = ?", withdraw.UserID).
		ToSql()
	if err != nil {
		return fmt.Errorf("repo - UpdateUserBalance - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	withdraw.Status = "PROCESSED"
	withdraw.ProcessedAt = time.Now()

	sql, args, err = r.Builder.
		Insert("withdrawals").
		Columns("user_id, order_num, sum, status, processed_at").
		Values(withdraw.UserID, withdraw.OrderID, withdraw.Sum, withdraw.Status, withdraw.ProcessedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("repo - StoreWithdraw - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return usecase.ErrDBDuplicatedEntry
		}
	}

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetWithdrawals(ctx context.Context, userID string) ([]entity.Withdraw, error) {
	var withdrawals []entity.Withdraw

	sql, args, err := r.Builder.
		Select("order_num, sum, processed_at").
		From("withdrawals").
		Where("user_id = ?", userID).
		ToSql()
	if err != nil {
		return withdrawals, fmt.Errorf("repo - GetWithdrawals - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return withdrawals, err
	}

	for rows.Next() {
		var withdraw entity.Withdraw

		err = rows.Scan(&withdraw.OrderID, &withdraw.Sum, &withdraw.ProcessedAt)
		if err != nil {
			return withdrawals, err
		}

		withdrawals = append(withdrawals, withdraw)
	}

	if err = rows.Err(); err != nil {
		return withdrawals, err
	}

	return withdrawals, nil
}

func (r *Repository) getUserIDByOrder(ctx context.Context, order string) (string, error) {
	var userID string

	sql, args, err := r.Builder.
		Select("user_id").
		From("orders").
		Where("number = ?", order).
		ToSql()
	if err != nil {
		return userID, fmt.Errorf("repo - getUserIDByOrder - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

func (r *Repository) getUserByID(ctx context.Context, userID string) (entity.User, error) {
	var user entity.User

	sql, args, err := r.Builder.
		Select("id, balance").
		From("users").
		Where("id = ?", userID).
		ToSql()
	if err != nil {
		return user, fmt.Errorf("repo - getUserByID - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	err = row.Scan(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, order entity.Order) error {
	userID, err := r.getUserIDByOrder(ctx, order.Number)
	if err != nil {
		return err
	}

	user, err := r.getUserByID(ctx, userID)
	if err != nil {
		return err
	}

	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql, args, err := r.Builder.
		Update("orders").
		Set("accrual", order.Accrual).
		Set("status", order.Status).
		Where("number = ?", order.Number).
		ToSql()
	if err != nil {
		return fmt.Errorf("repo - UpdateOrderStatus - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	sql, args, err = r.Builder.
		Update("users").
		Set("balance", user.UserBalance.Balance+order.Accrual).
		Where("id = ?", user.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("repo - UpdateOrderStatus / UserBalance - r.Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
